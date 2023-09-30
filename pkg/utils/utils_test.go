package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RobDoan/go-docker-template/pkg/utils"
)

func TestSafeSprintf(t *testing.T) {
	str, err := utils.SafeSprintf("hello %s", "world")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", str)
}

func TestSafeSprintfWithMissingPlaceholder(t *testing.T) {
	testStr := "hello %s %s"
	str, err := utils.SafeSprintf(testStr, "world")
	assert.Error(t, err, "missing or extra placeholders in the format string")
	assert.Equal(t, "", str)
}
