package agent

import (
	"errors"

	capnp "github.com/foohq/foojank-proto/go/agent/capnp"
)

// StopWorkerResponse is a response to a StopWorkerRequest.
type StopWorkerResponse struct {
	Error error
}

func marshalStopWorkerResponse(message *capnp.Message, data StopWorkerResponse) error {
	m, err := capnp.NewStopWorkerResponse(message.Segment())
	if err != nil {
		return err
	}

	if data.Error != nil {
		err = m.SetError(data.Error.Error())
		if err != nil {
			return err
		}
	}

	err = message.Content().SetStopWorkerResponse(m)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalStopWorkerResponse(message *capnp.Message) (StopWorkerResponse, error) {
	v, err := message.Content().StopWorkerResponse()
	if err != nil {
		return StopWorkerResponse{}, err
	}

	errMsg, err := v.Error()
	if err != nil {
		return StopWorkerResponse{}, err
	}

	var respErr error
	if errMsg != "" {
		respErr = errors.New(errMsg)
	}

	return StopWorkerResponse{
		Error: respErr,
	}, nil
}
