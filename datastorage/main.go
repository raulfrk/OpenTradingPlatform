package main

import (
	"tradingplatform/datastorage/command/local"
)

func main() {
	local.NewRootCmd().Execute()
}
