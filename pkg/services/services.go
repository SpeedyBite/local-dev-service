package services

import "github.com/RobDoan/go-docker-template/pkg/utils"

func ListServices(dir string) ([]string, error) {
	subfolders, err := utils.ListSubfolders(dir)
	if err != nil {
		return nil, err
	}
	return subfolders, nil
}
