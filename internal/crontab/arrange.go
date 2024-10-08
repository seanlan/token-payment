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

func CronScanArrangeTransactions() {
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
		zap.S().Info("CronScanArrangeTransactions locked !!!")
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
			handler.ScanArrangeTransactions(ctx, &ch)
		}(ch)
	}
	wg.Wait()
}

func CronCheckArrangeTxFee() {
	var (
		timeout = time.Minute * 10
		ctx     = context.Background()
		chainQ  = sqlmodel.ChainColumns
		chains  []sqlmodel.Chain
		wg      sync.WaitGroup
	)
	// 获取锁
	if dao.Redis.GetLock(ctx, CheckArrangeTxFeeLockKey, timeout) {
		// 释放锁
		defer dao.Redis.ReleaseLock(ctx, CheckArrangeTxFeeLockKey)
	} else {
		// 未获取到锁
		zap.S().Info("CronCheckArrangeTxFee locked !!!")
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
			handler.CheckArrangeTxFee(ctx, &ch)
		}(ch)
	}
	wg.Wait()
}

func CronBuildArrangeTx() {
	var (
		timeout = time.Minute * 10
		ctx     = context.Background()
		chainQ  = sqlmodel.ChainColumns
		chains  []sqlmodel.Chain
		wg      sync.WaitGroup
	)
	// 获取锁
	if dao.Redis.GetLock(ctx, BuildArrangeTxLockKey, timeout) {
		// 释放锁
		defer dao.Redis.ReleaseLock(ctx, BuildArrangeTxLockKey)
	} else {
		// 未获取到锁
		zap.S().Info("CronBuildArrangeTx locked !!!")
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
			handler.BuildArrangeTxs(ctx, &ch)
		}(ch)
	}
	wg.Wait()
}

func CronBuildArrangeFeeTx() {
	var (
		timeout = time.Minute * 10
		ctx     = context.Background()
		chainQ  = sqlmodel.ChainColumns
		chains  []sqlmodel.Chain
		wg      sync.WaitGroup
	)
	// 获取锁
	if dao.Redis.GetLock(ctx, BuildArrangeFeeTxLockKey, timeout) {
		// 释放锁
		defer dao.Redis.ReleaseLock(ctx, BuildArrangeFeeTxLockKey)
	} else {
		// 未获取到锁
		zap.S().Info("CronBuildArrangeFeeTx locked !!!")
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
			handler.BuildArrangeFeeTxs(ctx, &ch)
		}(ch)
	}
	wg.Wait()
}
