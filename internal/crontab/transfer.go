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

func CronGenerateTransactions() {
	var (
		timeout = time.Minute * 10
		ctx     = context.Background()
		chainQ  = sqlmodel.ChainColumns
		chains  []sqlmodel.Chain
		lockKey = "cron_generate_transactions"
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
