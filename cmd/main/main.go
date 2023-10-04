package main

import (
	"flag"

	"github.com/RobDoan/go-docker-template/pkg/engine"
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

// func readVault() {
// 	vault := vault.NewVault("https://vault11.raven.k8s.payfare.com")
// 	secret, err := vault.GetVaults()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(secret)
// 	// fmt.Printf("hello world %s %s \n", config.Environment, config.ConfigFile)
// }

//	func createEnvFile() {
//		dir := os.Args[1]
//		fmt.Println(dir)
//	}
func main() {
	// config := configParse()
	// startServiceOption := services.NewStartOptions([]string{"bishop"}, "qa2",
	// 	"https://vault11.raven.k8s.payfare.com", "paidapp_us", "us")
	// startServiceOption.CreateServiceEnv("bishop")
	// file := filepath.Join(os.Args[1], "docker-compose.yaml")
	// b, err := os.ReadFile(file)
	// if err != nil {
	// 	log.Fatal("Unable to read docker-compose.yml file", err)
	// }
	// aDict := yaml.MapSlice{}

	// if err := yaml.Unmarshal(b, &aDict); err != nil {
	// 	log.Fatal("Unable to unmarshal docker-compose file: ", err)
	// }
	// testYaml, err := yaml.Marshal(aDict["services"].(map[string]interface{})["mysql"])
	// utils.WriteFile(filepath.Join(os.Args[1], "test.yaml"), testYaml)

	// subfolders, _ := utils.ListSubfolders(os.Args[1])
	// log.Printf("subfolders: %v", subfolders)
	// services.Create(os.Args[1])
	engine := engine.NewEngine()
	engine.Render()
}
