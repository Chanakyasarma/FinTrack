package config

import (
	"fmt"
	"log"
	"strings"
)

// Validate checks that production-critical config values are set.
// Call this from main() before starting the server.
func (c *Config) Validate() error {
	var errs []string

	if c.JWTSecret == "super-secret-change-in-production" && c.Environment == "production" {
		errs = append(errs, "JWT_SECRET must be changed from the default in production")
	}
	if len(c.JWTSecret) < 32 {
		errs = append(errs, "JWT_SECRET should be at least 32 characters")
	}
	if c.DatabaseURL == "" {
		errs = append(errs, "DATABASE_URL is required")
	}
	if c.RedisURL == "" {
		errs = append(errs, "REDIS_URL is required")
	}

	if len(errs) > 0 {
		return fmt.Errorf("config validation failed:\n  - %s", strings.Join(errs, "\n  - "))
	}

	log.Printf("config loaded [env=%s port=%s]", c.Environment, c.Port)
	return nil
}
