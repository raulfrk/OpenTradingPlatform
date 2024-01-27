package logging

import (
	"io"

	"tradingplatform/shared/entities"
	"tradingplatform/shared/types"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

var log *zerolog.Logger

func SetLogger(l *zerolog.Logger) {
	log = l
}

func Log() *zerolog.Logger {
	if log == nil {
		panic("Logger not initialized")
	}
	return log
}

type NatsWriter struct {
	nc    *nats.Conn
	topic string
}

func NewNatsWriter(nc *nats.Conn, topic string) *NatsWriter {
	return &NatsWriter{
		nc:    nc,
		topic: topic,
	}
}

func (nw *NatsWriter) Write(p []byte) (n int, err error) {
	message := entities.Message{
		Topic:    nw.topic,
		Payload:  p,
		DataType: string(types.Log),
	}
	b, err := proto.Marshal(&message)
	if err != nil {
		return 0, err
	}

	err = nw.nc.Publish(nw.topic, b)

	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func NewMultiLevelLogger(component types.Component, writers ...io.Writer) zerolog.Logger {
	multi := zerolog.MultiLevelWriter(writers...)
	logger := zerolog.New(multi).With().Str("component", string(component)).Timestamp().Stack().Logger()
	return logger
}
