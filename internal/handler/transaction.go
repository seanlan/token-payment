package handler

import (
	"context"
	"math"
	"strings"
	"sync"
	"token-payment/internal/chain"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/types"
)

// CheckTransactions
//
//	@Description: 检查交易
//	@param ctx
//	@param txs
//	@return err
func CheckTransactions(ctx context.Context, ch *sqlmodel.Chain, txs []*chain.Transaction) (err error) {
	for _, tx := range txs {
		if tx == nil {
			continue
		}
		err = CheckRechargeTransaction(ctx, ch, tx) // 检查充值交易
		if err != nil {
			return
		}
		err = CheckWithdrawTransaction(ctx, ch, tx) // 检查提现交易
		if err != nil {
			return
		}
	}
	return
}

// CheckRechargeTransaction
//
//	@Description: 检查充值交易
//	@param ctx
//	@param ch
//	@param tx
//	@return err
func CheckRechargeTransaction(ctx context.Context, ch *sqlmodel.Chain, tx *chain.Transaction) (err error) {
	var (
		addressQ = sqlmodel.ChainAddressColumns
		tokenQ   = sqlmodel.ChainTokenColumns
		bills    = make([]sqlmodel.ChainTx, 0)
	)
	for _, bill := range tx.Bills {
		var (
			address sqlmodel.ChainAddress
			token   sqlmodel.ChainToken
			toAddr  = strings.ToLower(bill.To)
		)
		_ = dao.FetchChainAddress(ctx, &address, dao.And(
			addressQ.ChainSymbol.Eq(ch.ChainSymbol),
			addressQ.Watch.Eq(1), // 是否是监控地址
			addressQ.Address.Eq(toAddr),
		))
		if address.ID == 0 {
			// 不是监控地址
			continue
		}
		_ = dao.FetchChainToken(ctx, &token, dao.And(
			tokenQ.ChainSymbol.Eq(ch.ChainSymbol),
			tokenQ.ContractAddress.Eq(bill.ContractAddress),
		))
		if token.ID == 0 {
			// 不是监控token
			continue
		}
		bills = append(bills, sqlmodel.ChainTx{
			ApplicationID:   address.ApplicationID,
			ChainSymbol:     ch.ChainSymbol,
			BlockNumber:     tx.BlockNumber,
			BlockHash:       tx.BlockHash,
			TxHash:          tx.Hash,
			FromAddress:     strings.ToLower(bill.From),
			ToAddress:       strings.ToLower(bill.To),
			ContractAddress: bill.ContractAddress,
			Symbol:          token.Symbol,
			Value:           bill.Value.String(),
			TxIndex:         int64(bill.Index),
			BatchIndex:      int64(bill.BatchIndex),
			TransferType:    int32(types.TransferTypeIn),
			CreateAt:        tx.Time.Unix(),
		})
	}
	if len(bills) > 0 {
		_, err = dao.AddsChainTx(ctx, &bills)
	}
	return
}

// CheckWithdrawTransaction
//
//	@Description: 检查提现交易
//	@param ctx
//	@param ch
//	@param tx
//	@return err
func CheckWithdrawTransaction(ctx context.Context, ch *sqlmodel.Chain, tx *chain.Transaction) (err error) {
	// TODO: 检查提现交易
	return
}

// UpdateTransactionsConfirm
//
//	@Description: 更新交易确认数
//	@param ctx
//	@param ch
//	@param currentBlockNumber
//	@return err
func UpdateTransactionsConfirm(ctx context.Context, ch *sqlmodel.Chain) (err error) {
	var (
		transQ = sqlmodel.ChainTxColumns
		txList = make([]sqlmodel.ChainTx, 0)
		wg     sync.WaitGroup
	)
	err = dao.FetchAllChainTx(ctx, &txList, dao.And(
		transQ.ChainSymbol.Eq(ch.ChainSymbol),
		transQ.Removed.Eq(0),
		transQ.BlockNumber.Lte(ch.RebaseBlock),
		transQ.Confirm.Lt(ch.Confirm)), 0, 0)
	if err != nil {
		return
	}
	for _, tx := range txList {
		wg.Add(1)
		go func(tx sqlmodel.ChainTx) {
			defer wg.Done()
			_ = UpdateTransactionConfirm(ctx, ch, &tx)
		}(tx)
	}
	return
}

// UpdateTransactionConfirm
//
//	@Description: 更新交易确认数
//	@param ctx
//	@param ch
//	@param tx
//	@return err
func UpdateTransactionConfirm(ctx context.Context, ch *sqlmodel.Chain, tx *sqlmodel.ChainTx) (err error) {
	var (
		blockQ = sqlmodel.ChainBlockColumns
		count  int64
	)
	count, err = dao.CountChainBlock(ctx, blockQ.BlockHash.Eq(tx.BlockHash))
	if err != nil || count == 0 { // 未找到区块
		return
	}
	if count == 0 {
		tx.Removed = 1
	} else {
		confirm := ch.RebaseBlock - tx.BlockNumber
		tx.Confirm = int32(math.Min(float64(ch.Confirm), float64(confirm)))
		tx.NotifyNextTime = 0
		tx.NotifySuccess = 0
		tx.NotifyFailedTimes = 0
	}
	_, err = dao.UpdateChainTx(ctx, tx)
	return
}
