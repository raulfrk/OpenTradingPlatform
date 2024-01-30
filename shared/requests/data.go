package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"tradingplatform/shared/communication"
	"tradingplatform/shared/communication/command"
	"tradingplatform/shared/communication/subscriber"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/utils"

	"tradingplatform/shared/types"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type DataRequest struct {
	Fingerprint string              `json:"fingerprint"`
	Source      types.Source        `json:"source" validate:"required,min=3,isValidDataSource"`
	AssetClass  types.AssetClass    `json:"assetClass" validate:"required,min=3,isValidAssetClass"`
	Symbol      string              `json:"symbol" validate:"required,min=1"`
	Operation   types.DataRequestOp `json:"operation" validate:"required,min=3,isValidOperation"`
	DataType    types.DataType      `json:"dataType" validate:"required,min=3,isValidDataType"`
	Account     Account             `json:"account" validate:"required,min=3,isValidAccount"`
	StartTime   int64               `json:"startTime" validate:"required,min=0"`
	EndTime     int64               `json:"endTime" validate:"required,min=0,isValidEndTime"`
	TimeFrame   types.TimeFrame     `json:"timeFrame" validate:"required,min=3,isValidDataFrame"`
	NoConfirm   bool                `json:"noConfirm"`
}

func (d *DataRequest) Validate() error {
	v := validator.New()
	v.RegisterValidation("isValidDataSource", IsValidDataSource)
	v.RegisterValidation("isValidAssetClass", IsValidAssetClass)
	v.RegisterValidation("isValidOperation", IsValidOperation)
	v.RegisterValidation("isValidDataType", IsValidDataType)
	v.RegisterValidation("isValidAccount", IsValidAccount)
	v.RegisterValidation("isValidDataFrame", IsValidDataFrame)
	v.RegisterValidation("isValidEndTime", IsValidEndTime)

	err := v.Struct(d)
	return SummarizeError(err)
}

func (r *DataRequest) JSON() []byte {
	js, err := json.Marshal(r)
	if err != nil {
		logging.Log().Error().
			Err(err).
			Msg("marshalling data request to json")
		return []byte{}
	}
	return js
}

func (d *DataRequest) GetSource() types.Source {
	return d.Source
}

func (d *DataRequest) GetAssetClass() types.AssetClass {
	return d.AssetClass
}

func (d *DataRequest) GetSymbol() string {
	return d.Symbol
}

func (d *DataRequest) GetOperation() types.DataRequestOp {
	return d.Operation
}

func (d *DataRequest) GetDataType() types.DataType {
	return d.DataType
}

func (d *DataRequest) GetAccount() Account {
	return d.Account
}

func (d *DataRequest) GetStartTime() int64 {
	return d.StartTime
}

func (d *DataRequest) GetEndTime() int64 {
	return d.EndTime
}

func (d *DataRequest) GetTimeFrame() types.TimeFrame {
	return d.TimeFrame
}

func (d *DataRequest) GetNoConfirm() bool {
	return d.NoConfirm
}

func (d *DataRequest) GetFingerprint() string {
	return d.Fingerprint
}

func NewDataRequest(source types.Source,
	assetClass types.AssetClass,
	symbol string,
	operation types.DataRequestOp,
	dataType types.DataType,
	account Account,
	startTime int64,
	endTime int64,
	timeFrame types.TimeFrame,
	noConfirm bool) DataRequest {

	return DataRequest{
		Source:     source,
		AssetClass: assetClass,
		Symbol:     symbol,
		Operation:  operation,
		DataType:   dataType,
		Account:    account,
		StartTime:  startTime,
		EndTime:    endTime,
		TimeFrame:  timeFrame,
		NoConfirm:  noConfirm,
	}
}

func NewDataRequestFromRaw(source string,
	assetClass string,
	symbol string,
	operation string,
	dataType string,
	account string,
	startTime int64,
	endTime int64,
	timeFrame string,
	noConfirm bool, defaultingFunc func(*DataRequest)) (DataRequest, error) {

	dataRequest := NewDataRequest(types.Source(source),
		types.AssetClass(assetClass),
		symbol,
		types.DataRequestOp(operation),
		types.DataType(dataType),
		Account(account),
		startTime,
		endTime,
		types.TimeFrame(timeFrame),
		noConfirm,
	)

	defaultingFunc(&dataRequest)
	err := dataRequest.Validate()
	return dataRequest, err
}

func NewDataRequestFromExisting(dataRequest *DataRequest, defaultingFunc func(*DataRequest)) (DataRequest, error) {
	return NewDataRequestFromRaw(string(dataRequest.Source),
		string(types.AssetClass(dataRequest.AssetClass)),
		dataRequest.Symbol,
		string(dataRequest.Operation),
		string(dataRequest.DataType),
		string(dataRequest.Account),
		dataRequest.StartTime,
		dataRequest.EndTime,
		string(dataRequest.TimeFrame),
		dataRequest.NoConfirm, defaultingFunc)
}

func RequestData(ctx context.Context, topic utils.Topic, dataRequest DataRequest, onData func(*entities.Message)) error {
	nc, _ := nats.Connect(communication.GetNatsURL())
	defer nc.Close()

	// Ensure that confirmation is required
	dataRequest.NoConfirm = false

	rawReq := command.JSONCommand{
		RootOperation: command.JSONOperationData,
		Request:       dataRequest.JSON(),
	}

	dataRequestJSON := rawReq.JSONWithHeader()

	// Request data
	msg, err := nc.RequestWithContext(ctx, topic.Generate(), []byte(dataRequestJSON))
	if err != nil {
		return fmt.Errorf("error while requesting data %v (topic: %s)", err, topic.Generate())
	}

	var res types.DataResponse
	json.Unmarshal(msg.Data, &res)

	if res.Err != "" {
		return fmt.Errorf(res.Err)
	}

	// Get elements count
	_, count := subscriber.GetQueueComponents(res.ResponseTopic)

	// Prepare counter
	var wg sync.WaitGroup
	wg.Add(count)

	// Subscribe to response topic
	ch := make(chan *nats.Msg, count)
	sub, err := nc.ChanSubscribe(res.ResponseTopic, ch)
	defer sub.Unsubscribe()
	if err != nil {
		return fmt.Errorf("error while subscribing to data response stream %v (topic: %s)", err, res.ResponseTopic)
	}

	// Create subcontext
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			select {
			case m := <-ch:
				var msg entities.Message
				// Empty message is just the receiver telling us that it's ready to receive data
				// ignore it
				if string(m.Data) == "" {
					continue
				}
				proto.Unmarshal(m.Data, &msg)
				onData(&msg)
				wg.Done()
			case <-subCtx.Done():
				return
			}
		}
	}()

	// Publish empty message to start receiving data
	nc.Publish(res.ResponseTopic, []byte(""))
	wg.Wait()
	return nil
}

func DrainChannel[T any](ch chan *T) {
	for range ch {
		continue
	}
}
