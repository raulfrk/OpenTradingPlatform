package data

import (
	"tradingplatform/shared/data"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"

	"gorm.io/gorm/clause"
)

type DataProviderStream struct {
	DataSource types.Source     `gorm:"uniqueIndex:idx_data_source_account_data_type_symbol"`
	Account    requests.Account `gorm:"uniqueIndex:idx_data_source_account_data_type_symbol"`
	DataType   types.DataType   `gorm:"uniqueIndex:idx_data_source_account_data_type_symbol"`
	AssetClass types.AssetClass `gorm:"uniqueIndex:idx_data_source_account_data_type_symbol"`
	Symbol     string           `gorm:"uniqueIndex:idx_data_source_account_data_type_symbol"`
}

// Get all active streams of the dataprovider from local database
func GetDataProviderStreams() []DataProviderStream {
	var activeStreams []DataProviderStream
	data.LocalDBLock.Lock()
	result := data.LocalDB.Find(&activeStreams, DataProviderStream{})
	data.LocalDBLock.Unlock()
	if result.Error != nil {
		logging.Log().Error().
			Err(result.Error).
			Msg("getting active streams from local database")
	}
	return activeStreams
}

// Get all active streams of the dataprovider from local database for a given asset class
func GetDataProviderStreamsAssetClass(assetClass types.AssetClass) []DataProviderStream {
	var activeStreams []DataProviderStream
	data.LocalDBLock.Lock()
	result := data.LocalDB.Where("asset_class = ?", assetClass).
		Find(&activeStreams, DataProviderStream{})
	data.LocalDBLock.Unlock()
	if result.Error != nil {
		logging.Log().Error().
			Err(result.Error).
			Msg("getting active streams for given asset class from local database")
	}
	return activeStreams
}

// Add active stream to local database
func AddDataProviderStreamForDType(req requests.StreamRequest, dataType types.DataType) {
	data.LocalDBLock.Lock()
	tx := data.LocalDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			logging.Log().Error().
				RawJSON("request", req.JSON()).
				Msg("panicked while adding active stream to local database")
			tx.Rollback()
		}
	}()

	for _, symbol := range req.GetSymbols() {
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&DataProviderStream{
			DataSource: req.GetSource(),
			Account:    req.GetAccount(),
			AssetClass: req.GetAssetClass(),
			DataType:   dataType,
			Symbol:     symbol,
		})
		if res.Error != nil {
			tx.Rollback()
			logging.Log().Error().
				Err(res.Error).
				RawJSON("request", req.JSON()).
				Str("symbol", symbol).
				Msg("adding active stream to local database")
			data.LocalDBLock.Unlock()
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("committing transaction used to add active stream to local database")
	}
	data.LocalDBLock.Unlock()
}

// Remove active stream from local database
func RemoveDataProviderStreamForDType(req requests.StreamRequest, dataType types.DataType) {
	data.LocalDBLock.Lock()
	tx := data.LocalDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			logging.Log().Error().
				RawJSON("request", req.JSON()).
				Msg("panicked while removing active stream from local database")
			tx.Rollback()
		}
	}()

	if err := tx.Where("data_source = ? AND account = ? AND asset_class = ? AND data_type = ? AND symbol IN ?",
		req.GetSource(),
		req.GetAccount(),
		req.GetAssetClass(),
		dataType,
		req.GetSymbols()).
		Delete(&DataProviderStream{}).
		Error; err != nil {
		tx.Rollback()
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("removing active stream from local database")
		data.LocalDBLock.Unlock()
		return
	}

	if err := tx.Commit().Error; err != nil {
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("committing transaction used to remove active stream to local database")
		return
	}
	data.LocalDBLock.Unlock()
}
