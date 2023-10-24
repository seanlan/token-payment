/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"go.uber.org/zap"
	"token-payment/internal/chain"
	_ "token-payment/internal/config"

	"github.com/spf13/cobra"
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
		ch, err := chain.NewChain(chain.Config{
			Name:        "eth",
			ChainType:   "evm",
			ChainID:     56,
			ChainSymbol: "bsc",
			Currency:    "bnb",
			RpcURLs:     []string{"https://bsc.rpc.blxrbdn.com"},
			GasPrice:    0,
		})
		if err != nil {
			panic(err)
		}
		b, err := ch.GetBlock(context.Background(),
			32533471)
		if err != nil {
			panic(err)
		}
		zap.S().Infof("%#v", b)
		//s := chain.EvmChain{
		//	Name:        "eth",
		//	ChainType:   "evm",
		//	ChainID:     80001,
		//	ChainSymbol: "eth",
		//	Currency:    "eth",
		//	RpcURLs:     []string{"https://polygon-bor.publicnode.com"},
		//	GasPrice:    0,
		//}
		//s.GetTransaction(context.Background(), "0x507ceead1c8a44806cc62a531e0c684d17ebcf2e5fe443c831adadc0ba2cc4dd")
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
