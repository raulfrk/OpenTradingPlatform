package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"tradingplatform/sentimentanalyzer/llmproviders/ollama"
	"tradingplatform/sentimentanalyzer/sentiment"
	"tradingplatform/shared/communication/producer"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"

	"google.golang.org/protobuf/proto"
)

// HandleAnalysisRequest handles a sentiment analysis request
func HandleAnalysisRequest(ctx context.Context, req *requests.SentimentAnalysisRequest, och chan<- types.DataResponse) {
	switch req.ModelProvider {
	case types.Ollama:
		och <- HandleAnalysisNewsFromDB(ctx, req, ollama.HandleAnalysisNews)
	default:
		err := fmt.Errorf("provided model provider \"%s\" is not currently supported", req.ModelProvider)
		logging.Log().Debug().Err(err).RawJSON("request", req.JSON()).Msg("while handling sentiment analysis request")
		och <- types.NewDataError(err)
	}
}

func isValidJSON(s string) bool {
	var js map[string]interface{}
	err := json.Unmarshal([]byte(s), &js)
	return err == nil
}

func castToJSON(s string) map[string]interface{} {
	var js map[string]interface{}
	json.Unmarshal([]byte(s), &js)
	return js
}

func worker(ctx context.Context,
	jobsCh <-chan *entities.News,
	resultsCh chan<- *entities.News,
	errCh chan<- error,
	wg *sync.WaitGroup, req *requests.SentimentAnalysisRequest, analysisF func(context.Context, *entities.News, *requests.SentimentAnalysisRequest) (string, error)) {
	defer wg.Done()
	for n := range jobsCh {
		// Check if news already has sentiment with given LLM, symbol and process or if it already failed previously
		alreadyPresent := false
		failedShouldRetry := false

		for _, s := range n.Sentiments {
			matchesExisting := s.LLM == req.GetFormattedModelProvider() &&
				s.Symbol == req.GetSymbol() &&
				s.SentimentAnalysisProcess == string(req.SentimentAnalysisProcess) && s.SystemPrompt == req.SystemPrompt

			failedShouldRetry = s.Failed && req.RetryFailed
			if matchesExisting {
				alreadyPresent = true
				break
			}
		}
		if alreadyPresent && !failedShouldRetry {
			resultsCh <- n
			continue
		}

		if n.Headline == "" {
			logging.Log().Debug().RawJSON("request", req.JSON()).Str("newsFingerprint", req.Fingerprint).Msg("news headline is empty")
			continue
		}
		childCtx, cancel := context.WithCancel(ctx)
		analyzedSentiment, err := analysisF(childCtx, n, req)
		if err != nil {
			errCh <- err
			cancel()
			continue
		}
		if isValidJSON(analyzedSentiment) {
			handleJSONResponse(analyzedSentiment, n, req, resultsCh, errCh)
		} else {
			handlePlainResponse(analyzedSentiment, n, req, resultsCh, errCh)
		}
		cancel()

	}
}

func findMatchingSentiment(n *entities.News, newSentiment *entities.NewsSentiment) *entities.NewsSentiment {
	for _, s := range n.Sentiments {
		if s.LLM == newSentiment.LLM &&
			s.Symbol == newSentiment.Symbol &&
			s.SentimentAnalysisProcess == string(newSentiment.SentimentAnalysisProcess) && s.SystemPrompt == newSentiment.SystemPrompt {
			return s
		}
	}
	return nil
}

func handleJSONResponse(analyzedSentiment string, n *entities.News, req *requests.SentimentAnalysisRequest, resultsCh chan<- *entities.News, errCh chan<- error) {
	semanticSentiments := castToJSON(analyzedSentiment)
	if req.SentimentAnalysisProcess == types.Plain && req.FailFastOnBadSentiment {
		err := fmt.Errorf("produced sentiment resulted in JSON format when expected plain response for news; consider adjusting the prompt")
		logging.Log().Debug().
			Str("newsHeadline", n.Headline).
			Str("newsFingerprint", n.Fingerprint).
			RawJSON("request", req.JSON()).
			Str("sentiment", analyzedSentiment).
			Err(err).Msg("while handling JSON response from LLM")
		errCh <- err
		return
	}
	failed := req.SentimentAnalysisProcess == types.Plain

	// Assumes the JSON is in the format {"symbol": "sentiment"}
	for k, v := range semanticSentiments {
		vString, ok := v.(string)
		if !ok {
			err := fmt.Errorf("error casting sentiment to string")
			logging.Log().Info().Interface("sentiment", v).
				RawJSON("request", req.JSON()).
				Str("newsHeadline", n.Headline).
				Str("newsFingerprint", n.Fingerprint).
				Str("symbol", k).
				Str("sentiment", analyzedSentiment).
				Err(err).
				Msg("while casting sentiment to string in JSON response handler")
			failed = true
		}
		if failed {
			vString = ""
		}
		sentiment := entities.NewsSentiment{
			Timestamp:                time.Now().Unix(),
			Sentiment:                vString,
			SentimentAnalysisProcess: string(req.SentimentAnalysisProcess),
			News:                     n,
			LLM:                      req.GetFormattedModelProvider(),
			Symbol:                   k,
			SystemPrompt:             req.SystemPrompt,
			Failed:                   failed,
		}
		existingSentiment := findMatchingSentiment(n, &sentiment)
		if existingSentiment != nil {
			sentiment.Fingerprint = existingSentiment.Fingerprint
		} else {
			sentiment.SetFingerprint()
		}
		// Reset news field to avoid circular dependency
		sentiment.News = nil
		n.Sentiments = append(n.Sentiments, &sentiment)
	}
	resultsCh <- n
}

func handlePlainResponse(analyzedSentiment string, n *entities.News, req *requests.SentimentAnalysisRequest, resultsCh chan<- *entities.News, errCh chan<- error) {
	extractedSentiment, err := sentiment.ExtractSentimentFromLLMAnswer(analyzedSentiment)
	if req.SentimentAnalysisProcess == types.Semantic && req.FailFastOnBadSentiment {
		err := fmt.Errorf("produced sentiment resulted in plain format when expected JSON response for news; consider adjusting the prompt to ensure plain sentiment is produced")
		logging.Log().Debug().
			Str("newsHeadline", n.Headline).
			Str("newsFingerprint", n.Fingerprint).
			RawJSON("request", req.JSON()).
			Str("sentiment", analyzedSentiment).
			Err(err).Msg("while handling plain response from LLM")
		errCh <- err
		return
	}
	failed := req.SentimentAnalysisProcess == types.Semantic

	if err != nil {
		logging.Log().Debug().
			Str("newsHeadline", n.Headline).
			Str("newsFingerprint", n.Fingerprint).
			RawJSON("request", req.JSON()).
			Str("sentiment", analyzedSentiment).
			Err(err).Msg("while extracting sentiment from plain response from LLM")
		failed = true
	}

	if failed {
		extractedSentiment = ""
	}

	sentiment := entities.NewsSentiment{
		Timestamp:                time.Now().Unix(),
		Sentiment:                extractedSentiment,
		SentimentAnalysisProcess: string(req.SentimentAnalysisProcess),
		News:                     n,
		LLM:                      req.GetFormattedModelProvider(),
		Symbol:                   req.GetSymbol(),
		SystemPrompt:             req.SystemPrompt,
		Failed:                   failed,
	}

	existingSentiment := findMatchingSentiment(n, &sentiment)
	if existingSentiment != nil {
		sentiment.Fingerprint = existingSentiment.Fingerprint
	} else {
		sentiment.SetFingerprint()
	}
	// Reset news field to avoid circular dependency
	sentiment.News = nil
	n.Sentiments = append(n.Sentiments, &sentiment)
	resultsCh <- n
}

// Handle analysis request for news in the database
func HandleAnalysisNewsFromDB(ctx context.Context, req *requests.SentimentAnalysisRequest, analysisF func(context.Context, *entities.News, *requests.SentimentAnalysisRequest) (string, error)) types.DataResponse {
	var news []*entities.News
	logging.Log().Debug().RawJSON("request", req.JSON()).Msg("fetching and analyzing news from database")

	err := requests.RequestData(ctx, utils.NewCommandTopic(types.DataStorage), req.DataRequest, func(msg *entities.Message) {
		var entity entities.News
		err := proto.Unmarshal(msg.Payload, &entity)
		if err != nil {
			logging.Log().Debug().Err(err).Msg("while unmarshalling news from database")
			return
		}
		news = append(news, &entity)
	})

	if err != nil {
		logging.Log().Debug().Err(err).Msg("while requesting data from datastorage")
		return types.NewDataError(fmt.Errorf("error while requesting data from datastorage %v", err))
	}

	// Create a buffered channel for jobs and results
	jobs := make(chan *entities.News, len(news))
	results := make(chan *entities.News, len(news))

	// Add jobs to the channel
	for _, n := range news {
		jobs <- n
	}

	close(jobs)
	logging.Log().Debug().RawJSON("request", req.JSON()).Msg("starting sentiment analysis")

	errCh := make(chan error)
	done := make(chan struct{})

	var wg sync.WaitGroup
	// Start workers
	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go worker(ctx, jobs, results, errCh, &wg, req, analysisF)
	}

	go func() {
		// Wait for all workers to finish
		wg.Wait()
		close(done)
		close(results)
	}()

	select {
	case <-done:
		logging.Log().Debug().RawJSON("request", req.JSON()).Msg("finished processing news, generating response queue")
		var news []*entities.News
		for n := range results {
			news = append(news, n)
		}

		sort.Slice(news, func(i, j int) bool {
			return news[i].UpdatedAt < news[j].UpdatedAt
		})

		responseTopic := utils.NewDataTopic(types.SentimentAnalyzer, types.Internal, types.News, types.NewsWithSentiment, req.GetSymbol(), news[len(news)-1].Fingerprint, len(news)).Generate()
		handler, handlerResponse := producer.GetQueueHandler(responseTopic, req.NoConfirm)
		if handlerResponse.Err != "" {
			return handlerResponse
		}
		var messages []*entities.Message
		for _, n := range news {
			message := entities.GenerateMessage(n, types.NewsWithSentiment, responseTopic)
			messages = append(messages, message)
		}

		handler.Ch <- &messages
		return types.NewDataResponse(
			types.Success,
			fmt.Sprintf("successfully processed %d news", len(news)),
			nil,
			responseTopic,
		)

	case err := <-errCh:
		requests.DrainChannel(jobs)
		logging.Log().Debug().RawJSON("request", req.JSON()).Err(err).Msg("failed processing news")
		return types.NewDataError(err)
	}

}
