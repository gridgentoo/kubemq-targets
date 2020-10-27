package openfaas

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.serverless.openfaas").
		SetDescription("Openfaas Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("gateway").
				SetDescription("Set Openfaas gateway address").
				SetMust(true).
				SetDefault("localhost:27017"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set Openfaas username").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set Openfaas password").
				SetMust(true).
				SetDefault(""),
		)
}
