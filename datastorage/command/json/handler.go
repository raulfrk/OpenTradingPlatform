package json

import (
	JSON "encoding/json"
	"fmt"

	"tradingplatform/datastorage/handler"
	"tradingplatform/shared/communication/command"
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
		return handler.HandleStreamRequest(jsonCommand.StreamSubscribeRequest).Respond()
	}

	if jsonCommand.RootOperation == command.JSONOperationData {
		var och chan types.DataResponse = make(chan types.DataResponse)
		go handler.HandleDataRequest(jsonCommand.DataRequest.ApplyDefault(), och)

		select {
		case response := <-och:
			return response.Respond()
		case <-command.GetCommandHandler().Ctx().Done():
			return ""
		}
	}
	return ""
}
