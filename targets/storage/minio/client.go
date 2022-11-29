package minio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	log      *logger.Logger
	opts     options
	s3Client *minio.Client
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
	c.s3Client, err = minio.New(c.opts.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.opts.accessKeyId, c.opts.secretAccessKey, ""),
		Secure: c.opts.useSSL,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "make_bucket":
		return c.MakeBucket(ctx, meta)
	case "list_buckets":
		return c.ListBuckets(ctx)
	case "bucket_exists":
		return c.BucketExist(ctx, meta)
	case "remove_bucket":
		return c.RemoveBucket(ctx, meta)
	case "list_objects":
		return c.ListObjects(ctx, meta)
	case "put":
		return c.Put(ctx, meta, req.Data)
	case "get":
		return c.Get(ctx, meta)
	case "remove":
		return c.Remove(ctx, meta)
	}
	return nil, nil
}

func (c *Client) MakeBucket(ctx context.Context, meta metadata) (*types.Response, error) {
	bucketOptions := minio.MakeBucketOptions{
		Region:        meta.param2,
		ObjectLocking: false,
	}
	err := c.s3Client.MakeBucket(ctx, meta.param1, bucketOptions)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) ListBuckets(ctx context.Context) (*types.Response, error) {
	buckets, err := c.s3Client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(&buckets)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok").
		SetData(data), nil
}

func (c *Client) BucketExist(ctx context.Context, meta metadata) (*types.Response, error) {
	found, err := c.s3Client.BucketExists(ctx, meta.param1)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("exist", fmt.Sprintf("%t", found)).
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) RemoveBucket(ctx context.Context, meta metadata) (*types.Response, error) {
	err := c.s3Client.RemoveBucket(ctx, meta.param1)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) ListObjects(ctx context.Context, meta metadata) (*types.Response, error) {
	var objects []minio.ObjectInfo
	for object := range c.s3Client.ListObjects(ctx, meta.param1, minio.ListObjectsOptions{Recursive: true}) {
		objects = append(objects, object)
	}
	data, err := json.Marshal(&objects)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok").
		SetData(data), nil
}

func (c *Client) Get(ctx context.Context, meta metadata) (*types.Response, error) {
	object, err := c.s3Client.GetObject(ctx, meta.param1, meta.param2, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = object.Close()
	}()
	data, err := ioutil.ReadAll(object)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok").
		SetData(data), nil
}

func (c *Client) Put(ctx context.Context, meta metadata, value []byte) (*types.Response, error) {
	r := bytes.NewReader(value)
	_, err := c.s3Client.PutObject(ctx, meta.param1, meta.param2, r, int64(r.Len()), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) Remove(ctx context.Context, meta metadata) (*types.Response, error) {
	err := c.s3Client.RemoveObject(ctx, meta.param1, meta.param2, minio.RemoveObjectOptions{
		GovernanceBypass: false,
		VersionID:        "",
		Internal:         minio.AdvancedRemoveOptions{},
	})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) Stop() error {
	return nil
}
