package cloudfunctions

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"strings"

	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	gf "github.com/kubemq-io/kubemq-targets/targets/gcp/cloudfunctions/functions/apiv1"
	"github.com/kubemq-io/kubemq-targets/types"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	functionspb "google.golang.org/genproto/googleapis/cloud/functions/v1"
)

type Client struct {
	log            *logger.Logger
	opts           options
	client         *gf.CloudFunctionsClient
	parrantProject string
	list           []string
	//nameFunctions  map[string]string
}

func New() *Client {
	return &Client{}

}
func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) Init(ctx context.Context, cfg config.Spec, log *logger.Logger) error {
	c.log = log
	if c.log == nil {
		c.log = logger.NewLogger(cfg.Kind)
	}

	c.log = logger.NewLogger(cfg.Name)
	var err error

	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}

	b := []byte(c.opts.credentials)

	client, err := gf.NewCloudFunctionsClient(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return err
	}
	c.client = client
	c.parrantProject = c.opts.parentProject

	if c.opts.locationMatch {
		it := client.ListFunctions(ctx, &functionspb.ListFunctionsRequest{
			Parent: fmt.Sprintf("projects/%s/locations/-", c.opts.parentProject),
		})

		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			if resp != nil {
				c.list = append(c.list, resp.GetName())
			}
		}
	}
	return nil
}

func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	m, err := parseMetadata(request.Metadata, c.opts)
	if err != nil {
		return nil, err
	}

	if m.project == "" {
		m.project = c.parrantProject
	}

	name := fmt.Sprintf("projects/%s/locations/%s/functions/%s", m.project, m.location, m.name)
	if m.location == "" {
		for _, n := range c.list {
			if strings.Contains(n, m.name) && strings.Contains(n, m.project) {
				m.location = "added from match"
				name = n
				break
			}
		}
	}
	if m.location == "" {
		return nil, fmt.Errorf("no location found for function")
	}

	cfo := &functionspb.CallFunctionRequest{
		Name: name,
		Data: string(request.Data),
	}

	res, err := c.client.CallFunction(ctx, cfo)
	if err != nil {
		return nil, err
	}
	if res.Error != "" {
		return nil, fmt.Errorf(res.Error)
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", res.Result).
		SetMetadataKeyValue("execution_id", res.ExecutionId).
		SetData([]byte(res.Result)), nil

}

func (c *Client) Stop() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
