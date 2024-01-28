package communication

import (
	"github.com/nats-io/nats.go"
)

var natsURL = nats.DefaultURL

func SetNatsURL(url string) {
	natsURL = url
}

func GetNatsURL() string {
	return natsURL
}
