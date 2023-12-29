package data

import (
	"log"
	"os"
	"time"
	"tradingplatform/shared/logging"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var dsn string

func SetDSN(d string) {
	dsn = d
}

// Initialize the database connection and returns the database object and cancel function
func InitializeDatabase() (*gorm.DB, func()) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // Slow SQL threshold
			LogLevel:      logger.Error, // Log level
			Colorful:      false,        // Disable color
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		logging.Log().Error().Err(err).Msg("failed to connect to remote database")
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

	db.AutoMigrate(
		&Log{},
		&Bar{},
		&Trade{},
		&TradeCondition{},
		&Quote{},
		&QuoteCondition{},
		&LULD{},
		&Orderbook{},
		&AsksOrderbookEntry{},
		&BidsOrderbookEntry{},
		&TradingStatus{},
		&News{},
		&NewsSymbol{},
		&DailyBar{},
		&LLM{},
		&Sentiment{},
	)
	DB = db

	return db, cleanup
}

// Generic function to insert entities into the database in batches
func InsertBatchEntity[I any](entities []I) {
	logging.Log().Debug().
		Int("count", len(entities)).
		Type("entity", entities[0]).
		Msg("started inserting batch of entities to db")
	tx := DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(entities, 3000)

	if tx.Error != nil {
		logging.Log().Error().Err(tx.Error).Msg("failed to insert batch of entities to db")
		return
	}
	logging.Log().Debug().
		Int("count", len(entities)).
		Type("entity", entities[0]).
		Msg("finished inserting batch of entities to db")
}
