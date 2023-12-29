package requests

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"
)

type Account string

const (
	AnyAccount     Account = "any"
	DefaultAccount Account = "default"
)
const (
	StreamDefaultSource     = "alpaca"
	StreamDefaultAssetClass = "stock"
	StreamDefaultAccount    = "default"
)

// Datastructure to represent a stream request
type StreamRequest struct {
	Source     types.Source          `json:"source"`
	AssetClass types.AssetClass      `json:"assetClass"`
	Symbols    []string              `json:"symbols"`
	Operation  types.StreamRequestOp `json:"operation"`
	DataTypes  []types.DataType      `json:"dataTypes"`
	Account    Account               `json:"account"`
}

type StreamSubscribeAgents struct {
	AgentCount int    `json:"agentCount"`
	Topic      string `json:"topic"`
}

type StreamSubscribeRequest struct {
	StreamSubscribeWithAgents []StreamSubscribeAgents `json:"streamSubscribeWithAgents"`
	Operation                 types.StreamRequestOp   `json:"operation"`
}

func (sr StreamRequest) ApplyDefault() StreamRequest {
	if sr.Source == "" {
		sr.Source = StreamDefaultSource
	}
	if sr.AssetClass == "" {
		sr.AssetClass = StreamDefaultAssetClass
	}
	if sr.Account == "" {
		sr.Account = StreamDefaultAccount
	}
	return sr
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
	return streamSubscribeRequest, nil
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

func NewStreamRequestFromRaw(iSource string,
	iAssetClass string,
	iSymbols []string,
	iOperation string,
	iDataTypes []string,
	iAccount string) (StreamRequest, error) {

	source, exists := GetDataSourceMap()[iSource]

	if !exists {
		return StreamRequest{}, errors.New("Invalid data source: " + iSource)
	}

	assetClass, exists := types.GetAssetClassMap()[iAssetClass]

	if !exists {
		return StreamRequest{}, errors.New("Invalid asset class: " + iAssetClass)
	}

	operation, exists := types.GetStreamRequestOpMap()[iOperation]

	if !exists {
		return StreamRequest{}, errors.New("Invalid operation: " + iOperation)
	}

	dataTypes := make([]types.DataType, len(iDataTypes))

	if len(iDataTypes) == 0 {
		return StreamRequest{}, errors.New("no data types specified")
	}

	if len(iSymbols) == 0 {
		return StreamRequest{}, errors.New("no symbols specified")
	}

	for i, dType := range iDataTypes {
		dataTypeMap := GetDataTypeMap()[source]
		dataType, exists := dataTypeMap(assetClass)[types.DataType(dType)]

		if !exists {
			return StreamRequest{}, errors.New("Invalid data type: " + dType)
		}

		dataTypes[i] = dataType
	}

	account, exists := GetAccount()[iAccount]

	if !exists {
		return StreamRequest{}, errors.New("Invalid account type: " + iAccount)
	}

	streamRequest := NewStreamRequest(source,
		assetClass,
		iSymbols,
		operation,
		dataTypes,
		account,
	)
	return streamRequest, nil
}

func (sr *StreamRequest) GetSource() types.Source {
	return sr.Source
}

func (sr *StreamRequest) GetAssetClass() types.AssetClass {
	return sr.AssetClass
}

func (sr *StreamRequest) GetSymbols() []string {
	return sr.Symbols
}

func (sr *StreamRequest) GetOperation() types.StreamRequestOp {
	return sr.Operation
}

func (sr *StreamRequest) GetDataTypes() []types.DataType {
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
