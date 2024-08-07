package handler

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"math/big"
	"sync"
	"token-payment/internal/chain"
	"token-payment/internal/config"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/types"
	"token-payment/internal/utils"
)

type AddressArrange struct {
	ChainSymbol   string  `json:"chain_symbol"`
	ApplicationID int64   `json:"application_id"`
	Address       string  `json:"address"`
	Amount        float64 `json:"amount"`
}

// ScanArrangeTransactions
//
//	@Description: 检查需要整理的交易
//	@param ctx
//	@param ch
func ScanArrangeTransactions(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		tokenQ = sqlmodel.ChainTokenColumns
		tokens []sqlmodel.ChainToken
		wg     sync.WaitGroup
	)
	err := dao.FetchAllChainToken(ctx, &tokens, dao.And(
		tokenQ.ChainSymbol.Eq(ch.ChainSymbol),
		tokenQ.ArrangeSwitch.Eq(1),
	), 0, 0)
	if err != nil {
		return
	}
	for _, token := range tokens {
		// 分不同的token进行扫描
		wg.Add(1)
		go func(token sqlmodel.ChainToken) {
			defer wg.Done()
			ScanTokenArrangeTransaction(ctx, ch, &token)
		}(token)
	}
	wg.Wait()
	return
}

// ScanTokenArrangeTransaction
//
//	@Description: 扫描需要整理交易
//	@param ctx
//	@param ch
//	@param token
func ScanTokenArrangeTransaction(ctx context.Context, ch *sqlmodel.Chain, token *sqlmodel.ChainToken) {
	var (
		txQ         = sqlmodel.ChainTxColumns
		arrangeList = make([]AddressArrange, 0)
		err         error
		wg          sync.WaitGroup
	)
	err = dao.GetDB(ctx).Model(&sqlmodel.ChainTx{}).Select(
		"application_id, to_address as address, sum(value) as amount").
		Where(
			txQ.ChainSymbol.Eq(token.ChainSymbol),
			txQ.Symbol.Eq(token.Symbol),
			txQ.TransferType.Eq(types.TransferTypeIn),
			txQ.Confirmed.Eq(1), // 已确认的交易
			txQ.Arranged.Eq(0),  // 未整理的交易
			txQ.Removed.Eq(0),
		).
		Group("application_id, address").
		Order("amount desc").
		Limit(int(ch.Concurrent)).
		Find(&arrangeList).Error
	if err != nil {
		return
	}
	for _, arrange := range arrangeList {
		if arrange.Amount < token.Threshold { // 低于阈值 不生成交易
			continue
		}
		wg.Add(1)
		go func(arrange AddressArrange) {
			defer wg.Done()
			InsertArrangeTransaction(ctx, ch, token, &arrange)
		}(arrange)
	}
	wg.Wait()
}

// InsertArrangeTransaction
//
//	@Description: 插入新的整理交易
//	@param ctx
//	@param ch
//	@param token
//	@param arrange
func InsertArrangeTransaction(ctx context.Context, ch *sqlmodel.Chain, token *sqlmodel.ChainToken, arrange *AddressArrange) {
	var (
		txQ     = sqlmodel.ChainTxColumns
		finalTx sqlmodel.ChainTx
		amount  float64
		effect  int64
		err     error
	)
	err = dao.FetchChainTx(ctx, &finalTx, dao.And(
		txQ.ChainSymbol.Eq(token.ChainSymbol),
		txQ.Symbol.Eq(token.Symbol),
		txQ.TransferType.Eq(types.TransferTypeIn),
		txQ.Confirm.Eq(ch.Confirm), // 确认数符合链的确认数量
		txQ.ToAddress.Eq(arrange.Address),
		txQ.Removed.Eq(0), // 未删除的交易
	), txQ.ID.Desc())
	if err != nil || finalTx.ID == 0 {
		return
	}
	// 计算交易额
	amount, err = dao.SumChainTx(ctx, txQ.Value, dao.And(
		txQ.ChainSymbol.Eq(ch.ChainSymbol),
		txQ.Symbol.Eq(token.Symbol),
		txQ.ToAddress.Eq(arrange.Address),
		txQ.TransferType.Eq(types.TransferTypeIn),
		txQ.ID.Lte(finalTx.ID),
	))
	if err != nil {
		return
	}
	_ = dao.GetDB(ctx).Transaction(func(tx *gorm.DB) (txErr error) {
		c := dao.CtxWithTransaction(ctx, tx)
		// 更改交易状态
		effect, txErr = dao.UpdatesChainTx(c, dao.And(
			txQ.ChainSymbol.Eq(ch.ChainSymbol),
			txQ.Symbol.Eq(token.Symbol),
			txQ.ToAddress.Eq(arrange.Address),
			txQ.TransferType.Eq(types.TransferTypeIn),
			txQ.ID.Lte(finalTx.ID),
		), dao.M{
			txQ.Arranged.FieldName: 1,
		})
		if txErr != nil || effect == 0 {
			err = errors.New("update chain tx error")
			return
		}
		// 生成整理交易
		_, txErr = dao.AddApplicationArrangeTx(c, &sqlmodel.ApplicationArrangeTx{
			ApplicationID:   arrange.ApplicationID,
			ChainSymbol:     ch.ChainSymbol,
			ContractAddress: token.ContractAddress,
			Symbol:          token.Symbol,
			FromAddress:     arrange.Address,
			Value:           amount,
		})
		return
	})
}

// CheckArrangeTxFee
//
//	@Description: 检测整理交易手续费
//	@param ctx
//	@param ch
//	@return {}
func CheckArrangeTxFee(ctx context.Context, ch *sqlmodel.Chain) {
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
		CheckArrangeTxFeeByApplication(ctx, ch, &appChain)
	}
}

// CheckArrangeTxFeeByApplication
//
//	@Description: 检测整理交易手续费
//	@param ctx
//	@param ch
//	@param appChain
func CheckArrangeTxFeeByApplication(ctx context.Context, ch *sqlmodel.Chain, appChain *sqlmodel.ApplicationChain) {
	var (
		arrangeQ       = sqlmodel.ApplicationArrangeTxColumns
		arrangeTxs     = make([]sqlmodel.ApplicationArrangeTx, 0)
		arrangeTxIDs   = make([]int64, 0)
		tokenQ         = sqlmodel.ChainTokenColumns
		tokenContracts = make([]string, 0)
		tokens         = make([]sqlmodel.ChainToken, 0)
		tokenMap       = make(map[string]sqlmodel.ChainToken)
		addressFee     = make(map[string]float64)
		sumFee         float64
		balanceEnough  = true
	)
	err := dao.FetchAllApplicationArrangeTx(ctx, &arrangeTxs, dao.And(
		arrangeQ.ChainSymbol.Eq(ch.ChainSymbol),
		arrangeQ.ApplicationID.Eq(appChain.ApplicationID),
		arrangeQ.ArrangeFeeTxID.Eq(0), // 未计算手续费的交易
	), 0, int(ch.Concurrent), arrangeQ.ID.Asc())
	if err != nil || len(arrangeTxs) == 0 {
		return
	}
	for _, order := range arrangeTxs {
		tokenContracts = append(tokenContracts, order.ContractAddress)
		arrangeTxIDs = append(arrangeTxIDs, order.ID)
	}
	tokenContracts = append(tokenContracts, "")
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
	mainToken, _ok := tokenMap[""]
	if !_ok {
		return
	}
	// 计算每个合约的消耗
	for _, order := range arrangeTxs {
		token, ok := tokenMap[order.ContractAddress]
		if !ok {
			continue
		}
		fee := float64(ch.GasPrice) * token.GasFee
		addressFee[order.FromAddress] += fee
		sumFee += fee
	}
	sumFee += float64(ch.GasPrice) * mainToken.GasFee * float64(len(addressFee))
	// 查询热钱包余额
	client, err := GetChainRpcClient(ctx, ch)
	if err != nil {
		return
	}
	balance, err := client.GetBalance(ctx, appChain.FeeWallet, "")
	if err != nil {
		return
	}
	consume := big.NewInt(int64(sumFee))
	zap.S().Infof("balance: %s, consume: %s", balance.String(), consume.String())
	if balance.Cmp(consume) < 0 { // 余额不足
		zap.S().Warnf("balance not enough, feeWallet: %s, contract: %s, balance: %s, consume: %s",
			appChain.FeeWallet, "", balance.String(), consume.String())
		balanceEnough = false
		return
	}
	if balanceEnough {
		for address := range addressFee {
			err = dao.GetDB(ctx).Transaction(func(tx *gorm.DB) (txErr error) {
				c := dao.CtxWithTransaction(ctx, tx)
				feeTx := &sqlmodel.ApplicationArrangeFeeTx{
					ApplicationID: appChain.ApplicationID,
					ChainSymbol:   ch.ChainSymbol,
					Symbol:        mainToken.Symbol,
					FromAddress:   appChain.FeeWallet,
					ToAddress:     address,
					Value:         addressFee[address] / math.Pow10(int(mainToken.Decimals)),
				}
				_, txErr = dao.AddApplicationArrangeFeeTx(c, feeTx)
				if txErr != nil {
					return
				}
				_, txErr = dao.UpdatesApplicationArrangeTx(c, dao.And(
					arrangeQ.ID.In(arrangeTxIDs),
					arrangeQ.FromAddress.Eq(address),
				), dao.M{
					arrangeQ.ArrangeFeeTxID.FieldName: feeTx.ID,
				})
				return
			})
		}
	}
}

func BuildArrangeTxs(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		arrangeQ   = sqlmodel.ApplicationArrangeTxColumns
		feeQ       = sqlmodel.ApplicationArrangeFeeTxColumns
		arrangeTxs = make([]sqlmodel.ApplicationArrangeTx, 0)
		err        error
	)
	err = dao.GetDB(ctx).Table(sqlmodel.TableNameApplicationArrangeTx).
		Select("application_arrange_tx.*").
		Joins("join application_arrange_fee_tx on application_arrange_tx.`arrange_fee_tx_id` = application_arrange_fee_tx.id").
		Where(dao.And(
			arrangeQ.Generated.Eq(0),
			feeQ.Confirmed.Eq(1))).
		Limit(int(ch.Concurrent)).
		Find(&arrangeTxs).Error
	if err != nil {
		return
	}
	for _, arrange := range arrangeTxs {
		err = BuildArrangeTx(ctx, ch, &arrange)
		zap.S().Warnf("BuildArrangeTx: %v", err)
	}

}

func BuildArrangeTx(ctx context.Context, ch *sqlmodel.Chain, arrangeTx *sqlmodel.ApplicationArrangeTx) (err error) {
	var (
		tokenQ   = sqlmodel.ChainTokenColumns
		token    sqlmodel.ChainToken
		addressQ = sqlmodel.ChainAddressColumns
		address  sqlmodel.ChainAddress
	)
	err = dao.FetchChainToken(ctx, &token,
		dao.And(
			tokenQ.ChainSymbol.Eq(ch.ChainSymbol),
			tokenQ.ContractAddress.Eq(arrangeTx.ContractAddress),
			tokenQ.Symbol.Eq(arrangeTx.Symbol)))
	if err != nil {
		return
	}
	err = dao.FetchChainAddress(ctx, &address, dao.And(
		addressQ.ChainSymbol.Eq(ch.ChainSymbol),
		addressQ.Address.Eq(arrangeTx.FromAddress),
		addressQ.ApplicationID.Eq(arrangeTx.ApplicationID),
	))
	if err != nil {
		return
	}
	pk, err := utils.AesDecrypt(address.EncKey, config.C.Secret)
	if err != nil {
		return
	}
	nonce, err := GetTransferNonce(ctx, ch, arrangeTx.FromAddress)
	if err != nil {
		return
	}
	transferOrder := chain.TransferOrder{
		From:            arrangeTx.FromAddress,
		FromPrivateKey:  pk,
		To:              arrangeTx.ToAddress,
		ContractAddress: arrangeTx.ContractAddress,
		Value:           big.NewInt(int64(arrangeTx.Value * math.Pow10(int(token.Decimals)))),
		TokenID:         big.NewInt(arrangeTx.TokenID),
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
			ApplicationID:   arrangeTx.ApplicationID,
			SerialNo:        arrangeTx.SerialNo,
			ChainSymbol:     ch.ChainSymbol,
			ContractAddress: arrangeTx.ContractAddress,
			Symbol:          arrangeTx.Symbol,
			FromAddress:     arrangeTx.FromAddress,
			ToAddress:       arrangeTx.ToAddress,
			Value:           arrangeTx.Value,
			GasPrice:        transferOrder.GasPrice.Int64(),
			TokenID:         arrangeTx.TokenID,
			TxHash:          transferOrder.TxHash,
			Nonce:           int64(transferOrder.Nonce),
			Hook:            arrangeTx.Hook,
			TransferType:    int32(types.TransferTypeArrange),
			CreateAt:        arrangeTx.CreateAt,
		}
		_, txErr = dao.AddChainSendTx(ctx, sendTx)
		if txErr != nil {
			return
		}
		arrangeTx.SendTxID = sendTx.ID
		arrangeTx.Generated = 1
		_, txErr = dao.UpdateApplicationArrangeTx(c, arrangeTx)
		return
	})
	return
}

func BuildArrangeFeeTxs(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		feeQ   = sqlmodel.ApplicationArrangeFeeTxColumns
		feeTxs = make([]sqlmodel.ApplicationArrangeFeeTx, 0)
		err    error
	)
	err = dao.FetchAllApplicationArrangeFeeTx(ctx, &feeTxs, dao.And(
		feeQ.ChainSymbol.Eq(ch.ChainSymbol),
	), 0, int(ch.Concurrent))
	if err != nil {
		return
	}
	for _, feeTx := range feeTxs {
		err = BuildArrangeFeeTx(ctx, ch, &feeTx)
		zap.S().Warnf("BuildArrangeTx: %v", err)
	}
}

func BuildArrangeFeeTx(ctx context.Context, ch *sqlmodel.Chain, feeTx *sqlmodel.ApplicationArrangeFeeTx) (err error) {
	var (
		tokenQ   = sqlmodel.ChainTokenColumns
		token    sqlmodel.ChainToken
		addressQ = sqlmodel.ChainAddressColumns
		address  sqlmodel.ChainAddress
	)
	err = dao.FetchChainToken(ctx, &token,
		dao.And(
			tokenQ.ChainSymbol.Eq(ch.ChainSymbol),
			tokenQ.ContractAddress.Eq(feeTx.ContractAddress),
			tokenQ.Symbol.Eq(feeTx.Symbol)))
	if err != nil {
		return
	}
	err = dao.FetchChainAddress(ctx, &address, dao.And(
		addressQ.ChainSymbol.Eq(ch.ChainSymbol),
		addressQ.Address.Eq(feeTx.FromAddress),
		addressQ.ApplicationID.Eq(feeTx.ApplicationID),
	))
	if err != nil {
		return
	}
	pk, err := utils.AesDecrypt(address.EncKey, config.C.Secret)
	if err != nil {
		return
	}
	nonce, err := GetTransferNonce(ctx, ch, feeTx.FromAddress)
	if err != nil {
		return
	}
	transferOrder := chain.TransferOrder{
		From:            feeTx.FromAddress,
		FromPrivateKey:  pk,
		To:              feeTx.ToAddress,
		ContractAddress: feeTx.ContractAddress,
		Value:           big.NewInt(int64(feeTx.Value * math.Pow10(int(token.Decimals)))),
		TokenID:         big.NewInt(feeTx.TokenID),
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
			ApplicationID:   feeTx.ApplicationID,
			SerialNo:        feeTx.SerialNo,
			ChainSymbol:     ch.ChainSymbol,
			ContractAddress: feeTx.ContractAddress,
			Symbol:          feeTx.Symbol,
			FromAddress:     feeTx.FromAddress,
			ToAddress:       feeTx.ToAddress,
			Value:           feeTx.Value,
			GasPrice:        transferOrder.GasPrice.Int64(),
			TokenID:         feeTx.TokenID,
			TxHash:          transferOrder.TxHash,
			Nonce:           int64(transferOrder.Nonce),
			Hook:            feeTx.Hook,
			TransferType:    int32(types.TransferTypeFee),
			CreateAt:        feeTx.CreateAt,
		}
		_, txErr = dao.AddChainSendTx(ctx, sendTx)
		if txErr != nil {
			return
		}
		feeTx.SendTxID = sendTx.ID
		feeTx.Generated = 1
		_, txErr = dao.UpdateApplicationArrangeFeeTx(c, feeTx)
		return
	})
	return
}
