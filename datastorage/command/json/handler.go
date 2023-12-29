package json

import (
	JSON "encoding/json"
	"fmt"

	"tradingplatform/datastorage/handler"
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"
)

// Handle a JSON command
func HandleJSONCommand(jsonStr string) string {
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
		return handler.HandleStreamRequest(streamSubscribeRequest).Respond()
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
