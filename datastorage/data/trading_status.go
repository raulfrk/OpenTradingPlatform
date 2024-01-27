package data

import (
	"time"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
)

type TradingStatus struct {
	Symbol      string
	StatusCode  string
	StatusMsg   string
	ReasonCode  string
	ReasonMsg   string
	Timestamp   time.Time `gorm:"index"`
	Tape        string
	Fingerprint string `gorm:"primaryKey"`
	Source      string
	AssetClass  string `gorm:"not null"`
}

func TradingStatusFromEntity(entity *entities.TradingStatus) TradingStatus {
	return TradingStatus{
		Symbol:      entity.Symbol,
		StatusCode:  entity.StatusCode,
		StatusMsg:   entity.StatusMsg,
		ReasonCode:  entity.ReasonCode,
		ReasonMsg:   entity.ReasonMsg,
		Timestamp:   time.Unix(entity.Timestamp, 0),
		Tape:        entity.Tape,
		Fingerprint: entity.Fingerprint,
		Source:      entity.Source,
		AssetClass:  entity.AssetClass,
	}
}

func TradingStatusToEntity(tradingStatus TradingStatus) *entities.TradingStatus {
	return &entities.TradingStatus{
		Symbol:      tradingStatus.Symbol,
		StatusCode:  tradingStatus.StatusCode,
		StatusMsg:   tradingStatus.StatusMsg,
		ReasonCode:  tradingStatus.ReasonCode,
		ReasonMsg:   tradingStatus.ReasonMsg,
		Timestamp:   tradingStatus.Timestamp.Unix(),
		Tape:        tradingStatus.Tape,
		Fingerprint: tradingStatus.Fingerprint,
		Source:      tradingStatus.Source,
		AssetClass:  tradingStatus.AssetClass,
	}
}

func TradingStatusesToEntities(tradingStatuses []TradingStatus) []*entities.TradingStatus {
	entities := make([]*entities.TradingStatus, len(tradingStatuses))
	for i, tradingStatus := range tradingStatuses {
		entities[i] = TradingStatusToEntity(tradingStatus)
	}
	return entities
}

func TradingStatusesFromEntities(entities []*entities.TradingStatus) []TradingStatus {
	tradingStatuses := make([]TradingStatus, len(entities))
	for i, entity := range entities {
		tradingStatuses[i] = TradingStatusFromEntity(entity)
	}
	return tradingStatuses
}

func InsertTradingStatus(tradingStatus *entities.TradingStatus) {
	dbTradingStatus := TradingStatusFromEntity(tradingStatus)
	tx := DB.Create(dbTradingStatus)
	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			RawJSON("tradingStatus", entities.GenerateJson(tradingStatus)).
			Msg("inserting tradingStatus")
	}
}

func GetTradingStatusesFromRequest(symbol string, req requests.DataRequest) ([]*entities.TradingStatus, error) {
	return GetTradingStatuses(string(req.GetSource()),
		symbol,
		string(req.AssetClass),
		req.GetStartTime(),
		req.GetEndTime()), nil
}

func GetTradingStatuses(source, symbol, assetClass string, startTime, endTime int64) []*entities.TradingStatus {
	var ts []TradingStatus
	DB.Where("source = ? AND symbol = ? AND asset_class = ? AND timestamp >= ? AND timestamp <= ?",
		source,
		symbol,
		assetClass,
		time.Unix(startTime, 0),
		time.Unix(endTime, 0)).Find(&ts)
	return TradingStatusesToEntities(ts)
}
