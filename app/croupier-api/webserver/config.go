package webserver

import (
	"github.com/7onn/croupier/internal/croupier"
)

// Config for the Croupier API.
type Config struct {
	Croupier croupier.Config
	Port     int
}

func validateConfig(cfg Config) error {
	// This function would be used to validate a hydrated configuration; return an error if its invalid.
	return nil
}
