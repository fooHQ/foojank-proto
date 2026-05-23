package proto

import (
	capnp "github.com/foohq/foojank-proto/go/agent"
)

// UpdateClientInfo contains information about a client.
type UpdateClientInfo struct {
	Username string
	Hostname string
	System   string
	Address  string
}

func marshalUpdateClientInfo(message *capnp.Message, data UpdateClientInfo) error {
	m, err := capnp.NewUpdateClientInfo(message.Segment())
	if err != nil {
		return err
	}

	err = m.SetUsername(data.Username)
	if err != nil {
		return err
	}

	err = m.SetHostname(data.Hostname)
	if err != nil {
		return err
	}

	err = m.SetSystem(data.System)
	if err != nil {
		return err
	}

	err = m.SetAddress(data.Address)
	if err != nil {
		return err
	}

	err = message.Content().SetUpdateClientInfo(m)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalUpdateClientInfo(message *capnp.Message) (UpdateClientInfo, error) {
	v, err := message.Content().UpdateClientInfo()
	if err != nil {
		return UpdateClientInfo{}, err
	}

	username, err := v.Username()
	if err != nil {
		return UpdateClientInfo{}, err
	}

	hostname, err := v.Hostname()
	if err != nil {
		return UpdateClientInfo{}, err
	}

	system, err := v.System()
	if err != nil {
		return UpdateClientInfo{}, err
	}

	address, err := v.Address()
	if err != nil {
		return UpdateClientInfo{}, err
	}

	return UpdateClientInfo{
		Username: username,
		Hostname: hostname,
		System:   system,
		Address:  address,
	}, nil
}
