package json

import (
	JSON "encoding/json"
	"tradingplatform/dataprovider/handler"
	"tradingplatform/dataprovider/requests"
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/types"
)

type JSONOperation string

const (
	JSONOperationQuit   JSONOperation = "quit"
	JSONOperationStream JSONOperation = "stream"
	JSONOperationData   JSONOperation = "data"
)

type JSONCommand struct {
	RootOperation JSONOperation          `json:"operation"`
	StreamRequest requests.StreamRequest `json:"streamRequest"`
	DataRequest   requests.DataRequest   `json:"dataRequest"`
}

func HandleJSONCommand(jsonStr string) string {
	var jsonCommand JSONCommand
	err := JSON.Unmarshal([]byte(jsonStr), &jsonCommand)
	//TODO: Handle error
	if err != nil {
		return err.Error()
	}

	if jsonCommand.RootOperation == JSONOperationQuit {
		command.GetCommandHandler().Cancel()
		return types.NewResponse(
			types.Success,
			"Gracefully shutting down DataProvider",
			nil,
		).Respond()
	}

	if jsonCommand.RootOperation == JSONOperationStream {
		return handler.HandleStreamRequest(jsonCommand.StreamRequest.ApplyDefault())
	}

	if jsonCommand.RootOperation == JSONOperationData {
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
