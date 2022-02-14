package lambda

import (
	"context"
	"encoding/json"

	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"

	"testing"
	"time"
)

type testStructure struct {
	awsKey       string
	awsSecretKey string
	region       string
	token        string

	zipFileName  string
	functionName string
	handlerName  string
	role         string
	runtime      string
	description  string

	lambdaExp []byte
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/aws/awsKey.txt")
	if err != nil {
		return nil, err
	}
	t.awsKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/awsSecretKey.txt")
	if err != nil {
		return nil, err
	}
	t.awsSecretKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/region.txt")
	if err != nil {
		return nil, err
	}
	t.region = string(dat)
	t.token = ""

	dat, err = ioutil.ReadFile("./../../../credentials/aws/lambda/zipFileName.txt")
	if err != nil {
		return nil, err
	}
	t.zipFileName = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/lambda/functionName.txt")
	if err != nil {
		return nil, err
	}
	t.functionName = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/lambda/handlerName.txt")
	if err != nil {
		return nil, err
	}
	t.handlerName = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/lambda/role.txt")
	if err != nil {
		return nil, err
	}
	t.role = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/lambda/runtime.txt")
	if err != nil {
		return nil, err
	}
	t.runtime = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/lambda/description.txt")
	if err != nil {
		return nil, err
	}
	t.description = string(dat)
	contents, err := ioutil.ReadFile("./../../../credentials/aws/lambda/lambdaCode.zip")
	if err != nil {
		return nil, err
	}
	t.lambdaExp = contents
	return t, nil
}

func TestClient_Init(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "aws-lambda",
				Kind: "aws.lambda",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: false,
		}, {
			name: "invalid init - missing aws_key",
			cfg: config.Spec{
				Name: "aws-lambda",
				Kind: "aws.lambda",
				Properties: map[string]string{
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing region",
			cfg: config.Spec{
				Name: "aws-lambda",
				Kind: "aws.lambda",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing aws_secret_key",
			cfg: config.Spec{
				Name: "aws-lambda",
				Kind: "aws.lambda",
				Properties: map[string]string{
					"aws_key": dat.awsKey,
					"region":  dat.region,
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

			err := c.Init(ctx, tt.cfg, nil)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestClient_List(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-lambda",
		Kind: "aws.lambda",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()

	err = c.Init(ctx, cfg, nil)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid list",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestClient_Create(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-lambda",
		Kind: "aws.lambda",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()

	err = c.Init(ctx, cfg, nil)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid create",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("zip_file_name", dat.zipFileName).
				SetMetadataKeyValue("description", dat.description).
				SetMetadataKeyValue("handler_name", dat.handlerName).
				SetMetadataKeyValue("memorySize", "256").
				SetMetadataKeyValue("timeout", "15").
				SetMetadataKeyValue("role", dat.role).
				SetMetadataKeyValue("function_name", dat.functionName).
				SetMetadataKeyValue("runtime", dat.runtime).
				SetData(dat.lambdaExp),
			wantErr: false,
		},
		{
			name: "invalid create- already exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("zip_file_name", dat.zipFileName).
				SetMetadataKeyValue("description", dat.description).
				SetMetadataKeyValue("handler_name", dat.handlerName).
				SetMetadataKeyValue("role", dat.role).
				SetMetadataKeyValue("function_name", dat.functionName).
				SetMetadataKeyValue("runtime", dat.runtime).
				SetData(dat.lambdaExp),
			wantErr: true,
		},
		{
			name: "invalid create- missing data",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("zip_file_name", dat.zipFileName).
				SetMetadataKeyValue("description", dat.description).
				SetMetadataKeyValue("handler_name", dat.handlerName).
				SetMetadataKeyValue("role", dat.role).
				SetMetadataKeyValue("function_name", dat.functionName).
				SetMetadataKeyValue("runtime", dat.runtime),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestClient_Delete(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-lambda",
		Kind: "aws.lambda",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()

	err = c.Init(ctx, cfg, nil)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid delete",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("function_name", dat.functionName),
			wantErr: false,
		},
		{
			name: "invalid delete - does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("function_name", dat.functionName),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestClient_Run(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-lambda",
		Kind: "aws.lambda",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()

	err = c.Init(ctx, cfg, nil)
	require.NoError(t, err)
	b, err := json.Marshal("my object")
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid run",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "run").
				SetMetadataKeyValue("function_name", dat.functionName).
				SetData(b),
			wantErr: false,
		},
		{
			name: "invalid run - function does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "run").
				SetMetadataKeyValue("function_name", "not_a_real_function").
				SetData(b),
			wantErr: true,
		},
		{
			name: "invalid run - missing function name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "run").
				SetData(b),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestClient_isJson(t *testing.T) {

	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "empty",
			data: nil,
			want: true,
		},
		{
			name: "not-valid",
			data: []byte("bXkgb2JqZWN0"),
			want: false,
		},
		{
			name: "not-valid",
			data: []byte("eyJ0ZXN0IjoidGVzdCJ9"),
			want: false,
		},
		{
			name: "valid",
			data: []byte("eyJ0ZXN0IjogInRlc3QifQ=="),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := isJson(tt.data); got != tt.want {
				t.Errorf("isJson() = %v, want %v", got, tt.want)
			}
		})
	}
}
