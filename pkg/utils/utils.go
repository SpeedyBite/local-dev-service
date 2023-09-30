package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func SafeSprintf(format string, args ...interface{}) (string, error) {
	// Check if the format string contains the same number of placeholders
	// as the number of arguments.
	numPlaceholders := strings.Count(format, "%")
	numArgs := len(args)

	if numPlaceholders != numArgs {
		return "", errors.New("missing or extra placeholders in the format string")
	}

	// If the number of placeholders and arguments match, use Sprintf to format the string.
	formattedString := fmt.Sprintf(format, args...)

	return formattedString, nil
}

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
