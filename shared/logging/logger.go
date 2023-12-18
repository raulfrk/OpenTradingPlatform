package logging

import (
	"io"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
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
	err = nw.nc.Publish(nw.topic, p)

	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func NewMultiLevelLogger(writers ...io.Writer) zerolog.Logger {
	multi := zerolog.MultiLevelWriter(writers...)
	logger := zerolog.New(multi).With().Timestamp().Stack().Logger()
	return logger
}
