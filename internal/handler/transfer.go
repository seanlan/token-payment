package handler

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"math"
	"math/big"
	"time"
	"token-payment/internal/chain"
	"token-payment/internal/config"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/utils"
)

// LoadTransferNonce
//
//	@Description: 加载转账nonce
//	@param ctx
//	@param ch
//	@param appChain
//	@return err
func LoadTransferNonce(ctx context.Context, ch *sqlmodel.Chain, appChain *sqlmodel.ApplicationChain) (err error) {
	//获取nonce
	var (
		nonceCacheKey        = fmt.Sprintf("transfer-nonce:%s:%s", appChain.ChainSymbol, appChain.HotWallet)
		duration             time.Duration
		orderQ               = sqlmodel.ApplicationWithdrawOrderColumns
		lastOrder            sqlmodel.ApplicationWithdrawOrder
		chainNonce, newNonce uint64
	)
	duration, err = dao.Redis.TTL(ctx, nonceCacheKey).Result()
	if err != nil || duration > time.Minute*5 { // 5分钟内有效
		return
	}
	err = dao.FetchApplicationWithdrawOrder(ctx, &lastOrder, dao.And(
		orderQ.ChainSymbol.Eq(appChain.ChainSymbol),
		orderQ.ApplicationID.Eq(appChain.ApplicationID),
	), orderQ.Nonce.Desc())
	if err != nil && !errors.Is(err, dao.ErrNotFound) {
		return
	}
	client, err := GetChainRpcClient(ctx, ch)
	if err != nil {
		return
	}
	chainNonce, err = client.GetNonce(ctx, appChain.HotWallet)
	if err != nil {
		return
	}
	if lastOrder.Nonce > int64(chainNonce) {
		newNonce = uint64(lastOrder.Nonce)
	} else {
		newNonce = chainNonce
	}
	err = dao.Redis.Set(ctx, nonceCacheKey, newNonce, time.Minute*30).Err() // 30分钟过期
	return
}

// GenerateTransactions
//
//	@Description: 生成交易
//	@param ctx
//	@param ch
func GenerateTransactions(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		appChainQ = sqlmodel.ApplicationChainColumns
		appChains = make([]sqlmodel.ApplicationChain, 0)
		err       error
	)
	err = dao.FetchAllApplicationChain(ctx, &appChains, dao.And(
		appChainQ.ChainSymbol.Eq(ch.ChainSymbol),
	), 0, 0)
	if err != nil {
		return
	}
	for _, appChain := range appChains {
		err = LoadTransferNonce(ctx, ch, &appChain)
		if err != nil {
			continue
		}
		GenerateAppChainTransactions(ctx, ch, &appChain)
	}
}

// GenerateAppChainTransactions
//
//	@Description: 按照指定链，指定应用链生成交易
//	@param ctx
//	@param ch
//	@param appChain
func GenerateAppChainTransactions(ctx context.Context, ch *sqlmodel.Chain, appChain *sqlmodel.ApplicationChain) {
	var (
		orderQ          = sqlmodel.ApplicationWithdrawOrderColumns
		orders          = make([]sqlmodel.ApplicationWithdrawOrder, 0)
		tokenQ          = sqlmodel.ChainTokenColumns
		tokenContracts  = make([]string, 0)
		tokens          = make([]sqlmodel.ChainToken, 0)
		tokenMap        = make(map[string]sqlmodel.ChainToken)
		tokenConsumeMap = make(map[string]float64)
		balanceEnough   = true
	)
	err := dao.FetchAllApplicationWithdrawOrder(ctx, &orders, dao.And(
		orderQ.ChainSymbol.Eq(ch.ChainSymbol),
		orderQ.ApplicationID.Eq(appChain.ApplicationID),
		orderQ.Generated.Eq(0), // 未生成交易
	), 0, 0, orderQ.ID.Asc())
	if err != nil || len(orders) == 0 {
		return
	}
	for _, order := range orders {
		tokenContracts = append(tokenContracts, order.ContractAddress)
	}
	err = dao.FetchAllChainToken(ctx, &tokens, dao.And(
		tokenQ.ChainSymbol.Eq(ch.ChainSymbol),
		tokenQ.ContractAddress.In(tokenContracts),
	), 0, 0)
	if err != nil || len(tokens) == 0 {
		return
	}
	for _, token := range tokens {
		tokenMap[token.ContractAddress] = token
	}
	// 计算每个合约的消耗
	for _, order := range orders {
		token, ok := tokenMap[order.ContractAddress]
		if !ok {
			continue
		}
		tokenConsumeMap[order.ContractAddress] += order.Value * math.Pow10(int(token.Decimals))
		tokenConsumeMap[""] += float64(ch.GasPrice) * token.GasFee
	}
	// 查询热钱包余额
	client, err := GetChainRpcClient(ctx, ch)
	if err != nil {
		return
	}
	for contract, consume := range tokenConsumeMap {
		balance, _err := client.GetBalance(ctx, appChain.HotWallet, contract)
		if _err != nil {
			return
		}
		_consume := big.NewInt(int64(consume))
		zap.S().Infof("balance: %s, consume: %s", balance.String(), _consume.String())
		if balance.Cmp(_consume) < 0 { // 余额不足
			zap.S().Warnf("balance not enough, hotWallet: %s, contract: %s, balance: %s, consume: %s",
				appChain.HotWallet, contract, balance.String(), _consume.String())
			balanceEnough = false
			continue
		}
	}
	if balanceEnough {
		for _, order := range orders {
			_ = GenerateTransaction(ctx, ch, &order, appChain)
		}
	}
}

// GenerateTransaction
//
//	@Description: 生成交易
//	@param ctx
//	@param ch
//	@param order
//	@param appChain
//	@return err
func GenerateTransaction(ctx context.Context, ch *sqlmodel.Chain, order *sqlmodel.ApplicationWithdrawOrder, appChain *sqlmodel.ApplicationChain) (err error) {
	var (
		nonceCacheKey = fmt.Sprintf("transfer-nonce:%s:%s", appChain.ChainSymbol, appChain.HotWallet)
		tokenQ        = sqlmodel.ChainTokenColumns
		token         sqlmodel.ChainToken
		addressQ      = sqlmodel.ChainAddressColumns
		address       sqlmodel.ChainAddress
	)
	nonce, err := dao.Redis.Incr(ctx, nonceCacheKey).Result()
	if err != nil {
		// 生成nonce失败 重置缓存
		dao.Redis.Del(ctx, nonceCacheKey)
		return
	}
	err = dao.FetchChainToken(ctx, &token,
		dao.And(
			tokenQ.ChainSymbol.Eq(ch.ChainSymbol),
			tokenQ.ContractAddress.Eq(order.ContractAddress),
			tokenQ.Symbol.Eq(order.Symbol)))
	if err != nil {
		return
	}
	err = dao.FetchChainAddress(ctx, &address, dao.And(
		addressQ.ChainSymbol.Eq(ch.ChainSymbol),
		addressQ.Address.Eq(appChain.HotWallet),
		addressQ.ApplicationID.Eq(order.ApplicationID),
	))
	if err != nil {
		return
	}
	pk, err := utils.AesDecrypt(address.EncKey, config.C.Secret)
	if err != nil {
		return
	}
	transferOrder := chain.TransferOrder{
		From:            appChain.HotWallet,
		FromPrivateKey:  pk,
		To:              order.ToAddress,
		ContractAddress: order.ContractAddress,
		Value:           big.NewInt(int64(order.Value * math.Pow10(int(token.Decimals)))),
		TokenID:         big.NewInt(order.TokenID),
		Gas:             uint64(ch.Gas),
		Nonce:           uint64(nonce - 1),
	}
	client, err := GetChainRpcClient(ctx, ch)
	if err != nil {
		return
	}
	err = client.GenerateTransaction(ctx, &transferOrder)
	if err != nil {
		return
	}
	order.TxHash = transferOrder.TxHash
	order.GasPrice = transferOrder.GasPrice.Int64()
	order.Nonce = int64(transferOrder.Nonce)
	order.Generated = 1
	order.TransferNextTime = time.Now().Unix()
	if order.GasPrice < ch.GasPrice { // 不能低于预设的gasPrice
		order.GasPrice = ch.GasPrice
	}
	_, err = dao.UpdateApplicationWithdrawOrder(ctx, order)
	return
}

// SendTransactions
//
//	@Description: 发送交易
//	@param ctx
//	@param ch
func SendTransactions(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		orderQ      = sqlmodel.ApplicationWithdrawOrderColumns
		orders      = make([]sqlmodel.ApplicationWithdrawOrder, 0)
		appIDs      = make([]int64, 0)
		appChainQ   = sqlmodel.ApplicationChainColumns
		appChains   = make([]sqlmodel.ApplicationChain, 0)
		appChainMap = make(map[int64]sqlmodel.ApplicationChain)
	)
	err := dao.FetchAllApplicationWithdrawOrder(ctx, &orders, dao.And(
		orderQ.ChainSymbol.Eq(ch.ChainSymbol),
		orderQ.Generated.Eq(1),                         // 已生成交易
		orderQ.TransferSuccess.Eq(0),                   // 未发送成功
		orderQ.TransferNextTime.Lte(time.Now().Unix()), // 下次发送时间小于当前时间
	), 0, int(ch.Concurrent), orderQ.TransferFailedTimes.Asc(), orderQ.ID.Asc())
	if len(orders) == 0 {
		return
	}
	for _, order := range orders {
		appIDs = append(appIDs, order.ApplicationID)
	}
	if len(appIDs) == 0 {
		return
	}
	err = dao.FetchAllApplicationChain(ctx, &appChains, dao.And(
		appChainQ.ChainSymbol.Eq(ch.ChainSymbol),
		appChainQ.ApplicationID.In(appIDs),
	), 0, 0)
	if err != nil {
		return
	}
	for _, appChain := range appChains {
		appChainMap[appChain.ApplicationID] = appChain
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	for _, order := range orders {
		appChain, ok := appChainMap[order.ApplicationID]
		if !ok {
			continue
		}
		_ = SendTransaction(ctx, ch, &order, &appChain)
	}
}

// SendTransaction
//
//	@Description: 发送交易
//	@param ctx
//	@param ch
//	@param order
//	@param appChain
//	@return err
func SendTransaction(ctx context.Context, ch *sqlmodel.Chain, order *sqlmodel.ApplicationWithdrawOrder, appChain *sqlmodel.ApplicationChain) (err error) {
	var (
		tokenQ   = sqlmodel.ChainTokenColumns
		token    sqlmodel.ChainToken
		addressQ = sqlmodel.ChainAddressColumns
		address  sqlmodel.ChainAddress
	)
	err = dao.FetchChainToken(ctx, &token,
		dao.And(
			tokenQ.ChainSymbol.Eq(ch.ChainSymbol),
			tokenQ.ContractAddress.Eq(order.ContractAddress),
			tokenQ.Symbol.Eq(order.Symbol)))
	if err != nil {
		return
	}
	err = dao.FetchChainAddress(ctx, &address, dao.And(
		addressQ.ChainSymbol.Eq(ch.ChainSymbol),
		addressQ.Address.Eq(appChain.HotWallet),
		addressQ.ApplicationID.Eq(order.ApplicationID),
	))
	if err != nil {
		return
	}
	pk, err := utils.AesDecrypt(address.EncKey, config.C.Secret)
	if err != nil {
		return
	}
	transferOrder := chain.TransferOrder{
		TxHash:          order.TxHash,
		From:            appChain.HotWallet,
		FromPrivateKey:  pk,
		To:              order.ToAddress,
		ContractAddress: order.ContractAddress,
		Value:           big.NewInt(int64(order.Value * math.Pow10(int(token.Decimals)))),
		TokenID:         big.NewInt(order.TokenID),
		Gas:             uint64(ch.Gas),
		GasPrice:        big.NewInt(order.GasPrice),
		Nonce:           uint64(order.Nonce),
	}
	client, err := GetChainRpcClient(ctx, ch)
	if err != nil {
		return
	}
	txHash, _err := client.Transfer(ctx, &transferOrder)
	if _err != nil {
		zap.S().Warnf("transfer failed, err: %v", _err)
		order.TransferSuccess = 0
		order.TransferFailedTimes++
		order.TransferNextTime = time.Now().Unix() + 60*5 // 5分钟后重试
	} else {
		order.TransferSuccess = 1
	}
	order.TxHash = txHash                // 更新txHash
	order.TransferAt = time.Now().Unix() // 记录发送时间 方便查询pending时间过长的交易 防止交易卡住 后面可以定时检索并处理
	_, err = dao.UpdateApplicationWithdrawOrder(ctx, order)
	return
}
