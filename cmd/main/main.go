package main

import (
	"flag"
	"os"

	"github.com/RobDoan/go-docker-template/pkg/service"
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

//	func readVault() {
//		vault := vault.NewVault("https://vault11.raven.k8s.payfare.com")
//		secret, err := vault.GetVaults()
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println(secret)
//		// fmt.Printf("hello world %s %s \n", config.Environment, config.ConfigFile)
//	}
//
//	func createEnvFile() {
//		dir := os.Args[1]
//		fmt.Println(dir)
//	}
func main() {
	// config := configParse()
	service.CreateLocalEnv(os.Args[1])
}
