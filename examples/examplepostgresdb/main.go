package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"tradingplatform/shared/communication"
	"tradingplatform/shared/entities"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

const (
	// Replace with your Supabase database URL
	// It should look like: "postgres://username:password@hostname:port/database"
	dsn = "postgres://postgres:example@localhost:5432/test"
)

type BarDB struct {
	Fingerprint string `gorm:"primaryKey"`
	Symbol      string `gorm:"not null"`
	Exchange    string
	Open        float64
	High        float64
	Low         float64
	Close       float64
	Volume      float64
	VWAP        float64
	Timestamp   time.Time `gorm:"index"`
	TradeCount  uint64
	Source      string
}

const writeTimeout = time.Second * 5
const bufferSize = 1000

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

var buffers = make(map[string][]BarDB)
var lock = sync.RWMutex{}

func main() {
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
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	_, err = db.DB()
	db.AutoMigrate(&BarDB{})

	// Connect to nats and subscribe to a channel
	nc, _ := nats.Connect(communication.GetNatsURL())
	// Use a buffered channel with an appropriate size
	nc.Subscribe("dataprovider.*.alpaca.*.bar.>", func(m *nats.Msg) {
		var newBar entities.Bar
		var newMsg entities.Message
		err := proto.Unmarshal(m.Data, &newMsg)
		if err != nil {
			fmt.Println("Failed to unmarshal message:", err)
			return
		}

		proto.Unmarshal(newMsg.Payload, &newBar)

		// Create a new BarDB object
		newBarDB := BarDB{
			Fingerprint: newBar.Fingerprint,
			Symbol:      newBar.Symbol,
			Exchange:    newBar.Exchange,
			Open:        newBar.Open,
			High:        newBar.High,
			Low:         newBar.Low,
			Close:       newBar.Close,
			Volume:      newBar.Volume,
			VWAP:        newBar.VWAP,
			Timestamp:   time.Unix(newBar.Timestamp, 0),
			TradeCount:  newBar.TradeCount,
			Source:      newBar.Source,
		}
		topicParts := strings.Split(newMsg.Topic, ".")
		if isNumber(topicParts[len(topicParts)-1]) {
			lock.Lock()
			if buffers[newMsg.Topic] == nil {
				buffers[newMsg.Topic] = make([]BarDB, 0, bufferSize)
			}
			count, _ := strconv.ParseInt(topicParts[len(topicParts)-1], 10, 64)
			buffers[newMsg.Topic] = append(buffers[newMsg.Topic], newBarDB)

			lock.Unlock()
			if int(count) == len(buffers[newMsg.Topic]) {
				go processBuffer(newMsg.Topic, db)
			}
		} else {
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&newBarDB)
			fmt.Printf("Wrote to database: %v\n", newBarDB)
		}
	})

	// Writer

	select {}

}

func processBuffer(topic string, db *gorm.DB) {
	// db.Exec("ALTER TABLE bar_dbs DISABLE TRIGGER ALL")
	now := time.Now()
	lock.Lock()
	buffer := buffers[topic]
	fmt.Printf("Writing to database %d bars\n", len(buffer))

	lock.Unlock()
	db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(buffer, 5000)
	// db.Exec("ALTER TABLE bar_dbs ENABLE TRIGGER ALL")
	fmt.Printf("Finished writing to database %d bars\n", len(buffer))
	lock.Lock()
	delete(buffers, topic)
	lock.Unlock()
	var count int64
	db.Table("bar_dbs").Count(&count)
	fmt.Printf("Total rows in database: %d\n", count)
	fmt.Printf("Time taken: %s\n", time.Since(now))
}

func worker(barChan <-chan BarDB, db *gorm.DB) {
	buffer := make([]BarDB, 0, bufferSize)
	ticker := time.NewTicker(writeTimeout)
	defer ticker.Stop()

	for {
		select {
		case bar := <-barChan:
			buffer = append(buffer, bar)

			if len(buffer) == bufferSize {
				db.Exec("ALTER TABLE bar_dbs DISABLE TRIGGER ALL")
				// fmt.Printf("Writing to database %d bars\n", bufferSize)
				db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(buffer, len(buffer)).Commit()
				buffer = buffer[:0] // Reset the buffer
				db.Exec("ALTER TABLE bar_dbs ENABLE TRIGGER ALL")

			}

		case <-ticker.C:
			if len(buffer) > 0 {
				// fmt.Printf("Writing to database %d bars\n", bufferSize)
				db.Exec("ALTER TABLE bar_dbs DISABLE TRIGGER ALL")

				db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(buffer, len(buffer)).Commit()
				buffer = buffer[:0] // Reset the buffer
				db.Exec("ALTER TABLE bar_dbs ENABLE TRIGGER ALL")

			}
		}
	}
}
