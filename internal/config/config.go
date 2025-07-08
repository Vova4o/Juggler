package config

import (
	"fmt"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	WebPort int
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		WebPort: 8080,
	}
}

// LoadFromArgs parses command line arguments and returns configuration
func LoadFromArgs(args []string) (*Config, error) {
	config := &Config{
		WebPort: 8080, 
	}

	// Optional first argument for custom port
	if len(args) >= 2 {
		port, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, fmt.Errorf("wrong port format: %v", err)
		}
		config.WebPort = port
	}

	return config, nil
}

// PrintUsage prints the usage information
func PrintUsage() {
	fmt.Println("Usage: go run main.go [port]")
	fmt.Println("Example: go run main.go")
	fmt.Println("Example with port: go run main.go 8080")
	fmt.Println()
	fmt.Println("Parameters:")
	fmt.Println("  port                - (optional) port for the web server (default 8080)")
	fmt.Println()
	fmt.Println("Juggling settings (number of balls, time) are set via the web interface.")
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.WebPort <= 0 || c.WebPort > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	return nil
}

// String returns a string representation of the configuration
func (c *Config) String() string {
	return fmt.Sprintf("Порт: %d", c.WebPort)
}
