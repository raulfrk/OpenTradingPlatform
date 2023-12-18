package producer

import (
	"strings"
	"sync"
	sharedent "tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/utils"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

var streams = make(map[string]*utils.Handler[sharedent.Message, sharedent.Message])
var streamsMutex sync.RWMutex

func GetStreamHandler(topic string) *utils.Handler[sharedent.Message, sharedent.Message] {
	splitTopic := strings.Split(topic, ".")
	formatted_topic := strings.Join(splitTopic[:len(splitTopic)-1], ".")

	streamsMutex.RLock()
	handler, ok := streams[formatted_topic]
	streamsMutex.RUnlock()

	if !ok {
		newTopicHandler := utils.NewHandler[sharedent.Message, sharedent.Message]()
		streamsMutex.Lock()
		streams[formatted_topic] = newTopicHandler
		streamsMutex.Unlock()
		logging.Log().Debug().Str("topic", formatted_topic).Msg("starting new stream handler")
		StartTopicHandler(newTopicHandler)
		handler = newTopicHandler
	}

	return handler
}

func StartTopicHandler(handler *utils.Handler[sharedent.Message, sharedent.Message]) {
	ich := make(chan *sharedent.Message)
	handler.SetInputChannel(ich)
	go handleTopic(handler)
}

func handleTopic(handler *utils.Handler[sharedent.Message, sharedent.Message]) {
	nc, _ := nats.Connect(nats.DefaultURL)

	defer nc.Close()
	for msg := range handler.Ich {
		messagePayload, _ := proto.Marshal(msg)
		err := nc.Publish(msg.Topic, messagePayload)
		if err != nil {
			logging.Log().Error().
				Err(err).
				RawJSON("message", messagePayload).
				Msg("publishing message")
		}
	}
	<-handler.Ctx().Done()
}
