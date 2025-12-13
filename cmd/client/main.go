package main

import (
	"go-gophkeeper/internal/client/cli"
	"go-gophkeeper/internal/config"
	"go-gophkeeper/internal/logger"
	"go-gophkeeper/internal/utils/building"
)

var buildVersion, buildDate, buildCommit string

func main() {
	building.PrintBuildVersion(buildVersion, buildDate, buildCommit)

	cfg, err := config.NewAgentConfig()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	client, err := cli.NewCli(cfg)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	err = client.Run()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

}
