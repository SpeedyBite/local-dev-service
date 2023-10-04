package services

import (
	"log"
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
	Descritpion      string   `yaml:"description"`
	Links            []string `yaml:"links"`
	ports            []string
	DebugPorts       map[string]int `yaml:"debugPorts,omitempty"`
	requiredDatabase bool           `yaml:"requiredDatabase"`
	requiredRedis    bool           `yaml:"requiredRedis"`
	requiredRabbitMQ bool           `yaml:"requiredRabbitMQ"`
}

func CreateSericeValues(sc *ServiceConfig) *ServiceValues {
	serviceValues := ServiceValues{
		Links:       sc.Links,
		Descritpion: "Service Description",
		ports:       sc.Ports,
	}
	serviceValues.DebugPorts = serviceValues.getDebugPorts()
	serviceValues.requiredDatabase = serviceValues.isRequiredDatabase()
	serviceValues.requiredRedis = serviceValues.isRequiredRedis()
	serviceValues.requiredRabbitMQ = serviceValues.isRequiredRabbitMQ()
	return &serviceValues
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

func (s *ServiceValues) isRequiredDatabase() bool {
	return arrayutils.Some(s.Links, func(link string) bool { return link == "mysql" })
}

func (s *ServiceValues) isRequiredRedis() bool {
	return arrayutils.Some(s.Links, func(link string) bool { return link == "redis" })
}

func (s *ServiceValues) isRequiredRabbitMQ() bool {
	return arrayutils.Some(s.Links, func(link string) bool { return link == "rabbitmq" })
}
