package requests

import (
	"fmt"
	"tradingplatform/shared/types"

	"github.com/go-playground/validator/v10"
)

func IsValidDataSource(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, exists := GetDataSourceMap()[value]
	return exists
}

func IsValidAssetClass(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, exists := types.GetAssetClassMap()[value]
	return exists
}

func IsValidOperation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, exists := types.GetDataRequestOpMap()[value]
	return exists
}

func IsValidOperationStream(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, exists := types.GetStreamRequestOpMap()[value]
	return exists
}

func IsValidOperationStreamSubscribe(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, exists := types.GetStreamSubscribeRequestOpMap()[value]
	return exists
}

func IsValidDataType(fl validator.FieldLevel) bool {
	source := fl.Parent().FieldByName("Source").String()
	assetClass := fl.Parent().FieldByName("AssetClass").String()
	value := fl.Field().String()
	dataTypeMap := GetDataTypeMap()[types.Source(source)]
	_, exists := dataTypeMap(types.AssetClass(assetClass))[types.DataType(value)]
	return exists
}

func IsValidMultiDataType(fl validator.FieldLevel) bool {
	source := fl.Parent().FieldByName("Source").String()
	assetClass := fl.Parent().FieldByName("AssetClass").String()
	for i := 0; i < fl.Field().Len(); i++ {
		// Get the value of the current element
		value := fl.Field().Index(i).Interface().(types.DataType)
		dataTypeMap := GetDataTypeMap()[types.Source(source)]
		_, exists := dataTypeMap(types.AssetClass(assetClass))[types.DataType(value)]

		if !exists {
			return false
		}
	}
	return true
}
func IsValidMultiDataTypeStream(fl validator.FieldLevel) bool {
	operation := fl.Parent().FieldByName("Operation").String()
	if operation == "get" {
		return true
	}
	return fl.Field().Len() > 0 && IsValidMultiDataType(fl)
}

func IsValidMultiSymbolStream(fl validator.FieldLevel) bool {
	operation := fl.Parent().FieldByName("Operation").String()
	if operation == "get" {
		return true
	}
	return fl.Field().Len() > 0
}

func IsValidAccount(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, exists := GetAccount()[value]
	return exists
}

func IsValidDataFrame(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, exists := types.GetTimeFrameMap()[value]
	return exists
}

func IsValidEndTime(fl validator.FieldLevel) bool {
	startTime := fl.Parent().FieldByName("StartTime").Int()
	endTime := fl.Field().Int()
	return endTime > startTime
}

func IsValidSentimentAnalysisProcess(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, exists := types.GetSentimentAnalysisProcessMap()[value]
	return exists
}

func IsValidModelProvider(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, exists := GetModelProviderMap()[value]
	return exists
}

func SummarizeError(err error) error {
	var errMsg string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsg += fmt.Sprintf("Field: %s, Value: %s, Error: %s\n", err.Field(), err.Value(), err.Tag())
		}
		return fmt.Errorf(errMsg)
	}
	return nil
}
