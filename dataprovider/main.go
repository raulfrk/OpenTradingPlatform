package main

import (
	"os"
	"os/signal"
	"syscall"
	"tradingplatform/dataprovider/command/cli"
	"tradingplatform/dataprovider/command/json"
	"tradingplatform/dataprovider/data"
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer nc.Close()
	mlLogger := logging.NewMultiLevelLogger(os.Stdout, logging.NewNatsWriter(nc, "dataprovider.logging"))
	logging.SetLogger(&mlLogger)

	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)

	// Register the channel to receive SIGINT and SIGTERM signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	_, cleanup := data.InitializeDatabase()
	defer cleanup()
	command.StartCommandHandler(types.DataProvider, cli.NewRootCmd, json.HandleJSONCommand)
	handler := command.GetCommandHandler()
	go func() {
		<-sigs
		handler.Cancel()
	}()
	<-handler.Ctx().Done()
	handler.Wg.Wait()
}
