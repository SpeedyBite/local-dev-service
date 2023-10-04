package utils

import (
	"os"
	"path/filepath"
)

func ListSubfolders(rootFolder string) ([]string, error) {
	subfolders := []string{}

	err := filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != rootFolder {
			subfolderName := filepath.Base(path)
			subfolders = append(subfolders, subfolderName)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return subfolders, nil
}
