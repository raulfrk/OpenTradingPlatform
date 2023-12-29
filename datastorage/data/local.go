package data

import (
	"tradingplatform/shared/data"
)

// InitializeDataStorageLocalDatabase initializes the local database and returns a function to close the database.
func InitializeDataStorageLocalDatabase() func() {
	_, cancel := data.InitializeLocalDatabase(&data.SubscribedTopic{})
	return cancel
}
