package utils

import (
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func HandleEntityQueue[T protoreflect.ProtoMessage](msg []*entities.Message, sample T, f func([]T)) {
	entity := make([]T, len(msg))
	for i, m := range msg {
		clone := proto.Clone(sample)
		err := proto.Unmarshal(m.Payload, clone)
		if err != nil {
			logging.Log().Error().Err(err).Msg("unmarshalling entity")
			return
		}
		entity[i] = clone.(T)
	}

	f(entity)
}

func HandleEntity[T protoreflect.ProtoMessage](msg *entities.Message, sample T, f func(T)) {
	err := proto.Unmarshal(msg.Payload, sample)
	if err != nil {
		logging.Log().Error().Err(err).Msg("unmarshalling entity")
		return
	}

	f(sample)
}

func HandleEntityWithConversion[T protoreflect.ProtoMessage, I any](msg *entities.Message, sample T, f func(I), converter func(T) I) {
	entity := proto.Clone(sample)
	err := proto.Unmarshal(msg.Payload, entity)
	if err != nil {
		logging.Log().Error().Err(err).Msg("unmarshalling entity")
		return
	}

	f(converter(entity.(T)))
}

// HandleEntityQueueWithConversion handles a queue of messages and performs conversion on the entities.
// It takes in a slice of messages, a sample entity of type T, a function f to handle the converted entities,
// and a converter function to convert the entities from type T to type I.
// The function iterates over the messages, unmarshals each message into a clone of the sample entity,
// and stores the converted entities in a slice. Finally, it calls the function f with the converted entities.
func HandleEntityQueueWithConversion[T protoreflect.ProtoMessage, I any](msg []*entities.Message, sample T, f func([]I), converter func([]T) []I) {
	entity := make([]T, len(msg))
	for i, m := range msg {
		clone := proto.Clone(sample)
		err := proto.Unmarshal(m.Payload, clone)
		if err != nil {
			logging.Log().Error().Err(err).Msg("unmarshalling entity")
			return
		}
		entity[i] = clone.(T)
	}

	f(converter(entity))
}
