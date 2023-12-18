package types

type TimeFrame string
type DataRequestOp string

const (
	OneMin    TimeFrame     = "1min"
	OneHour                 = "1hour"
	OneDay                  = "1day"
	OneWeek                 = "1week"
	OneMonth                = "1month"
	DataGetOp DataRequestOp = "get"
)

func GetTimeFrameMap() map[string]TimeFrame {
	return map[string]TimeFrame{
		"1min":   OneMin,
		"1hour":  OneHour,
		"1day":   OneDay,
		"1week":  OneWeek,
		"1month": OneMonth,
	}
}

func GetDataRequestOpMap() map[string]DataRequestOp {
	return map[string]DataRequestOp{
		"get": DataGetOp,
	}
}
