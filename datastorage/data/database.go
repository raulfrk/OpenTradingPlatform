package data

import (
	"fmt"
	"log"
	"os"
	"strings"
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

const DatabaseConnectionAttempts = 3

func hasDatabaseName(dsn string) bool {
	splitDsn := strings.Split(dsn, " ")
	for _, element := range splitDsn {
		if strings.Contains(element, "dbname") {
			return true
		}
	}
	return false
}

func getDatabaseName(dsn string) (string, bool) {
	if !hasDatabaseName(dsn) {
		return "", false
	}

	splitDsn := strings.Split(dsn, " ")
	for _, element := range splitDsn {
		if strings.Contains(element, "dbname") {
			return strings.Split(element, "=")[1], true
		}
	}
	return "", false
}

func getDsnRoot(dsn string) string {
	splitDsn := strings.Split(dsn, " ")
	var dsnRoot []string
	for _, element := range splitDsn {
		if !strings.Contains(element, "dbname") {
			dsnRoot = append(dsnRoot, element)
		}
	}
	return strings.Join(dsnRoot, " ")
}

// NOTE: Only the following dsn format is supported:
// "host=somehost user=someuser password=somepassword dbname=somedb port=someport"
func createDatabaseIfNotExists(dsn string, attempts int) {
	var db *gorm.DB
	var err error
	currentAttempt := attempts
	dbName, exists := getDatabaseName(dsn)
	logging.Log().Debug().Str("dbName", dbName).Msg("checking if database exists")
	if !exists {
		logging.Log().Debug().Msg("no database name found in dsn")
		logging.Log().Warn().Msg(`Only the following dsn format is supported:
		host=somehost user=someuser password=somepassword dbname=somedb port=someport`)
		return
	}
	dsnRoot := getDsnRoot(dsn)
	for currentAttempt != 0 && db == nil {
		db, err = gorm.Open(postgres.Open(dsnRoot), &gorm.Config{})
		if err != nil {
			if currentAttempt == 1 {
				logging.Log().Error().Str("dbName", dbName).Err(err).Msg("failed to connect to the database in the database creation step")
				panic("failed to connect to the database in the database creation step: " + err.Error())
			}
			db = nil
		}
		currentAttempt--
		time.Sleep(time.Second * 1)
	}
	var count int64
	db.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", dbName).Scan(&count)
	if count > 0 {
		logging.Log().Debug().Str("dbName", dbName).Msg("database already exists")
		return
	}
	result := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName))
	if result.Error != nil {
		logging.Log().Error().Str("dbName", dbName).Err(err).Msg("failed to create the database")
		panic("failed to create the database: " + result.Error.Error())
	}

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

	createDatabaseIfNotExists(dsn, DatabaseConnectionAttempts)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logging.Log().Error().Err(err).Msg("failed to connect to remote database")
		panic(err)
	}
	logging.Log().Debug().Str("dsn", dsn).Msg("connected to database successfully")

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

func PaginateRequest[T any](tx *gorm.DB, _ T) []T {
	pageSize := 65000
	page := 1
	ents := make([]T, 0)
	for {
		var batch []T
		tx.Offset((page - 1) * pageSize).Limit(pageSize).Find(&batch)
		ents = append(ents, batch...)
		if len(batch) < pageSize {
			break
		}
		page++
	}

	return ents

}
