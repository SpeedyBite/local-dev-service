package engine

import (
	"text/template"

	"github.com/RobDoan/go-docker-template/pkg/service"
)

type Engine struct {
	DockerComposeDir string
	Template         *template.Template
}

func NewEngine() Engine {
	t := template.New("docker-compose.yml")
	return Engine{
		Template: t,
	}
}

func Render(service service.Service) error {
	return nil
}
