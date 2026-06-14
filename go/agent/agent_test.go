package agent_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/foohq/foojank-proto/go/agent"
)

func TestMarshalUnmarshal(t *testing.T) {
	testError := errors.New("test error")

	tests := []struct {
		name        string
		input       agent.Envelope
		want        agent.Envelope
		wantMarshal bool
		wantErr     error
	}{
		{
			name: "StartWorkerRequest",
			input: agent.Envelope{
				Subject: agent.CmdStartWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StartWorkerRequest{
					Command: "cmd",
					Args:    []string{"arg1", "arg2"},
					Env:     []string{"KEY1=val1", "KEY2=val2"},
				},
			},
			want: agent.Envelope{
				Subject: agent.CmdStartWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StartWorkerRequest{
					Command: "cmd",
					Args:    []string{"arg1", "arg2"},
					Env:     []string{"KEY1=val1", "KEY2=val2"},
				},
			},
			wantMarshal: true,
		},
		{
			name: "StartWorkerRequest with empty slices",
			input: agent.Envelope{
				Subject: agent.CmdStartWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StartWorkerRequest{
					Command: "cmd",
					Args:    []string{},
					Env:     []string{},
				},
			},
			want: agent.Envelope{
				Subject: agent.CmdStartWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StartWorkerRequest{
					Command: "cmd",
					Args:    nil,
					Env:     nil,
				},
			},
			wantMarshal: true,
		},
		{
			name: "StartWorkerResponse without error",
			input: agent.Envelope{
				Subject: agent.EvtStartWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StartWorkerResponse{
					Error: nil,
				},
			},
			want: agent.Envelope{
				Subject: agent.EvtStartWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StartWorkerResponse{
					Error: nil,
				},
			},
			wantMarshal: true,
		},
		{
			name: "StartWorkerResponse with error",
			input: agent.Envelope{
				Subject: agent.EvtStartWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StartWorkerResponse{
					Error: testError,
				},
			},
			want: agent.Envelope{
				Subject: agent.EvtStartWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StartWorkerResponse{
					Error: testError,
				},
			},
			wantMarshal: true,
		},
		{
			name: "StopWorkerRequest",
			input: agent.Envelope{
				Subject: agent.CmdStopWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StopWorkerRequest{},
			},
			want: agent.Envelope{
				Subject: agent.CmdStopWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StopWorkerRequest{},
			},
			wantMarshal: true,
		},
		{
			name: "StopWorkerResponse without error",
			input: agent.Envelope{
				Subject: agent.EvtStopWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StopWorkerResponse{
					Error: nil,
				},
			},
			want: agent.Envelope{
				Subject: agent.EvtStopWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StopWorkerResponse{
					Error: nil,
				},
			},
			wantMarshal: true,
		},
		{
			name: "StopWorkerResponse with error",
			input: agent.Envelope{
				Subject: agent.EvtStopWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StopWorkerResponse{
					Error: testError,
				},
			},
			want: agent.Envelope{
				Subject: agent.EvtStopWorkerSubject("gateway1", "agent1", "worker1"),
				Payload: agent.StopWorkerResponse{
					Error: testError,
				},
			},
			wantMarshal: true,
		},
		{
			name: "UpdateWorkerStatus",
			input: agent.Envelope{
				Subject: agent.EvtWorkerStatusSubject("gateway1", "agent1", "worker1"),
				Payload: agent.UpdateWorkerStatus{
					Status: 42,
				},
			},
			want: agent.Envelope{
				Subject: agent.EvtWorkerStatusSubject("gateway1", "agent1", "worker1"),
				Payload: agent.UpdateWorkerStatus{
					Status: 42,
				},
			},
			wantMarshal: true,
		},
		{
			name: "UpdateWorkerStdio",
			input: agent.Envelope{
				Subject: agent.EvtWorkerStdoutSubject("gateway1", "agent1", "worker1"),
				Payload: agent.UpdateWorkerStdio{
					Data: []byte("Hello, World!"),
				},
			},
			want: agent.Envelope{
				Subject: agent.EvtWorkerStdoutSubject("gateway1", "agent1", "worker1"),
				Payload: agent.UpdateWorkerStdio{
					Data: []byte("Hello, World!"),
				},
			},
			wantMarshal: true,
		},
		{
			name: "UpdateClientInfo",
			input: agent.Envelope{
				Subject: agent.EvtAgentInfoSubject("gateway1", "agent1"),
				Payload: agent.UpdateClientInfo{
					Username: "testuser",
					Hostname: "testhost",
					System:   "linux",
					Address:  "192.168.1.1",
				},
			},
			want: agent.Envelope{
				Subject: agent.EvtAgentInfoSubject("gateway1", "agent1"),
				Payload: agent.UpdateClientInfo{
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
			input: agent.Envelope{
				Subject: agent.EvtAgentInfoSubject("gateway1", "agent1"),
				Payload: agent.UpdateClientInfo{
					Username: "",
					Hostname: "",
					System:   "",
					Address:  "",
				},
			},
			want: agent.Envelope{
				Subject: agent.EvtAgentInfoSubject("gateway1", "agent1"),
				Payload: agent.UpdateClientInfo{
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
			input: agent.Envelope{
				Subject: "",
				Payload: struct{}{},
			},
			wantErr: agent.ErrUnknownMessage,
		},
		{
			name: "Nil input",
			input: agent.Envelope{
				Subject: "",
				Payload: nil,
			},
			wantErr: agent.ErrUnknownMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Marshal
			marshaled, err := agent.Marshal(tt.input)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, marshaled)

			// Test Unmarshal
			unmarshaled, err := agent.Unmarshal(marshaled)
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
			_, err := agent.Unmarshal(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCmdStartWorkerSubject(t *testing.T) {
	got := agent.CmdStartWorkerSubject("gateway1", "agent1", "worker1")
	require.Equal(t, "FJ.GATEWAY.gateway1.AGENT.agent1.CMD.WORKER.worker1.START", got)
}

func TestCmdStopWorkerSubject(t *testing.T) {
	got := agent.CmdStopWorkerSubject("gateway1", "agent1", "worker1")
	require.Equal(t, "FJ.GATEWAY.gateway1.AGENT.agent1.CMD.WORKER.worker1.STOP", got)
}

func TestCmdWriteStdinSubject(t *testing.T) {
	got := agent.CmdWriteStdinSubject("gateway1", "agent1", "worker1")
	require.Equal(t, "FJ.GATEWAY.gateway1.AGENT.agent1.CMD.WORKER.worker1.STDIN", got)
}

func TestEvtStartWorkerSubject(t *testing.T) {
	got := agent.EvtStartWorkerSubject("gateway1", "agent1", "worker1")
	require.Equal(t, "FJ.GATEWAY.gateway1.AGENT.agent1.EVT.WORKER.worker1.START", got)
}

func TestEvtStopWorkerSubject(t *testing.T) {
	got := agent.EvtStopWorkerSubject("gateway1", "agent1", "worker1")
	require.Equal(t, "FJ.GATEWAY.gateway1.AGENT.agent1.EVT.WORKER.worker1.STOP", got)
}

func TestEvtWorkerStatusSubject(t *testing.T) {
	got := agent.EvtWorkerStatusSubject("gateway1", "agent1", "worker1")
	require.Equal(t, "FJ.GATEWAY.gateway1.AGENT.agent1.EVT.WORKER.worker1.STATUS", got)
}

func TestEvtWorkerStdoutSubject(t *testing.T) {
	got := agent.EvtWorkerStdoutSubject("gateway1", "agent1", "worker1")
	require.Equal(t, "FJ.GATEWAY.gateway1.AGENT.agent1.EVT.WORKER.worker1.STDOUT", got)
}

func TestEvtAgentInfoSubject(t *testing.T) {
	got := agent.EvtAgentInfoSubject("gateway1", "agent1")
	require.Equal(t, "FJ.GATEWAY.gateway1.AGENT.agent1.EVT.INFO", got)
}
