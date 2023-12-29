package json

import (
	JSON "encoding/json"
	"tradingplatform/dataprovider/handler"
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
		return handler.HandleStreamRequest(jsonCommand.StreamRequest.ApplyDefault())
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
