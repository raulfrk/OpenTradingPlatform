package cli

import (
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/types"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "dataprovider",
		Short: "DataProvider for the data pipeline",
		Long:  `Provides market data and handles re-distribution.`,
	}

	rootCmd.AddCommand(NewStreamCmd())
	rootCmd.AddCommand(NewQuitCommand())
	rootCmd.AddCommand(NewDataCmd())

	return &rootCmd
}

func NewQuitCommand() *cobra.Command {
	quitCmd := cobra.Command{
		Use:   "quit",
		Short: "Gracefully shuts down the DataProvider",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Print(types.NewResponse(
				types.Success,
				"Gracefully shutting down DataProvider",
				nil,
			).Respond())
			command.GetCommandHandler().Cancel()
		},
	}
	return &quitCmd
}
