package handler

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"math/big"
	"token-payment/internal/chain"
	"token-payment/internal/config"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/types"
	"token-payment/internal/utils"
)

// BuildWithdrawTransactions
//
//	@Description: 生成交易
//	@param ctx
//	@param ch
func BuildWithdrawTransactions(ctx context.Context, ch *sqlmodel.Chain) {
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
		BuildWithdrawSendTxList(ctx, ch, &appChain)
	}
}

func BuildWithdrawSendTxList(ctx context.Context, ch *sqlmodel.Chain, appChain *sqlmodel.ApplicationChain) {
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
	), 0, int(ch.Concurrent), orderQ.ID.Asc())
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
			_ = BuildWithdrawSendTx(ctx, ch, &order, appChain)
		}
	}
}

func BuildWithdrawSendTx(ctx context.Context, ch *sqlmodel.Chain, order *sqlmodel.ApplicationWithdrawOrder, appChain *sqlmodel.ApplicationChain) (err error) {
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
	nonce, err := GetTransferNonce(ctx, ch, appChain.HotWallet)
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
		Nonce:           nonce,
	}
	client, err := GetChainRpcClient(ctx, ch)
	if err != nil {
		return
	}
	err = client.GenerateTransaction(ctx, &transferOrder)
	if err != nil {
		return
	}
	if transferOrder.GasPrice.Int64() < ch.GasPrice { // 不能低于预设的gasPrice
		transferOrder.GasPrice = big.NewInt(ch.GasPrice)
	}
	err = dao.GetDB(ctx).Transaction(func(tx *gorm.DB) (txErr error) {
		c := dao.CtxWithTransaction(ctx, tx)
		sendTx := &sqlmodel.ChainSendTx{
			ApplicationID:   order.ApplicationID,
			SerialNo:        order.SerialNo,
			ChainSymbol:     ch.ChainSymbol,
			ContractAddress: order.ContractAddress,
			Symbol:          order.Symbol,
			FromAddress:     appChain.HotWallet,
			ToAddress:       order.ToAddress,
			Value:           order.Value,
			GasPrice:        transferOrder.GasPrice.Int64(),
			TokenID:         order.TokenID,
			TxHash:          transferOrder.TxHash,
			Nonce:           int64(transferOrder.Nonce),
			Hook:            order.Hook,
			TransferType:    int32(types.TransferTypeOut),
			CreateAt:        order.CreateAt,
		}
		_, txErr = dao.AddChainSendTx(ctx, sendTx)
		if txErr != nil {
			return
		}
		order.SendTxID = sendTx.ID
		order.Generated = 1
		_, txErr = dao.UpdateApplicationWithdrawOrder(c, order)
		return
	})
	return
}
