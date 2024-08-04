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

func CronCheckArrangeTransactions() {
	var (
		timeout = time.Minute * 10
		ctx     = context.Background()
		chainQ  = sqlmodel.ChainColumns
		chains  []sqlmodel.Chain
		wg      sync.WaitGroup
	)
	// 获取锁
	if dao.Redis.GetLock(ctx, CheckArrangeTxLockKey, timeout) {
		// 释放锁
		defer dao.Redis.ReleaseLock(ctx, CheckArrangeTxLockKey)
	} else {
		// 未获取到锁
		zap.S().Info("CronCheckArrangeTransactions locked !!!")
		return
	}
	// 获取所有的链
	err := dao.FetchAllChain(ctx, &chains, chainQ.Watch.Eq(1), 0, 0)
	if err != nil {
		zap.S().Errorf("fetch all chain error: %#v", err)
		return
	}
	// 读取区块
	for _, ch := range chains {
		wg.Add(1)
		go func(ch sqlmodel.Chain) {
			defer wg.Done()
			handler.CheckArrangeTransactions(ctx, &ch)
		}(ch)
	}
	wg.Wait()
}
