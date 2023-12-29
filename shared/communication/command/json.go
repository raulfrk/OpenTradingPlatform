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

	JSONOperationData                    JSONOperation = "data"
	JSONOperationSentimentAnalysisFromDB JSONOperation = "sentiment-analysis-from-db"
)

type JSONCommand struct {
	RootOperation JSONOperation   `json:"operation"`
	Request       json.RawMessage `json:"request"`
}

func (r *JSONCommand) JSONWithHeader() string {
	j, _ := json.Marshal(r)
	return fmt.Sprintf("json%s", j)
}
