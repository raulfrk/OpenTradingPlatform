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

var streams = make(map[string]*utils.Handler[sharedent.Message])
var streamsMutex sync.RWMutex

// GetStreamHandler returns the stream handler for a topic
func GetStreamHandler(topic string) *utils.Handler[sharedent.Message] {
	splitTopic := strings.Split(topic, ".")
	formattedTopic := strings.Join(splitTopic[:len(splitTopic)-1], ".")

	streamsMutex.RLock()
	handler, ok := streams[formattedTopic]
	streamsMutex.RUnlock()

	if !ok {
		newTopicHandler := utils.NewHandler[sharedent.Message]()
		streamsMutex.Lock()
		streams[formattedTopic] = newTopicHandler
		streamsMutex.Unlock()
		logging.Log().Debug().Str("topic", formattedTopic).Msg("starting new stream handler")
		StartTopicHandler(newTopicHandler)
		handler = newTopicHandler
	}

	return handler
}

// StartTopicHandler starts a new stream handler for a topic
func StartTopicHandler(handler *utils.Handler[sharedent.Message]) {
	ich := make(chan *sharedent.Message)
	handler.SetChannel(ich)
	go handleTopic(handler)
}

// StopTopicHandler stops a stream handler for a topic
func StopTopicHandler(topic string) {
	streamsMutex.RLock()
	handler, ok := streams[topic]
	streamsMutex.RUnlock()
	if !ok {
		return
	}
	handler.Cancel()
	close(handler.Ch)
	streamsMutex.Lock()
	delete(streams, topic)
	streamsMutex.Unlock()
}

func handleTopic(handler *utils.Handler[sharedent.Message]) {
	nc, _ := nats.Connect(nats.DefaultURL)

	defer nc.Close()
	for msg := range handler.Ch {
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
