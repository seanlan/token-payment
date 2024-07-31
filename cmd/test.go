/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"token-payment/internal/crontab"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//crontab.CronGenerateTransactions()
		crontab.CronSendTransactions()
		//toAddress := common.HexToAddress("0x0d49ea539217d011faec8c48ec864941aab1cf17")
		//tx := types.NewTx(&types.LegacyTx{
		//	Nonce:    1,
		//	To:       &toAddress,
		//	Value:    big.NewInt(100),
		//	Gas:      20000,
		//	GasPrice: big.NewInt(100),
		//	Data:     nil,
		//})
		//zap.S().Infof("tx hash: %s", tx.Hash().String())
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
