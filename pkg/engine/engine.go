package engine

import (
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/RobDoan/go-docker-template/pkg/utils"
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

type rendereredService struct {
	Name        string
	YamlContent string
}

func (e Engine) Render() error {
	localServices := []string{"mysql", "rabbitmq"}
	dir := "/Users/quytrandoan/payfare/dev/local-dev-services/testdata/docker-services"
	referenceTpls := []rendereredService{}
	for _, s := range localServices {
		templateFile := filepath.Join(dir, s, "service.yaml.tpl")
		t := template.New(templateFile)
		t, err := t.ParseFiles(templateFile)
		if err != nil {
			log.Fatal("could not parse template file: ", err)
			return err
		}
		log.Println(t.Name())
		var str strings.Builder
		t.ExecuteTemplate(&str, "service.yaml.tpl", make(map[string]interface{}))
		referenceTpls = append(referenceTpls, rendereredService{
			Name:        s,
			YamlContent: str.String(),
		})
	}
	dockerComposeTemplate := template.New("docker-compose.yml")
	dockerComposeTemplate, err := dockerComposeTemplate.Funcs(sprig.FuncMap()).ParseFiles("/Users/quytrandoan/payfare/dev/local-dev-services/testdata/docker-services/docker-compose.yaml.tpl")
	if err != nil {
		log.Fatal("could not parse template file: ", err)
		return err
	}
	var str strings.Builder
	dockerComposeTemplate.ExecuteTemplate(&str, "docker-compose.yaml.tpl", referenceTpls)
	final := strings.ReplaceAll(str.String(), "<no value>", "")

	return utils.WriteFile(filepath.Join(dir, "docker-compose.yaml"), []byte(final))
}
