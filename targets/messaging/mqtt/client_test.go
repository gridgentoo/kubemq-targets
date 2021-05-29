package mqtt

import (
	"context"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "messaging.mqtt",
				Kind: "messaging.mqtt",
				Properties: map[string]string{
					"host":      "localhost:1883",
					"username":  "",
					"password":  "",
					"client_id": "",
				},
			},
			wantErr: false,
		},
		{
			name: "init - bad host",
			cfg: config.Spec{
				Name: "messaging.mqtt",
				Kind: "messaging.mqtt",
				Properties: map[string]string{
					"host":      "localhost:6000",
					"username":  "",
					"password":  "",
					"client_id": "",
				},
			},
			wantErr: true,
		}, {
			name: "init - no host",
			cfg: config.Spec{
				Name: "messaging.mqtt",
				Kind: "messaging.mqtt",
				Properties: map[string]string{
					"username":  "",
					"password":  "",
					"client_id": "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			if err := c.Init(ctx, tt.cfg, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantExecErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config.Spec
		request      *types.Request
		wantResponse *types.Response
		wantErr      bool
	}{
		{
			name: "valid publish request with confirmation",
			cfg: config.Spec{
				Name: "messaging.mqtt",
				Kind: "messaging.mqtt",
				Properties: map[string]string{
					"host":      "localhost:1883",
					"username":  "",
					"password":  "",
					"client_id": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("topic", "some-queue").
				SetMetadataKeyValue("qos", "0").
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantErr: false,
		},
		{
			name: "invalid publish request - no topic",
			cfg: config.Spec{
				Name: "messaging.mqtt",
				Kind: "messaging.mqtt",
				Properties: map[string]string{
					"host":      "localhost:1883",
					"username":  "",
					"password":  "",
					"client_id": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("qos", "0").
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantErr: true,
		},
		{
			name: "invalid publish request - bad qos",
			cfg: config.Spec{
				Name: "messaging.mqtt",
				Kind: "messaging.mqtt",
				Properties: map[string]string{
					"host":      "localhost:1883",
					"username":  "",
					"password":  "",
					"client_id": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("topic", "some-queue").
				SetMetadataKeyValue("qos", "-1").
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			gotResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotResponse)
			require.EqualValues(t, tt.wantResponse, gotResponse)
		})
	}
}
