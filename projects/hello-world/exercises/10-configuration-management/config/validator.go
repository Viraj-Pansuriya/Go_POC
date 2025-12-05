package config

import (
	"fmt"
)

// Validate checks if the configuration is valid
// Like Spring's @Validated on @ConfigurationProperties
func (c *Config) Validate() error {
	// Validate server config
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	// Validate database config
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config: %w", err)
	}

	// Validate JWT config
	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("jwt config: %w", err)
	}

	return nil
}

func (s *ServerConfig) Validate() error {
	if s.Port <= 0 || s.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", s.Port)
	}
	if s.Host == "" {
		return fmt.Errorf("host is required")
	}
	if s.ReadTimeout <= 0 {
		return fmt.Errorf("read_timeout must be positive")
	}
	if s.WriteTimeout <= 0 {
		return fmt.Errorf("write_timeout must be positive")
	}
	return nil
}

func (d *DatabaseConfig) Validate() error {
	if d.Host == "" {
		return fmt.Errorf("host is required")
	}
	if d.Port <= 0 || d.Port > 65535 {
		return fmt.Errorf("invalid port: %d", d.Port)
	}
	if d.DBName == "" {
		return fmt.Errorf("dbname is required")
	}
	return nil
}

func (j *JWTConfig) Validate() error {
	if j.Secret == "" {
		return fmt.Errorf("secret is required (set JWT_SECRET env var)")
	}
	if len(j.Secret) < 32 {
		return fmt.Errorf("secret must be at least 32 characters for security")
	}
	if j.Expiration <= 0 {
		return fmt.Errorf("expiration must be positive")
	}
	return nil
}

