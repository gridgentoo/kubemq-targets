package sqs

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.sqs").
		SetDescription("AWS SQS Target").
		SetName("SQS").
		SetProvider("AWS").
		SetCategory("Messaging").
		SetTags("queue", "cloud", "managed").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set SQS aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set SQS aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set SQS aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set SQS token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_receive").
				SetDescription("Set SQS max receive").
				SetMust(false).
				SetDefault("0").
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("retries").
				SetDescription("Set SQS number of retries on failed send request").
				SetMust(false).
				SetDefault("0").
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("default_delay").
				SetDescription("Set SQS default delay in seconds").
				SetMust(false).
				SetDefault("10").
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("dead_letter").
				SetTitle("Dead Letter Queue").
				SetDescription("Set SQS dead letter queue").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("default_queue").
				SetTitle("Default Queue").
				SetDescription("Set SQS Default Queue").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("queue").
				SetKind("string").
				SetDescription("Set EventHubs queue name").
				SetDefault("").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("tags").
				SetKind("string").
				SetDescription("Set EventHubs tags").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("partition_key").
				SetKind("string").
				SetDescription("Set EventHubs partition key").
				SetDefault("").
				SetMust(false),
		)
}
