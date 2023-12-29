package types

import (
	"encoding/json"
)

type DataResponse struct {
	Err           string
	Message       string
	Status        OpStatus
	ResponseTopic string
	// TODO: handle open queues
}

func NewDataError(err error) DataResponse {
	return NewDataResponse(Failure, "", err, "")
}

func NewDataResponse(status OpStatus, message string, err error, responseTopic string) DataResponse {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	newResponse := DataResponse{
		Err:           errStr,
		Message:       message,
		Status:        status,
		ResponseTopic: responseTopic,
	}

	return newResponse
}

type Response struct {
	Err     string
	Message string
	Status  OpStatus
}

type StreamResponse struct {
	Response
	// TODO: Implement response that provides list of topics
	Topics  string
	Streams string
}

func NewError(err error) Response {
	return NewResponse(Failure, "", err)
}

func NewResponse(status OpStatus, message string, err error) Response {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	newResponse := Response{
		Err:     errStr,
		Message: message,
		Status:  status,
	}

	return newResponse
}

func NewStreamError(err error, streams string) StreamResponse {
	return NewStreamResponse(Failure, "", err, streams)
}

func NewStreamResponse(status OpStatus, message string, err error, streams string) StreamResponse {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	newResponse := StreamResponse{

		Streams: streams,
	}
	newResponse.Err = errStr
	newResponse.Message = message
	newResponse.Status = status

	return newResponse
}

func NewStreamResponseTopic(status OpStatus, message string, err error, topics string) StreamResponse {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	newResponse := StreamResponse{

		Topics: topics,
	}
	newResponse.Err = errStr
	newResponse.Message = message
	newResponse.Status = status

	return newResponse
}

func NewStreamTopicError(err error, topics string) StreamResponse {
	return NewStreamResponseTopic(Failure, "", err, topics)
}

func (r Response) Respond() string {
	response, err := json.Marshal(r)
	if err != nil {
		marshalError := &Response{
			Err:     "Error marshalling response",
			Message: "",
			Status:  Failure,
		}
		response, _ = json.Marshal(marshalError)
	}
	return string(response)
}

func (r StreamResponse) Respond() string {
	response, err := json.Marshal(r)
	if err != nil {
		marshalError := &StreamResponse{
			Response: Response{
				Err:     "Error marshalling response",
				Message: "",
				Status:  Failure,
			},
		}
		response, _ = json.Marshal(marshalError)
	}
	return string(response)
}

func (r DataResponse) Respond() string {
	response, err := json.Marshal(r)
	if err != nil {
		marshalError := &StreamResponse{
			Response: Response{
				Err:     "Error marshalling response",
				Message: "",
				Status:  Failure,
			},
		}
		response, _ = json.Marshal(marshalError)
	}
	return string(response)
}
