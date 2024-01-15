package main

import (
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
	"tradingplatform/dataprovider/data"
	shcomm "tradingplatform/shared/communication"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer nc.Close()
	loggingTopic := utils.NewLoggingTopic(types.DataProvider).Generate()
	mlLogger := logging.NewMultiLevelLogger(types.DataProvider,
		os.Stdout, logging.NewNatsWriter(nc, loggingTopic))
	logging.SetLogger(&mlLogger)

	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)

	// Register the channel to receive SIGINT and SIGTERM signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cleanup := data.InitializeDataProviderLocalDatabase()
	defer cleanup()

	server := grpc.NewServer()
	shcomm.StartCommunicationHandler(server)
	//command.StartCommandHandler(types.DataProvider, cli.NewRootCmd, json.HandleJSONCommand)
	//handler := command.GetCommandHandler()
	//go func() {
	//	<-sigs
	//	handler.Cancel()
	//}()
	//<-handler.Ctx().Done()
	//handler.Wg.Wait()
}
