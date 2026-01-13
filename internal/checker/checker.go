package checker

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"envwatch/internal/model"
)

// Checker runs a single check for a given Service and returns the result.
type Checker interface {
	Check(ctx context.Context, svc model.Service) (model.CheckResult, error)
}

// DefaultChecker is our concrete implementation of Checker.
type DefaultChecker struct {
	httpClient *http.Client
	timeout    time.Duration
}

// NewDefaultChecker creates a checker with a sane default timeout.
func NewDefaultChecker() *DefaultChecker {
	return &DefaultChecker{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		timeout: 5 * time.Second,
	}
}

// Ensure DefaultChecker implements Checker at compile time.
var _ Checker = (*DefaultChecker)(nil)

func (c *DefaultChecker) Check(ctx context.Context, svc model.Service) (model.CheckResult, error) {
	switch svc.Type {
	case model.ServiceHTTP:
		return c.checkHTTP(ctx, svc)
	case model.ServiceTCP:
		return c.checkTCP(ctx, svc)
	default:
		return model.CheckResult{}, errors.New("unsupported service type")
	}
}

func (c *DefaultChecker) checkHTTP(ctx context.Context, svc model.Service) (model.CheckResult, error) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, svc.Target, nil)
	if err != nil {
		return model.CheckResult{}, err
	}

	resp, err := c.httpClient.Do(req)
	latency := time.Since(start)

	res := model.CheckResult{
		ServiceID: svc.ID,
		CheckedAt: time.Now(),
		LatencyMs: int(latency.Milliseconds()), // convert int64 to int
	}

	if err != nil {
		msg := err.Error()
		res.Success = false
		res.ErrorMsg = &msg
		return res, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		res.Success = true
		return res, nil
	}

	msg := "unexpected status code: " + resp.Status
	res.Success = false
	res.ErrorMsg = &msg
	return res, nil
}

func (c *DefaultChecker) checkTCP(ctx context.Context, svc model.Service) (model.CheckResult, error) {
	start := time.Now()

	dialer := net.Dialer{
		Timeout: c.timeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", svc.Target)
	latency := time.Since(start)

	res := model.CheckResult{
		ServiceID: svc.ID,
		CheckedAt: time.Now(),
		LatencyMs: int(latency.Milliseconds()), // convert int64 to int
	}

	if err != nil {
		msg := err.Error()
		res.Success = false
		res.ErrorMsg = &msg
		return res, nil
	}
	_ = conn.Close()

	res.Success = true
	return res, nil
}
