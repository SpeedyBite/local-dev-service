package service

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RobDoan/go-docker-template/pkg/utils"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

var (
	supportedYamlExtensions = []string{".yml", ".yaml"}
)

const (
	dest = "docker-services"
)

type ServiceConfig struct {
	Image       string   `yaml:"image"`
	Environment []string `yaml:"environment,omitempty"`
	Ports       []string `yaml:"ports,omitempty"`
	Volumes     []string `yaml:"volumes,omitempty"`
	Networks    []string `yaml:"networks,omitempty"`
	Links       []string `yaml:"links,omitempty"`
}

type DockerComposeConfig struct {
	Services map[string]ServiceConfig `yaml:"services"`
	Version  string                   `yaml:"version"`
}

// Find docker-compose.ya?ml file in path
// if not found, return error
func findDockerComposeFile(path string) (string, error) {

	for _, extension := range supportedYamlExtensions {
		dockerComposeFile := filepath.Join(path, "docker-compose"+extension)
		if _, err := os.Stat(dockerComposeFile); err == nil {
			return dockerComposeFile, nil
		}
	}
	return "", errors.Errorf("no docker-compose file found in %s", path)
}

func writeFile(name string, content []byte) error {
	if _, err := os.Stat(name); err == nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	return os.WriteFile(name, content, 0644)
}

// Read docker-compose.yml file and create all local environment files
func CreateLocalEnv(dir string) (string, error) {

	path, err := utils.GetAbsDirectory(dir)

	if err != nil {
		log.Fatal("Unable to get absolute path: ", err)
		return path, err
	}

	yamlFile, err := findDockerComposeFile(path)
	if err != nil {
		log.Fatal("Unable to find docker-compose.yml file: ", err)
		return path, err
	}

	b, err := os.ReadFile(yamlFile)
	if err != nil {
		log.Fatal("Unable to read docker-compose.yml file", err)
		return path, err
	}

	dockerComposeConfig := DockerComposeConfig{}
	if err := yaml.Unmarshal(b, &dockerComposeConfig); err != nil {
		log.Fatal("Unable to unmarshal docker-compose file: ", err)
	}

	for serviceName, service := range dockerComposeConfig.Services {
		environmentFilePath := filepath.Join(path, dest, ".env", serviceName, "local.env")

		content := strings.Join(service.Environment, "\n")

		if err := writeFile(environmentFilePath, []byte(content)); err != nil {
			log.Fatal("count not write file: ", err)
			continue
		}
	}
	return path, nil
}
