package service

import (
	"fmt"
)

var ()

type Service struct {
	Name  string
	Env   string
	Ports []string
}

func CreateService(service string, env string, dir string) {
	fmt.Println("CreateService")
}
