package services

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RobDoan/go-docker-template/pkg/utils"
	"github.com/judedaryl/go-arrayutils"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

var (
	supportedYamlExtensions = []string{".yml", ".yaml"}
)

const (
	dest                          = "docker-services"
	ServiceDefinedFileNames       = "services.yaml"
	ValuesfileName                = "values.yaml"
	HelperfileName                = "helpers.tpl"
	ServiceTemplateFileName       = "service.yaml.tpl"
	ServiceDockerTemplateFileName = "docker-compose.yaml.tpl"
	ServiceConnectionsFileName    = "service-connections.yaml"
)

const defaultServiceHelperTemplate = `
{{- define "gibson.environment" -}}
{{- range $k, $v := .Environment }}
{{ print "- " $k "=" $v | indent 6}}
{{- end -}}
{{- end -}}
`

const defaultServiceTemplate = `
version: "3.7"
services:
  {{ .Name }}:
    enviroment:
      {{- template "gibson.environment" .}}
    ports:
		{{- range $_, $v := .DebugPorts }}
      - {{ print $v ":" $v | quote }}
		{{- end }}
`
const couldNotWriteFileMessage = "could not write file:"

func couldNotWriteFile(err error) {
	log.Fatal(couldNotWriteFileMessage, err)
	utils.PrintError(couldNotWriteFileMessage)
}

type ServiceConfig struct {
	Image       string   `yaml:"image,omitempty"`
	Environment []string `yaml:"environment,omitempty"` // Environments variables
	Ports       []string `yaml:"ports,omitempty"`
	Volumes     []string `yaml:"volumes,omitempty"`
	Networks    []string `yaml:"networks,omitempty"`
	Links       []string `yaml:"links,omitempty"`
}

type DockerComposeConfig struct {
	Services       map[string]ServiceConfig `yaml:"services"`
	Version        string                   `yaml:"version"`
	dockerFilePath string
	directory      string
}

func LoadServiceConfig(dir string) (*DockerComposeConfig, error) {
	path, err := utils.GetAbsDirectory(dir)
	dockerComposeConfig := DockerComposeConfig{directory: path}
	if err != nil {
		log.Fatal("Unable to get absolute path: ", err)
		return nil, err
	}

	yamlFile, err := findDockerComposeFile(path)
	if err != nil {
		log.Fatal("Unable to find docker-compose.yml file: ", err)
		return nil, err
	}
	dockerComposeConfig.dockerFilePath = yamlFile

	b, err := os.ReadFile(yamlFile)
	if err != nil {
		log.Fatal("Unable to read docker-compose.yml file", err)
		return nil, err
	}

	if err := yaml.Unmarshal(b, &dockerComposeConfig); err != nil {
		log.Fatal("Unable to unmarshal docker-compose file: ", err)
	}
	return &dockerComposeConfig, nil
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

func (dc *DockerComposeConfig) createLocalEnv() (string, error) {
	path := dc.directory
	for serviceName, service := range dc.Services {
		environmentFilePath := filepath.Join(path, dest, serviceName, "local.env")

		content := strings.Join(service.Environment, "\n")

		if err := utils.WriteFile(environmentFilePath, []byte(content)); err != nil {
			couldNotWriteFile(err)
			continue
		}
		helpFilePath := filepath.Join(path, dest, serviceName, HelperfileName)
		helpContent := transform(defaultServiceHelperTemplate, serviceName)
		if err := utils.WriteFile(helpFilePath, []byte(helpContent)); err != nil {
			couldNotWriteFile(err)
			continue
		}
		serviceFilePath := filepath.Join(path, dest, serviceName, ServiceTemplateFileName)
		serviceTemplateContent := transform(defaultServiceTemplate, serviceName)
		if err := utils.WriteFile(serviceFilePath, []byte(serviceTemplateContent)); err != nil {
			couldNotWriteFile(err)
			continue
		}
	}
	return path, nil
}

func (dc *DockerComposeConfig) CreateValues() error {
	path := dc.directory
	for serviceName, service := range dc.Services {
		serviceValueFilePath := filepath.Join(path, dest, serviceName, ValuesfileName)
		serviceValues := CreateSericeValues(&service)

		if (serviceValues) == nil {
			log.Fatal("could not create service values: ", serviceName)
			return errors.Errorf("could not create service values: %s", serviceName)
		}

		if err := serviceValues.DumpYamlTo(serviceValueFilePath); (err) != nil {
			log.Fatal("could not marshal service values: ", err)
			return err
		}
	}
	return nil
}

func (dc *DockerComposeConfig) createServiceConnections() error {
	serviceConnection := make(map[string]interface{})
	for serviceName := range dc.Services {
		// Skip non-payfare services such as mysql, redis, rabbitmq
		if arrayutils.Some(StorageServices, func(service string) bool { return service == serviceName }) {
			continue
		}
		for _, environment := range Environments {
			serviceDomains, ok := serviceConnection[serviceName].(map[string]string)
			if !ok {
				serviceConnection[serviceName] = make(map[string]string)
				serviceDomains = serviceConnection[serviceName].(map[string]string)
			}
			serviceDomains[environment] = BuildServiceDomain(environment, serviceName)
			serviceConnection[serviceName] = serviceDomains
		}
	}
	path := dc.directory
	serviceConnectionsFilePath := filepath.Join(path, dest, ServiceConnectionsFileName)

	content, err := yaml.Marshal(serviceConnection)
	if (err) != nil {
		log.Fatal("could not marshal service connections: ", err)
		return err
	}

	return utils.WriteFile(serviceConnectionsFilePath, content)
}

func transform(src, replacement string) []byte {
	return []byte(strings.ReplaceAll(src, "<SERVICENAME>", replacement))
}

// Read docker-compose.yml file and create all local environment files
func Create(dir string) (string, error) {
	dockerComposeConfig, err := LoadServiceConfig(dir)
	if err != nil {
		return "", err
	}

	if err := dockerComposeConfig.CreateValues(); err != nil {
		return "", err
	}

	if err := dockerComposeConfig.createServiceConnections(); err != nil {
		return "", err
	}
	return dockerComposeConfig.createLocalEnv()
}
