package mqtt

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/uuid"
)

type options struct {
	host     string
	username string
	password string
	clientId string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, err = cfg.Properties.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.username = cfg.Properties.ParseString("username", "")
	o.password = cfg.Properties.ParseString("password", "")
	o.clientId = cfg.Properties.ParseString("client_id", uuid.New().String())
	return o, nil
}
