package handler

import (
	"context"
	"strings"
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
