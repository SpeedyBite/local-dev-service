package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/RobDoan/go-docker-template/pkg/utils"
	"github.com/judedaryl/go-arrayutils"
	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"
)

var (
	supportPortMappingRegex = regexp.MustCompile(`\d+:\d+`)
)

type ServiceValues struct {
	Descritpion         string            `yaml:"description"`
	Links               []string          `yaml:"links"`
	Command             interface{}       `yaml:"command,omitempty"`
	DebugPorts          map[string]int    `yaml:"debugPorts,omitempty"`
	LinkEnvironmentName map[string]string `yaml:"linkEnvironmentName,omitempty"`
	ports               []string
}

func CreateSericeValues(sc *ServiceConfig) *ServiceValues {
	serviceValues := ServiceValues{
		Links:       sc.Links,
		Descritpion: "Service Description",
		ports:       sc.Ports,
	}
	serviceValues.DebugPorts = serviceValues.getDebugPorts()
	serviceValues.LinkEnvironmentName = serviceValues.getConnectionEnvironmentName(sc.Environment)
	// serviceValues.requiredDatabase = serviceValues.isRequiredDatabase()
	// serviceValues.requiredRedis = serviceValues.isRequiredRedis()
	// serviceValues.requiredRabbitMQ = serviceValues.isRequiredRabbitMQ()
	return &serviceValues
}

func (s *ServiceValues) getConnectionEnvironmentName(environmentVariables []string) map[string]string {
	linkEvnName := map[string]string{}
	for _, serviceName := range s.Links {
		isStorageService := arrayutils.Some(StorageServices, func(service string) bool { return service == serviceName })

		// Skipp all storage service
		if isStorageService {
			continue
		}
		domainPattern := fmt.Sprintf("^(https?://)?%s", serviceName)
		re := regexp.MustCompile(domainPattern)
		for _, envString := range environmentVariables {
			if !strings.Contains(envString, "=") {
				utils.PrintError(fmt.Sprintf("[%s] Invalid environment variable %s", serviceName, envString))
				continue
			}
			parts := strings.SplitN(envString, "=", 2)
			if (len(parts) != 2) || (parts[0] == "") || (parts[1] == "") {
				utils.PrintError(fmt.Sprintf("[%s] Invalid environment variable %s", serviceName, envString))
			}
			envName := strings.Trim(parts[0], " ")
			envValue := strings.Trim(parts[1], " ")
			if re.Match([]byte(envValue)) {
				linkEvnName[serviceName] = envName
			}
		}
	}
	return linkEvnName
}

func (s *ServiceValues) DumpYamlTo(filepath string) error {
	data, err := yaml.Marshal(s)
	if err != nil {
		log.Fatal("Unable to marshal service values: ", err)
		return err
	}
	return utils.WriteFile(filepath, data)
}

// Get the first host post which is bigger than 10000
func (s *ServiceValues) getAppProxyPort() (int, error) {
	for _, portMapping := range s.ports {
		if !supportPortMappingRegex.MatchString(portMapping) {
			log.Printf("Port mapping %s is not supported", portMapping)
			continue
		}
		hostPort, err := strconv.Atoi(strings.Split(portMapping, ":")[0])
		if err != nil {
			return 0, err
		}
		if hostPort > 10000 {
			return hostPort, nil
		}
	}
	return 0, errors.New("no port mapping found")
}

func (s *ServiceValues) getDebugPorts() map[string]int {
	config := make(map[string]int)
	hostPort, err := s.getAppProxyPort()
	if err != nil {
		return config
	}
	config["appDebugPort"] = hostPort + 10000
	config["jobDebugPort"] = hostPort + 20000
	return config
}

// func (s *ServiceValues) isRequiredDatabase() bool {
// 	return arrayutils.Some(s.Links, func(link string) bool { return link == "mysql" })
// }

// func (s *ServiceValues) isRequiredRedis() bool {
// 	return arrayutils.Some(s.Links, func(link string) bool { return link == "redis" })
// }

// func (s *ServiceValues) isRequiredRabbitMQ() bool {
// 	return arrayutils.Some(s.Links, func(link string) bool { return link == "rabbitmq" })
// }

func ReadServiceValue(dir string, service string) (ServiceValues, error) {
	valueFile := filepath.Join(dir, service, "values.yaml")
	utils.PrintInfo(fmt.Sprintf("Reading values.yaml file: %s", valueFile))
	if fi, err := os.Stat(valueFile); err != nil && fi.IsDir() {
		infoMessage := fmt.Sprintf("%s does not have values configuration file", service)
		utils.PrintInfo(infoMessage)
	}
	values := ServiceValues{}
	err := utils.ReadYamlFile(valueFile, &values)
	if err != nil {
		log.Fatal("Unable to read values.yaml", err)
		return ServiceValues{}, err
	}
	return values, nil
}
