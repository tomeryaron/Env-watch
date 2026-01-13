package main

import (
	"context"
	"log"
	"time"

	"envwatch/internal/checker"
	"envwatch/internal/model"
	"envwatch/internal/scheduler"
	"envwatch/internal/store"
)

func main() {
	log.Println("envwatch – scheduler test (Phase 1)")

	// Create in-memory store (implements both ServiceStore and ResultStore)
	mem := store.NewMemoryStore()

	// Create a test HTTP service
	svc := model.Service{
		Name:        "Example HTTP",
		Type:        model.ServiceHTTP,
		Target:      "https://example.com",
		IntervalSec: 30,
	}

	if err := svc.Validate(); err != nil {
		log.Fatalf("service validation failed: %v", err)
	}

	// Save the service in the store (ID will be assigned)
	if err := mem.CreateService(context.Background(), &svc); err != nil {
		log.Fatalf("CreateService failed: %v", err)
	}
	log.Printf("created service with ID=%d", svc.ID)

	// Build checker and scheduler
	chk := checker.NewDefaultChecker()

	sched := &scheduler.Scheduler{
		Services:     mem,
		Results:      mem,
		Checker:      chk,
		TickInterval: 10 * time.Second, // run checks every 10s
	}

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start scheduler in background
	go sched.Start(ctx)

	// Let it run for ~35 seconds then exit
	time.Sleep(35 * time.Second)
	log.Println("main: done (exiting after test run)")
}
