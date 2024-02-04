package json

import (
	"context"
	JSON "encoding/json"
	"fmt"

	"tradingplatform/datastorage/handler"
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
			"Gracefully shutting down DataStorage",
			nil,
		).Respond()
	}

	if jsonCommand.RootOperation == command.JSONOperationStream {
		return types.NewError(
			fmt.Errorf("operation %s not supported", jsonCommand.RootOperation),
		).Respond()
	}
	if jsonCommand.RootOperation == command.JSONOperationStreamSubscribe {
		var streamSubscribeRequest requests.StreamSubscribeRequest
		err := JSON.Unmarshal(jsonCommand.Request, &streamSubscribeRequest)
		if err != nil {
			return types.NewError(err).Respond()
		}
		validatedStreamSubscribeRequest, err := requests.NewStreamSubscribeRequestFromExisting(&streamSubscribeRequest)
		if err != nil {
			return types.NewError(err).Respond()
		}
		return handler.HandleStreamRequest(validatedStreamSubscribeRequest).Respond()
	}

	if jsonCommand.RootOperation == command.JSONOperationData {
		var dataRequest requests.DataRequest
		err := JSON.Unmarshal(jsonCommand.Request, &dataRequest)
		if err != nil {
			return types.NewError(err).Respond()
		}
		validatedDataRequest, err := requests.NewDataRequestFromExisting(&dataRequest, requests.DefaultForEmptyDataRequest)
		if err != nil {
			return types.NewError(err).Respond()
		}
		var och chan types.DataResponse = make(chan types.DataResponse)
		go handler.HandleDataRequest(validatedDataRequest, och)

		select {
		case response := <-och:
			return response.Respond()
		case <-command.GetCommandHandler().Ctx().Done():
			return ""
		}
	}
	return ""
}
