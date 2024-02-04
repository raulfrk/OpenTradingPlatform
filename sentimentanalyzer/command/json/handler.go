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

	// Register cancel function
	if jsonCommand.CancelKey != "" && jsonCommand.RootOperation != command.JSONOperationCancel {
		cancelKey := jsonCommand.CancelKey
		err := command.AddCancelFunc(cancelKey, ctx.Value(command.CancelKey{}).(context.CancelFunc))
		if err != nil {
			logging.Log().Error().Str("key", cancelKey).Err(err).Msg("adding cancel function")
			return types.NewError(err).Respond()
		}
		logging.Log().Info().Str("key", cancelKey).Msg("added cancel function")
		defer command.RemoveCancelFunc(cancelKey)
	}

	if jsonCommand.RootOperation == command.JSONOperationCancel {
		cancelFunc, found := command.GetCancelFunc(jsonCommand.CancelKey)
		cancelKey := jsonCommand.CancelKey
		if !found {
			err := fmt.Errorf("cancel function not found for key %s", cancelKey)
			logging.Log().Error().Str("key", cancelKey).Err(err).Msg("getting cancel function")
			return types.NewError(err).Respond()
		}
		cancelFunc()
		logging.Log().Info().Str("key", cancelKey).Msg("called cancel function")
		return types.NewResponse(
			types.Success,
			"Cancelled operation",
			nil,
		).Respond()
	}

	if jsonCommand.RootOperation == command.JSONOperationQuit {
		command.GetCommandHandler().Cancel()
		return types.NewResponse(
			types.Success,
			"Gracefully shutting down SentimentAnalyzer",
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
		// Create a new request that is validated
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
