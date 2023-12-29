package utils

import (
	"fmt"
	"tradingplatform/shared/types"
)

type Topic struct {
	Component     types.Component
	Functionality types.Functionality
	Source        types.Source
	AssetClass    types.AssetClass
	DataType      types.DataType
	TimeFrame     types.TimeFrame
	QueueID       string
	QueueCount    int
	Symbol        string
}

func (t Topic) Generate() string {
	base := fmt.Sprintf("%s.%s", t.Component, t.Functionality)

	if t.Functionality == types.Command {
		return base
	}
	if t.Source != "" {
		base = fmt.Sprintf("%s.%s", base, t.Source)
	}

	if t.Functionality == types.Stream {
		stream := base
		stream = fmt.Sprintf("%s.%s.%s.%s", stream, t.AssetClass, t.DataType, t.Symbol)
		return stream
	}

	if t.Functionality == types.Data {
		data := fmt.Sprintf("%s.%s.%s", base, t.AssetClass, t.DataType)
		if t.DataType == types.Bar {
			data = fmt.Sprintf("%s.%s", data, t.TimeFrame)
		}
		data = fmt.Sprintf("%s.%s.%s.%d", data, t.Symbol, t.QueueID, t.QueueCount)
		return data
	}
	return base
}

func NewDataTopic(component types.Component, source types.Source, assetClass types.AssetClass, dataType types.DataType, symbol string, queueID string, queueCount int) Topic {
	return Topic{
		Component:     component,
		Functionality: types.Data,
		Source:        source,
		AssetClass:    assetClass,
		DataType:      dataType,
		QueueID:       queueID,
		QueueCount:    queueCount,
		Symbol:        symbol,
	}
}

func NewBarDataTopic(component types.Component, source types.Source, assetClass types.AssetClass, timeFrame types.TimeFrame, symbol string, queueID string, queueCount int) Topic {
	return Topic{
		Component:     component,
		Functionality: types.Data,
		Source:        source,
		AssetClass:    assetClass,
		DataType:      types.Bar,
		TimeFrame:     timeFrame,
		QueueID:       queueID,
		QueueCount:    queueCount,
		Symbol:        symbol,
	}
}

func NewStreamTopic(component types.Component, source types.Source, assetClass types.AssetClass, dataType types.DataType, symbol string) Topic {
	return Topic{
		Component:     component,
		Functionality: types.Stream,
		Source:        source,
		AssetClass:    assetClass,
		DataType:      dataType,
		Symbol:        symbol,
	}
}

func NewCommandTopic(component types.Component) Topic {
	return Topic{
		Component:     component,
		Functionality: types.Command,
	}
}

func NewLoggingTopic(component types.Component) Topic {
	return Topic{
		Component:     component,
		Functionality: types.Logging,
	}
}
