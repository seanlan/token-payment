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
		crontab.CronReadNextBlock()
		//ch, err := chain.NewChain(chain.Config{
		//	Name:        "Polygon",
		//	ChainType:   "polygon",
		//	ChainID:     137,
		//	ChainSymbol: "polygon",
		//	Currency:    "matic",
		//	RpcURLs:     []string{"https://poly-rpc.gateway.pokt.network"},
		//	GasPrice:    0,
		//})
		//if err != nil {
		//	panic(err)
		//}
		//b, err := ch.GetBlock(context.Background(),
		//	49165088)
		//if err != nil {
		//	panic(err)
		//}
		//zap.S().Infof("%#v", b)
		//s := chain.PolygonChain{
		//	Name:        "eth",
		//	ChainType:   "evm",
		//	ChainID:     80001,
		//	ChainSymbol: "eth",
		//	Currency:    "eth",
		//	RpcURLs:     []string{"https://gateway.tenderly.co/public/polygon-mumbai"},
		//	GasPrice:    0,
		//}
		//t, err := s.GetTransaction(context.Background(), "0x71c19da19a75611aa6cef7484dbac3b000a2735887f58f5e3d265589370d839f")
		//if err != nil {
		//	panic(err)
		//}
		//zap.S().Infof("%#v", t)
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
