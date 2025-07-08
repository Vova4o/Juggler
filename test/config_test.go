package test

import (
	"testing"

	"juggler/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.WebPort != 8080 {
		t.Errorf("Expected WebPort to be 8080, got %d", cfg.WebPort)
	}
}

func TestLoadFromArgs(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedPort int
		expectError  bool
	}{
		{
			name:         "No arguments",
			args:         []string{"program"},
			expectedPort: 8080,
			expectError:  false,
		},
		{
			name:         "Valid port",
			args:         []string{"program", "9000"},
			expectedPort: 9000,
			expectError:  false,
		},
		{
			name:         "Invalid port format",
			args:         []string{"program", "abc"},
			expectedPort: 0,
			expectError:  true,
		},
		{
			name:         "Multiple arguments",
			args:         []string{"program", "3000", "extra"},
			expectedPort: 3000,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.LoadFromArgs(tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if cfg.WebPort != tt.expectedPort {
				t.Errorf("Expected WebPort to be %d, got %d", tt.expectedPort, cfg.WebPort)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
	}{
		{
			name: "Valid config",
			config: &config.Config{
				WebPort: 8080,
			},
			expectError: false,
		},
		{
			name: "Invalid port - zero",
			config: &config.Config{
				WebPort: 0,
			},
			expectError: true,
		},
		{
			name: "Invalid port - negative",
			config: &config.Config{
				WebPort: -1,
			},
			expectError: true,
		},
		{
			name: "Invalid port - too high",
			config: &config.Config{
				WebPort: 70000,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestConfigString(t *testing.T) {
	cfg := &config.Config{
		WebPort: 9000,
	}

	result := cfg.String()
	expected := "Порт: 9000"

	if result != expected {
		t.Errorf("Expected string '%s', got '%s'", expected, result)
	}
}
