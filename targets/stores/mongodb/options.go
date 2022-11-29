package mongodb

import (
	"fmt"
	"math"
	"time"

	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	host             string
	username         string
	password         string
	database         string
	collection       string
	writeConcurrency string
	readConcurrency  string
	params           string
	operationTimeout time.Duration
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
	o.database, err = cfg.Properties.MustParseString("database")
	if err != nil {
		return options{}, fmt.Errorf("error parsing database name, %w", err)
	}
	o.collection, err = cfg.Properties.MustParseString("collection")
	if err != nil {
		return options{}, fmt.Errorf("error parsing collection name, %w", err)
	}
	o.writeConcurrency = cfg.Properties.ParseString("write_concurrency", "")
	o.readConcurrency = cfg.Properties.ParseString("read_concurrency", "")

	o.params = cfg.Properties.ParseString("params", "")
	operationTimeoutSeconds, err := cfg.Properties.ParseIntWithRange("operation_timeout_seconds", 60, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error operation timeout seconds, %w", err)
	}
	o.operationTimeout = time.Duration(operationTimeoutSeconds) * time.Second
	return o, nil
}
