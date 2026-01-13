# Env-watch

## Overview

EnvWatch is a service monitoring tool that continuously pings important services (HTTP endpoints, TCP ports, and Kubernetes services), measures uptime and latency, keeps history in a database, and exposes an API to view current status and SLO compliance.

## What EnvWatch Does

### Service Registration

You register services to monitor, for example:

- **API Gateway**: `https://api.example.com/health`
- **Payments Service**: `http://payments.namespace.svc.cluster.local:8080/health`
- **Redis**: `redis.namespace.svc.cluster.local:6379` (TCP)

### Monitoring Process

EnvWatch performs the following operations:

1. **Periodic Checks**: Runs health checks every N seconds
2. **Metrics Collection**: Measures:
   - Status (up/down)
   - Latency (in milliseconds)
3. **Data Storage**: Stores check results in SQLite database
4. **SLO Computation**: Calculates service level objectives, such as:
   - 99.9% availability over the last 7 days
   - P95 latency over the last 1 hour

## API Endpoints

EnvWatch exposes an HTTP API (and optionally a simple HTML page) with the following endpoints:

### Service Management

- `GET /services` - List all services and their current status
- `POST /services` - Register a new service to monitor

### Service Details

- `GET /services/{id}/history` - Get last X check results for a service
- `GET /services/{id}/slo` - Get computed SLOs for a service
- `POST /services/{id}/check` - Trigger an on-demand check for a service

## Deployment

EnvWatch is designed to run inside Kubernetes:

- **Deployment**: Runs as a Kubernetes Deployment
- **Service**: Exposes a Kubernetes Service
- **Health Probes**: Configured with liveness and readiness probes
- **Monitoring Scope**: Monitors other services in the same cluster or external URLs
