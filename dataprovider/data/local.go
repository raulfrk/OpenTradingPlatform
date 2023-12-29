package data

import "tradingplatform/shared/data"

// Initialize local database for the DataProvider
// Returns cancel function
func InitializeDataProviderLocalDatabase() func() {
	_, cancel := data.InitializeLocalDatabase(&DataProviderStream{})
	return cancel
}
