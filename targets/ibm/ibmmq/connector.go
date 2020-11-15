package ibmmq

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("ibm.ibmmq").
		SetDescription("IBM-MQ Messaging Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("queue_manager_name").
				SetDescription("Set IBM-MQ queue manager name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host_name").
				SetDescription("Set IBM-MQ host name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("channel_name").
				SetDescription("Set IBM-MQ channel name the queue is under").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set IBM-MQ username").
				SetDefault("").
				SetMust(true),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("certificate_label").
				SetDescription("Set IBM-MQ certificate_label for requests").
				SetDefault("").
				SetMust(false),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("ttl").
				SetDescription("Sets IBM-MQ message time to live (milliseconds)").
				SetDefault("1000000").
				SetMax(1000000000).
				SetMin(0).
				SetMust(false),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("queue_name").
				SetDescription("Sets IBM-MQ queue name").
				SetDefault("").
				SetMust(true),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("transport_type").
				SetDescription("Set IBM-MQ Transport type").
				SetDefault("0").
				SetMax(1).
				SetMin(0).
				SetMust(false),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("tls_client_auth").
				SetDescription("Set IBM-MQ tls_client_auth").
				SetDefault("NONE").
				SetMust(false),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("port_number").
				SetDescription("Set IBM-MQ server port_number").
				SetDefault("1414").
				SetMax(1000000000).
				SetMin(0).
				SetMust(false),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set IBM-MQ password").
				SetDefault("").
				SetMust(false),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("key_repository").
				SetDescription("Set IBM-MQ key_repository a certificate store").
				SetDefault("").
				SetMust(false),
		)
}
