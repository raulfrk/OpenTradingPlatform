package command

import "encoding/json"

type JSONOperation string

const (
	JSONOperationQuit            JSONOperation = "quit"
	JSONOperationStream          JSONOperation = "stream"
	JSONOperationStreamSubscribe JSONOperation = "stream-subscribe"

	JSONOperationData JSONOperation = "data"
)

type JSONCommand struct {
	RootOperation JSONOperation   `json:"operation"`
	Request       json.RawMessage `json:"request"`
}
