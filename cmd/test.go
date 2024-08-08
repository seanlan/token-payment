/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"token-payment/internal/tokenpay"
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
		sdk := tokenpay.NewClient(
			"f4399f851e984405aa1eba51ecbce790",
			"f4399f851e984405aa1eba51ecbce790",
			"http://127.0.0.1:8080")
		resp, err := sdk.Withdraw(tokenpay.WithdrawReqData{
			Chain:           "amoy",
			SerialNo:        "0b2695c3-a487-4102-8a9a-75b2562c660b4",
			Symbol:          "matic",
			ContractAddress: "",
			TokenID:         0,
			Value:           "0.01",
			ToAddress:       "0x92304f29496ded1314d124c1d7905be44608b709",
			NotifyUrl:       "",
		})
		if err != nil {
			zap.S().Infof("err: %v", err)
		}
		zap.S().Infof("resp: %v", resp)
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
