/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"token-payment/internal/crontab"
)

type CronZapLogger struct {
}

func (*CronZapLogger) Info(msg string, keysAndValues ...interface{}) {
	zap.S().Info(msg, keysAndValues)
}

func (*CronZapLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	zap.S().Error(err, msg, keysAndValues)
}

func cronFunc(cmd *cobra.Command, args []string) {
	crontab.ClearRedisLock(cmd.Context())
	var err error
	c := cron.New(cron.WithLogger(&CronZapLogger{}), cron.WithChain(cron.Recover(&CronZapLogger{})))
	// 索引新的区块(优先生成交易高度，并做好排序） 方便后面可以并发的读取区块
	_, err = c.AddFunc("@every 1s", crontab.CronReadNextBlock)
	if err != nil {
		zap.S().Fatalf("cron add func error: %#v", err)
	}
	// 检测区块（检索交易） 由于前面已经排序，所以可以并发的读取区块
	_, err = c.AddFunc("@every 1s", crontab.CronCheckBlock)
	if err != nil {
		zap.S().Fatalf("cron add func error: %#v", err)
	}
	// 检查是否有分叉 同时刷新交易确认数
	_, err = c.AddFunc("@every 1s", crontab.CronCheckRebase)
	if err != nil {
		zap.S().Fatalf("cron add func error: %#v", err)
	}
	// 更新rebase区块(修复分叉)
	_, err = c.AddFunc("@every 1s", crontab.CronRebaseBlock)
	if err != nil {
		zap.S().Fatalf("cron add func error: %#v", err)
	}
	// 发送交易
	_, err = c.AddFunc("@every 1s", crontab.CronSendTransactions)
	if err != nil {
		zap.S().Fatalf("cron add func error: %#v", err)
	}
	// 扫描整理交易
	_, err = c.AddFunc("@every 1s", crontab.CronScanArrangeTransactions)
	if err != nil {
		zap.S().Fatalf("cron add func error: %#v", err)
	}
	// 检测整理交易手续费
	_, err = c.AddFunc("@every 1s", crontab.CronCheckArrangeTxFee)
	if err != nil {
		zap.S().Fatalf("cron add func error: %#v", err)
	}
	// 交易通知
	_, err = c.AddFunc("@every 1s", crontab.CronTransactionNotify)
	if err != nil {
		zap.S().Fatalf("cron add func error: %#v", err)
	}
	// 构建提币交易
	_, err = c.AddFunc("@every 1s", crontab.CronBuildWithdrawTransactions)
	if err != nil {
		zap.S().Fatalf("cron add func error: %#v", err)
	}
	// 构建零钱整理交易
	_, err = c.AddFunc("@every 1s", crontab.CronBuildArrangeTx)
	if err != nil {
		zap.S().Fatalf("cron add func error: %#v", err)
	}
	// 构建零钱整理手续费交易
	_, err = c.AddFunc("@every 1s", crontab.CronBuildArrangeFeeTx)
	if err != nil {
		zap.S().Fatalf("cron add func error: %#v", err)
	}
	c.Start()
	select {}
}

// cronCmd represents the crontab command
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: cronFunc,
}

func init() {
	rootCmd.AddCommand(cronCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cronCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cronCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
