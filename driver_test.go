package d1sql

import (
	"testing"
)

func TestParseDSN(t *testing.T) {
	tests := []struct {
		name        string
		dsn         string
		wantConfig  *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "valid DSN",
			dsn:  "d1://account123:token456@database789",
			wantConfig: &Config{
				AccountID:  "account123",
				APIToken:   "token456",
				DatabaseID: "database789",
				Timeout:    DefaultTimeout,
			},
			wantErr: false,
		},
		{
			name: "valid DSN with timeout",
			dsn:  "d1://account123:token456@database789?timeout=60s",
			wantConfig: &Config{
				AccountID:  "account123",
				APIToken:   "token456",
				DatabaseID: "database789",
				Timeout:    60000000000, // 60 seconds in nanoseconds
			},
			wantErr: false,
		},
		{
			name:        "empty DSN",
			dsn:         "",
			wantConfig:  nil,
			wantErr:     true,
			errContains: "empty DSN",
		},
		{
			name:        "invalid scheme",
			dsn:         "mysql://account123:token456@database789",
			wantConfig:  nil,
			wantErr:     true,
			errContains: "invalid scheme",
		},
		{
			name:        "missing account ID",
			dsn:         "d1://:token456@database789",
			wantConfig:  nil,
			wantErr:     true,
			errContains: "missing account ID",
		},
		{
			name:        "missing API token",
			dsn:         "d1://account123:@database789",
			wantConfig:  nil,
			wantErr:     true,
			errContains: "missing API token",
		},
		{
			name:        "missing database ID",
			dsn:         "d1://account123:token456@",
			wantConfig:  nil,
			wantErr:     true,
			errContains: "missing database ID",
		},
		{
			name:        "invalid timeout",
			dsn:         "d1://account123:token456@database789?timeout=invalid",
			wantConfig:  nil,
			wantErr:     true,
			errContains: "invalid timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseDSN(tt.dsn)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseDSN() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ParseDSN() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseDSN() unexpected error = %v", err)
				return
			}

			if config.AccountID != tt.wantConfig.AccountID {
				t.Errorf("AccountID = %v, want %v", config.AccountID, tt.wantConfig.AccountID)
			}
			if config.APIToken != tt.wantConfig.APIToken {
				t.Errorf("APIToken = %v, want %v", config.APIToken, tt.wantConfig.APIToken)
			}
			if config.DatabaseID != tt.wantConfig.DatabaseID {
				t.Errorf("DatabaseID = %v, want %v", config.DatabaseID, tt.wantConfig.DatabaseID)
			}
			if config.Timeout != tt.wantConfig.Timeout {
				t.Errorf("Timeout = %v, want %v", config.Timeout, tt.wantConfig.Timeout)
			}
		})
	}
}

func TestConfigFormatDSN(t *testing.T) {
	config := &Config{
		AccountID:  "account123",
		APIToken:   "token456",
		DatabaseID: "database789",
		Timeout:    DefaultTimeout,
	}

	dsn := config.FormatDSN()
	expected := "d1://account123:token456@database789"

	if dsn != expected {
		t.Errorf("FormatDSN() = %v, want %v", dsn, expected)
	}

	// Test round-trip
	parsedConfig, err := ParseDSN(dsn)
	if err != nil {
		t.Errorf("ParseDSN() error = %v", err)
		return
	}

	if parsedConfig.AccountID != config.AccountID {
		t.Errorf("Round-trip AccountID = %v, want %v", parsedConfig.AccountID, config.AccountID)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config",
			config: &Config{
				AccountID:  "account123",
				APIToken:   "token456",
				DatabaseID: "database789",
				Timeout:    DefaultTimeout,
			},
			wantErr: false,
		},
		{
			name: "missing account ID",
			config: &Config{
				AccountID:  "",
				APIToken:   "token456",
				DatabaseID: "database789",
				Timeout:    DefaultTimeout,
			},
			wantErr:     true,
			errContains: "account ID is required",
		},
		{
			name: "missing API token",
			config: &Config{
				AccountID:  "account123",
				APIToken:   "",
				DatabaseID: "database789",
				Timeout:    DefaultTimeout,
			},
			wantErr:     true,
			errContains: "API token is required",
		},
		{
			name: "missing database ID",
			config: &Config{
				AccountID:  "account123",
				APIToken:   "token456",
				DatabaseID: "",
				Timeout:    DefaultTimeout,
			},
			wantErr:     true,
			errContains: "database ID is required",
		},
		{
			name: "invalid timeout",
			config: &Config{
				AccountID:  "account123",
				APIToken:   "token456",
				DatabaseID: "database789",
				Timeout:    0,
			},
			wantErr:     true,
			errContains: "timeout must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("Validate() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

