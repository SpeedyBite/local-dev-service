package serviceconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RobDoan/go-docker-template/pkg/serviceconfig"
)

func TestServiceConfig_GetDomain(t *testing.T) {
	tests := []struct {
		name          string
		env           string
		domainFormat  string
		expected      string
		expectedError error
	}{
		{
			name:         "valid domain format",
			env:          "production",
			domainFormat: "https://%s.example.com",
			expected:     "https://production.example.com",
		},
		{
			name:         "invalid domain format",
			env:          "staging",
			domainFormat: "https://example.com",
			expected:     "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := serviceconfig.New("test-service", tt.domainFormat)
			actual, err := config.GetDomain(tt.env)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}

// END: 5f8b3c1d7f6e
