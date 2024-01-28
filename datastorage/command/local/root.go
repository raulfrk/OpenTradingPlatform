package local

import (
	JSON "encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"tradingplatform/datastorage/command/cli"
	"tradingplatform/datastorage/command/json"
	"tradingplatform/datastorage/data"
	"tradingplatform/datastorage/handler"

	"tradingplatform/shared/communication"
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "datastorage",
		Short: "DataStorage for the data pipeline",
		Run: func(cmd *cobra.Command, args []string) {
			dns, _ := cmd.Flags().GetString("dsn")
			natsURL, _ := cmd.Flags().GetString("nats-url")
			startupConfig, _ := cmd.Flags().GetString("startup-commands")
			data.SetDSN(dns)
			communication.SetNatsURL(natsURL)
			nc, err := nats.Connect(communication.GetNatsURL())
			if err != nil {
				panic(err)
			}
			defer nc.Close()
			loggingTopic := utils.NewLoggingTopic(types.DataStorage).Generate()
			mlLogger := logging.NewMultiLevelLogger(types.DataStorage,
				os.Stdout, logging.NewNatsWriter(nc, loggingTopic))
			logging.SetLogger(&mlLogger)

			logging.Log().Info().
				Str("natsUrl", communication.GetNatsURL()).
				Msg("starting datastorage, remote logging enabled")

			// Create a channel to receive OS signals
			sigs := make(chan os.Signal, 1)

			// Register the channel to receive SIGINT and SIGTERM signals
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			_, cleanup := data.InitializeDatabase()
			localDbCleanup := data.InitializeDataStorageLocalDatabase()
			defer localDbCleanup()
			defer cleanup()
			command.StartCommandHandler(types.DataStorage, cli.NewRootCmd, json.HandleJSONCommand)
			cmdHandler := command.GetCommandHandler()

			if startupConfig != "" {
				var commands []command.JSONCommand
				// Open config file
				file, err := os.Open(startupConfig)
				if err != nil {
					panic(err)
				}
				defer file.Close()
				fileContent, err := io.ReadAll(file)
				if err != nil {
					panic(err)
				}

				err = JSON.Unmarshal(fileContent, &commands)
				if err != nil {
					panic(err)
				}
				for _, cmd := range commands {
					if cmd.RootOperation != command.JSONOperationStreamSubscribe {
						panic(fmt.Errorf("startup config can only contain stream subscribe requests for now"))
					}
					var streamSubscribeRequest requests.StreamSubscribeRequest
					err := JSON.Unmarshal(cmd.Request, &streamSubscribeRequest)
					if err != nil {
						panic(err)
					}
					streamSubscribeRequest, err = requests.NewStreamSubscribeRequestFromExisting(&streamSubscribeRequest)
					if err != nil {
						panic(err)
					}
					handler.HandleStreamRequest(streamSubscribeRequest)
				}
			}

			go func() {
				<-sigs
				cmdHandler.Cancel()
			}()
			<-cmdHandler.Ctx().Done()
			cmdHandler.Wg.Wait()
		},
	}
	rootCmd.Flags().StringP("dsn", "d", "postgres://postgres:otc@localhost:5432/otc", "DSN for the database")
	rootCmd.Flags().StringP("nats-url", "n", communication.GetNatsURL(), "NATS server URL")

	rootCmd.Flags().StringP("startup-commands", "c", "", "Path to the config startup commands")
	return &rootCmd
}
