package main

import "tradingplatform/dataprovider/command/local"

func main() {
	local.NewRootCmd().Execute()
}
