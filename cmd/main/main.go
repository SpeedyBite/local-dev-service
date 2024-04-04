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
	// command := "docker-compose"
	// args := []string{"run", "-it", "--rm", "container", "/bin/bash"}

	// // Create a new command and attach it to the standard input/output
	// cmd := exec.Command(command, args...)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// // Run the command
	// if err := cmd.Run(); err != nil {
	// 	fmt.Println("Error running command:", err)
	// }
	// config := configParse()
	jobs.GetJob("gibson")
	return
	startServiceOption := engine.NewBuildOptions([]string{"bishop"}, "qa2",
		"https://vault11.raven.k8s.payfare.com", "paidapp_us", "us", os.Args[1])
	startServiceOption.Build()

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
	// engine := engine.NewEngine()
	// engine.Render()
}
