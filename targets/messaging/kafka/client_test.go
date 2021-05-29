package kafka

import (
	"context"
	b64 "encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
)

func replaceHeaderValues(req *types.Request) {
	r := strings.NewReplacer(
		"_replaceHK_", b64.StdEncoding.EncodeToString([]byte("header1")),
		"_replaceHV_", b64.StdEncoding.EncodeToString([]byte("headervalue1")))
	req.Metadata["Headers"] = r.Replace(req.Metadata["Headers"])

}

func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "messaging-kafka",
				Kind: "messaging.kafka",
				Properties: map[string]string{
					"brokers": "localhost:9092",
					"topic":   "TestTopic",
				},
			},
			wantErr: false,
		}, {
			name: "invalid init - missing brokers",
			cfg: config.Spec{
				Name: "messaging-kafka",
				Kind: "messaging.kafka",
				Properties: map[string]string{
					"topic": "TestTopic",
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing topic",
			cfg: config.Spec{
				Name: "messaging-kafka",
				Kind: "messaging.kafka",
				Properties: map[string]string{
					"brokers": "localhost:9092",
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
			name: "valid publish request ",
			cfg: config.Spec{
				Name: "messaging-kafka",
				Kind: "messaging.kafka",
				Properties: map[string]string{
					"brokers": "localhost:9092",
					"topic":   "TestTopic",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("Key", "S2V5").
				SetData([]byte("new-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("partition", "0").
				SetMetadataKeyValue("offset", "1"),
			wantErr: false,
		},
		{
			name: "valid publish request with headers",
			cfg: config.Spec{
				Name: "messaging-kafka",
				Kind: "messaging.kafka",
				Properties: map[string]string{
					"brokers": "localhost:9092",
					"topic":   "TestTopic",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("Key", "S2V5").
				SetData([]byte("new-data")).SetMetadataKeyValue(
				"Headers", `[{"Key": "_replaceHK_","Value": "_replaceHV_"}]`),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("partition", "0").
				SetMetadataKeyValue("offset", "2"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			replaceHeaderValues(tt.request)
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
