package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func GetAbsDirectory(dir string) (string, error) {
	path, err := filepath.Abs(dir)
	if err != nil {
		return path, err
	}

	if fi, err := os.Stat(path); err != nil {
		return path, err
	} else if !fi.IsDir() {
		return path, errors.Errorf("no such directory %s", path)
	}

	return path, nil
}

func WriteFile(name string, content []byte) error {
	if _, err := os.Stat(name); err == nil {
		// TODO: truncate file name before asking for confirmation
		overrided := AskForConfirmation(fmt.Sprintf(" \"%s\" already exists. Do you want to overwrite it?", name))
		if !overrided {
			return nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		log.Printf("could not create directory %s", filepath.Dir(name))
		return err
	}

	return os.WriteFile(name, content, 0644)
}

func ReadYamlFile(yamlFile string, out interface{}) error {
	b, err := os.ReadFile(yamlFile)
	if err != nil {
		log.Fatal("Unable to read docker-compose.yml file", err)
		return err
	}

	if err := yaml.Unmarshal(b, out); err != nil {
		log.Fatal("Unable to unmarshal docker-compose file: ", err)
	}
	return nil
}
