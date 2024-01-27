package cli

import (
	"fmt"
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/types"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "sentalyzer",
		Short: "SentimentAnalyzer for the data pipeline",
	}

	rootCmd.AddCommand(NewQuitCommand())
	rootCmd.AddCommand(NewDataCmd())
	rootCmd.AddCommand(NewCancelCommand())

	return &rootCmd
}

func NewQuitCommand() *cobra.Command {
	quitCmd := cobra.Command{
		Use:   "quit",
		Short: "Gracefully shuts down the SentimentAnalyzer",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Print(types.NewResponse(
				types.Success,
				"Gracefully shutting down SentimentAnalyzer",
				nil,
			).Respond())
			command.GetCommandHandler().Cancel()
		},
	}

	return &quitCmd
}

func NewCancelCommand() *cobra.Command {
	cancelCmd := cobra.Command{
		Use:   "cancel",
		Short: "Cancel a running command",
		Run: func(cmd *cobra.Command, args []string) {
			key, _ := cmd.Flags().GetString("key")
			cancelFunc, ok := command.GetCancelFunc(key)
			if !ok {
				cmd.PrintErr(types.NewError(fmt.Errorf("no command with cancel key %s", key)).Respond())
				return
			}
			cancelFunc()
			command.RemoveCancelFunc(key)
			cmd.Print(types.NewResponse(
				types.Success,
				"Command canceled",
				nil,
			).Respond())
		},
	}

	cancelCmd.Flags().StringP("key", "k", "",
		"Key of the command to cancel")
	cancelCmd.MarkFlagRequired("key")

	return &cancelCmd
}
