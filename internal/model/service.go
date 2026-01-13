package model

import "errors"

// Phase 1: we’ll fill this with Service struct in the next step.

type ServiceType string

const (
	ServiceHTTP ServiceType = "http"
	ServiceTCP  ServiceType = "tcp"
)

// Service represents a monitored service, like a HTTP endpoint or TCP port.
// It contains the configuration for the service and its current status.
type Service struct {
	ID          int64       `json:"id"`                   // will be set by storage later
	Name        string      `json:"name"`                 // human-friendly name
	Type        ServiceType `json:"type"`                 // "http" or "tcp"
	Target      string      `json:"target"`               // URL or host:port
	IntervalSec int         `json:"interval_sec"`         // how often to check
	SLOTarget   string      `json:"slo_target,omitempty"` // optional SLO target
}

// Validate validates the service configuration, it is used to ensure the service configuration is valid before storing it in the database.
func (s Service) Validate() error { // validates the service configuration
	if s.Name == "" {
		return errors.New("Name is required")
	}
	if s.Type != ServiceHTTP && s.Type != ServiceTCP {
		return errors.New("invalid service type, must be http or tcp")
	}
	if s.Target == "" {
		return errors.New("target URL or host:port is required")
	}
	if s.IntervalSec <= 0 {
		return errors.New("interval must be greater than 0")
	}
	return nil
}
