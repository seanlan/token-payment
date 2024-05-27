/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
		//crontab.CronReadNextBlock()
		client, err := ethclient.Dial("https://rpc.ankr.com/polygon/f60b6a29d8551b2156461783d5ebc4b00983609c846db245e42bf3c5aa51af5c")
		if err != nil {
			panic(err)
		}
		//b, err := client.BlockByNumber(context.Background(), big.NewInt(57175074))
		//if err != nil {
		//	panic(err)
		//}
		ts, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithNumber(57175074))
		for _, t := range ts {
			zap.S().Infof("tx: %+v", t)
		}
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
