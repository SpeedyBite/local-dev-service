package serviceconfig

import (
	"github.com/RobDoan/go-docker-template/pkg/utils"
)

type ServiceConfig struct {
	Name          string
	DomaninFormat string
	Channel       string
	CountryCode   string
}

func New(name string, domainFormat string) *ServiceConfig {
	return &ServiceConfig{
		Name:          name,
		DomaninFormat: domainFormat,
	}
}

func (s *ServiceConfig) GetDomain(env string) (string, error) {
	return utils.SafeSprintf(s.DomaninFormat, env)
}
