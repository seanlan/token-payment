/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"token-payment/cmd"
	_ "token-payment/internal/config"
)

func main() {
	cmd.Execute()
}
