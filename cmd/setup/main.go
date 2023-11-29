package main

import (
  "github.com/bugfixes/go-bugfixes/logs"
  "github.com/joho/godotenv"
  "github.com/k8sdeploy/queue-manager-service/internal/master"
  ConfigBuilder "github.com/keloran/go-config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		_ = logs.Errorf("unable to load env: %v", err)
		return
	}

	cfg, err := ConfigBuilder.Build(ConfigBuilder.Rabbit)
	if err != nil {
		_ = logs.Errorf("unable to build config: %v", err)
		return
	}

	ad, err := master.NewMaster(*cfg).Build()
	if err != nil {
		_ = logs.Errorf("unable to build master: %v", err)
		return
	}

	logs.Local().Infof("Username: %s, Password: %s, VHost: %s", ad.Username, ad.Password, ad.VHost)
}
