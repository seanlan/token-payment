package crontab

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"time"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/handler"
)

var (
	notifyMaxCount = 100 // 通知最大数量
)

// CronTransactionNotify
//
//	@Description: 交易通知
func CronTransactionNotify() {
	var (
		timeout = time.Minute * 10
		ctx     = context.Background()
		txQ     = sqlmodel.ChainTxColumns
		txs     []sqlmodel.ChainTx
		wg      sync.WaitGroup
	)
	// 获取锁
	if dao.Redis.GetLock(ctx, NotifyLockKey, timeout) {
		// 释放锁
		defer dao.Redis.ReleaseLock(ctx, NotifyLockKey)
	} else {
		// 未获取到锁
		zap.S().Info("CronReadNextBlock locked !!!")
		return
	}
	// 获取所有待通知的交易
	err := dao.FetchAllChainTx(ctx, &txs, dao.And(
		txQ.NotifySuccess.Eq(0),
		txQ.NotifyNextTime.Lte(time.Now().Unix()),
	), 0, notifyMaxCount)
	if err != nil {
		zap.S().Errorf("fetch all chain tx error: %#v", err)
		return
	}
	// 通知交易
	for _, tx := range txs {
		wg.Add(1)
		go func(tx sqlmodel.ChainTx) {
			defer wg.Done()
			// TODO: 通知交易
			zap.S().Infow("notify transaction", "tx", tx.TxHash)
			_ = handler.NotifyTransaction(ctx, &tx)
		}(tx)
	}
	wg.Wait()
}
