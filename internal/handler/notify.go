package handler

import (
	"context"
	"errors"
	"math"
	"time"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/tokenpay"
	"token-payment/internal/types"
)

// NotifyTransaction
//
//	@Description: 通知交易
//	@param ctx
//	@param tx
//	@return err
func NotifyTransaction(ctx context.Context, tx *sqlmodel.ChainTx) (err error) {
	var (
		notifyUrl   string
		addressQ    = sqlmodel.ChainAddressColumns
		appQ        = sqlmodel.ApplicationColumns
		chainQ      = sqlmodel.ChainColumns
		tokenQ      = sqlmodel.ChainTokenColumns
		orderQ      = sqlmodel.ApplicationWithdrawOrderColumns
		address     sqlmodel.ChainAddress
		application sqlmodel.Application
		chain       sqlmodel.Chain
		token       sqlmodel.ChainToken
		order       sqlmodel.ApplicationWithdrawOrder
	)
	switch types.TransferType(tx.TransferType) {
	case types.TransferTypeIn:
		// 充值
		err = dao.FetchChainAddress(ctx, &address, dao.And(
			addressQ.ChainSymbol.Eq(tx.ChainSymbol),
			addressQ.Address.Eq(tx.ToAddress),
			addressQ.Watch.Eq(1),
		))
		if err != nil && !errors.Is(err, dao.ErrNotFound) {
			return err
		}
		notifyUrl = address.Hook
	case types.TransferTypeOut:
		// 提现
		err = dao.FetchApplicationWithdrawOrder(ctx, &order, dao.And(
			orderQ.ChainSymbol.Eq(tx.ChainSymbol),
			orderQ.TxHash.Eq(tx.TxHash),
		))
		if err != nil && !errors.Is(err, dao.ErrNotFound) {
			return err
		}
		notifyUrl = order.Hook
	case types.TransferTypeFee:
		// 手续费转账
		notifyUrl = ""
	}
	if notifyUrl == "" { // 不需要通知
		tx.NotifySuccess = 1
		_, err = dao.UpdateChainTx(ctx, tx)
		return
	}
	// 获取应用信息
	err = dao.FetchApplication(ctx, &application, appQ.ID.Eq(tx.ApplicationID))
	if err != nil {
		return
	}
	// 获取链信息
	err = dao.FetchChain(ctx, &chain, chainQ.ChainSymbol.Eq(tx.ChainSymbol))
	if err != nil {
		return
	}
	// 获取token信息
	err = dao.FetchChainToken(ctx, &token, dao.And(
		tokenQ.ChainSymbol.Eq(tx.ChainSymbol),
		tokenQ.Symbol.Eq(tx.Symbol),
		tokenQ.ContractAddress.Eq(tx.ContractAddress),
	))
	if err != nil {
		return
	}
	// 进行http通知
	ntx := tokenpay.NotifyTx{
		ApplicationID:   tx.ApplicationID,
		ChainSymbol:     tx.ChainSymbol,
		TxHash:          tx.TxHash,
		FromAddress:     tx.FromAddress,
		ToAddress:       tx.ToAddress,
		ContractAddress: tx.ContractAddress,
		Symbol:          tx.Symbol,
		Decimals:        token.Decimals,
		TokenID:         tx.TokenID,
		Value:           tx.Value,
		TxIndex:         tx.TxIndex,
		BatchIndex:      tx.BatchIndex,
		Confirm:         tx.Confirm,
		MaxConfirm:      chain.Confirm,
		TransferType:    tx.TransferType,
		SerialNo:        tx.SerialNo,
		CreateAt:        tx.CreateAt,
	}
	client := tokenpay.NewClient(application.AppName, application.AppKey, "")
	success, notifyErr := client.NotifyTransaction(ntx, notifyUrl)
	if notifyErr != nil {
		success = false
	}
	if success || tx.NotifyFailedTimes >= 10 {
		tx.NotifySuccess = 1
	} else {
		tx.NotifyFailedTimes++
		scale := math.Pow(2, float64(tx.NotifyFailedTimes))
		tx.NotifyNextTime = time.Now().Unix() + int64(scale)*30
	}
	_, err = dao.UpdateChainTx(ctx, tx)
	return
}
