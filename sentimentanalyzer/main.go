package main

import "tradingplatform/sentimentanalyzer/command/local"

func main() {
	local.NewRootCmd().Execute()
}
