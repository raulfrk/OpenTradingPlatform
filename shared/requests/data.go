package requests

import (
	"encoding/json"
	"errors"
	"tradingplatform/shared/logging"

	"tradingplatform/shared/types"
)

const (
	DataDefaultSource     = "alpaca"
	DataDefaultSymbol     = ""
	DataDefaultAssetClass = "stock"
	DataDefaultDataType   = ""
	DataDefaultAccount    = "default"
	DataDefaultStartTime  = 0
	DataDefaultEndTime    = 0
	DataDefaultTimeFrame  = "1min"
	DataDefaultNoConfirm  = false
)

type DataRequest struct {
	Source     DataSource          `json:"source"`
	AssetClass types.AssetClass    `json:"assetClass"`
	Symbols    []string            `json:"symbols"`
	Operation  types.DataRequestOp `json:"operation"`
	DataTypes  []types.DataType    `json:"dataTypes"`
	Account    Account             `json:"account"`
	StartTime  int64               `json:"startTime"`
	EndTime    int64               `json:"endTime"`
	TimeFrame  types.TimeFrame     `json:"timeFrame"`
	NoConfirm  bool                `json:"noConfirm"`
}

func (r *DataRequest) JSON() []byte {
	js, err := json.Marshal(r)
	if err != nil {
		logging.Log().Error().
			Err(err).
			Msg("marshalling data request to json")
		return []byte{}
	}
	return js
}
func (r DataRequest) ApplyDefault() DataRequest {
	if r.Source == "" {
		r.Source = DataDefaultSource
	}
	if r.AssetClass == "" {
		r.AssetClass = DataDefaultAssetClass
	}
	if r.Account == "" {
		r.Account = DataDefaultAccount
	}
	if r.StartTime == 0 {
		r.StartTime = DataDefaultStartTime
	}
	if r.EndTime == 0 {
		r.EndTime = DataDefaultEndTime
	}
	if r.TimeFrame == "" {
		r.TimeFrame = DataDefaultTimeFrame
	}
	return r
}

func (d *DataRequest) GetSource() DataSource {
	return d.Source
}

func (d *DataRequest) GetAssetClass() types.AssetClass {
	return d.AssetClass
}

func (d *DataRequest) GetSymbols() []string {
	return d.Symbols
}

func (d *DataRequest) GetOperation() types.DataRequestOp {
	return d.Operation
}

func (d *DataRequest) GetDataTypes() []types.DataType {
	return d.DataTypes
}

func (d *DataRequest) GetAccount() Account {
	return d.Account
}

func (d *DataRequest) GetStartTime() int64 {
	return d.StartTime
}

func (d *DataRequest) GetEndTime() int64 {
	return d.EndTime
}

func (d *DataRequest) GetTimeFrame() types.TimeFrame {
	return d.TimeFrame
}

func (d *DataRequest) GetNoConfirm() bool {
	return d.NoConfirm
}

func NewDataRequest(source DataSource,
	assetClass types.AssetClass,
	symbols []string,
	operation types.DataRequestOp,
	dataTypes []types.DataType,
	account Account,
	startTime int64,
	endTime int64,
	timeFrame types.TimeFrame,
	noConfirm bool) DataRequest {

	return DataRequest{
		Source:     source,
		AssetClass: assetClass,
		Symbols:    symbols,
		Operation:  operation,
		DataTypes:  dataTypes,
		Account:    account,
		StartTime:  startTime,
		EndTime:    endTime,
		TimeFrame:  timeFrame,
		NoConfirm:  noConfirm,
	}
}

func NewDataRequestFromRaw(iSource string,
	iAssetClass string,
	iSymbols []string,
	iOperation string,
	iDataTypes []string,
	iAccount string,
	iStartTime int64,
	iEndTime int64,
	iTimeFrame string,
	iNoConfirm bool) (DataRequest, error) {

	source, exists := GetDataSourceMap()[iSource]

	if !exists {
		return DataRequest{}, errors.New("Invalid data source: " + iSource)
	}

	assetClass, exists := types.GetAssetClassMap()[iAssetClass]

	if !exists {
		return DataRequest{}, errors.New("Invalid asset type: " + iAssetClass)
	}

	operation, exists := types.GetDataRequestOpMap()[iOperation]

	if !exists {
		return DataRequest{}, errors.New("Invalid operation: " + iOperation)
	}

	dataTypes := make([]types.DataType, len(iDataTypes))

	if len(iDataTypes) == 0 {
		return DataRequest{}, errors.New("no data types specified")
	}
	if len(iSymbols) == 0 {
		return DataRequest{}, errors.New("no symbols specified")
	}

	if assetClass == types.News {
		dataTypes = []types.DataType{types.RawText}
	} else {
		for i, dType := range iDataTypes {
			dataTypeMap := GetDataTypeMap()[source]
			dataType, exists := dataTypeMap(assetClass)[types.DataType(dType)]

			if !exists {
				return DataRequest{}, errors.New("Invalid data type: " + dType)
			}

			dataTypes[i] = dataType
		}
	}

	account, exists := GetAccount()[iAccount]

	if !exists {
		return DataRequest{}, errors.New("Invalid account type: " + iAccount)
	}

	timeFrame, exists := types.GetTimeFrameMap()[iTimeFrame]

	if !exists {
		return DataRequest{}, errors.New("Invalid time frame: " + iTimeFrame)
	}

	dataRequest := NewDataRequest(source,
		assetClass,
		iSymbols,
		operation,
		dataTypes,
		account,
		iStartTime,
		iEndTime,
		timeFrame,
		iNoConfirm,
	)
	return dataRequest, nil
}
