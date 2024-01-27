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

func HandleAnalysisRequest(ctx context.Context, req *requests.SentimentAnalysisRequest, och chan<- types.DataResponse) {
	switch req.ModelProvider {
	case types.Ollama:
		och <- HandleAnalysisNewsFromDB(ctx, req, ollama.HandleAnalysisNews)
	default:
		och <- types.NewDataError(
			fmt.Errorf("provided model provider \"%s\" is not currently supported", req.ModelProvider),
		)
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
		logging.Log().Debug().Str("news", n.Fingerprint).Msg("processing news")

		if n.Headline == "" {
			logging.Log().Debug().Msg("news headline is empty")
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
			logging.Log().Debug().Str("news", n.Headline).Str("sentiment", analyzedSentiment).Msg("received plain response from LLM")
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
	mmap := castToJSON(analyzedSentiment)
	if req.SentimentAnalysisProcess == types.Plain && req.FailFastOnBadSentiment {
		err := fmt.Errorf("for news \"%s\" got sentiment \"%s\"; %v", n.Headline, analyzedSentiment, "received semantic response when expected plain response")
		logging.Log().Debug().Err(err).Msg("while handling JSON response from LLM")
		errCh <- err
		return
	}
	failed := req.SentimentAnalysisProcess == types.Plain

	// Assumes the JSON is in the format {"symbol": "sentiment"}
	for k, v := range mmap {
		vString, ok := v.(string)
		if !ok {
			logging.Log().Debug().Str("symbol", k).
				Interface("value", v).Msg("error while casting json value to string (expected structure: {\"symbol\": \"sentiment\"})")
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
		err := fmt.Errorf("for news \"%s\" got sentiment \"%s\"; %v", n.Headline, analyzedSentiment, "received plain response when expected semantic response in JSON format")
		logging.Log().Debug().Err(err).Msg("while handling response from LLM")
		errCh <- err
		return
	}
	failed := req.SentimentAnalysisProcess == types.Semantic

	if err != nil {
		err = fmt.Errorf("for news \"%s\"; %v", n.Headline, err)
		logging.Log().Debug().Err(err).Msg("while extracting sentiment")
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

	err := requests.RequestData(ctx, utils.NewCommandTopic(types.DataStorage), req.DataRequest, func(msg *entities.Message) {
		var entity entities.News
		err := proto.Unmarshal(msg.Payload, &entity)
		if err != nil {
			logging.Log().Debug().Err(err).Msg("error while unmarshalling news")
			return
		}
		news = append(news, &entity)
	})

	if err != nil {
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
		logging.Log().Debug().Err(err).Msg("error while processing news")
		return types.NewDataError(err)
	}

}