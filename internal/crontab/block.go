// Package crontab
// 区块更新
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

// CronReadNextBlock 读取下一个区块
//
//	@Description: 读取下一个区块
func CronReadNextBlock() {
	var (
		timeout = time.Minute * 10
		ctx     = context.Background()
		chainQ  = sqlmodel.ChainColumns
		chains  []sqlmodel.Chain
		lockKey = "cron_read_block"
		wg      sync.WaitGroup
	)

	// 获取锁
	if dao.Redis.GetLock(ctx, lockKey, timeout) {
		// 释放锁
		defer dao.Redis.ReleaseLock(ctx, lockKey)
	} else {
		// 未获取到锁
		zap.S().Info("CronReadNextBlock locked !!!")
		return
	}
	defer func() {
		if _err := recover(); _err != nil {
			zap.S().Errorw("cron read block error", "error", _err)
		}
	}()
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
			// TODO: 读取下一个区块
			zap.S().Infow("read block", "chain", ch.ChainSymbol)
			if ch.RebaseBlock > 0 { // 有区块需要更新
				zap.S().Infow("update rebase block", "chain", ch.ChainSymbol, "block", ch.RebaseBlock)
			} else {
				handler.ReadNextBlock(ctx, &ch)
			}
		}(ch)
	}
	wg.Wait()
}

// CronRebaseBlock 更新区块
//
//	@Description: 更新区块
func CronRebaseBlock() {
	var (
		timeout = time.Minute * 10
		ctx     = context.TODO()
		chainQ  = sqlmodel.ChainColumns
		chains  []sqlmodel.Chain
		lockKey = "cron_rebase_block"
	)
	// 获取锁
	if dao.Redis.GetLock(ctx, lockKey, timeout) {
		// 结束后释放锁
		defer dao.Redis.ReleaseLock(ctx, lockKey)

	} else {
		return
	}
	// 获取所有的链
	err := dao.FetchAllChain(ctx, &chains, chainQ.Watch.Eq(1), 0, 0)
	if err != nil {
		zap.S().Errorf("fetch all chain error: %#v", err)
		return
	}
	for _, ch := range chains {
		// TODO: 更新区块
		if ch.RebaseBlock > 0 { // 有区块需要更新
			zap.S().Infow("read block", "chain", ch.ChainSymbol)
			handler.RebaseBlock(ctx, &ch)
		}
	}
}

// CronCheckBlock
//
//	@Description: 检查区块
func CronCheckBlock() {
	var (
		timeout = time.Minute * 10
		ctx     = context.Background()
		chainQ  = sqlmodel.ChainColumns
		chains  []sqlmodel.Chain
		lockKey = "cron_check_block"
		wg      sync.WaitGroup
	)

	// 获取锁
	if dao.Redis.GetLock(ctx, lockKey, timeout) {
		// 释放锁
		defer dao.Redis.ReleaseLock(ctx, lockKey)
	} else {
		// 未获取到锁
		zap.S().Info("CronCheckBlock locked !!!")
		return
	}
	defer func() {
		if _err := recover(); _err != nil {
			zap.S().Errorw("cron check block error", "error", _err)
		}
	}()
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
			// TODO: 检测区块
			zap.S().Infow("check block", "chain", ch.ChainSymbol)
			handler.CheckBlock(ctx, &ch)
		}(ch)
	}
	wg.Wait()
}
