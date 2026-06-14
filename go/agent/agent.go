// Package agent provides functions for marshaling and unmarshaling messages.
package agent

import (
	"errors"

	capnplib "capnproto.org/go/capnp/v3"

	capnp "github.com/foohq/foojank-proto/go/agent/capnp"
)

var (
	ErrUnknownMessage = errors.New("unknown message")
)

type Envelope struct {
	Subject string
	Payload any
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
	case StartWorkerRequest:
		err = marshalStartWorkerRequest(&payload, v)
	case StartWorkerResponse:
		err = marshalStartWorkerResponse(&payload, v)
	case StopWorkerRequest:
		err = marshalStopWorkerRequest(&payload, v)
	case StopWorkerResponse:
		err = marshalStopWorkerResponse(&payload, v)
	case UpdateWorkerStatus:
		err = marshalUpdateWorkerStatus(&payload, v)
	case UpdateWorkerStdio:
		err = marshalUpdateWorkerStdio(&payload, v)
	case UpdateClientInfo:
		err = marshalUpdateClientInfo(&payload, v)
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
	case payload.Content().HasStartWorkerRequest():
		v, err = unmarshalStartWorkerRequest(&payload)
	case payload.Content().HasStartWorkerResponse():
		v, err = unmarshalStartWorkerResponse(&payload)
	case payload.Content().HasStopWorkerRequest():
		v, err = unmarshalStopWorkerRequest(&payload)
	case payload.Content().HasStopWorkerResponse():
		v, err = unmarshalStopWorkerResponse(&payload)
	case payload.Content().HasUpdateWorkerStatus():
		v, err = unmarshalUpdateWorkerStatus(&payload)
	case payload.Content().HasUpdateWorkerStdio():
		v, err = unmarshalUpdateWorkerStdio(&payload)
	case payload.Content().HasUpdateClientInfo():
		v, err = unmarshalUpdateClientInfo(&payload)
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

// CmdStartWorkerSubject returns the NATS subject for sending a start worker command to an agent.
func CmdStartWorkerSubject(gatewayID, agentID, workerID string) string {
	return replaceStringPlaceholders(capnp.CmdStartWorkerT, gatewayID, agentID, workerID)
}

// CmdStopWorkerSubject returns the NATS subject for sending a stop worker command to an agent.
func CmdStopWorkerSubject(gatewayID, agentID, workerID string) string {
	return replaceStringPlaceholders(capnp.CmdStopWorkerT, gatewayID, agentID, workerID)
}

// CmdWriteStdinSubject returns the NATS subject for sending stdin to a worker via an agent.
func CmdWriteStdinSubject(gatewayID, agentID, workerID string) string {
	return replaceStringPlaceholders(capnp.CmdWriteStdinT, gatewayID, agentID, workerID)
}

// EvtStartWorkerSubject returns the NATS subject for a worker start event.
func EvtStartWorkerSubject(gatewayID, agentID, workerID string) string {
	return replaceStringPlaceholders(capnp.EvtStartWorkerT, gatewayID, agentID, workerID)
}

// EvtStopWorkerSubject returns the NATS subject for a worker stop event.
func EvtStopWorkerSubject(gatewayID, agentID, workerID string) string {
	return replaceStringPlaceholders(capnp.EvtStopWorkerT, gatewayID, agentID, workerID)
}

// EvtWorkerStatusSubject returns the NATS subject for a worker status event.
func EvtWorkerStatusSubject(gatewayID, agentID, workerID string) string {
	return replaceStringPlaceholders(capnp.EvtWorkerStatusT, gatewayID, agentID, workerID)
}

// EvtWorkerStdoutSubject returns the NATS subject for a worker stdout event.
func EvtWorkerStdoutSubject(gatewayID, agentID, workerID string) string {
	return replaceStringPlaceholders(capnp.EvtWorkerStdoutT, gatewayID, agentID, workerID)
}

// EvtAgentInfoSubject returns the NATS subject for an agent info event.
func EvtAgentInfoSubject(gatewayID, agentID string) string {
	return replaceStringPlaceholders(capnp.EvtAgentInfoT, gatewayID, agentID)
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

func newTextList(segment *capnplib.Segment, ss []string) (capnplib.TextList, error) {
	tl, err := capnplib.NewTextList(segment, int32(len(ss)))
	if err != nil {
		return capnplib.TextList{}, err
	}

	for i, s := range ss {
		err := tl.Set(i, s)
		if err != nil {
			return capnplib.TextList{}, err
		}
	}

	return tl, nil
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

func textListToStringSlice(list capnplib.TextList) ([]string, error) {
	if list.Len() == 0 {
		return nil, nil
	}
	result := make([]string, 0, list.Len())
	for i := 0; i < list.Len(); i++ {
		v, err := list.At(i)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}
