package models_test

import (
	"testing"

	"github.com/ecumenos/orbis-socius/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateAccountUniqueName(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected bool
	}{
		"should be ok for lower case letters only": {
			input:    "qwerty",
			expected: true,
		},
		"should be not ok for upper case letters only": {
			input:    "QWERTY",
			expected: false,
		},
		"should be ok for lower & upper case letters only": {
			input:    "QWErty",
			expected: true,
		},
		"should be ok for lower case letters & digits only": {
			input:    "qwerty123",
			expected: true,
		},
		"should be ok for lower case letters & underscore only": {
			input:    "qwerty_qwerty",
			expected: true,
		},
		"should be ok for lower case & upper case letters & underscore only": {
			input:    "qweRTY_qwerty",
			expected: true,
		},
		"should be ok for lower case letters & digits & underscore only": {
			input:    "qwerty_qwerty123",
			expected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, models.ValidateAccountUniqueName(tc.input))
		})
	}
}
