package main

import (
  "github.com/bugfixes/go-bugfixes/logs"
  "github.com/k8sdeploy/queue-manager-service/internal/service"
  ConfigBuilder "github.com/keloran/go-config"
)

var (
	BuildVersion = "dev"
	BuildHash    = "none"
	ServiceName  = "queue-manager-service"
)

func main() {
	logs.Local().Infof("Starting %s", ServiceName)
	logs.Local().Infof("Version: %s, Hash: %s", BuildVersion, BuildHash)

	cfg, err := ConfigBuilder.Build(ConfigBuilder.Rabbit)
	if err != nil {
		_ = logs.Errorf("unable to build config: %v", err)
		return
	}

	if err := service.NewService(*cfg).Start(); err != nil {
		_ = logs.Errorf("unable to start service: %v", err)
		return
	}
}
