package db

import (
	"fmt"
	"testing"
)

func TestGenerateParamPlaceholders(t *testing.T) {
	testCases := []struct {
		count            int
		expectedResult   string
		expectedErrorMsg string
	}{
		{
			count:            0,
			expectedResult:   "",
			expectedErrorMsg: "",
		},
		{
			count:            3,
			expectedResult:   "$1, $2, $3",
			expectedErrorMsg: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Count%d", testCase.count), func(t *testing.T) {
			result := generateParamPlaceholders(testCase.count)

			// Check if the result matches the expected value
			if result != testCase.expectedResult {
				t.Errorf("Expected: %s, Got: %s", testCase.expectedResult, result)
			}
		})
	}
}
