package main

import (
	"fmt"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/types"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()
	nc.Subscribe("dataprovider.stream.>", func(msg *nats.Msg) {
		var message entities.Message
		proto.Unmarshal(msg.Data, &message)

		switch message.DataType {
		case string(types.Bar):
			var bar entities.Bar
			proto.Unmarshal(message.Payload, &bar)
			fmt.Printf("Received bar: %v\n", &bar)
		case string(types.Orderbook):
			var orderbook entities.Orderbook
			proto.Unmarshal(message.Payload, &orderbook)
			fmt.Printf("Received orderbook: %v\n", &orderbook)
		case string(types.DailyBars):
			var dailyBars entities.Bar
			proto.Unmarshal(message.Payload, &dailyBars)
			fmt.Printf("Received daily bars: %v\n", &dailyBars)
		case string(types.Quotes):
			var quotes entities.Quote
			proto.Unmarshal(message.Payload, &quotes)
			fmt.Printf("Received quotes: %v\n", &quotes)
		case string(types.Trades):
			var trades entities.Trade
			proto.Unmarshal(message.Payload, &trades)
			fmt.Printf("Received trades: %v\n", &trades)
		case string(types.LULD):
			var luld entities.LULD
			proto.Unmarshal(message.Payload, &luld)
			fmt.Printf("Received LULD: %v\n", &luld)
		case string(types.Status):
			var status entities.TradingStatus
			proto.Unmarshal(message.Payload, &status)
			fmt.Printf("Received status: %v\n", &status)
		case string(types.RawText):
			var originalText entities.News
			proto.Unmarshal(message.Payload, &originalText)
			fmt.Printf("Received news: %v\n", &originalText)
		}

	})

	select {}
}
