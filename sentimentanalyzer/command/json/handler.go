package json

import (
	"context"
	JSON "encoding/json"
	"fmt"

	"tradingplatform/sentimentanalyzer/handler"
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"
)

func HandleJSONCommand(ctx context.Context, jsonStr string) string {
	var jsonCommand command.JSONCommand
	err := JSON.Unmarshal([]byte(jsonStr), &jsonCommand)
	if err != nil {
		return types.NewError(err).Respond()
	}

	if jsonCommand.CancelKey != "" {
		cancelKey := jsonCommand.CancelKey
		err := command.AddCancelFunc(cancelKey, ctx.Value(command.CancelKey{}).(context.CancelFunc))
		if err != nil {
			logging.Log().Error().Err(err).Msg("Error adding cancel function")
			return types.NewError(err).Respond()
		}
		defer command.RemoveCancelFunc(cancelKey)
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
			fmt.Errorf("operation %s not yet implemented", jsonCommand.RootOperation),
		).Respond()
	}
	if jsonCommand.RootOperation == command.JSONOperationStreamSubscribe {
		// var streamSubscribeRequest requests.StreamSubscribeRequest
		// err := JSON.Unmarshal(jsonCommand.Request, &streamSubscribeRequest)
		// if err != nil {
		// 	return types.NewError(err).Respond()
		// }
		// return handler.HandleStreamRequest(streamSubscribeRequest).Respond()
		return types.NewError(
			fmt.Errorf("operation %s not yet implemented", jsonCommand.RootOperation),
		).Respond()
	}

	if jsonCommand.RootOperation == command.JSONOperationData {
		var dataRequest requests.SentimentAnalysisRequest
		err := JSON.Unmarshal(jsonCommand.Request, &dataRequest)
		if err != nil {
			return types.NewError(err).Respond()
		}
		validatedRequest, err := requests.NewSentimentAnalysisRequestFromExisting(&dataRequest,
			requests.DefaultForEmptySentimentAnalysisRequest)

		if err != nil {
			return types.NewError(err).Respond()
		}
		var och chan types.DataResponse = make(chan types.DataResponse)
		go handler.HandleAnalysisRequest(ctx, &validatedRequest, och)

		select {
		case response := <-och:
			return response.Respond()
		case <-command.GetCommandHandler().Ctx().Done():
			return ""
		}
	}

	return ""
}
