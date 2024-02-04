package requests

import (
	"encoding/json"
	"fmt"
	"strings"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"

	"github.com/go-playground/validator/v10"
)

type SentimentAnalysisRequest struct {
	DataRequest              `validate:"-"`
	SentimentAnalysisProcess types.SentimentAnalysisProcess `json:"sentimentAnalysisProcess" validate:"required,min=3,isValidSentimentAnalysisProcess"`
	Model                    string                         `json:"model" validate:"required,min=3"`
	ModelProvider            types.LLMProvider              `json:"modelProvider" validate:"required,min=3,isValidModelProvider"`
	SystemPrompt             string                         `json:"systemPrompt" validate:"required,min=3"`
	FailFastOnBadSentiment   bool                           `json:"failFastOnBadSentiment"`
	RetryFailed              bool                           `json:"retryFailed"`
}

func (sar *SentimentAnalysisRequest) Validate() error {
	v := validator.New()
	v.RegisterValidation("isValidSentimentAnalysisProcess", IsValidSentimentAnalysisProcess)
	v.RegisterValidation("isValidModelProvider", IsValidModelProvider)

	err := v.Struct(sar)
	return SummarizeError(err)

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
	sentimentAnalysisProcess string,
	model string,
	systemPrompt string,
	failFastOnBadSentiment bool,
	retryFailed bool,
	defaultingFunc func(*SentimentAnalysisRequest),
) (SentimentAnalysisRequest, error) {

	// Validate model
	providerV, modelV, err := extractProviderModel(model)
	if err != nil {
		return SentimentAnalysisRequest{}, err
	}

	req := SentimentAnalysisRequest{
		DataRequest:              dataRequest,
		SentimentAnalysisProcess: types.SentimentAnalysisProcess(sentimentAnalysisProcess),
		ModelProvider:            types.LLMProvider(providerV),
		Model:                    modelV,
		SystemPrompt:             systemPrompt,
		FailFastOnBadSentiment:   failFastOnBadSentiment,
		RetryFailed:              retryFailed,
	}

	defaultingFunc(&req)
	_, err = NewDataRequestFromExisting(&req.DataRequest, DefaultForEmptyDataRequest)
	if err != nil {
		return SentimentAnalysisRequest{}, err
	}
	err = req.Validate()
	return req, err

}

func NewSentimentAnalysisRequestFromExisting(req *SentimentAnalysisRequest, defaultingFunc func(*SentimentAnalysisRequest)) (SentimentAnalysisRequest, error) {

	return NewSentimentAnalysisRequestFromRaw(req.DataRequest,
		string(req.SentimentAnalysisProcess),
		req.GetFormattedModelProvider(),
		req.SystemPrompt,
		req.FailFastOnBadSentiment,
		req.RetryFailed,
		defaultingFunc,
	)
}

func GetModelProviderMap() map[string]types.LLMProvider {
	return map[string]types.LLMProvider{
		"ollama":  types.Ollama,
		"gpt4all": types.GPT4All,
	}
}
