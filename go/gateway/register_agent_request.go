package gateway

import (
	capnp "github.com/foohq/foojank-proto/go/gateway/capnp"
)

// RegisterAgentRequest is sent by a client to register an agent with a gateway.
type RegisterAgentRequest struct {
	Properties []Property
}

func marshalRegisterAgentRequest(message *capnp.Message, data RegisterAgentRequest) error {
	m, err := capnp.NewRegisterAgentRequest(message.Segment())
	if err != nil {
		return err
	}

	propertiesList, err := m.NewProperties(int32(len(data.Properties)))
	if err != nil {
		return err
	}

	for i, p := range data.Properties {
		item := propertiesList.At(i)
		err := item.SetKey(p.Key)
		if err != nil {
			return err
		}

		err = item.SetValue(p.Value)
		if err != nil {
			return err
		}
	}

	err = message.Content().SetRegisterAgentRequest(m)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalRegisterAgentRequest(message *capnp.Message) (RegisterAgentRequest, error) {
	v, err := message.Content().RegisterAgentRequest()
	if err != nil {
		return RegisterAgentRequest{}, err
	}

	vProperties, err := v.Properties()
	if err != nil {
		return RegisterAgentRequest{}, err
	}

	properties := make([]Property, 0, vProperties.Len())
	for i := range vProperties.Len() {
		item := vProperties.At(i)
		key, err := item.Key()
		if err != nil {
			return RegisterAgentRequest{}, err
		}

		value, err := item.Value()
		if err != nil {
			return RegisterAgentRequest{}, err
		}

		properties = append(properties, Property{
			Key:   key,
			Value: value,
		})
	}

	return RegisterAgentRequest{
		Properties: properties,
	}, nil
}
