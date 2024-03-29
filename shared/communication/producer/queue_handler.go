package producer

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
	"tradingplatform/shared/communication"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"

	sharedent "tradingplatform/shared/entities"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

var queues = make(map[string]*utils.Handler[[]*sharedent.Message])
var queuesMutex sync.RWMutex

// GetQueueHandler returns a queue handler for a topic
func GetQueueHandler(topic string, noConfirm bool) (*utils.Handler[[]*sharedent.Message], types.DataResponse) {

	queuesMutex.RLock()
	_, ok := queues[topic]
	queuesMutex.RUnlock()

	if !ok {
		newQueueHandler := utils.NewHandler[[]*sharedent.Message]()
		queuesMutex.Lock()
		queues[topic] = newQueueHandler
		queuesMutex.Unlock()
		StartQueueHandler(newQueueHandler, noConfirm)
		handler := newQueueHandler
		return handler, types.NewDataResponse(
			types.Success,
			"Created new queue handler",
			nil,
			"",
		)
	}

	return nil, types.NewDataError(
		fmt.Errorf("a queue with a response for topic %s already exists", topic),
	)
}

// StartQueueHandler starts a queue handler
func StartQueueHandler(handler *utils.Handler[[]*sharedent.Message], noConfirm bool) {
	ich := make(chan *[]*sharedent.Message)
	handler.SetChannel(ich)
	go handleQueue(handler, noConfirm)
}

func handleQueue(handler *utils.Handler[[]*sharedent.Message], noConfirm bool) {
	nc, err := nats.Connect(communication.GetNatsURL())
	var topic string
	if err != nil {
		logging.Log().Fatal().Err(err).Msg("connecting to NATS")
		return
	}
	defer nc.Close()
	msgs := <-handler.Ch
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)

	first := (*msgs)[0]
	topic = first.Topic
	logging.Log().Debug().Str("topic", topic).Msg("starting queue handler")
	if !noConfirm {
		sub, _ := nc.Subscribe(first.Topic, func(msg *nats.Msg) {
			cancel()
		})
		<-ctx.Done()
		sub.Unsubscribe()
	}
	for _, msg := range *msgs {
		messagePayload, _ := proto.Marshal(msg)
		err := nc.Publish(msg.Topic, messagePayload)
		if err != nil {
			logging.Log().Error().
				Err(err).
				Msg("publishing message to data queue")
		}

	}
	handler.Cancel()
	cancel()
	<-handler.Ctx().Done()
	queuesMutex.Lock()
	delete(queues, topic)
	queuesMutex.Unlock()
	logging.Log().Debug().Str("topic", topic).Msg("queue handler stopped")
}

func GenerateQueueID() string {
	return uuid.New().String()

}
