package bigtable

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	projectID   string
	instance    string
	credentials string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.projectID, err = cfg.Properties.MustParseString("project_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project_id, %w", err)
	}
	o.instance, err = cfg.Properties.MustParseString("instance")
	if err != nil {
		return options{}, fmt.Errorf("error parsing instance, %w", err)
	}
	o.credentials, err = cfg.Properties.MustParseString("credentials")
	if err != nil {
		return options{}, err
	}
	return o, nil
}
