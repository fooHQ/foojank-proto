// Package gateway provides functions for building gateway communication subjects.
package gateway

import (
	"errors"

	capnplib "capnproto.org/go/capnp/v3"

	capnp "github.com/foohq/foojank-proto/go/gateway/capnp"
)

var (
	ErrUnknownMessage = errors.New("unknown message")
)

type Envelope struct {
	Subject string
	Payload any
}

type Property struct {
	Key   string
	Value string
}

type Error struct {
	Code    int32
	Message string
}

// Marshal serializes the given data into a byte slice.
// It supports various request and response types defined in the proto package.
func Marshal(envelope Envelope) ([]byte, error) {
	env, err := newEnvelope()
	if err != nil {
		return nil, err
	}

	err = env.SetSubject(envelope.Subject)
	if err != nil {
		return nil, err
	}

	payload, err := env.NewPayload()
	if err != nil {
		return nil, err
	}

	switch v := envelope.Payload.(type) {
	case RegisterAgentRequest:
		err = marshalRegisterAgentRequest(&payload, v)
	case RegisterAgentResponse:
		err = marshalRegisterAgentResponse(&payload, v)
	case UnregisterAgentRequest:
		err = marshalUnregisterAgentRequest(&payload, v)
	case UnregisterAgentResponse:
		err = marshalUnregisterAgentResponse(&payload, v)
	default:
		err = ErrUnknownMessage
	}
	if err != nil {
		return nil, err
	}

	err = env.SetPayload(payload)
	if err != nil {
		return nil, err
	}

	return env.Message().Marshal()
}

// Unmarshal deserializes the given byte slice into a message object.
// It returns an interface{} which can be type-asserted to the specific message type.
func Unmarshal(b []byte) (Envelope, error) {
	envelope, err := unmarshalEnvelope(b)
	if err != nil {
		return Envelope{}, err
	}

	subject, err := envelope.Subject()
	if err != nil {
		return Envelope{}, err
	}

	payload, err := envelope.Payload()
	if err != nil {
		return Envelope{}, err
	}

	var v any
	switch {
	case payload.Content().HasRegisterAgentRequest():
		v, err = unmarshalRegisterAgentRequest(&payload)
	case payload.Content().HasRegisterAgentResponse():
		v, err = unmarshalRegisterAgentResponse(&payload)
	case payload.Content().HasUnregisterAgentRequest():
		v, err = unmarshalUnregisterAgentRequest(&payload)
	case payload.Content().HasUnregisterAgentResponse():
		v, err = unmarshalUnregisterAgentResponse(&payload)
	default:
		err = ErrUnknownMessage
	}
	if err != nil {
		return Envelope{}, err
	}

	return Envelope{
		Subject: subject,
		Payload: v,
	}, nil
}

// CmdRegisterAgentSubject returns the NATS subject for sending an agent registration command to a gateway.
func CmdRegisterAgentSubject(gatewayID, agentID string) string {
	return replaceStringPlaceholders(capnp.CmdRegisterT, gatewayID, agentID)
}

// CmdUnregisterAgentSubject returns the NATS subject for sending an agent unregistration command to a gateway.
func CmdUnregisterAgentSubject(gatewayID, agentID string) string {
	return replaceStringPlaceholders(capnp.CmdUnregisterT, gatewayID, agentID)
}

// EvtRegisterAgentSubject returns the NATS subject for an agent registration event.
func EvtRegisterAgentSubject(gatewayID, agentID string) string {
	return replaceStringPlaceholders(capnp.EvtRegisterT, gatewayID, agentID)
}

// EvtUnregisterAgentSubject returns the NATS subject for an agent unregistration event.
func EvtUnregisterAgentSubject(gatewayID, agentID string) string {
	return replaceStringPlaceholders(capnp.EvtUnregisterT, gatewayID, agentID)
}

func replaceStringPlaceholders(s string, values ...string) string {
	result := s
	valIndex := 0

	for valIndex < len(values) {
		found := false
		for i := 0; i < len(result)-1; i++ {
			if result[i] == '%' && result[i+1] == 's' {
				result = result[:i] + values[valIndex] + result[i+2:]
				valIndex++
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	return result
}

func newEnvelope() (capnp.Envelope, error) {
	arena := capnplib.SingleSegment(nil)
	_, seg, err := capnplib.NewMessage(arena)
	if err != nil {
		return capnp.Envelope{}, err
	}

	envelope, err := capnp.NewRootEnvelope(seg)
	if err != nil {
		return capnp.Envelope{}, err
	}

	return envelope, nil
}

func unmarshalEnvelope(b []byte) (capnp.Envelope, error) {
	capMsg, err := capnplib.Unmarshal(b)
	if err != nil {
		return capnp.Envelope{}, err
	}

	envelope, err := capnp.ReadRootEnvelope(capMsg)
	if err != nil {
		return capnp.Envelope{}, err
	}

	return envelope, nil
}
