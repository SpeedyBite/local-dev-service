package services

import "fmt"

var (
	Environments = []string{"local", "uat", "stg",
		"qa", "qa2", "qa3",
		"sandbox", "sandbox", "sandbox2", "sandbox3",
	}
	StorageServices = []string{"redis", "mysql", "rabbitmq", "elasticsearch"}
)

func BuildServiceDomain(env string, service string) string {
	if env == "local" {
		return service
	}
	return fmt.Sprintf("%s-%s.payfare.com", env, service)
}
