package service

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

type Dependencies map[string]string

type ServiceValues struct {
	Image            string `yaml:"image"`
	Ports            []string
	links            []string
	debugPorts       map[string]int          `yaml:"debugPorts,omitempty"`
	dependencies     map[string]Dependencies `yaml:"dependencies,omitempty"`
	requiredDatabase bool                    `yaml:"requiredDatabase,omitempty"`
	requiredRedis    bool                    `yaml:"requiredRedis,omitempty"`
	requiredRabbitMQ bool                    `yaml:"requiredRabbitMQ,omitempty"`
}

func CreateSericeValues(image string, ports []string, links []string) *ServiceValues {
	serviceValues := ServiceValues{
		Image: image,
		Ports: ports,
		links: links,
	}
	serviceValues.debugPorts = serviceValues.getDebugPorts()
	serviceValues.dependencies = serviceValues.getDependencies()
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
	for _, portMapping := range s.Ports {
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

func (s *ServiceValues) getDependencies() map[string]Dependencies {
	ignoreLinks := []string{"mysql", "redis", "rabbitmq"}
	config := make(map[string]Dependencies)
	for _, environment := range Environments {
		dependencies := Dependencies{}
		for _, link := range s.links {
			isIgnored := arrayutils.Some(ignoreLinks, func(ignoreLink string) bool { return ignoreLink == link })
			if !isIgnored {
				dependencies[link] = BuildServiceDomain(environment, link)
			}
		}
	}
	return config
}

func (s *ServiceValues) isRequiredDatabase() bool {
	return arrayutils.Some(s.links, func(link string) bool { return link == "mysql" })
}

func (s *ServiceValues) isRequiredRedis() bool {
	return arrayutils.Some(s.links, func(link string) bool { return link == "redis" })
}

func (s *ServiceValues) isRequiredRabbitMQ() bool {
	return arrayutils.Some(s.links, func(link string) bool { return link == "rabbitmq" })
}
