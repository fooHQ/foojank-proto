package gateway

import (
	capnp "github.com/foohq/foojank-proto/go/gateway/capnp"
)

// UnregisterAgentResponse is sent by a gateway in response to an UnregisterAgentRequest.
type UnregisterAgentResponse struct {
	Properties []Property
	Error      Error
}

func marshalUnregisterAgentResponse(message *capnp.Message, data UnregisterAgentResponse) error {
	m, err := capnp.NewUnregisterAgentResponse(message.Segment())
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

	errError, err := m.NewError()
	if err != nil {
		return err
	}

	errError.SetCode(data.Error.Code)
	err = errError.SetMessage_(data.Error.Message)
	if err != nil {
		return err
	}

	err = message.Content().SetUnregisterAgentResponse(m)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalUnregisterAgentResponse(message *capnp.Message) (UnregisterAgentResponse, error) {
	v, err := message.Content().UnregisterAgentResponse()
	if err != nil {
		return UnregisterAgentResponse{}, err
	}

	vProperties, err := v.Properties()
	if err != nil {
		return UnregisterAgentResponse{}, err
	}

	properties := make([]Property, 0, vProperties.Len())
	for i := range vProperties.Len() {
		item := vProperties.At(i)
		key, err := item.Key()
		if err != nil {
			return UnregisterAgentResponse{}, err
		}

		value, err := item.Value()
		if err != nil {
			return UnregisterAgentResponse{}, err
		}

		properties = append(properties, Property{
			Key:   key,
			Value: value,
		})
	}

	var respErr Error
	if v.HasError() {
		errStruct, err := v.Error()
		if err != nil {
			return UnregisterAgentResponse{}, err
		}

		msg, err := errStruct.Message_()
		if err != nil {
			return UnregisterAgentResponse{}, err
		}

		respErr = Error{
			Code:    errStruct.Code(),
			Message: msg,
		}
	}

	return UnregisterAgentResponse{
		Properties: properties,
		Error:      respErr,
	}, nil
}
