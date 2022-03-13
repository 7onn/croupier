package croupier

import (
	"fmt"
)

// Croupier exposes all functionalities of the Croupier service.
type Croupier interface {
	Ping() bool
}

// Broker manages the internal state of the Croupier service.
type Broker struct {
	cfg Config // the croupier's configuration
}

// New initializes a new Croupier service.
func New(cfg Config) (*Broker, error) {
	r := &Broker{}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return r, nil
}

// Ping checks to see if the croupier's database is responding.
func (brk *Broker) Ping() bool {
	// This function would check the croupier's dependencies (datastores and whatnot); useful for Kubernetes probes
	return true
}
