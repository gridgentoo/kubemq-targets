package logs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
	"sort"
)

type Client struct {
	log    *logger.Logger
	opts   options
	client *cloudwatchlogs.CloudWatchLogs
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

	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.opts.region),
		Credentials: credentials.NewStaticCredentials(c.opts.awsKey, c.opts.awsSecretKey, c.opts.token),
	})
	if err != nil {
		return err
	}

	svc := cloudwatchlogs.New(sess)
	c.client = svc

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "create_log_event_stream":
		return c.createLogEventStream(ctx, meta)
	case "describe_log_event_stream":
		return c.describeLogEventStream(ctx, meta)
	case "delete_log_event_stream":
		return c.deleteLogEventStream(ctx, meta)
	case "put_log_event":
		return c.putLogEvent(ctx, meta, req.Data)
	case "get_log_event":
		return c.getLogEvent(ctx, meta)
	case "create_log_group":
		return c.createLogEventGroup(ctx, meta, req.Data)
	case "describe_log_group":
		return c.describeLogGroup(ctx, meta)
	case "delete_log_group":
		return c.deleteLogEventGroup(ctx, meta)
	default:
		return nil, errors.New("invalid method type")
	}
}

func (c *Client) createLogEventStream(ctx context.Context, meta metadata) (*types.Response, error) {
	_, err := c.client.CreateLogStreamWithContext(ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(meta.logGroupName),
		LogStreamName: aws.String(meta.logStreamName),
	})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) describeLogEventStream(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.DescribeLogStreamsWithContext(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(meta.logGroupName),
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) deleteLogEventStream(ctx context.Context, meta metadata) (*types.Response, error) {
	_, err := c.client.DeleteLogStreamWithContext(ctx, &cloudwatchlogs.DeleteLogStreamInput{
		LogGroupName:  aws.String(meta.logGroupName),
		LogStreamName: aws.String(meta.logStreamName),
	})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) putLogEvent(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	var m map[int64]string
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, errors.New("failed to parse messages ,please verify data is map[int64]string int64:timestamp and string: message")
	}
	var inputLogs []*cloudwatchlogs.InputLogEvent
	for k, v := range m {
		i := cloudwatchlogs.InputLogEvent{
			Message:   aws.String(v),
			Timestamp: aws.Int64(k),
		}
		inputLogs = append(inputLogs, &i)
	}
	sort.Slice(inputLogs,
		func(i, j int) bool {
			return *inputLogs[i].Timestamp < *inputLogs[j].Timestamp
		})
	resp, err := c.client.PutLogEventsWithContext(ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(meta.logGroupName),
		LogStreamName: aws.String(meta.logStreamName),
		SequenceToken: aws.String(meta.sequenceToken),
		LogEvents:     inputLogs,
	})
	if err != nil {
		return nil, err
	}
	if resp.RejectedLogEventsInfo != nil {
		return nil, fmt.Errorf("%v", resp.RejectedLogEventsInfo)
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) getLogEvent(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.GetLogEventsWithContext(ctx, &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(meta.logGroupName),
		LogStreamName: aws.String(meta.logStreamName),
		Limit:         aws.Int64(meta.limit),
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) createLogEventGroup(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	m := make(map[string]*string)
	var err error
	if data != nil {
		err := json.Unmarshal(data, &m)
		if err != nil {
			return nil, err
		}
		_, err = c.client.CreateLogGroupWithContext(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(meta.logGroupName),
			Tags:         m,
		})
		if err != nil {
			return nil, err
		}
	} else {
		_, err = c.client.CreateLogGroupWithContext(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(meta.logGroupName),
		})
		if err != nil {
			return nil, err
		}
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) deleteLogEventGroup(ctx context.Context, meta metadata) (*types.Response, error) {
	_, err := c.client.DeleteLogGroupWithContext(ctx, &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(meta.logGroupName),
	})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) describeLogGroup(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.DescribeLogGroupsWithContext(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(meta.logGroupPrefix),
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) Stop() error {
	return nil
}
