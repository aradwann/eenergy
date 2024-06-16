package migrate

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetSQLFiles(t *testing.T) {
	testCases := []struct {
		name                    string
		migrationDir            string
		expectedFiles           []string
		expectedErrMsgSubstring string
	}{
		{
			name:          "Valid Directory",
			migrationDir:  "test_data",
			expectedFiles: []string{"test_data/file_1.sql", "test_data/subdir/file_2.sql"},
		},
		{
			name:                    "Non-Existent Directory",
			migrationDir:            "nonexistent_directory",
			expectedErrMsgSubstring: "no such file or directory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			files, err := getSQLFiles(tc.migrationDir)

			if tc.expectedErrMsgSubstring != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectedErrMsgSubstring) {
					t.Fatalf("Expected error containing '%s', but got '%v'", tc.expectedErrMsgSubstring, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(files, tc.expectedFiles) {
				t.Fatalf("Expected files %v, but got %v", tc.expectedFiles, files)
			}
		})
	}
}
