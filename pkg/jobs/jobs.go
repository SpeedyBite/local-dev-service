package jobs

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/RobDoan/go-docker-template/pkg/tui"
	"github.com/RobDoan/go-docker-template/pkg/utils"
	"github.com/pkg/errors"
)

type Job struct {
	Name       string
	Args       []string
	Names      []string
	Namespaces []string
}

func (j Job) command() string {
	return strings.Join(j.Args, " ")
}

func (j Job) Title() string {
	return j.Name
}

func (j Job) Description() string {
	return j.command()
}

func SelectAJob(serviceName string) (Job, error) {
	homedir, err := utils.CurrentUserHomeDir()
	if err != nil {
		log.Fatal(err)
		return Job{}, err
	}
	appDir := filepath.Join(homedir, ProjectFolder, serviceName)
	if err != nil {
		utils.PrintError("Unable to get absolute path")
		log.Fatal(err)
		return Job{}, err
	}
	jobs, err := readAstroFile(appDir)

	if err != nil {
		utils.PrintError("Unable to read astro file")
		log.Fatal(err)
		return Job{}, err
	}

	jobSelectOptions := make([]tui.Option, len(*jobs))
	i := 0
	for _, job := range *jobs {
		jobSelectOptions[i] = job
		i++
	}

	selectedIndex, err := tui.SelectList(jobSelectOptions, "Select a job to run")

	if err != nil {
		log.Fatal(err)
		return Job{}, err
	}

	if selectedIndex < 0 || selectedIndex+1 > len(jobSelectOptions) {
		log.Fatal(errors.New("No job to selected or selected job is invalid"))
		return Job{}, errors.New("No job to selected or selected job is invalid")
	}

	utils.PrintInfo(fmt.Sprintf("Selected job: %s", jobSelectOptions[selectedIndex].Title()))

	return (*jobs)[selectedIndex], nil
}
