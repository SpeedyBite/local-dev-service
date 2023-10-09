package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/RobDoan/go-docker-template/pkg/jobs"
)

func main() {
	// Get first command line argument as a directory or default to current directory
	port := flag.Int("port", 33120, "debug port")
	service := flag.String("service", "", "Service")
	flag.Parse()
	serviceName := *service
	if serviceName == "" {
		dir, _ := os.Getwd()
		absDir, _ := filepath.Abs(dir)
		serviceName = filepath.Base(absDir)
	}

	log.Printf("Service name: %s %d \n", serviceName, *port)
	job, err := jobs.SelectAJob(serviceName)
	if err != nil {
		panic(err)
	}

	jobs.Run(job, serviceName, *port)
}
