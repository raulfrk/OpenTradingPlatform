package requests

import (
	"encoding/json"
	"fmt"
	"strings"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"
)

type SentimentAnalysisRequest struct {
	DataRequest
	SentimentAnalysisProcess types.SentimentAnalysisProcess `json:"sentimentAnalysisProcess"`
	Model                    string                         `json:"model"`
	ModelProvider            types.LLMProvider              `json:"modelProvider"`
	IgnoreFailedParsing      bool                           `json:"ignoreFailedParsing"`
	SystemPrompt             string                         `json:"systemPrompt"`
	FailFastOnBadSentiment   bool                           `json:"failFastOnBadSentiment"`
	RetryFailed              bool                           `json:"retryFailed"`
}

func (r SentimentAnalysisRequest) ApplyDefault() SentimentAnalysisRequest {
	if r.Source == "" {
		r.Source = DataDefaultSource
	}

	r.AssetClass = types.News

	if r.Account == "" {
		r.Account = DataDefaultAccount
	}

	r.Operation = types.DataGetOp
	r.DataTypes = []types.DataType{types.RawText}

	if r.SentimentAnalysisProcess == "" {
		r.SentimentAnalysisProcess = DefaultSentimentAnalysisProcess
	}
	return r
}

func (r *SentimentAnalysisRequest) JSON() []byte {
	js, err := json.Marshal(r)
	if err != nil {
		logging.Log().Error().
			Err(err).
			Msg("marshalling sentiment analysis request to json")
		return []byte{}
	}
	return js
}

func (req *SentimentAnalysisRequest) GetFormattedModelProvider() string {
	return fmt.Sprintf("%s/%s", req.ModelProvider, req.Model)
}

func extractProviderModel(iModel string) (string, string, error) {
	mSplit := strings.Split(iModel, "/")

	if len(mSplit) != 2 {
		return "", "", fmt.Errorf(
			"invalid model format %s. Expecting {provider}/{model}", iModel)
	}

	return mSplit[0], mSplit[1], nil

}

func (req *SentimentAnalysisRequest) GetSystemPrompt() (string, error) {
	if req.SystemPrompt == "" {
		return "", fmt.Errorf("system prompt cannot be empty")
	}
	return req.SystemPrompt, nil
}

func NewSentimentAnalysisRequestFromRaw(
	dataRequest DataRequest,
	iSentimentAnalysisProcess string,
	iModel string,
	iIgnoreFailedParsing bool,
	iSystemPrompt string,
	iFailFastOnBadSentiment bool,
	iRetryFailed bool,
) (SentimentAnalysisRequest, error) {
	// Validate sentiment analysis process
	sentProcess, exists := types.GetSentimentAnalysisProcessMap()[iSentimentAnalysisProcess]
	if !exists {
		return SentimentAnalysisRequest{}, fmt.Errorf("sentiment analysis process %s is not currently supproted", iSentimentAnalysisProcess)
	}
	// Validate model
	provider, model, err := extractProviderModel(iModel)
	if err != nil {
		return SentimentAnalysisRequest{}, err
	}

	verifiedProvider, exists := GetModelProviderMap()[provider]
	if !exists {
		return SentimentAnalysisRequest{},
			fmt.Errorf("provider %s is not currently supported", provider)
	}

	if iSystemPrompt == "" {
		return SentimentAnalysisRequest{}, fmt.Errorf("system prompt cannot be empty")
	}

	return SentimentAnalysisRequest{
		DataRequest:              dataRequest,
		SentimentAnalysisProcess: sentProcess,
		ModelProvider:            verifiedProvider,
		Model:                    model,
		IgnoreFailedParsing:      iIgnoreFailedParsing,
		SystemPrompt:             iSystemPrompt,
		FailFastOnBadSentiment:   iFailFastOnBadSentiment,
		RetryFailed:              iRetryFailed,
	}, nil
}

func GetModelProviderMap() map[string]types.LLMProvider {
	return map[string]types.LLMProvider{
		"ollama": types.Ollama,
	}
}
