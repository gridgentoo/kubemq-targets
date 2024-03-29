package postgres

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/stretchr/testify/require"
)

type testStructure struct {
	instanceConnectionName string
	dbUser                 string
	dbPassword             string
	dbName                 string
	cred                   string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../../credentials/sql/postgresinstanceConnectionName.txt")
	if err != nil {
		return nil, err
	}
	t.instanceConnectionName = string(dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/sql/dbUser.txt")
	if err != nil {
		return nil, err
	}
	t.dbUser = string(dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/sql/dbPassword.txt")
	if err != nil {
		return nil, err
	}
	t.dbPassword = string(dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/sql/postGresDBName.txt")
	t.dbName = string(dat)
	if err != nil {
		return nil, err
	}
	dat, err = ioutil.ReadFile("./../../../../credentials/google_cred.json")
	if err != nil {
		return nil, err
	}
	t.cred = string(dat)
	return t, nil
}

type post struct {
	Id      int    `json:"id"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}
type posts []*post

func (p *posts) marshal() []byte {
	b, _ := json.Marshal(p)
	return b
}

func unmarshal(data []byte) *posts {
	if data == nil {
		return nil
	}
	p := &posts{}
	_ = json.Unmarshal(data, p)
	return p
}

var allPosts = posts{
	&post{
		Id:      1,
		Content: "Content One",
	},
	&post{
		Id:      2,
		Title:   "Title Two",
		Content: "Content Two",
	},
}

const (
	createPostTable = `
	DROP TABLE IF EXISTS post;
	       CREATE TABLE post (
	         ID serial,
	         TITLE varchar(40),
	         CONTENT varchar(255),
	         CONSTRAINT pk_post PRIMARY KEY(ID)
	       );
	       INSERT INTO post(ID,TITLE,CONTENT) VALUES
	                       (1,NULL,'Content One'),
	                       (2,'Title Two','Content Two');
	`
	selectPostTable = `SELECT id,title,content FROM post;`
)

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
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			wantErr: false,
		}, {
			name: "invalid init - missing db_user",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing db_name",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing connection",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"db_user":     dat.dbUser,
					"db_name":     dat.dbName,
					"db_password": dat.dbPassword,
					"credentials": dat.cred,
					"use_proxy":   "true",
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing db_password",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_name":                  dat.dbName,
					"db_user":                  dat.dbUser,
					"credentials":              dat.cred,
					"use_proxy":                "true",
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
			err = c.Init(ctx, tt.cfg, nil)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			err = c.Stop()
			require.NoError(t, err)
		})
	}
}

func TestClient_Query_Exec_Transaction(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name              string
		cfg               config.Spec
		execRequest       *types.Request
		queryRequest      *types.Request
		wantExecResponse  *types.Response
		wantQueryResponse *types.Response
		wantExecErr       bool
		wantQueryErr      bool
	}{
		{
			name: "valid exec query request",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetData([]byte(createPostTable)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetData([]byte(selectPostTable)),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok").
				SetData(allPosts.marshal()),
			wantExecErr:  false,
			wantQueryErr: false,
		},
		{
			name: "empty exec request",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec"),
			queryRequest:      nil,
			wantExecResponse:  nil,
			wantQueryResponse: nil,
			wantExecErr:       true,
			wantQueryErr:      false,
		},
		{
			name: "invalid exec request",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetData([]byte("bad statement")),
			queryRequest:      nil,
			wantExecResponse:  nil,
			wantQueryResponse: nil,
			wantExecErr:       true,
			wantQueryErr:      false,
		},
		{
			name: "valid exec empty query request",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetData([]byte(createPostTable)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetData([]byte("")),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: nil,
			wantExecErr:       false,
			wantQueryErr:      true,
		},
		{
			name: "valid exec bad query request",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetData([]byte(createPostTable)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetData([]byte("some bad query")),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: nil,
			wantExecErr:       false,
			wantQueryErr:      true,
		},
		{
			name: "valid exec valid query - no results",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetData([]byte(createPostTable)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetData([]byte("SELECT id,title,content FROM post where id=100")),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantExecErr:  false,
			wantQueryErr: false,
		},
		{
			name: "valid exec query request",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetData([]byte(createPostTable)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetData([]byte(selectPostTable)),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok").
				SetData(allPosts.marshal()),
			wantExecErr:  false,
			wantQueryErr: false,
		},
		{
			name: "empty transaction request",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "transaction"),
			queryRequest:      nil,
			wantExecResponse:  nil,
			wantQueryResponse: nil,
			wantExecErr:       true,
			wantQueryErr:      false,
		},
		{
			name: "invalid transaction request",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetData([]byte("bad statement")),
			queryRequest:      nil,
			wantExecResponse:  nil,
			wantQueryResponse: nil,
			wantExecErr:       true,
			wantQueryErr:      false,
		},
		{
			name: "valid transaction empty query request",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetData([]byte(createPostTable)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetData([]byte("")),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: nil,
			wantExecErr:       false,
			wantQueryErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			defer func() {
				err = c.Stop()
				require.NoError(t, err)
			}()
			gotSetResponse, err := c.Do(ctx, tt.execRequest)
			if tt.wantExecErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
			require.EqualValues(t, tt.wantExecResponse, gotSetResponse)
			gotGetResponse, err := c.Do(ctx, tt.queryRequest)
			if tt.wantQueryErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotGetResponse)

			if tt.wantQueryResponse != nil {
				wantPosts := unmarshal(tt.wantQueryResponse.Data)
				var gotPosts *posts
				if gotGetResponse != nil {
					gotPosts = unmarshal(gotGetResponse.Data)
				}
				require.EqualValues(t, wantPosts, gotPosts)
			} else {
				require.EqualValues(t, tt.wantQueryResponse, gotGetResponse)
			}
		})
	}
}

func TestClient_Do(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid request",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetMetadataKeyValue("isolation_level", "read_uncommitted").
				SetData([]byte(createPostTable)),
			wantErr: false,
		},
		{
			name: "valid request - 2",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetMetadataKeyValue("isolation_level", "read_committed").
				SetData([]byte(createPostTable)),
			wantErr: false,
		},
		{
			name: "valid request - 3",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetMetadataKeyValue("isolation_level", "repeatable_read").
				SetData([]byte(createPostTable)),
			wantErr: false,
		},
		{
			name: "valid request - 3",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetMetadataKeyValue("isolation_level", "serializable").
				SetData([]byte(createPostTable)),
			wantErr: false,
		},
		{
			name: "invalid request - bad method",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "bad-method"),
			wantErr: true,
		},
		{
			name: "invalid request - bad isolation level",
			cfg: config.Spec{
				Name: "gcp-stores-postgres",
				Kind: "gcp.stores.postgres",
				Properties: map[string]string{
					"instance_connection_name": dat.instanceConnectionName,
					"db_user":                  dat.dbUser,
					"db_name":                  dat.dbName,
					"db_password":              dat.dbPassword,
					"credentials":              dat.cred,
					"use_proxy":                "true",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetMetadataKeyValue("isolation_level", "bad_level").
				SetData([]byte(createPostTable)),
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
			defer func() {
				err = c.Stop()
				require.NoError(t, err)
			}()
			_, err = c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
