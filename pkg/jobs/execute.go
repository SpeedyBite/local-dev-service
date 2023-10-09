package jobs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type runningContainer struct {
	ContainerName string `json:"Name"`
	Project       string `json:"Project"`
	Service       string `json:"Service"`
	Status        string `json:"State"`
}

func findContainerName(service string) (string, error) {
	command := exec.Command("docker", "compose", "ps", "--format", "json", service)
	output, err := command.Output()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	var data []runningContainer

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(output, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return "", err
	}
	if (len(data)) <= 0 {
		return "", errors.New("Unable to find container name")
	}
	var container runningContainer
	for _, c := range data {
		if c.Status == "running" {
			container = c
			break
		}
	}

	if container.ContainerName == "" {
		log.Fatal("Unable to find running container")
		return "", errors.New("Unable to find running container")
	}

	return container.ContainerName, nil
}

func Run(job Job, service string, port int) error {
	containerName, err := findContainerName(service)
	if err != nil {
		return err
	}

	commands := []string{
		"cd /srv/root",
		"pip install --upgrade debugpy",
		fmt.Sprintf("python3 -m debugpy --wait-for-client --listen 0.0.0.0:%d %s", port, job.command()),
	}
	cmdStr := strings.Join(commands, " && ")
	cmd := exec.Command("docker", "exec", "-it", containerName, "bash", "-c", cmdStr)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running command")
		log.Fatal(err)
		return err
	}

	return nil
}
