package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"envwatch/internal/checker"
	"envwatch/internal/model"
	"envwatch/internal/store"
)

type Scheduler struct {
	Services store.ServiceStore
	Results  store.ResultStore
	Checker  checker.Checker

	// How often the scheduler wakes up to run checks.
	TickInterval time.Duration
}

// Start begins the periodic checking loop.
// It returns when ctx is cancelled.
func (s *Scheduler) Start(ctx context.Context) {
	if s.TickInterval <= 0 {
		s.TickInterval = 10 * time.Second // default
	}

	ticker := time.NewTicker(s.TickInterval)
	defer ticker.Stop()

	log.Printf("scheduler: starting, tick interval=%s", s.TickInterval)

	for {
		select {
		case <-ctx.Done():
			log.Println("scheduler: stopping due to context cancellation")
			return
		case <-ticker.C:
			s.runCycle(ctx)
		}
	}
}

// runCycle runs one round of checks for all services.
func (s *Scheduler) runCycle(ctx context.Context) {
	services, err := s.Services.ListServices(ctx)
	if err != nil {
		log.Printf("scheduler: ListServices error: %v", err)
		return
	}

	if len(services) == 0 {
		log.Println("scheduler: no services to check")
		return
	}

	log.Printf("scheduler: running checks for %d service(s)", len(services))

	var wg sync.WaitGroup

	for _, svc := range services {
		svc := svc // capture range variable
		wg.Add(1)

		go func(svc model.Service) {
			defer wg.Done()
			s.runCheckForService(ctx, svc)
		}(svc)
	}

	wg.Wait()
}

func (s *Scheduler) runCheckForService(ctx context.Context, svc model.Service) {
	res, err := s.Checker.Check(ctx, svc)
	if err != nil {
		log.Printf("scheduler: checker error for service %d (%s): %v", svc.ID, svc.Name, err)
		return
	}

	if err := s.Results.SaveResult(ctx, &res); err != nil {
		log.Printf("scheduler: SaveResult error for service %d (%s): %v", svc.ID, svc.Name, err)
		return
	}

	// Just a small log so we see activity
	if res.Success {
		log.Printf("scheduler: service=%s success latency=%dms", svc.Name, res.LatencyMs)
	} else {
		msg := ""
		if res.ErrorMsg != nil {
			msg = *res.ErrorMsg
		}
		log.Printf("scheduler: service=%s FAILED latency=%dms error=%q", svc.Name, res.LatencyMs, msg)
	}
}
