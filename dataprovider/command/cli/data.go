package cli

import (
	"tradingplatform/dataprovider/handler"
	"tradingplatform/dataprovider/provider"
	"tradingplatform/dataprovider/requests"
	"tradingplatform/shared/types"

	"github.com/spf13/cobra"
)

// Stream-related commands (does nothing by itself)
func NewDataCmd() *cobra.Command {
	dataCmd := cobra.Command{
		Use:   "data",
		Short: "Command to handle data",
		Long:  `This command allows the user to get data`,
	}

	dataCmd.AddCommand(NewDataGetCmd())

	return &dataCmd
}

// Get all active streams
func NewDataGetCmd() *cobra.Command {
	dataGetCmd := cobra.Command{
		Use:   "get",
		Short: "Get all active streams from a given source on a given account.",

		Run: func(cmd *cobra.Command, args []string) {
			// Get flags
			source, _ := cmd.Flags().GetString("source")
			assetClass, _ := cmd.Flags().GetString("assetClass")
			symbols, _ := cmd.Flags().GetString("symbol")
			operation := "get"
			dataTypes, _ := cmd.Flags().GetString("dataType")
			account, _ := cmd.Flags().GetString("account")
			startTime, _ := cmd.Flags().GetInt64("startTime")
			endTime, _ := cmd.Flags().GetInt64("endTime")
			timeFrame, _ := cmd.Flags().GetString("timeFrame")
			noConfirm, _ := cmd.Flags().GetBool("noConfirm")

			// Generate stream request from flags
			dataRequest, err := requests.NewDataRequestFromRaw(source, assetClass, []string{symbols}, operation, []string{dataTypes}, account, startTime, endTime, timeFrame, noConfirm)

			if err != nil {
				cmd.Print(provider.NewStreamError(err).Respond())
				return
			}
			var och chan types.DataResponse = make(chan types.DataResponse)
			go handler.HandleDataRequest(dataRequest, och)

			select {
			case response := <-och:
				cmd.Print(response.Respond())
				return
			case <-cmd.Context().Done():
				return
			}
		},
	}
	dataGetCmd.Flags().StringP("source", "s", requests.DataDefaultSource, "Source of the data stream")
	dataGetCmd.Flags().StringP("symbol", "y", requests.DataDefaultSymbol, "Symbols to stream")
	dataGetCmd.Flags().StringP("assetClass", "a", requests.DataDefaultAssetClass, "Type of asset to stream")
	dataGetCmd.Flags().StringP("dataType", "t", requests.DataDefaultDataType, "Type of data to stream")
	dataGetCmd.Flags().StringP("account", "c", requests.DataDefaultAccount, "Account to use for the stream")
	dataGetCmd.Flags().Int64P("startTime", "b", requests.DataDefaultStartTime, "Start time for the data")
	dataGetCmd.Flags().Int64P("endTime", "e", requests.DataDefaultEndTime, "End time for the data")
	dataGetCmd.Flags().StringP("timeFrame", "f", requests.DataDefaultTimeFrame, "Time frame")
	dataGetCmd.Flags().BoolP("noConfirm", "o", requests.DataDefaultNoConfirm, "Whether to confirm the operation")

	dataGetCmd.MarkFlagRequired("assetClass")
	dataGetCmd.MarkFlagRequired("symbol")
	dataGetCmd.MarkFlagRequired("startTime")
	dataGetCmd.MarkFlagRequired("endTime")

	return &dataGetCmd
}
