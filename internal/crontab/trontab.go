package crontab

import (
	"context"
	"token-payment/internal/dao"
)

const (
	ReadBlockLockKey        = "cron_read_block"
	CheckBlockLockKey       = "cron_check_block"
	CheckRebaseLockKey      = "cron_check_rebase"
	RebaseBlockLockKey      = "cron_rebase_block"
	NotifyLockKey           = "cron_transaction_notify"
	GenerateLockKey         = "cron_generate_transactions"
	SendLockKey             = "cron_send_transactions"
	CheckAddressPoolLockKey = "cron_check_address_pool"
)

func ClearRedisLock(ctx context.Context) {
	dao.Redis.Del(ctx, ReadBlockLockKey)
	dao.Redis.Del(ctx, CheckBlockLockKey)
	dao.Redis.Del(ctx, CheckRebaseLockKey)
	dao.Redis.Del(ctx, RebaseBlockLockKey)
	dao.Redis.Del(ctx, NotifyLockKey)
	dao.Redis.Del(ctx, GenerateLockKey)
	dao.Redis.Del(ctx, SendLockKey)
	dao.Redis.Del(ctx, CheckAddressPoolLockKey)
}
