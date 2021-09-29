// +build !container

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/api"
	"github.com/kubemq-io/kubemq-targets/binding"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/browser"
	"github.com/kubemq-io/kubemq-targets/pkg/builder"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/sources"
	"github.com/kubemq-io/kubemq-targets/targets"
	"io/ioutil"

	"os"
	"os/signal"
	"syscall"
)

var (
	version = ""
)

var (
	log              *logger.Logger
	generateManifest = flag.Bool("manifest", false, "generate targets connectors manifest")
	build            = flag.Bool("build", false, "build targets configuration")
	buildUrl         = flag.String("get", "", "get config file from url")
	configFile       = flag.String("config", "config.yaml", "set config file name")
	svcFlag          = flag.String("service", "", "control the kubemq-targets service")
	svcUsername      = flag.String("username", "", "kubemq-targets service username")
	svcPassword      = flag.String("password", "", "kubemq-targets service password")
)

func saveManifest() error {
	sourceConnectors := sources.Connectors()
	if err := sourceConnectors.Validate(); err != nil {
		return err
	}
	targetConnectors := targets.Connectors()
	if err := targetConnectors.Validate(); err != nil {
		return err
	}
	return common.NewManifest().
		SetSchema("targets").
		SetVersion(version).
		SetSourceConnectors(sourceConnectors).
		SetTargetConnectors(targetConnectors).
		Save()
}
func downloadUrl() error {
	c, err := builder.GetBuildManifest(*buildUrl)
	if err != nil {
		return err
	}
	cfg := &config.Config{}
	err = yaml.Unmarshal([]byte(c.Spec.Config), &cfg)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.yaml", data, 0644)
	if err != nil {
		return err
	}
	return nil
}
func runInteractive(serviceExit chan bool) error {
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)
	configCh := make(chan *config.Config)
	cfg, err := config.Load(configCh)
	if err != nil {
		return err
	}
	err = cfg.Validate()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bindingsService, err := binding.New()
	if err != nil {
		return err
	}
	err = bindingsService.Start(ctx, cfg)
	if err != nil {
		return err
	}

	apiServer, err := api.Start(ctx, cfg.ApiPort, bindingsService)
	if err != nil {
		return err
	}
	for {
		select {
		case newConfig := <-configCh:
			err = newConfig.Validate()
			if err != nil {
				return fmt.Errorf("error on validation new config file: %s", err.Error())
			}
			bindingsService.Stop()
			err = bindingsService.Start(ctx, newConfig)
			if err != nil {
				return fmt.Errorf("error on restarting service with new config file: %s", err.Error())
			}
			if apiServer != nil {
				err = apiServer.Stop()
				if err != nil {
					return fmt.Errorf("error on shutdown api server: %s", err.Error())
				}
			}

			apiServer, err = api.Start(ctx, newConfig.ApiPort, bindingsService)
			if err != nil {
				return fmt.Errorf("error on start api server: %s", err.Error())
			}
		case <-gracefulShutdown:
			_ = apiServer.Stop()
			bindingsService.Stop()
			return nil
		case <-serviceExit:
			_ = apiServer.Stop()
			bindingsService.Stop()
			return nil
		}
	}
}
func preRun() {
	if *generateManifest {
		err := saveManifest()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		log.Infof("generated manifest.json completed")
		os.Exit(0)
	}
	if *build {
		err := browser.OpenURL("https://build.kubemq.io/#/targets")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}
	if *buildUrl != "" {
		err := downloadUrl()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
}

func main() {
	log = logger.NewLogger("kubemq-targets")
	flag.Parse()
	config.SetConfigFile(*configFile)
	app := newAppService()
	if err := app.init(*svcFlag, *svcUsername, *svcPassword); err != nil {
		log.Error(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
