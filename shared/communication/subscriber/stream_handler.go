package subscriber

import (
	"strconv"
	"strings"
	"sync"

	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/utils"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

var streams = make(map[string]*utils.Handler[entities.Message])
var streamsMutex sync.RWMutex
var dataQueues = make(map[string][]*entities.Message)
var dataQueuesMutex sync.RWMutex

func GetStreamHandler(topic string) *utils.Handler[entities.Message] {
	streamsMutex.RLock()
	handler, ok := streams[topic]
	streamsMutex.RUnlock()

	if !ok {
		streamsMutex.Lock()
		streams[topic] = utils.NewHandler[entities.Message]()
		handler = streams[topic]
		streamsMutex.Unlock()
		StartTopicHandler(streams[topic], topic)
	}

	return handler
}

// TODO: separate into multiple functions
func GetStreamHandlerRoundRobin(topic string, maxAgents int) *utils.Handler[entities.Message] {
	streamsMutex.RLock()
	handler, ok := streams[topic]
	streamsMutex.RUnlock()

	if !ok {
		streamsMutex.Lock()
		streams[topic] = utils.NewRoundRobinHandler[entities.Message](maxAgents)
		handler = streams[topic]
		streamsMutex.Unlock()
		StartTopicHandler(streams[topic], topic)
	}

	return handler
}

func StartTopicHandler(handler *utils.Handler[entities.Message], topic string) {
	ich := make(chan *entities.Message)
	handler.SetChannel(ich)
	go handleTopic(handler, topic)
}

func StopTopicHandler(topic string) {
	streamsMutex.Lock()
	defer streamsMutex.Unlock()
	handler, ok := streams[topic]
	if !ok {
		return
	}
	handler.Cancel()
	handler.FunctionalityAttatched = false

	delete(streams, topic)
}

func handleTopic(handler *utils.Handler[entities.Message], topic string) {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()
	sub, _ := nc.QueueSubscribe(topic, "storage", func(m *nats.Msg) {
		go func() {
			var msg entities.Message
			err := proto.Unmarshal(m.Data, &msg)
			if err != nil {
				logging.Log().Error().Err(err).Msg("unmarshalling message")
			}
			if IsQueue(m.Subject) {
				accumulateData(m.Subject, &msg)
			}
			ch, err := handler.GetNextAgent()
			if err != nil {
				logging.Log().Error().Err(err).Msg("getting next agent")
				return
			}
			ch <- &msg
		}()
	})
	sub.SetPendingLimits(-1, -1)

	<-handler.Ctx().Done()
	sub.Unsubscribe()
	nc.Close()
}

// TODO: separate into multiple functions
func AttatchFunctionalityRoundRobin(topic string, f func(*entities.Message), numAgents int) {
	handler := GetStreamHandlerRoundRobin(topic, numAgents)
	handler.Lock.RLock()
	if handler.FunctionalityAttatched {
		logging.Log().Debug().Str("topic", topic).Int("numAgents", numAgents).Msg("functionality already attatched to all agents")
		handler.Lock.RUnlock()
		return
	}
	handler.Lock.RUnlock()
	handler.Lock.Lock()
	for i := 0; i < len(handler.RoundRobinChannels); i++ {
		go func(h *utils.Handler[entities.Message], chId int) {
			h.Lock.RLock()
			ch := h.RoundRobinChannels[chId]
			h.Lock.RUnlock()
			for {
				select {
				case msg := <-ch:
					f(msg)
				case <-h.Ctx().Done():
					logging.Log().Debug().Str("topic", topic).Int("agent", chId).Msg("functionality detatched from agent")
					return
				}
			}
		}(handler, i)
	}
	handler.FunctionalityAttatched = true
	logging.Log().Debug().Str("topic", topic).Int("numAgents", numAgents).Msg("attatched functionality to all agents")
	handler.Lock.Unlock()
}

func accumulateData(topic string, msg *entities.Message) {
	dataQueuesMutex.Lock()
	if len(dataQueues[topic]) == 0 {
		dataQueues[topic] = make([]*entities.Message, 0)
	}
	dataQueues[topic] = append(dataQueues[topic], msg)
	dataQueuesMutex.Unlock()
}

func DrainQueue(topic string) []*entities.Message {
	_, queueCount := GetQueueComponents(topic)
	dataQueuesMutex.Lock()
	if len(dataQueues[topic]) < queueCount {
		dataQueuesMutex.Unlock()
		return nil
	}
	logging.Log().Debug().Str("topic", topic).Msg("received all messages for queue")
	defer dataQueuesMutex.Unlock()
	queue := dataQueues[topic]
	dataQueues[topic] = make([]*entities.Message, 0)
	return queue
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func IsQueue(topic string) bool {
	topicParts := strings.Split(topic, ".")
	return isNumber(topicParts[len(topicParts)-1])
}

func GetQueueComponents(topic string) (string, int) {
	topicParts := strings.Split(topic, ".")
	queueID := topicParts[len(topicParts)-2]
	queueCount, _ := strconv.ParseInt(topicParts[len(topicParts)-1], 10, 64)
	return queueID, int(queueCount)
}
