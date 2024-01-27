package data

import (
	"sync"

	"tradingplatform/dataprovider/requests"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var streamDB *gorm.DB
var streamDBLock = sync.RWMutex{}

type ActiveStream struct {
	DataSource requests.DataSource `gorm:"uniqueIndex:idx_data_source_account_data_type_symbol"`
	Account    requests.Account    `gorm:"uniqueIndex:idx_data_source_account_data_type_symbol"`
	DataType   types.DataType      `gorm:"uniqueIndex:idx_data_source_account_data_type_symbol"`
	AssetClass types.AssetClass    `gorm:"uniqueIndex:idx_data_source_account_data_type_symbol"`
	Symbol     string              `gorm:"uniqueIndex:idx_data_source_account_data_type_symbol"`
}

func GetActiveStreams() []ActiveStream {
	var activeStreams []ActiveStream
	streamDBLock.Lock()
	result := streamDB.Find(&activeStreams, ActiveStream{})
	streamDBLock.Unlock()
	if result.Error != nil {
		logging.Log().Error().
			Err(result.Error).
			Msg("getting active streams from local database")
	}
	return activeStreams
}

func GetActiveStreamsAssetClass(assetClass types.AssetClass) []ActiveStream {
	var activeStreams []ActiveStream
	streamDBLock.Lock()
	result := streamDB.Where("asset_class = ?", assetClass).Find(&activeStreams, ActiveStream{})
	streamDBLock.Unlock()
	if result.Error != nil {
		logging.Log().Error().
			Err(result.Error).
			Msg("getting active streams for given asset class from local database")
	}
	return activeStreams
}

func AddActiveStreamForDType(req requests.StreamRequest, dataType types.DataType) {
	streamDBLock.Lock()
	tx := streamDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			logging.Log().Error().
				RawJSON("request", req.JSON()).
				Msg("panicked while adding active stream to local database")
			tx.Rollback()
		}
	}()

	for _, symbol := range req.GetSymbols() {
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&ActiveStream{
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
			streamDBLock.Unlock()
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("committing transaction used to add active stream to local database")
	}
	streamDBLock.Unlock()
}

func RemoveActiveStreamForDType(req requests.StreamRequest, dataType types.DataType) {
	streamDBLock.Lock()
	tx := streamDB.Begin()
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
		Delete(&ActiveStream{}).
		Error; err != nil {
		tx.Rollback()
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("removing active stream from local database")
		streamDBLock.Unlock()
		return
	}

	if err := tx.Commit().Error; err != nil {
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("committing transaction used to remove active stream to local database")
		return
	}
	streamDBLock.Unlock()
}

func InitializeDatabase() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	if err != nil {
		logging.Log().Error().Err(err).Msg("failed to connect to local database")
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logging.Log().Error().
			Err(err).
			Msg("failed to get sql database from gorm database")
		panic(err)
	}

	cleanup := func() {
		sqlDB.Close()
	}

	// Migrate the schema
	db.AutoMigrate(&ActiveStream{})

	streamDB = db

	return db, cleanup
}
