package handler

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"math/big"
	"sync"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/types"
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
			txQ.Confirm.Eq(ch.Confirm),
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

func CheckArrangeTxFeeByApplication(ctx context.Context, ch *sqlmodel.Chain, appChain *sqlmodel.ApplicationChain) {
	var (
		arrangeQ        = sqlmodel.ApplicationArrangeTxColumns
		arrangeTxs      = make([]sqlmodel.ApplicationArrangeTx, 0)
		tokenQ          = sqlmodel.ChainTokenColumns
		tokenContracts  = make([]string, 0)
		tokens          = make([]sqlmodel.ChainToken, 0)
		tokenMap        = make(map[string]sqlmodel.ChainToken)
		tokenConsumeMap = make(map[string]float64)
		balanceEnough   = true
	)
	err := dao.FetchAllApplicationArrangeTx(ctx, &arrangeTxs, dao.And(
		arrangeQ.ChainSymbol.Eq(ch.ChainSymbol),
		arrangeQ.ApplicationID.Eq(appChain.ApplicationID),
		arrangeQ.Generated.Eq(0), // 未生成交易
	), 0, int(ch.Concurrent), arrangeQ.ID.Asc())
	if err != nil || len(arrangeTxs) == 0 {
		return
	}
	for _, order := range arrangeTxs {
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
	for _, order := range arrangeTxs {
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
		balance, _err := client.GetBalance(ctx, appChain.FeeWallet, contract)
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
		for _, order := range arrangeTxs {
			_ = GenerateArrangeTx(ctx, ch, &order, appChain)
		}
	}
}

func GenerateArrangeTx(ctx context.Context, ch *sqlmodel.Chain, tx *sqlmodel.ApplicationArrangeTx, appChain *sqlmodel.ApplicationChain) error {
	return nil
}
