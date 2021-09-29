package msk

import (
	"testing"

	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/stretchr/testify/require"
)

func TestMetadata_parseOptions(t *testing.T) {
	tests := []struct {
		name     string
		meta     config.Spec
		wantOpts options
		wantErr  bool
	}{
		{
			name: "valid options",
			meta: config.Spec{
				Name: "Kafka config",
				Kind: "kafka",
				Properties: map[string]string{
					"brokers": "localhost:9092,localhost:9093",
					"topic":   "TestTopic",
				},
			},
			wantOpts: options{
				brokers: []string{"localhost:9092", "localhost:9093"},
				topic:   "TestTopic",
			},
			wantErr: false,
		}, {
			name: "valid options  with userpass",
			meta: config.Spec{
				Name: "Kafka options conf",
				Kind: "kafka",
				Properties: map[string]string{
					"brokers":      "localhost:9092,localhost:9093",
					"topic":        "TestTopic",
					"saslUsername": "admin",
					"saslPassword": "password",
				},
			},
			wantOpts: options{
				brokers:      []string{"localhost:9092", "localhost:9093"},
				topic:        "TestTopic",
				saslUsername: "admin",
				saslPassword: "password",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOpts, err := parseOptions(tt.meta)
			if tt.wantErr {
				require.Error(t, err)

			} else {
				require.NoError(t, err)

			}
			require.EqualValues(t, tt.wantOpts, gotOpts)
		})
	}
}
