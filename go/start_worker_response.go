package proto

import (
	"errors"

	capnp "github.com/foohq/foojank-proto/go/agent"
)

// StartWorkerResponse is a response to a StartWorkerRequest.
type StartWorkerResponse struct {
	Error error
}

func marshalStartWorkerResponse(message *capnp.Message, data StartWorkerResponse) error {
	m, err := capnp.NewStartWorkerResponse(message.Segment())
	if err != nil {
		return err
	}

	if data.Error != nil {
		err = m.SetError(data.Error.Error())
		if err != nil {
			return err
		}
	}

	err = message.Content().SetStartWorkerResponse(m)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalStartWorkerResponse(message *capnp.Message) (StartWorkerResponse, error) {
	v, err := message.Content().StartWorkerResponse()
	if err != nil {
		return StartWorkerResponse{}, err
	}

	errMsg, err := v.Error()
	if err != nil {
		return StartWorkerResponse{}, err
	}

	var respErr error
	if errMsg != "" {
		respErr = errors.New(errMsg)
	}

	return StartWorkerResponse{
		Error: respErr,
	}, nil
}
