package main

import (
	"flag"
	"os"

	"github.com/RobDoan/go-docker-template/pkg/engine"
	"github.com/RobDoan/go-docker-template/pkg/jobs"
)

type Config struct {
	Environment string
	ConfigFile  string
}

func configParse() Config {
	var configFile = flag.String("config", "config.json", "config file")
	var environment = flag.String("environment", "qa2", "environment")
	flag.Parse()
	config := Config{*environment, *configFile}
	return config
}


func main() {
	
}
