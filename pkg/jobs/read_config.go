package jobs

import (
	"log"
	"path/filepath"

	"github.com/RobDoan/go-docker-template/pkg/utils"
)

type container struct {
	Name string   `yaml:"name"`
	Args []string `yaml:"args"`
}

type instance struct {
	Name       string      `yaml:"name"`
	Namespace  string      `yaml:"namespace"`
	Containers []container `yaml:"containers"`
}

type astroconfig struct {
	Instances []instance `yaml:"instances"`
}

const ProjectFolder = "payfare/dev"

func readAstroFile(dir string) (*[]Job, error) {
	astroFilePath := filepath.Join(dir, "conf", "astro-deploy-config.yaml")
	astroConfig := astroconfig{}
	if err := utils.ReadYamlFile(astroFilePath, &astroConfig); err != nil {
		log.Fatal("Unable to read astro-deploy-config.yaml file", err)
		return nil, err
	}

	jobs := []Job{}
	appendJobs := map[string]bool{}

	for _, instance := range astroConfig.Instances {

		for _, container := range instance.Containers {
			if len(container.Args) <= 0 {
				continue
			}
			if _, ok := appendJobs[instance.Name]; ok {
				continue
			}

			name := instance.Name
			appendJobs[name] = true
			job := Job{
				Name: name,
				Args: container.Args,
			}
			jobs = append(jobs, job)
		}
	}
	return &jobs, nil
}
