package gateway

import (
	capnp "github.com/foohq/foojank-proto/go/gateway/capnp"
)

// RegisterAgentResponse is sent by a gateway in response to a RegisterAgentRequest.
type RegisterAgentResponse struct {
	Properties []Property
	Error      Error
}

func marshalRegisterAgentResponse(message *capnp.Message, data RegisterAgentResponse) error {
	m, err := capnp.NewRegisterAgentResponse(message.Segment())
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

	err = message.Content().SetRegisterAgentResponse(m)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalRegisterAgentResponse(message *capnp.Message) (RegisterAgentResponse, error) {
	v, err := message.Content().RegisterAgentResponse()
	if err != nil {
		return RegisterAgentResponse{}, err
	}

	vProperties, err := v.Properties()
	if err != nil {
		return RegisterAgentResponse{}, err
	}

	properties := make([]Property, 0, vProperties.Len())
	for i := range vProperties.Len() {
		item := vProperties.At(i)
		key, err := item.Key()
		if err != nil {
			return RegisterAgentResponse{}, err
		}

		value, err := item.Value()
		if err != nil {
			return RegisterAgentResponse{}, err
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
			return RegisterAgentResponse{}, err
		}

		msg, err := errStruct.Message_()
		if err != nil {
			return RegisterAgentResponse{}, err
		}

		respErr = Error{
			Code:    errStruct.Code(),
			Message: msg,
		}
	}

	return RegisterAgentResponse{
		Properties: properties,
		Error:      respErr,
	}, nil
}
