package data

import (
	"sync"
	"tradingplatform/shared/logging"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var LocalDB *gorm.DB
var LocalDBLock sync.RWMutex

func InitializeLocalDatabase(dst ...interface{}) (*gorm.DB, func()) {
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
	db.AutoMigrate(dst...)

	LocalDB = db

	return db, cleanup
}
