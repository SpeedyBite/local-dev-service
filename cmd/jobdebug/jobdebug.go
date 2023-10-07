package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RobDoan/go-docker-template/pkg/jobs"
)

type Config struct {
	Environment string
	ConfigFile  string
}

func main() {
	// Get first command line argument as a directory or default to current directory

	var dir string
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		dir, _ = os.Getwd()
	}
	absDir, err := filepath.Abs(dir)

	if err != nil {
		panic("Invalid directory: " + dir)
	}

	serviceName := filepath.Base(absDir)
	port := flag.String("port", "33120", "debug port")

	job, err := jobs.SelectAJob(serviceName)
	if err != nil {
		panic(err)
	}

	jobs.Run(job, serviceName, *port)

	fmt.Printf("Service name: %s %s \n", serviceName, *port)
}
