package cli

import (
	"tradingplatform/dataprovider/handler"
	"tradingplatform/dataprovider/provider"
	"tradingplatform/dataprovider/requests"

	"tradingplatform/shared/types"

	"github.com/spf13/cobra"
)

// Stream-related commands (does nothing by itself)
func NewStreamCmd() *cobra.Command {
	streamCmd := cobra.Command{
		Use:   "stream",
		Short: "Command to handle data streams",
		Long:  `This command allows the user to add, remove and see active data streams.`,
	}

	streamCmd.AddCommand(NewStreamAddCmd())
	streamCmd.AddCommand(NewStreamDeleteCmd())
	streamCmd.AddCommand(NewStreamGetCmd())

	return &streamCmd
}

// Add a new data stream
func NewStreamAddCmd() *cobra.Command {
	streamAddCmd := cobra.Command{
		Use:   "add",
		Short: "Add a new data stream",
		Long:  `This command allows the user to add a new data stream.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get flags
			source, _ := cmd.Flags().GetString("source")
			assetClass, _ := cmd.Flags().GetString("assetClass")
			symbols, _ := cmd.Flags().GetStringArray("symbols")
			operation := "add"
			dataTypes, _ := cmd.Flags().GetStringArray("dataTypes")
			account, _ := cmd.Flags().GetString("account")

			if len(dataTypes) == 0 {
				providerDatatypes := requests.GetDataTypeMap()[requests.DataSource(source)]
				for _, dtype := range providerDatatypes(types.AssetClass(assetClass)) {
					dataTypes = append(dataTypes, string(dtype))
				}
			}
			// Generate stream request from flags
			streamRequest, err := requests.NewStreamRequestFromRaw(source, assetClass, symbols, operation, dataTypes, account)

			if err != nil {
				cmd.Print(provider.NewStreamError(err).Respond())
				return
			}

			response := handler.HandleStreamRequest(streamRequest)
			cmd.Print(response)

		},
	}

	streamAddCmd.Flags().StringP("source", "s", requests.StreamDefaultSource, "Source of the data stream")
	streamAddCmd.Flags().StringArrayP("symbols", "y", []string{}, "Symbols to stream")
	streamAddCmd.Flags().StringP("assetClass", "a", requests.StreamDefaultAssetClass, "Type of asset to stream")
	streamAddCmd.Flags().StringArrayP("dataTypes", "t", []string{}, "Type of data to stream")
	streamAddCmd.Flags().StringP("account", "c", requests.StreamDefaultAccount, "Account to use for the stream")

	streamAddCmd.MarkFlagRequired("symbols")
	streamAddCmd.MarkFlagRequired("assetClass")

	return &streamAddCmd
}

// Get all active streams
func NewStreamGetCmd() *cobra.Command {
	streamGetCmd := cobra.Command{
		Use:   "get",
		Short: "Get all active streams from a given source on a given account.",

		Run: func(cmd *cobra.Command, args []string) {
			// Get flags
			source, _ := cmd.Flags().GetString("source")
			assetClass, _ := cmd.Flags().GetString("assetClass")
			operation := "get"
			account, _ := cmd.Flags().GetString("account")

			// Generate stream request from flags
			streamRequest, err := requests.NewStreamRequestFromRaw(source, assetClass, []string{}, operation, []string{}, account)

			if err != nil {
				cmd.Print(provider.NewStreamError(err).Respond())
				return
			}
			response := handler.HandleStreamRequest(streamRequest)
			cmd.Print(response)
		},
	}
	streamGetCmd.Flags().StringP("source", "s", requests.StreamDefaultSource, "Source of the data stream")
	streamGetCmd.Flags().StringP("assetClass", "a", requests.StreamDefaultAssetClass, "Type of asset to stream")
	streamGetCmd.Flags().StringP("account", "c", requests.StreamDefaultAccount, "Account to use for the stream")

	streamGetCmd.MarkFlagRequired("assetClass")
	return &streamGetCmd
}

func NewStreamDeleteCmd() *cobra.Command {
	streamAddCmd := cobra.Command{
		Use:   "remove",
		Short: "Remove a data stream",
		Long:  `This command allows the user to remove one or multiple data streams.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get flags
			source, _ := cmd.Flags().GetString("source")
			assetClass, _ := cmd.Flags().GetString("assetClass")
			symbols, _ := cmd.Flags().GetStringArray("symbols")
			operation := "remove"
			dataTypes, _ := cmd.Flags().GetStringArray("dataTypes")
			account, _ := cmd.Flags().GetString("account")

			if len(dataTypes) == 0 {
				providerDatatypes := requests.GetDataTypeMap()[requests.DataSource(source)]
				for _, dtype := range providerDatatypes(types.AssetClass(assetClass)) {
					dataTypes = append(dataTypes, string(dtype))
				}
			}
			// Generate stream request from flags
			streamRequest, err := requests.NewStreamRequestFromRaw(source, assetClass, symbols, operation, dataTypes, account)

			if err != nil {
				cmd.Print(provider.NewStreamError(err).Respond())
				return
			}

			response := handler.HandleStreamRequest(streamRequest)
			cmd.Print(response)

		},
	}

	streamAddCmd.Flags().StringP("source", "s", requests.StreamDefaultSource, "Source of the data stream")
	streamAddCmd.Flags().StringArrayP("symbols", "y", []string{}, "Symbols to stream")
	streamAddCmd.Flags().StringP("assetClass", "a", requests.StreamDefaultAssetClass, "Type of asset to stream")
	streamAddCmd.Flags().StringArrayP("dataTypes", "t", []string{}, "Type of data to stream")
	streamAddCmd.Flags().StringP("account", "c", requests.StreamDefaultAccount, "Account to use for the stream")

	streamAddCmd.MarkFlagRequired("symbols")
	streamAddCmd.MarkFlagRequired("assetClass")

	return &streamAddCmd
}
