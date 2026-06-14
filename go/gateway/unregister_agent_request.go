package gateway

import (
	capnp "github.com/foohq/foojank-proto/go/gateway/capnp"
)

// UnregisterAgentRequest is sent by a client to unregister an agent from a gateway.
type UnregisterAgentRequest struct {
	Properties []Property
}

func marshalUnregisterAgentRequest(message *capnp.Message, data UnregisterAgentRequest) error {
	m, err := capnp.NewUnregisterAgentRequest(message.Segment())
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

	err = message.Content().SetUnregisterAgentRequest(m)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalUnregisterAgentRequest(message *capnp.Message) (UnregisterAgentRequest, error) {
	v, err := message.Content().UnregisterAgentRequest()
	if err != nil {
		return UnregisterAgentRequest{}, err
	}

	vProperties, err := v.Properties()
	if err != nil {
		return UnregisterAgentRequest{}, err
	}

	properties := make([]Property, 0, vProperties.Len())
	for i := range vProperties.Len() {
		item := vProperties.At(i)
		key, err := item.Key()
		if err != nil {
			return UnregisterAgentRequest{}, err
		}

		value, err := item.Value()
		if err != nil {
			return UnregisterAgentRequest{}, err
		}

		properties = append(properties, Property{
			Key:   key,
			Value: value,
		})
	}

	return UnregisterAgentRequest{
		Properties: properties,
	}, nil
}
