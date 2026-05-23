package proto

import (
	"errors"

	capnp "github.com/foohq/foojank-proto/go/agent"
)

const (
	ExitSuccess         = capnp.ExitSuccess
	ExitFailure         = capnp.ExitFailure
	ExitCommandNotFound = capnp.ExitCommandNotFound
	ExitInterrupted     = capnp.ExitInterrupted
)

// UpdateWorkerStatus is used to update the status of a worker.
type UpdateWorkerStatus struct {
	Status int64
	Error  error
}

func marshalUpdateWorkerStatus(message *capnp.Message, data UpdateWorkerStatus) error {
	m, err := capnp.NewUpdateWorkerStatus(message.Segment())
	if err != nil {
		return err
	}

	m.SetStatus(data.Status)

	if data.Error != nil {
		err = m.SetError(data.Error.Error())
		if err != nil {
			return err
		}
	}

	err = message.Content().SetUpdateWorkerStatus(m)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalUpdateWorkerStatus(message *capnp.Message) (UpdateWorkerStatus, error) {
	v, err := message.Content().UpdateWorkerStatus()
	if err != nil {
		return UpdateWorkerStatus{}, err
	}

	errMsg, err := v.Error()
	if err != nil {
		return UpdateWorkerStatus{}, err
	}

	var respErr error
	if errMsg != "" {
		respErr = errors.New(errMsg)
	}

	return UpdateWorkerStatus{
		Status: v.Status(),
		Error:  respErr,
	}, nil
}
