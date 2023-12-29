package json

import (
	"context"
	JSON "encoding/json"
	"tradingplatform/dataprovider/handler"
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"
)

// Handle a JSON command
func HandleJSONCommand(ctx context.Context, jsonStr string) string {
	var jsonCommand command.JSONCommand
	err := JSON.Unmarshal([]byte(jsonStr), &jsonCommand)
	if err != nil {
		return types.NewError(err).Respond()
	}

	if jsonCommand.RootOperation == command.JSONOperationQuit {
		command.GetCommandHandler().Cancel()
		return types.NewResponse(
			types.Success,
			"Gracefully shutting down DataProvider",
			nil,
		).Respond()
	}

	if jsonCommand.RootOperation == command.JSONOperationStream {
		var streamRequest requests.StreamRequest
		err := JSON.Unmarshal(jsonCommand.Request, &streamRequest)
		if err != nil {
			return types.NewError(err).Respond()
		}
		return handler.HandleStreamRequest(streamRequest.ApplyDefault())
	}

	if jsonCommand.RootOperation == command.JSONOperationData {
		var dataRequest requests.DataRequest
		err := JSON.Unmarshal(jsonCommand.Request, &dataRequest)
		if err != nil {
			return types.NewError(err).Respond()
		}
		var och chan types.DataResponse = make(chan types.DataResponse)
		go handler.HandleDataRequest(dataRequest.ApplyDefault(), och)

		select {
		case response := <-och:
			return response.Respond()
		case <-command.GetCommandHandler().Ctx().Done():
			return ""
		}
	}
	return ""
}
