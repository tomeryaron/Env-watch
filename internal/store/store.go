package store

// Phase 1: we’ll define interfaces + in-memory implementation later.
import (
	"context"
	"envwatch/internal/model"
)

// Why we use context: It allows us to pass a cancellation signal to the function,
// and it allows us to pass a timeout to the function. It is a way to handle concurrent operations.

// ServiceStore handles storing and fetching monitored services.
// It is responsible for creating, updating, and retrieving services from the database.
// We use interface to allow for different implementations of the store.
// Meaning we can store services in a database, file, or in memory.
type ServiceStore interface {
	CreateService(ctx context.Context, service *model.Service) error
	GetService(ctx context.Context, id int64) (*model.Service, error)
	ListServices(ctx context.Context) ([]model.Service, error) // lists all services
}

type ResultStore interface {
	SaveResult(ctx context.Context, res *model.CheckResult) error
	GetRecentResults(ctx context.Context, serviceID int64, limit int) ([]model.CheckResult, error) // gets the most recent results for a service
}

// CheckResultStore handles storing and fetching check results.
// It is responsible for creating, updating, and retrieving check results from the database.
// We use interface to allow for different implementations of the store.
// Meaning we can store check results in a database, file, or in memory.
