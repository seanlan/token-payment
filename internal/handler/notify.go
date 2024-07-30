package handler

import (
	"context"
	"errors"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/types"
)

func NotifyTransaction(ctx context.Context, tx *sqlmodel.ChainTx) (err error) {
	var (
		notifyUrl   string
		addressQ    = sqlmodel.ChainAddressColumns
		appQ        = sqlmodel.ApplicationColumns
		address     sqlmodel.ChainAddress
		application sqlmodel.Application
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
	}
	if notifyUrl == "" { // 不需要通知
		tx.NotifySuccess = 1
		_, err = dao.UpdateChainTx(ctx, tx)
		return
	}
	err = dao.FetchApplication(ctx, &application, appQ.ID.Eq(tx.ApplicationID))
	if err != nil {
		return
	}
	// 进行http通知
	client := NewHttpClient()

	return nil
}
