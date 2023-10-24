package crontab

import (
	"context"
	"go.uber.org/zap"
	"time"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/handler"
)

func CronReadNextBlock() {
	var (
		timeout = time.Minute * 10
		ctx     = context.TODO()
		chains  []sqlmodel.Chain
		lockKey = "cron_read_block"
	)
	// 释放锁
	defer dao.Redis.ReleaseLock(ctx, lockKey)
	// 获取锁
	if !dao.Redis.GetLock(ctx, lockKey, timeout) {
		return
	}
	// 获取所有的链
	err := dao.FetchAllChain(ctx, &chains, nil, 0, 0)
	if err != nil {
		zap.S().Errorf("fetch all chain error: %#v", err)
		return
	}
	for _, ch := range chains {
		// TODO: 读取下一个区块
		zap.S().Infow("read block", "chain", ch.ChainSymbol)
		if ch.RebaseBlock > 0 { // 有区块需要更新
			zap.S().Infow("update rebase block", "chain", ch.ChainSymbol, "block", ch.RebaseBlock)
		} else {
			handler.ReadNextBlock(ctx, &ch)
		}
	}
}
