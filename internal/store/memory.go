package store

// Phase 1: we’ll add in-memory implementation later.
import (
	"context"
	"envwatch/internal/model"
	"errors"
	"sync"
)

var (
	ErrServiceNotFound = errors.New("service not found")
)

type MemoryStore struct {
	nextServiceID int64                   // It is a counter to generate unique IDs for services
	services      map[int64]model.Service // It is a map of services, the key is the service ID, the value is the service

	mu sync.RWMutex // what is this? it is a mutex to protect the services map from concurrent access

	nextResultID     int64                         // It is a counter to generate unique IDs for results
	resultsByService map[int64][]model.CheckResult // It is a map of results, the key is the service ID, the value is the results
}

// NewMemoryStore creates a new empty in-memory store.
func NewMemoryStore() *MemoryStore {
	return newMemoryStore()
}

func newMemoryStore() *MemoryStore {
	return &MemoryStore{
		services:         make(map[int64]model.Service),
		resultsByService: make(map[int64][]model.CheckResult),
		nextServiceID:    1,
		nextResultID:     1,
	}
}

var _ ServiceStore = (*MemoryStore)(nil) // ensure MemoryStore implements ServiceStore interface
var _ ResultStore = (*MemoryStore)(nil)

func (m *MemoryStore) CreateService(ctx context.Context, svc *model.Service) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	svc.ID = m.nextServiceID
	m.nextServiceID++
	m.services[svc.ID] = *svc
	return nil
}

func (m *MemoryStore) GetService(ctx context.Context, id int64) (*model.Service, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	svc, ok := m.services[id]
	if !ok {
		return nil, ErrServiceNotFound
	}
	return &svc, nil
}

func (m *MemoryStore) ListServices(ctx context.Context) ([]model.Service, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	services := make([]model.Service, 0, len(m.services))
	for _, svc := range m.services {
		services = append(services, svc)
	}
	return services, nil
}

func (m *MemoryStore) SaveResult(ctx context.Context, res *model.CheckResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	res.ID = m.nextResultID
	m.nextResultID++

	// append to slice for that service
	m.resultsByService[res.ServiceID] = append(m.resultsByService[res.ServiceID], *res)
	return nil
}

func (m *MemoryStore) GetRecentResults(ctx context.Context, serviceID int64, limit int) ([]model.CheckResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	all := m.resultsByService[serviceID]
	n := len(all)
	if n == 0 {
		return nil, nil
	}

	if limit <= 0 || limit > n {
		limit = n
	}

	// return last `limit` results (newest at the end)
	start := n - limit
	results := make([]model.CheckResult, 0, limit)
	results = append(results, all[start:]...)
	return results, nil
}
