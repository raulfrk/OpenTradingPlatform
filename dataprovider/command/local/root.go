package local

import (
	"os"
	"os/signal"
	"syscall"

	"tradingplatform/dataprovider/command/cli"
	"tradingplatform/dataprovider/command/json"
	"tradingplatform/dataprovider/data"
	"tradingplatform/shared/communication"
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "dataprovider",
		Short: "DataProvider for the data pipeline",
		Run: func(cmd *cobra.Command, args []string) {
			natsURL, _ := cmd.Flags().GetString("nats-url")
			communication.SetNatsURL(natsURL)
			nc, err := nats.Connect(communication.GetNatsURL())
			if err != nil {
				panic(err)
			}
			defer nc.Close()
			loggingTopic := utils.NewLoggingTopic(types.DataProvider).Generate()
			mlLogger := logging.NewMultiLevelLogger(types.DataProvider,
				os.Stdout, logging.NewNatsWriter(nc, loggingTopic))
			logging.SetLogger(&mlLogger)

			logging.Log().Info().
				Str("natsUrl", communication.GetNatsURL()).
				Msg("starting dataprovider, remote logging enabled")

			// Create a channel to receive OS signals
			sigs := make(chan os.Signal, 1)

			// Register the channel to receive SIGINT and SIGTERM signals
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

			cleanup := data.InitializeDataProviderLocalDatabase()
			defer cleanup()
			command.StartCommandHandler(types.DataProvider, cli.NewRootCmd, json.HandleJSONCommand)
			handler := command.GetCommandHandler()
			go func() {
				<-sigs
				handler.Cancel()
			}()
			<-handler.Ctx().Done()
			handler.Wg.Wait()
		},
	}
	rootCmd.Flags().StringP("nats-url", "n", communication.GetNatsURL(), "NATS server URL")
	return &rootCmd
}
