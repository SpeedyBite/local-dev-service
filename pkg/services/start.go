package services

import (
	"fmt"

	"github.com/RobDoan/go-docker-template/pkg/vault"
)

type StartOptions struct {
	LocalServices []string
	Environment   string
	VaultUrl      string
	Channel       string
	CountryCode   string
}

func NewStartOptions(services []string, environment string, vault string, channel string, countryCode string) *StartOptions {
	return &StartOptions{
		LocalServices: services,
		Environment:   environment,
		VaultUrl:      vault,
		Channel:       channel,
		CountryCode:   countryCode,
	}
}

func (s *StartOptions) Start() {
	fmt.Println("Start")
}

func (s *StartOptions) readVaults(service string) (map[string]interface{}, error) {
	vault := vault.NewVault(s.VaultUrl)
	secret, err := vault.GetVaults(s.Environment, service)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return secret.Data, nil
}

func (s *StartOptions) CreateServiceEnv(service string) {
	environmentData, err := s.readVaults(service)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, value := range environmentData {
		fmt.Printf("%s: %s \n", k, value)
	}
}
