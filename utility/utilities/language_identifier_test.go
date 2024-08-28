package utilities

import (
	"errors"
	"testing"
	"utility/internal/consts"

	"github.com/stretchr/testify/require"
)

func TestGenerateLangIdentifier(t *testing.T) {
	tests := []struct {
		testCaseDesc string
		input        string
		inputType    string
		identifier   []string
		expected     string
		expectedErr  error
	}{
		// Test case 1: Normal input with default identifier
		{
			testCaseDesc: "Normal input with default identifier",
			input:        "HelloWorld",
			inputType:    "Type",
			expected:     "helloworldType" + consts.DefaultIdentifier,
			expectedErr:  nil,
		},
		// Test case 2: Normal input with custom identifier
		{
			testCaseDesc: "Normal input with custom identifier",
			input:        "Hello123",
			inputType:    "Type",
			identifier:   []string{"CustomIdentifier"},
			expected:     "hello123TypeCustomIdentifier",
			expectedErr:  nil,
		},
		// Test case 3: Input with special characters
		{
			testCaseDesc: "Input with special characters",
			input:        "!@#$Hello%^& World*()",
			inputType:    "Type",
			identifier:   []string{""},
			expected:     "!@#$hello%^&World*()Type" + consts.DefaultIdentifier,
			expectedErr:  nil,
		},
		// Test case 4: Empty input
		{
			testCaseDesc: "Empty input",
			input:        "",
			inputType:    "Type",
			expected:     "",
			expectedErr:  errors.New("cannot create language label invalid input or input type"),
		},
	}

	for _, test := range tests {
		t.Run(test.testCaseDesc, func(t *testing.T) {
			actual, err := GenerateLangIdentifier(test.input, test.inputType, test.identifier...)
			require.Equal(t, test.expected, actual)
			require.Equal(t, test.expectedErr, err)
		})
	}
}
