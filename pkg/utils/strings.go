package utils

import (
	"fmt"
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
