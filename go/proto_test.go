package proto_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	proto "github.com/foohq/foojank-proto/go"
)

func TestMarshalUnmarshal(t *testing.T) {
	testError := errors.New("test error")

	tests := []struct {
		name        string
		input       proto.Envelope
		want        proto.Envelope
		wantMarshal bool
		wantErr     error
	}{
		{
			name: "StartWorkerRequest",
			input: proto.Envelope{
				Subject: proto.CmdStartWorkerSubject("agent1", "worker1"),
				Payload: proto.StartWorkerRequest{
					Command: "cmd",
					Args:    []string{"arg1", "arg2"},
					Env:     []string{"KEY1=val1", "KEY2=val2"},
				},
			},
			want: proto.Envelope{
				Subject: proto.CmdStartWorkerSubject("agent1", "worker1"),
				Payload: proto.StartWorkerRequest{
					Command: "cmd",
					Args:    []string{"arg1", "arg2"},
					Env:     []string{"KEY1=val1", "KEY2=val2"},
				},
			},
			wantMarshal: true,
		},
		{
			name: "StartWorkerRequest with empty slices",
			input: proto.Envelope{
				Subject: proto.CmdStartWorkerSubject("agent1", "worker1"),
				Payload: proto.StartWorkerRequest{
					Command: "cmd",
					Args:    []string{},
					Env:     []string{},
				},
			},
			want: proto.Envelope{
				Subject: proto.CmdStartWorkerSubject("agent1", "worker1"),
				Payload: proto.StartWorkerRequest{
					Command: "cmd",
					Args:    nil,
					Env:     nil,
				},
			},
			wantMarshal: true,
		},
		{
			name: "StartWorkerResponse without error",
			input: proto.Envelope{
				Subject: proto.EvtStartWorkerSubject("agent1", "worker1"),
				Payload: proto.StartWorkerResponse{
					Error: nil,
				},
			},
			want: proto.Envelope{
				Subject: proto.EvtStartWorkerSubject("agent1", "worker1"),
				Payload: proto.StartWorkerResponse{
					Error: nil,
				},
			},
			wantMarshal: true,
		},
		{
			name: "StartWorkerResponse with error",
			input: proto.Envelope{
				Subject: proto.EvtStartWorkerSubject("agent1", "worker1"),
				Payload: proto.StartWorkerResponse{
					Error: testError,
				},
			},
			want: proto.Envelope{
				Subject: proto.EvtStartWorkerSubject("agent1", "worker1"),
				Payload: proto.StartWorkerResponse{
					Error: testError,
				},
			},
			wantMarshal: true,
		},
		{
			name: "StopWorkerRequest",
			input: proto.Envelope{
				Subject: proto.CmdStopWorkerSubject("agent1", "worker1"),
				Payload: proto.StopWorkerRequest{},
			},
			want: proto.Envelope{
				Subject: proto.CmdStopWorkerSubject("agent1", "worker1"),
				Payload: proto.StopWorkerRequest{},
			},
			wantMarshal: true,
		},
		{
			name: "StopWorkerResponse without error",
			input: proto.Envelope{
				Subject: proto.EvtStopWorkerSubject("agent1", "worker1"),
				Payload: proto.StopWorkerResponse{
					Error: nil,
				},
			},
			want: proto.Envelope{
				Subject: proto.EvtStopWorkerSubject("agent1", "worker1"),
				Payload: proto.StopWorkerResponse{
					Error: nil,
				},
			},
			wantMarshal: true,
		},
		{
			name: "StopWorkerResponse with error",
			input: proto.Envelope{
				Subject: proto.EvtStopWorkerSubject("agent1", "worker1"),
				Payload: proto.StopWorkerResponse{
					Error: testError,
				},
			},
			want: proto.Envelope{
				Subject: proto.EvtStopWorkerSubject("agent1", "worker1"),
				Payload: proto.StopWorkerResponse{
					Error: testError,
				},
			},
			wantMarshal: true,
		},
		{
			name: "UpdateWorkerStatus",
			input: proto.Envelope{
				Subject: proto.EvtWorkerStatusSubject("agent1", "worker1"),
				Payload: proto.UpdateWorkerStatus{
					Status: 42,
				},
			},
			want: proto.Envelope{
				Subject: proto.EvtWorkerStatusSubject("agent1", "worker1"),
				Payload: proto.UpdateWorkerStatus{
					Status: 42,
				},
			},
			wantMarshal: true,
		},
		{
			name: "UpdateWorkerStdio",
			input: proto.Envelope{
				Subject: proto.EvtWorkerStdoutSubject("agent1", "worker1"),
				Payload: proto.UpdateWorkerStdio{
					Data: []byte("Hello, World!"),
				},
			},
			want: proto.Envelope{
				Subject: proto.EvtWorkerStdoutSubject("agent1", "worker1"),
				Payload: proto.UpdateWorkerStdio{
					Data: []byte("Hello, World!"),
				},
			},
			wantMarshal: true,
		},
		{
			name: "UpdateClientInfo",
			input: proto.Envelope{
				Subject: proto.EvtAgentInfoSubject("agent1"),
				Payload: proto.UpdateClientInfo{
					Username: "testuser",
					Hostname: "testhost",
					System:   "linux",
					Address:  "192.168.1.1",
				},
			},
			want: proto.Envelope{
				Subject: proto.EvtAgentInfoSubject("agent1"),
				Payload: proto.UpdateClientInfo{
					Username: "testuser",
					Hostname: "testhost",
					System:   "linux",
					Address:  "192.168.1.1",
				},
			},
			wantMarshal: true,
		},
		{
			name: "UpdateClientInfo with empty fields",
			input: proto.Envelope{
				Subject: proto.EvtAgentInfoSubject("agent1"),
				Payload: proto.UpdateClientInfo{
					Username: "",
					Hostname: "",
					System:   "",
					Address:  "",
				},
			},
			want: proto.Envelope{
				Subject: proto.EvtAgentInfoSubject("agent1"),
				Payload: proto.UpdateClientInfo{
					Username: "",
					Hostname: "",
					System:   "",
					Address:  "",
				},
			},
			wantMarshal: true,
		},
		{
			name: "Unsupported type",
			input: proto.Envelope{
				Subject: "",
				Payload: struct{}{},
			},
			wantErr: proto.ErrUnknownMessage,
		},
		{
			name: "Nil input",
			input: proto.Envelope{
				Subject: "",
				Payload: nil,
			},
			wantErr: proto.ErrUnknownMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Marshal
			marshaled, err := proto.Marshal(tt.input)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, marshaled)

			// Test Unmarshal
			unmarshaled, err := proto.Unmarshal(marshaled)
			require.NoError(t, err)
			require.Equal(t, tt.want, unmarshaled)
		})
	}
}

func TestUnmarshalInvalidData(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "Empty input",
			input:   []byte{},
			wantErr: true,
		},
		{
			name:    "Invalid data",
			input:   []byte("invalid data"),
			wantErr: true,
		},
		{
			name:    "Corrupt Cap'n Proto message",
			input:   []byte{0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: true,
		},
		{
			name:    "Nil input",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := proto.Unmarshal(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCmdStartWorkerSubject(t *testing.T) {
	got := proto.CmdStartWorkerSubject("agent1", "worker1")
	require.Equal(t, "FJ.AGENT.agent1.CMD.WORKER.worker1.START", got)
}

func TestCmdStopWorkerSubject(t *testing.T) {
	got := proto.CmdStopWorkerSubject("agent1", "worker1")
	require.Equal(t, "FJ.AGENT.agent1.CMD.WORKER.worker1.STOP", got)
}

func TestCmdWriteStdinSubject(t *testing.T) {
	got := proto.CmdWriteStdinSubject("agent1", "worker1")
	require.Equal(t, "FJ.AGENT.agent1.CMD.WORKER.worker1.STDIN", got)
}

func TestEvtStartWorkerSubject(t *testing.T) {
	got := proto.EvtStartWorkerSubject("agent1", "worker1")
	require.Equal(t, "FJ.AGENT.agent1.EVT.WORKER.worker1.START", got)
}

func TestEvtStopWorkerSubject(t *testing.T) {
	got := proto.EvtStopWorkerSubject("agent1", "worker1")
	require.Equal(t, "FJ.AGENT.agent1.EVT.WORKER.worker1.STOP", got)
}

func TestEvtWorkerStatusSubject(t *testing.T) {
	got := proto.EvtWorkerStatusSubject("agent1", "worker1")
	require.Equal(t, "FJ.AGENT.agent1.EVT.WORKER.worker1.STATUS", got)
}

func TestEvtWorkerStdoutSubject(t *testing.T) {
	got := proto.EvtWorkerStdoutSubject("agent1", "worker1")
	require.Equal(t, "FJ.AGENT.agent1.EVT.WORKER.worker1.STDOUT", got)
}

func TestEvtAgentInfoSubject(t *testing.T) {
	got := proto.EvtAgentInfoSubject("agent1")
	require.Equal(t, "FJ.AGENT.agent1.EVT.INFO", got)
}
