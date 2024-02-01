package json

import (
	"context"
	JSON "encoding/json"
	"tradingplatform/dataprovider/handler"
	shcommand "tradingplatform/shared/communication/command"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"
)

// Handle a JSON command
func HandleJSONCommand(ctx context.Context, jsonStr string) string {
	var jsonCommand shcommand.JSONCommand
	err := JSON.Unmarshal([]byte(jsonStr), &jsonCommand)
	if err != nil {
		return types.NewError(err).Respond()
	}

	if jsonCommand.RootOperation == shcommand.JSONOperationQuit {
		shcommand.GetCommandHandler().Cancel()
		return types.NewResponse(
			types.Success,
			"Gracefully shutting down DataProvider",
			nil,
		).Respond()
	}

	if jsonCommand.RootOperation == shcommand.JSONOperationStream {
		var streamRequest requests.StreamRequest
		err := JSON.Unmarshal(jsonCommand.Request, &streamRequest)
		if err != nil {
			return types.NewError(err).Respond()
		}
		validatedStreamRequest, err := requests.NewStreamRequestFromExisting(&streamRequest, requests.DefaultForEmptyStreamAddDeleteRequest)
		if err != nil {
			return types.NewError(err).Respond()
		}
		return handler.HandleStreamRequest(validatedStreamRequest)
	}

	if jsonCommand.RootOperation == shcommand.JSONOperationData {
		var dataRequest requests.DataRequest
		err := JSON.Unmarshal(jsonCommand.Request, &dataRequest)
		if err != nil {
			return types.NewError(err).Respond()
		}
		validDataRequest, err := requests.NewDataRequestFromExisting(&dataRequest, requests.DefaultForEmptyDataRequest)
		if err != nil {
			return types.NewError(err).Respond()
		}
		var och chan types.DataResponse = make(chan types.DataResponse)
		go handler.HandleDataRequest(validDataRequest, och)

		select {
		case response := <-och:
			return response.Respond()
		case <-shcommand.GetCommandHandler().Ctx().Done():
			return ""
		}
	}
	return ""
}
