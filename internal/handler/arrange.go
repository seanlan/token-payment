package handler

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

func CheckArrangeTransactions(ctx context.Context, ch *sqlmodel.Chain) {
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
		wg.Add(1)
		go func(token sqlmodel.ChainToken) {
			defer wg.Done()
			CheckArrangeTransaction(ctx, ch, &token)
		}(token)
	}
	wg.Wait()
	return
}

func CheckArrangeTransaction(ctx context.Context, ch *sqlmodel.Chain, token *sqlmodel.ChainToken) {
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
			txQ.TransferType.Eq(types.TransferTypeIn)).
		Group("application_id, address").
		Order("amount desc").
		Limit(int(ch.Concurrent)).
		Find(&arrangeList).Error
	if err != nil {
		return
	}
	zap.S().Infof("arrangeList: %#v", arrangeList)
	for _, arrange := range arrangeList {
		if arrange.Amount < token.Threshold { // 低于阈值 不生成交易
			continue
		}
		wg.Add(1)
		go func(arrange AddressArrange) {
			defer wg.Done()
			GenerateArrangeTransaction(ctx, ch, token, &arrange)
		}(arrange)
	}
	wg.Wait()
}

func GenerateArrangeTransaction(ctx context.Context, ch *sqlmodel.Chain, token *sqlmodel.ChainToken, arrange *AddressArrange) {
	var (
		txQ     = sqlmodel.ChainTxColumns
		finalTx sqlmodel.ChainTx
		amount  float64
		effect  int64
		err     error
	)
	err = dao.FetchChainTx(ctx, &finalTx, dao.And(
		txQ.ChainSymbol.Eq(ch.ChainSymbol),
		txQ.Symbol.Eq(token.Symbol),
		txQ.ToAddress.Eq(arrange.Address),
	), txQ.ID.Desc())
	if err != nil || finalTx.ID == 0 {
		return
	}
	// 生成交易
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
