package crontab

import (
	"context"
	"token-payment/internal/dao"
)

const (
	ReadBlockLockKey                = "cron_read_block"
	CheckBlockLockKey               = "cron_check_block"
	CheckRebaseLockKey              = "cron_check_rebase"
	RebaseBlockLockKey              = "cron_rebase_block"
	NotifyLockKey                   = "cron_transaction_notify"
	SendLockKey                     = "cron_send_transactions"
	CheckAddressPoolLockKey         = "cron_check_address_pool"
	CheckArrangeTxLockKey           = "cron_check_arrange_tx"
	CheckArrangeTxFeeLockKey        = "cron_check_arrange_tx_fee"
	BuildWithdrawTransactionLockKey = "cron_build_withdraw_transactions"
	BuildArrangeTxLockKey           = "cron_build_arrange_tx"
	BuildArrangeFeeTxLockKey        = "cron_build_arrange_fee_tx"
)

func ClearRedisLock(ctx context.Context) {
	dao.Redis.Del(ctx, ReadBlockLockKey)
	dao.Redis.Del(ctx, CheckBlockLockKey)
	dao.Redis.Del(ctx, CheckRebaseLockKey)
	dao.Redis.Del(ctx, RebaseBlockLockKey)
	dao.Redis.Del(ctx, NotifyLockKey)
	dao.Redis.Del(ctx, SendLockKey)
	dao.Redis.Del(ctx, CheckAddressPoolLockKey)
	dao.Redis.Del(ctx, CheckArrangeTxLockKey)
	dao.Redis.Del(ctx, CheckArrangeTxFeeLockKey)
	dao.Redis.Del(ctx, BuildWithdrawTransactionLockKey)
	dao.Redis.Del(ctx, BuildArrangeTxLockKey)
	dao.Redis.Del(ctx, BuildArrangeFeeTxLockKey)
}
