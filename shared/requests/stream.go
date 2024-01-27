package requests

import (
	"encoding/json"
	"strconv"
	"strings"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"

	"github.com/go-playground/validator/v10"
)

type Account string

const (
	AnyAccount     Account = "any"
	DefaultAccount Account = "default"
)

// Datastructure to represent a stream request
type StreamRequest struct {
	Source     types.Source          `json:"source" validate:"required,min=3,isValidDataSource"`
	AssetClass types.AssetClass      `json:"assetClass" validate:"required,min=3,isValidAssetClass"`
	Symbols    []string              `json:"symbols" validate:"required,isValidSymbols"`
	Operation  types.StreamRequestOp `json:"operation" validate:"required,min=3,isValidOperation"`
	DataTypes  []types.DataType      `json:"dataTypes" validate:"required,isValidMultiDataType"`
	Account    Account               `json:"account"`
}

func (sr *StreamRequest) Validate() error {
	v := validator.New()
	v.RegisterValidation("isValidDataSource", IsValidDataSource)
	v.RegisterValidation("isValidAssetClass", IsValidAssetClass)
	v.RegisterValidation("isValidOperation", IsValidOperationStream)
	v.RegisterValidation("isValidMultiDataType", IsValidMultiDataTypeStream)
	v.RegisterValidation("isValidAccount", IsValidAccount)
	v.RegisterValidation("isValidSymbols", IsValidMultiSymbolStream)

	err := v.Struct(sr)
	return SummarizeError(err)
}

type StreamSubscribeAgents struct {
	AgentCount int    `json:"agentCount" validate:"required,min=1"`
	Topic      string `json:"topic" validate:"required,min=3"`
}

func (ssa *StreamSubscribeAgents) Validate() error {
	v := validator.New()
	err := v.Struct(ssa)
	return SummarizeError(err)
}

type StreamSubscribeRequest struct {
	StreamSubscribeWithAgents []StreamSubscribeAgents `json:"streamSubscribeWithAgents" validate:"-"`
	Operation                 types.StreamRequestOp   `json:"operation" validate:"required,min=3,isValidOperation"`
}

func (ssr *StreamSubscribeRequest) Validate() error {
	for _, ssa := range ssr.StreamSubscribeWithAgents {
		err := ssa.Validate()
		if err != nil {
			return err
		}
	}
	v := validator.New()
	v.RegisterValidation("isValidOperation", IsValidOperationStreamSubscribe)
	err := v.Struct(ssr)
	return SummarizeError(err)
}

func (sr *StreamRequest) JSON() []byte {
	js, err := json.Marshal(sr)
	if err != nil {
		logging.Log().Error().
			Err(err).
			Msg("marshalling stream request to json")
		return []byte{}
	}
	return js
}

func (sr *StreamSubscribeRequest) JSON() []byte {
	js, err := json.Marshal(sr)
	if err != nil {
		logging.Log().Error().
			Err(err).
			Msg("marshalling stream subscribe request to json")
		return []byte{}
	}
	return js
}

func NewStreamSubscribeRequestFromRaw(iTopicAgents []string,
	operation types.StreamRequestOp) (StreamSubscribeRequest, error) {
	var streamSubscribeWithAgents []StreamSubscribeAgents
	for _, t := range iTopicAgents {
		// Split topic in topic and agents
		splitT := strings.Split(t, ",")
		agents := 5
		topic := splitT[0]
		if len(splitT) == 2 {
			agentsStr := splitT[1]
			agentsInt, err := strconv.Atoi(agentsStr)
			if err != nil {
				return StreamSubscribeRequest{}, err
			}
			agents = agentsInt
		}
		streamSubscribeWithAgents = append(streamSubscribeWithAgents, StreamSubscribeAgents{
			AgentCount: agents,
			Topic:      topic,
		})
	}
	streamSubscribeRequest := StreamSubscribeRequest{
		StreamSubscribeWithAgents: streamSubscribeWithAgents,
		Operation:                 operation,
	}
	err := streamSubscribeRequest.Validate()
	return streamSubscribeRequest, err
}

func NewStreamSubscribeRequestFromExisting(req *StreamSubscribeRequest) (StreamSubscribeRequest, error) {
	agents := make([]string, len(req.StreamSubscribeWithAgents))
	for i, ssa := range req.StreamSubscribeWithAgents {
		if ssa.AgentCount == 0 {
			agents[i] = ssa.Topic
		} else {
			agents[i] = ssa.Topic + "," + strconv.Itoa(ssa.AgentCount)
		}
	}
	return NewStreamSubscribeRequestFromRaw(agents, req.Operation)

}

func NewStreamRequest(source types.Source,
	assetClass types.AssetClass,
	symbols []string,
	operation types.StreamRequestOp,
	dataTypes []types.DataType,
	account Account) StreamRequest {

	return StreamRequest{
		Source:     source,
		AssetClass: assetClass,
		Symbols:    symbols,
		Operation:  operation,
		DataTypes:  dataTypes,
		Account:    account,
	}
}

func NewStreamRequestFromRaw(source string,
	assetClass string,
	symbols []string,
	operation string,
	dataTypes []string,
	account string, defaultingFunc func(*StreamRequest)) (StreamRequest, error) {

	var convertedDataTypes = make([]types.DataType, len(dataTypes))
	for i, dataType := range dataTypes {
		convertedDataTypes[i] = types.DataType(dataType)
	}
	streamRequest := NewStreamRequest(types.Source(source),
		types.AssetClass(assetClass),
		symbols,
		types.StreamRequestOp(operation),
		convertedDataTypes,
		Account(account),
	)
	defaultingFunc(&streamRequest)
	err := streamRequest.Validate()
	return streamRequest, err
}

func NewStreamRequestFromExisting(req *StreamRequest, defaultingFunc func(*StreamRequest)) (StreamRequest, error) {
	return NewStreamRequestFromRaw(string(req.Source),
		string(req.AssetClass),
		req.Symbols,
		string(req.Operation),
		req.GetStrDataTypes(),
		string(req.Account),
		defaultingFunc)
}

func (sr *StreamRequest) GetSource() types.Source {
	return sr.Source
}

func (sr *StreamRequest) GetAssetClass() types.AssetClass {
	return sr.AssetClass
}

func (sr *StreamRequest) GetSymbol() []string {
	return sr.Symbols
}

func (sr *StreamRequest) GetOperation() types.StreamRequestOp {
	return sr.Operation
}

func (sr *StreamRequest) GetDataType() []types.DataType {
	return sr.DataTypes
}
func (sr *StreamRequest) GetStrDataTypes() []string {
	var strDataTypes []string
	for _, dataType := range sr.DataTypes {
		strDataTypes = append(strDataTypes, string(dataType))
	}
	return strDataTypes
}

func (sr *StreamRequest) GetAccount() Account {
	return sr.Account
}

func GetAccount() map[string]Account {
	return map[string]Account{
		"any":     AnyAccount,
		"default": DefaultAccount,
	}
}

func GetDataTypeMap() map[types.Source]func(types.AssetClass) map[types.DataType]types.DataType {
	return map[types.Source]func(types.AssetClass) map[types.DataType]types.DataType{
		types.Alpaca: getAlpacaDataType,
	}
}

func getAlpacaDataType(assetClass types.AssetClass) map[types.DataType]types.DataType {
	switch assetClass {
	case types.Stock:
		return map[types.DataType]types.DataType{
			"bar":          types.Bar,
			"daily-bars":   types.DailyBars,
			"quotes":       types.Quotes,
			"trades":       types.Trades,
			"updated-bars": types.UpdatedBars,
			"luld":         types.LULD,
			"status":       types.Status,
		}
	case types.Crypto:
		return map[types.DataType]types.DataType{
			"bar":          types.Bar,
			"orderbook":    types.Orderbook,
			"daily-bars":   types.DailyBars,
			"quotes":       types.Quotes,
			"trades":       types.Trades,
			"updated-bars": types.UpdatedBars,
		}
	case types.News:
		return map[types.DataType]types.DataType{
			"raw-text": types.RawText,
		}
	default:
		return nil
	}
}

func GetDataSourceMap() map[string]types.Source {
	return map[string]types.Source{
		"alpaca":   types.Alpaca,
		"internal": types.Internal,
	}
}
