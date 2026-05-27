package agent

import (
	capnp "github.com/foohq/foojank-proto/go/agent/capnp"
)

// StopWorkerRequest is a request to stop a running worker.
type StopWorkerRequest struct{}

func marshalStopWorkerRequest(message *capnp.Message, _ StopWorkerRequest) error {
	m, err := capnp.NewStopWorkerRequest(message.Segment())
	if err != nil {
		return err
	}

	err = message.Content().SetStopWorkerRequest(m)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalStopWorkerRequest(message *capnp.Message) (StopWorkerRequest, error) {
	_, err := message.Content().StopWorkerRequest()
	if err != nil {
		return StopWorkerRequest{}, err
	}

	return StopWorkerRequest{}, nil
}
