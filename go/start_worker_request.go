package proto

import (
	capnp "github.com/foohq/foojank-proto/go/agent"
)

// StartWorkerRequest is a request to start a new worker.
type StartWorkerRequest struct {
	Command string
	Args    []string
	Env     []string
}

func marshalStartWorkerRequest(message *capnp.Message, data StartWorkerRequest) error {
	m, err := capnp.NewStartWorkerRequest(message.Segment())
	if err != nil {
		return err
	}

	err = m.SetCommand(data.Command)
	if err != nil {
		return err
	}

	argsList, err := newTextList(message.Segment(), data.Args)
	if err != nil {
		return err
	}

	err = m.SetArgs(argsList)
	if err != nil {
		return err
	}

	envList, err := newTextList(message.Segment(), data.Env)
	if err != nil {
		return err
	}

	err = m.SetEnv(envList)
	if err != nil {
		return err
	}

	err = message.Content().SetStartWorkerRequest(m)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalStartWorkerRequest(message *capnp.Message) (StartWorkerRequest, error) {
	v, err := message.Content().StartWorkerRequest()
	if err != nil {
		return StartWorkerRequest{}, err
	}

	command, err := v.Command()
	if err != nil {
		return StartWorkerRequest{}, err
	}

	vArgs, err := v.Args()
	if err != nil {
		return StartWorkerRequest{}, err
	}

	args, err := textListToStringSlice(vArgs)
	if err != nil {
		return StartWorkerRequest{}, err
	}

	vEnv, err := v.Env()
	if err != nil {
		return StartWorkerRequest{}, err
	}

	env, err := textListToStringSlice(vEnv)
	if err != nil {
		return StartWorkerRequest{}, err
	}

	return StartWorkerRequest{
		Command: command,
		Args:    args,
		Env:     env,
	}, nil
}
