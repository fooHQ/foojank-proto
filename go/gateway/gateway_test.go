package gateway_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/foohq/foojank-proto/go/gateway"
)

func TestMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		input       gateway.Envelope
		want        gateway.Envelope
		wantMarshal bool
		wantErr     error
	}{
		{
			name: "RegisterAgentRequest with properties",
			input: gateway.Envelope{
				Subject: gateway.CmdRegisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.RegisterAgentRequest{
					Properties: []gateway.Property{
						{Key: "key1", Value: "val1"},
						{Key: "key2", Value: "val2"},
					},
				},
			},
			want: gateway.Envelope{
				Subject: gateway.CmdRegisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.RegisterAgentRequest{
					Properties: []gateway.Property{
						{Key: "key1", Value: "val1"},
						{Key: "key2", Value: "val2"},
					},
				},
			},
			wantMarshal: true,
		},
		{
			name: "RegisterAgentRequest with empty properties",
			input: gateway.Envelope{
				Subject: gateway.CmdRegisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.RegisterAgentRequest{
					Properties: []gateway.Property{},
				},
			},
			want: gateway.Envelope{
				Subject: gateway.CmdRegisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.RegisterAgentRequest{
					Properties: []gateway.Property{},
				},
			},
			wantMarshal: true,
		},
		{
			name: "RegisterAgentRequest with nil properties",
			input: gateway.Envelope{
				Subject: gateway.CmdRegisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.RegisterAgentRequest{
					Properties: nil,
				},
			},
			want: gateway.Envelope{
				Subject: gateway.CmdRegisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.RegisterAgentRequest{
					Properties: []gateway.Property{},
				},
			},
			wantMarshal: true,
		},
		{
			name: "RegisterAgentResponse without error",
			input: gateway.Envelope{
				Subject: gateway.EvtRegisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.RegisterAgentResponse{
					Properties: []gateway.Property{
						{Key: "key1", Value: "val1"},
					},
					Error: nil,
				},
			},
			want: gateway.Envelope{
				Subject: gateway.EvtRegisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.RegisterAgentResponse{
					Properties: []gateway.Property{
						{Key: "key1", Value: "val1"},
					},
					Error: nil,
				},
			},
			wantMarshal: true,
		},
		{
			name: "RegisterAgentResponse with error",
			input: gateway.Envelope{
				Subject: gateway.EvtRegisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.RegisterAgentResponse{
					Properties: []gateway.Property{},
					Error:      errors.New("registration failed"),
				},
			},
			want: gateway.Envelope{
				Subject: gateway.EvtRegisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.RegisterAgentResponse{
					Properties: []gateway.Property{},
					Error:      errors.New("registration failed"),
				},
			},
			wantMarshal: true,
		},
		{
			name: "UnregisterAgentRequest with properties",
			input: gateway.Envelope{
				Subject: gateway.CmdUnregisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.UnregisterAgentRequest{
					Properties: []gateway.Property{
						{Key: "reason", Value: "shutdown"},
					},
				},
			},
			want: gateway.Envelope{
				Subject: gateway.CmdUnregisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.UnregisterAgentRequest{
					Properties: []gateway.Property{
						{Key: "reason", Value: "shutdown"},
					},
				},
			},
			wantMarshal: true,
		},
		{
			name: "UnregisterAgentRequest with nil properties",
			input: gateway.Envelope{
				Subject: gateway.CmdUnregisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.UnregisterAgentRequest{
					Properties: nil,
				},
			},
			want: gateway.Envelope{
				Subject: gateway.CmdUnregisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.UnregisterAgentRequest{
					Properties: []gateway.Property{},
				},
			},
			wantMarshal: true,
		},
		{
			name: "UnregisterAgentResponse without error",
			input: gateway.Envelope{
				Subject: gateway.EvtUnregisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.UnregisterAgentResponse{
					Properties: []gateway.Property{
						{Key: "key1", Value: "val1"},
					},
					Error: nil,
				},
			},
			want: gateway.Envelope{
				Subject: gateway.EvtUnregisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.UnregisterAgentResponse{
					Properties: []gateway.Property{
						{Key: "key1", Value: "val1"},
					},
					Error: nil,
				},
			},
			wantMarshal: true,
		},
		{
			name: "UnregisterAgentResponse with error",
			input: gateway.Envelope{
				Subject: gateway.EvtUnregisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.UnregisterAgentResponse{
					Properties: []gateway.Property{},
					Error:      errors.New("internal error"),
				},
			},
			want: gateway.Envelope{
				Subject: gateway.EvtUnregisterAgentSubject("gateway1", "agent1"),
				Payload: gateway.UnregisterAgentResponse{
					Properties: []gateway.Property{},
					Error:      errors.New("internal error"),
				},
			},
			wantMarshal: true,
		},
		{
			name: "Unsupported type",
			input: gateway.Envelope{
				Subject: "",
				Payload: struct{}{},
			},
			wantErr: gateway.ErrUnknownMessage,
		},
		{
			name: "Nil input",
			input: gateway.Envelope{
				Subject: "",
				Payload: nil,
			},
			wantErr: gateway.ErrUnknownMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshaled, err := gateway.Marshal(tt.input)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, marshaled)

			unmarshaled, err := gateway.Unmarshal(marshaled)
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
			_, err := gateway.Unmarshal(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCmdRegisterAgentSubject(t *testing.T) {
	got := gateway.CmdRegisterAgentSubject("gateway1", "agent1")
	require.Equal(t, "FJ.GATEWAY.gateway1.CMD.AGENT.agent1.REGISTER", got)
}

func TestCmdUnregisterAgentSubject(t *testing.T) {
	got := gateway.CmdUnregisterAgentSubject("gateway1", "agent1")
	require.Equal(t, "FJ.GATEWAY.gateway1.CMD.AGENT.agent1.UNREGISTER", got)
}

func TestEvtRegisterAgentSubject(t *testing.T) {
	got := gateway.EvtRegisterAgentSubject("gateway1", "agent1")
	require.Equal(t, "FJ.GATEWAY.gateway1.EVT.AGENT.agent1.REGISTER", got)
}

func TestEvtUnregisterAgentSubject(t *testing.T) {
	got := gateway.EvtUnregisterAgentSubject("gateway1", "agent1")
	require.Equal(t, "FJ.GATEWAY.gateway1.EVT.AGENT.agent1.UNREGISTER", got)
}
