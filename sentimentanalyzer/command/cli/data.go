package cli

import (
	"context"
	"tradingplatform/sentimentanalyzer/handler"
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"

	"github.com/spf13/cobra"
)

func NewDataCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "data",
		Short: "Data commands",
	}

	rootCmd.AddCommand(NewQuitCommand())
	rootCmd.AddCommand(NewAnalyzeFromDBCmd())

	return &rootCmd
}

func NewAnalyzeFromDBCmd() *cobra.Command {
	analyzeFromDBCmd := cobra.Command{
		Use:   "analyze",
		Short: "Analyzes news from the database",
		Run: func(cmd *cobra.Command, args []string) {
			symbol, _ := cmd.Flags().GetString("symbol")
			source, _ := cmd.Flags().GetString("source")
			systemPrompt, _ := cmd.Flags().GetString("system-prompt")
			cancelKey, _ := cmd.Flags().GetString("with-cancel-key")
			retryFailed, _ := cmd.Flags().GetBool("retry-failed")

			startTime, _ := cmd.Flags().GetInt64("start-time")
			endTime, _ := cmd.Flags().GetInt64("end-time")
			noConfirm, _ := cmd.Flags().GetBool("no-confirm")
			failFastOnBadSentiment, _ := cmd.Flags().GetBool("fail-fast-bad-sentiment")

			err := command.AddCancelFunc(cancelKey, cmd.Context().Value(command.CancelKey{}).(context.CancelFunc))
			if err != nil {
				logging.Log().Error().Err(err).Msg("Error adding cancel function")
				cmd.PrintErr(types.NewError(err).Respond())
				return
			}
			defer command.RemoveCancelFunc(cancelKey)

			model, _ := cmd.Flags().GetString("model")

			sentimentAnalysisProcess, _ := cmd.Flags().GetString("process")
			dataReq, err := requests.NewDataRequestFromRaw(
				source,
				string(types.News),
				symbol,
				"get",
				string(types.RawText),
				"",
				startTime,
				endTime,
				types.NoTimeFrame,
				noConfirm,
				requests.DefaultForEmptyDataRequest,
			)
			if err != nil {
				logging.Log().Error().Err(err).Msg("Error creating data request")
			}
			logging.Log().Debug().Msgf("System prompt: %s", systemPrompt)
			req, _ := requests.NewSentimentAnalysisRequestFromRaw(
				dataReq,
				sentimentAnalysisProcess,
				model,
				systemPrompt,
				failFastOnBadSentiment,
				retryFailed,
				requests.DefaultForEmptySentimentAnalysisRequest,
			)
			och := make(chan types.DataResponse)
			go handler.HandleAnalysisRequest(cmd.Context(), &req, och)

			select {
			case response := <-och:
				cmd.Print(response.Respond())
				return
			case <-cmd.Context().Done():
				return
			}

		},
	}
	analyzeFromDBCmd.Flags().StringP("source", "s", "",
		"Source of the news data")
	analyzeFromDBCmd.Flags().StringP("symbol", "y", "",
		"Symbols")
	analyzeFromDBCmd.Flags().StringP("system-prompt", "t", "",
		"System prompt for sentiment analysis")

	analyzeFromDBCmd.Flags().StringP("model", "m", "", `LLM to use for sentiment analysis. Format: 
	{provider}/{model} (e.g. ollama/llama2)`)
	analyzeFromDBCmd.Flags().Int64P("start-time", "b", 0,
		"Start time for the data")
	analyzeFromDBCmd.Flags().Int64P("end-time", "e", 0,
		"End time for the data")
	analyzeFromDBCmd.Flags().StringP("process", "p", "", "Sentiment analysis process")
	analyzeFromDBCmd.Flags().BoolP("no-confirm", "o", false,
		"Setting this flag will make so that data is streamed as soon as ready")
	analyzeFromDBCmd.Flags().BoolP("fail-fast-bad-sentiment", "f", false,
		"Whether to fail fast when invalid sentiment is detected")
	analyzeFromDBCmd.Flags().BoolP("retry-failed", "r", false,
		"Whether to retry sentiment analysis for news that failed previously")
	analyzeFromDBCmd.Flags().StringP("with-cancel-key", "c", "",
		"Set the cancellation key")

	return &analyzeFromDBCmd
}
