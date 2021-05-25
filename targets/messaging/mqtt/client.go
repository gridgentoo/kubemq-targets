package mqtt

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"time"
)

const (
	defaultConnectTimeout = 5 * time.Second
)

type Client struct {
	name   string
	opts   options
	client mqtt.Client
}

func New() *Client {
	return &Client{}
}
func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", c.opts.host))
	opts.SetUsername(c.opts.username)
	opts.SetPassword(c.opts.password)
	opts.SetClientID(c.opts.clientId)
	opts.SetConnectTimeout(defaultConnectTimeout)
	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error connecting to mqtt broker, %w", token.Error())
	}

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, ok := c.opts.defaultMetadata()
	if !ok {
		var err error
		meta, err = parseMetadata(req.Metadata)
		if err != nil {
			return nil, err
		}
	}
	token := c.client.Publish(meta.topic, byte(meta.qos), false, req.Data)
	token.Wait()
	if token.Error() != nil {
		return nil, token.Error()
	}
	return types.NewResponse().SetMetadataKeyValue("result", "ok"), nil
}


func (c *Client) Stop() error {
	if c.client != nil {
		c.client.Disconnect(0)
	}
	return nil
}
