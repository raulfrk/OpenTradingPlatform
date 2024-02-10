package data

import (
	"time"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"

	"gorm.io/gorm/clause"
)

type OrderbookEntry struct {
	Price  float64 `gorm:"uniqueIndex:idx_orderbook_entry"`
	Size   float64 `gorm:"uniqueIndex:idx_orderbook_entry"`
	Source string  `gorm:"uniqueIndex:idx_orderbook_entry"`
}

type AsksOrderbookEntry struct {
	OrderbookEntry
	OrderbookFingerprint string `gorm:"primaryKey;uniqueIndex:idx_orderbook_entry"`
}

type BidsOrderbookEntry struct {
	OrderbookEntry
	OrderbookFingerprint string `gorm:"primaryKey;uniqueIndex:idx_orderbook_entry"`
}

type Orderbook struct {
	Symbol      string
	Exchange    string
	Timestamp   time.Time            `gorm:"index"`
	Asks        []AsksOrderbookEntry `gorm:"foreignKey:OrderbookFingerprint"`
	Bids        []BidsOrderbookEntry `gorm:"foreignKey:OrderbookFingerprint"`
	Reset_      bool
	Fingerprint string `gorm:"primaryKey"`
	Source      string
	AssetClass  string `gorm:"not null"`
}

func OrderbookFromEntity(entity *entities.Orderbook) Orderbook {
	asks := make([]AsksOrderbookEntry, len(entity.Asks))
	for i, a := range entity.Asks {
		asks[i] = AsksOrderbookEntry{
			OrderbookEntry: OrderbookEntry{
				Price:  a.Price,
				Size:   a.Size,
				Source: a.Source,
			},
		}

	}
	bids := make([]BidsOrderbookEntry, len(entity.Bids))
	for i, b := range entity.Bids {
		bids[i] = BidsOrderbookEntry{
			OrderbookEntry: OrderbookEntry{
				Price:  b.Price,
				Size:   b.Size,
				Source: b.Source,
			},
		}
	}
	return Orderbook{
		Symbol:      entity.Symbol,
		Exchange:    entity.Exchange,
		Timestamp:   time.Unix(entity.Timestamp, 0),
		Asks:        asks,
		Bids:        bids,
		Reset_:      entity.Reset_,
		Fingerprint: entity.Fingerprint,
		Source:      entity.Source,
		AssetClass:  entity.AssetClass,
	}
}

func OrderbookToEntity(orderbook Orderbook) *entities.Orderbook {
	asks := make([]*entities.OrderbookEntry, len(orderbook.Asks))
	for i, a := range orderbook.Asks {
		asks[i] = &entities.OrderbookEntry{
			Price:  a.Price,
			Size:   a.Size,
			Source: a.Source,
		}
	}
	bids := make([]*entities.OrderbookEntry, len(orderbook.Bids))
	for i, b := range orderbook.Bids {
		bids[i] = &entities.OrderbookEntry{
			Price:  b.Price,
			Size:   b.Size,
			Source: b.Source,
		}
	}
	return &entities.Orderbook{
		Symbol:      orderbook.Symbol,
		Exchange:    orderbook.Exchange,
		Timestamp:   orderbook.Timestamp.Unix(),
		Asks:        asks,
		Bids:        bids,
		Reset_:      orderbook.Reset_,
		Fingerprint: orderbook.Fingerprint,
		Source:      orderbook.Source,
		AssetClass:  orderbook.AssetClass,
	}
}

func OrderbooksFromEntities(entities []*entities.Orderbook) []Orderbook {
	orderbooks := make([]Orderbook, len(entities))
	for i, entity := range entities {
		orderbooks[i] = OrderbookFromEntity(entity)
	}
	return orderbooks
}

func OrderbooksToEntities(orderbooks []Orderbook) []*entities.Orderbook {
	entities := make([]*entities.Orderbook, len(orderbooks))
	for i, orderbook := range orderbooks {
		entities[i] = OrderbookToEntity(orderbook)
	}
	return entities
}

func InsertOrderbook(orderbook *entities.Orderbook) {
	dbOrderbook := OrderbookFromEntity(orderbook)
	tx := DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(dbOrderbook)
	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			RawJSON("orderbook", entities.GenerateJson(orderbook)).
			Msg("inserting orderbook")
	}
}

func GetOrderbookFromRequest(symbol string, req requests.DataRequest) ([]*entities.Orderbook, error) {
	return GetOrderbook(string(req.GetSource()),
		symbol,
		string(req.AssetClass),
		req.GetStartTime(),
		req.GetEndTime()), nil
}

func GetOrderbook(source, symbol, assetClass string, startTime, endTime int64) []*entities.Orderbook {
	var orderbooks []Orderbook

	tx := DB.Preload("Asks").Preload("Bids").Where("source = ? AND symbol = ? AND asset_class = ? AND timestamp >= ? AND timestamp < ?",
		source,
		symbol,
		assetClass,
		time.Unix(startTime, 0),
		time.Unix(endTime, 0))
	orderbooks = PaginateRequest(tx, Orderbook{})
	return OrderbooksToEntities(orderbooks)
}
