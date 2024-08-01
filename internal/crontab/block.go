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
		wg      sync.WaitGroup
	)

	// 获取锁
	if dao.Redis.GetLock(ctx, ReadBlockLockKey, timeout) {
		// 释放锁
		defer dao.Redis.ReleaseLock(ctx, ReadBlockLockKey)
	} else {
		// 未获取到锁
		zap.S().Info("CronReadNextBlock locked !!!")
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
			// TODO: 读取下一个区块
			zap.S().Infow("read next block", "chain", ch.ChainSymbol)
			if ch.HasBranch > 0 { // 有区块分叉需要更新
				zap.S().Infow("update rebase block", "chain", ch.ChainSymbol, "block", ch.RebaseBlock)
			} else {
				handler.ReadNextBlock(ctx, &ch)
			}
		}(ch)
	}
	wg.Wait()
}

// CronCheckRebase
//
//	@Description: 检查区块
func CronCheckRebase() {
	var (
		timeout = time.Minute * 10
		ctx     = context.Background()
		chainQ  = sqlmodel.ChainColumns
		chains  []sqlmodel.Chain
		wg      sync.WaitGroup
	)
	// 获取锁
	if dao.Redis.GetLock(ctx, CheckRebaseLockKey, timeout) {
		// 释放锁
		defer dao.Redis.ReleaseLock(ctx, CheckRebaseLockKey)
	} else {
		// 未获取到锁
		zap.S().Info("CronReadNextBlock locked !!!")
		return
	}
	// 获取所有的链
	err := dao.FetchAllChain(ctx, &chains, chainQ.Watch.Eq(1), 0, 0)
	if err != nil {
		zap.S().Errorf("fetch all chain error: %#v", err)
		return
	}
	// 检查区块
	for _, ch := range chains {
		wg.Add(1)
		go func(ch sqlmodel.Chain) {
			defer wg.Done()
			// TODO: 读取下一个区块
			zap.S().Infow("check rebase block", "chain", ch.ChainSymbol)
			if ch.HasBranch == 0 { // 没有区块分叉需要更新
				handler.CheckRebase(ctx, &ch)
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
		wg      sync.WaitGroup
	)
	// 获取锁
	if dao.Redis.GetLock(ctx, RebaseBlockLockKey, timeout) {
		// 结束后释放锁
		defer dao.Redis.ReleaseLock(ctx, RebaseBlockLockKey)
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
		wg.Add(1)
		go func(ch sqlmodel.Chain) {
			defer wg.Done()
			zap.S().Infow("rebase block", "chain", ch.ChainSymbol)
			if ch.HasBranch > 0 { // 有区块分叉需要更新
				handler.RebaseBlock(ctx, &ch)
			}
		}(ch)
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
		wg      sync.WaitGroup
	)

	// 获取锁
	if dao.Redis.GetLock(ctx, CheckBlockLockKey, timeout) {
		// 释放锁
		defer dao.Redis.ReleaseLock(ctx, CheckBlockLockKey)
	} else {
		// 未获取到锁
		zap.S().Info("CronCheckBlock locked !!!")
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
			// TODO: 检测区块
			zap.S().Infow("check block", "chain", ch.ChainSymbol)
			handler.CheckBlocks(ctx, &ch)
		}(ch)
	}
	wg.Wait()
}
