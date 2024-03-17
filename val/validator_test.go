package val

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFullName(t *testing.T) {
	tests := []struct {
		name      string
		fullName  string
		expectErr bool
	}{
		{"Valid Full Name", "Valid Name", false},
		{"Invalid Characters", "Invalid_Name123", true},
		{"Empty String", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFullName(tt.fullName)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		expectErr bool
	}{
		{"Valid Password", "password123", false},
		{"Short Password", "short", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		expectErr bool
	}{
		{"Valid Email", "valid@example.com", false},
		{"Invalid Email", "invalid.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateID(t *testing.T) {
	tests := []struct {
		name      string
		emailId   int64
		expectErr bool
	}{
		{"Valid Email ID", 1, false},
		{"Invalid Email ID", 0, true},
		{"Invalid Email ID", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateID(tt.emailId)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateIntNotNegative(t *testing.T) {
	tests := []struct {
		name      string
		emailId   int64
		expectErr bool
	}{
		{"zero", 0, false},
		{"postive", 1, false},
		{"negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIntNotNegative(tt.emailId)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateSecretCode(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		expectErr bool
	}{
		{"Valid Code", "12345678901234567890123456789012", false},
		{"Short Code", "short", true},
		{"Long Code", "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSecretCode(tt.code)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		expectErr bool
	}{
		{"Too Short", "ab", true},
		{"Too Long", repeat("a", 101), true}, // Make this string 101 characters long
		{"Valid Simple", "valid_username123", false},
		{"Invalid Characters", "Invalid-User", true},
		{"Empty", "", true},
		{"Upper Case", "UserName", true},
		{"Valid with Underscore", "valid_user_name", false},
		{"Contains Space", "user name", true},
		{"Numeric Username", "123456", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Helper function to repeat a string n times.
func repeat(s string, n int) string {
	str := ""
	for i := 0; i < n; i++ {
		str += s
	}
	return str
}
