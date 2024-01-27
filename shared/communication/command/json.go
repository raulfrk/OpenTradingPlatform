package command

import (
	"encoding/json"
	"fmt"
)

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
	CancelKey     string          `json:"cancelKey"`
}

func (r *JSONCommand) JSONWithHeader() string {
	j, _ := json.Marshal(r)
	return fmt.Sprintf("json%s", j)
}
