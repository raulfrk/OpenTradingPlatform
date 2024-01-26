package cli

import (
	"tradingplatform/datastorage/handler"

	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"

	"github.com/spf13/cobra"
)

// Data-related commands (does nothing by itself)
func NewDataCmd() *cobra.Command {
	dataCmd := cobra.Command{
		Use:   "data",
		Short: "Command to handle data",
	}

	dataCmd.AddCommand(NewDataGetCmd())

	return &dataCmd
}

// Get all active streams
func NewDataGetCmd() *cobra.Command {
	dataGetCmd := cobra.Command{
		Use:   "get",
		Short: "Get data from datastorage.",

		Run: func(cmd *cobra.Command, args []string) {
			// Get flags
			source, _ := cmd.Flags().GetString("source")
			assetClass, _ := cmd.Flags().GetString("asset-class")
			symbol, _ := cmd.Flags().GetString("symbol")
			operation := "get"
			dataType, _ := cmd.Flags().GetString("data-type")
			account, _ := cmd.Flags().GetString("account")
			startTime, _ := cmd.Flags().GetInt64("start-time")
			endTime, _ := cmd.Flags().GetInt64("end-time")
			timeFrame, _ := cmd.Flags().GetString("time-frame")
			noConfirm, _ := cmd.Flags().GetBool("no-confirm")

			// Generate stream request from flags
			dataRequest, err := requests.NewDataRequestFromRaw(source,
				assetClass,
				symbol,
				operation,
				dataType,
				account,
				startTime,
				endTime,
				timeFrame,
				noConfirm)

			if err != nil {
				cmd.Print(types.NewDataError(err).Respond())
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

	dataGetCmd.Flags().StringP("source", "s", requests.DataDefaultSource,
		"Source of the data stream")
	dataGetCmd.Flags().StringP("symbol", "y", requests.DataDefaultSymbol,
		"Symbols")
	dataGetCmd.Flags().StringP("asset-class", "a", requests.DataDefaultAssetClass,
		"Asset class")
	dataGetCmd.Flags().StringP("data-type", "t", requests.DataDefaultDataType,
		"Type of data")
	dataGetCmd.Flags().StringP("account", "c", requests.DataDefaultAccount,
		"Account (not doing anything currently)")
	dataGetCmd.Flags().Int64P("start-time", "b", requests.DataDefaultStartTime,
		"Start time for the data")
	dataGetCmd.Flags().Int64P("end-time", "e", requests.DataDefaultEndTime,
		"End time for the data")
	dataGetCmd.Flags().StringP("time-frame", "f", requests.DataDefaultTimeFrame,
		"Time frame (only for bar data)")
	dataGetCmd.Flags().BoolP("no-confirm", "o", requests.DataDefaultNoConfirm,
		"Setting this flag will make so that data is streamed as soon as ready")

	dataGetCmd.MarkFlagRequired("asset-class")
	dataGetCmd.MarkFlagRequired("symbol")
	dataGetCmd.MarkFlagRequired("start-time")
	dataGetCmd.MarkFlagRequired("end-time")

	return &dataGetCmd
}
