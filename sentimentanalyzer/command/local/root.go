package local

import (
	JSON "encoding/json"
	"io"
	"os"
	"os/signal"
	"syscall"

	"tradingplatform/sentimentanalyzer/command/cli"
	"tradingplatform/sentimentanalyzer/command/json"
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "sentalyzer",
		Short: "SentimentAnalyzer startup command",
		Run: func(cmd *cobra.Command, args []string) {
			startupConfig, _ := cmd.Flags().GetString("startup-commands")
			nc, err := nats.Connect(nats.DefaultURL)
			if err != nil {
				panic(err)
			}
			defer nc.Close()

			loggingTopic := utils.NewLoggingTopic(types.SentimentAnalyzer).Generate()
			mlLogger := logging.NewMultiLevelLogger(types.SentimentAnalyzer,
				os.Stdout, logging.NewNatsWriter(nc, loggingTopic))
			logging.SetLogger(&mlLogger)

			// Create a channel to receive OS signals
			sigs := make(chan os.Signal, 1)

			// Register the channel to receive SIGINT and SIGTERM signals
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			// TODO: Update this for this particular component
			// localDbCleanup := data.InitializeDataStorageLocalDatabase()
			// defer localDbCleanup()
			command.StartCommandHandler(types.SentimentAnalyzer, cli.NewRootCmd, json.HandleJSONCommand)
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

				JSON.Unmarshal(fileContent, &commands)
				// TODO: Handle the startup config
				// for _, cmd := range commands {
				// 	if cmd.RootOperation != command.JSONOperationStreamSubscribe {
				// 		panic(fmt.Errorf("startup config can only contain stream subscribe requests for now"))
				// 	}
				// 	handler.HandleStreamRequest(cmd.StreamSubscribeRequest)
				// }
			}

			go func() {
				<-sigs
				cmdHandler.Cancel()
			}()
			<-cmdHandler.Ctx().Done()
			cmdHandler.Wg.Wait()
		},
	}
	rootCmd.Flags().StringP("startup-commands", "c", "", "Path to the config startup commands")
	return &rootCmd
}
