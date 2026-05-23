package proto

import (
	capnp "github.com/foohq/foojank-proto/go/agent"
)

// UpdateWorkerStdio contains data from a worker's stdout or stdin.
type UpdateWorkerStdio struct {
	Data []byte
}

func marshalUpdateWorkerStdio(message *capnp.Message, data UpdateWorkerStdio) error {
	m, err := capnp.NewUpdateWorkerStdio(message.Segment())
	if err != nil {
		return err
	}

	err = m.SetData(data.Data)
	if err != nil {
		return err
	}

	err = message.Content().SetUpdateWorkerStdio(m)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalUpdateWorkerStdio(message *capnp.Message) (UpdateWorkerStdio, error) {
	v, err := message.Content().UpdateWorkerStdio()
	if err != nil {
		return UpdateWorkerStdio{}, err
	}

	data, err := v.Data()
	if err != nil {
		return UpdateWorkerStdio{}, err
	}

	return UpdateWorkerStdio{
		Data: data,
	}, nil
}
