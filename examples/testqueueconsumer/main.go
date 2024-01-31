package main

import (
	"sort"
	"sync"
	"tradingplatform/shared/communication"
	"tradingplatform/shared/entities"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Subscribe to a queue
	// Receive messages from the queue
	// Unmarshal the message
	// Print the message
	nc, err := nats.Connect(communication.GetNatsURL())
	if err != nil {

		log.Fatal().Err(err).Msg("connecting to JetStream")
	}
	defer nc.Close()

	var bars []*entities.Bar
	var barslock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(1_483_546)
	nc.QueueSubscribe("dataprovider.data.alpaca.crypto.bar.47b25753d24a3e4eb7888178a1b892dd9e73f4e424ec8bd418dc4b8c6845ef6a.1483545", "test-queue", func(m *nats.Msg) {
		var msg entities.Message
		proto.Unmarshal(m.Data, &msg)

		var bar entities.Bar
		proto.Unmarshal(msg.Payload, &bar)
		barslock.Lock()
		bars = append(bars, &bar)
		barslock.Unlock()
		wg.Done()
	})
	nc.Publish("dataprovider.data.alpaca.crypto.bar.47b25753d24a3e4eb7888178a1b892dd9e73f4e424ec8bd418dc4b8c6845ef6a.1483545", []byte(""))
	wg.Wait()
	sort.Slice(bars, func(i, j int) bool {
		return bars[i].Timestamp < bars[j].Timestamp
	})

	if bars[len(bars)-1].Fingerprint == "47b25753d24a3e4eb7888178a1b892dd9e73f4e424ec8bd418dc4b8c6845ef6a" {
		log.Info().Msg("Success")
	} else {
		log.Error().Msg("Failure")
	}

}
