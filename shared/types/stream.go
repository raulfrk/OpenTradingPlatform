package types

type StreamRequestOp string

const (
	StreamAddOp    StreamRequestOp = "add"
	StreamRemoveOp StreamRequestOp = "remove"

	StreamGetOp StreamRequestOp = "get"
)

func GetStreamRequestOpMap() map[string]StreamRequestOp {
	return map[string]StreamRequestOp{
		"add":    StreamAddOp,
		"remove": StreamRemoveOp,
		"stop":   StreamAddOp,
		"get":    StreamGetOp,
	}
}

func GetStreamSubscribeRequestOpMap() map[string]StreamRequestOp {
	return map[string]StreamRequestOp{
		"add":    StreamAddOp,
		"remove": StreamRemoveOp,
	}
}
