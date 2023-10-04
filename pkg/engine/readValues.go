package engine

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/RobDoan/go-docker-template/pkg/utils"
)

type Values struct {
	Image      string         `yaml:"image"`
	Ports      []string       `yaml:"ports"`
	Links      []string       `yaml:"links"`
	DebugPorts map[string]int `yaml:"debugPorts,omitempty"`
}

func ReadServiceValue(dir string, service string) error {
	valueFile := filepath.Join(dir, service, "values.yaml")
	if fi, err := os.Stat(valueFile); err != nil && fi.IsDir() {
		infoMessage := fmt.Sprintf("%s does not have values configuration file", service)
		utils.PrintInfo(infoMessage)
	}
	values := Values{}
	err := utils.ReadYamlFile(valueFile, &values)
	if err != nil {
		log.Fatal("Unable to read values.yaml file: ", err)
		return err
	}
	return nil
}
