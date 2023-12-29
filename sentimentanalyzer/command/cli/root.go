package cli

import (
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
