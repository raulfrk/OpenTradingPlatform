package command

import "tradingplatform/shared/requests"

type JSONOperation string

const (
	JSONOperationQuit            JSONOperation = "quit"
	JSONOperationStream          JSONOperation = "stream"
	JSONOperationStreamSubscribe JSONOperation = "stream-subscribe"

	JSONOperationData JSONOperation = "data"
)

type JSONCommand struct {
	RootOperation          JSONOperation                   `json:"operation"`
	StreamRequest          requests.StreamRequest          `json:"streamRequest"`
	DataRequest            requests.DataRequest            `json:"dataRequest"`
	StreamSubscribeRequest requests.StreamSubscribeRequest `json:"streamSubscribeRequest"`
}
