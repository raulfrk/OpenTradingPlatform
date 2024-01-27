package cli

import (
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/types"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "datastorage",
		Short: "DataStorage for the data pipeline",
	}

	rootCmd.AddCommand(NewQuitCommand())
	rootCmd.AddCommand(NewStreamCommand())
	rootCmd.AddCommand(NewDataCmd())

	return &rootCmd
}

func NewQuitCommand() *cobra.Command {
	quitCmd := cobra.Command{
		Use:   "quit",
		Short: "Gracefully shuts down the DataStorage",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Print(types.NewResponse(
				types.Success,
				"Gracefully shutting down DataStorage",
				nil,
			).Respond())
			command.GetCommandHandler().Cancel()
		},
	}

	return &quitCmd
}
