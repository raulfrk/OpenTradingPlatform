package requests

import "tradingplatform/shared/types"

func DefaultForEmptyDataRequest(dr *DataRequest) {
	if dr.Source == "" {
		dr.Source = types.Alpaca
	}
	if dr.Account == "" {
		dr.Account = DefaultAccount
	}
	if dr.TimeFrame == "" {
		dr.TimeFrame = types.OneMin
	}
}

func DefaultForEmptyStreamRequest(sr *StreamRequest) {
	if sr.Source == "" {
		sr.Source = types.Alpaca
	}
	if sr.Account == "" {
		sr.Account = DefaultAccount
	}
}

func DefaultForEmptySentimentAnalysisRequest(sr *SentimentAnalysisRequest) {
	if sr.SentimentAnalysisProcess == "" {
		sr.SentimentAnalysisProcess = types.Plain
	}

	if sr.DataType == "" {
		sr.DataType = types.RawText
	}
	if sr.AssetClass == "" {
		sr.AssetClass = types.News
	}

	if sr.Operation == "" {
		sr.Operation = types.DataGetOp
	}
}

func DefaultForEmptyStreamAddDeleteRequest(sr *StreamRequest) {
	if sr.Source == "" {
		sr.Source = types.Alpaca
	}
	if sr.Account == "" {
		sr.Account = DefaultAccount
	}

	if len(sr.DataTypes) == 0 {
		providerDatatypes := GetDataTypeMap()[types.Source(sr.Source)]
		for _, dtype := range providerDatatypes(types.AssetClass(sr.AssetClass)) {
			sr.DataTypes = append(sr.DataTypes, types.DataType(dtype))
		}
	}
}
