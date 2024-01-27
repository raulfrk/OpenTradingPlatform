package cli

import (
	"tradingplatform/dataprovider/handler"

	"tradingplatform/shared/logging"

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

// Create new data get command
func NewDataGetCmd() *cobra.Command {
	dataGetCmd := cobra.Command{
		Use:   "get",
		Short: "Get data from dataprovider.",

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
				noConfirm, requests.DefaultForEmptyDataRequest)
			logging.Log().Info().
				RawJSON("dataRequest", dataRequest.JSON()).
				Msg("receiving data request")
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
	dataGetCmd.Flags().StringP("source", "s", "",
		"Source of the data")
	dataGetCmd.Flags().StringP("symbol", "y", "",
		"Symbols")
	dataGetCmd.Flags().StringP("asset-class", "a", "",
		"Asset class")
	dataGetCmd.Flags().StringP("data-type", "t", "",
		"Type of data (e.g. bar, trade...)")
	dataGetCmd.Flags().StringP("account", "c", "",
		"Account to use for the stream")
	dataGetCmd.Flags().Int64P("start-time", "b", 0,
		"Start time for the data")
	dataGetCmd.Flags().Int64P("end-time", "e", 0,
		"End time for the data")
	dataGetCmd.Flags().StringP("time-frame", "f", "",
		"Time frame (only available for bar data)")
	dataGetCmd.Flags().BoolP("no-confirm", "o", false,
		"Setting this flag will make so that data is streamed as soon as ready")

	return &dataGetCmd
}
