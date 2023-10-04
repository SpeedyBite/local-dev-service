// package service

// import (
// 	"os"
// 	"path/filepath"

// 	"github.com/pkg/errors"
// )

// type Service struct {
// 	Name      string `yaml:"name"`
// 	Image     string
// 	Port      int
// 	IsGateWay bool
// }

// const defaultHelpers = `{{/*
//   Expand helper for templates
//   */}}
//   {{- define "environmentFiles" -}}
//   {{- range $key, $value := .Values.environmentFiles }}
//   - {{ $value }}
//   {{- end }}
//   {{- end }}
// `
// const ServiceTemplate = `
// {{service.Name}}:
//   image: {{service.Image}}
//   ports:
//     - {{service.Port}}:{{service.Port}}
//   environmentFile:
//     - {{service.EnvironmentFile}}
//   environment:
//     - SERVICE_NAME={{service.Name}}
//   volumes:
//     - {{service.Volume}}:{{service.Volume}}
//   networks:
//     - {{service.Network}}
// `

// func Create(name string, path string) (string, error) {
// 	// This function is currently empty and needs to be implemented.
// 	// It is intended to create a new service with the given name and path.
// 	// The function should return the path of the created service and any errors encountered.
// 	path, err := filepath.Abs(path)
// 	if err != nil {
// 		return path, err
// 	}

// 	if fi, err := os.Stat(path); err != nil {
// 		return path, err
// 	} else if !fi.IsDir() {
// 		return path, errors.Errorf("no such directory %s", path)
// 	}

// 	serviceFile := filepath.Join(path, name)
// 	if fi, err := os.Stat(serviceFile); err == nil && fi.IsDir() {
// 		return serviceFile, errors.Errorf("file %s already exists and is a directory", serviceFile)
// 	}

// 	return "", nil
// }


