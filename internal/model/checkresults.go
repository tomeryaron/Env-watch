package model

// Phase 1: we’ll fill this with CheckResult struct soon.
import "time"

type CheckResult struct {
	ID        int64     `json:"id"`                  // will be set by storage later
	ServiceID int64     `json:"service_id"`          // the ID of the service that was checked
	CheckedAt time.Time `json:"checked_at"`          // when the check was performed
	Success   bool      `json:"success"`             // true if the check was successful, false otherwise
	LatencyMs int       `json:"latency_ms"`          // in milliseconds
	ErrorMsg  *string   `json:"error_msg,omitempty"` // error message if the check failed
}
