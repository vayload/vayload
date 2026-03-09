/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Services
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package vayload

import "context"

type ServiceStatus string

const (
	ServiceStopped  ServiceStatus = "stopped"
	ServiceStarting ServiceStatus = "starting"
	ServiceRunning  ServiceStatus = "running"
	ServiceError    ServiceStatus = "error"
)

// For other services reading
type PublicService interface {
	// Get service name
	Name() string

	// Get service status
	Status() ServiceStatus

	// Indicate if the service is running
	IsRunning() bool

	// Indicate if the kernel should fail if this service fails to start
	ShouldFailOnError() bool

	// Get service container
	Container() Container

	// Get required services (other services that need to be running before this service)
	RequiredServices() []ServiceName
}

// For kernel services
type Service interface {
	PublicService

	// Bootstrap the service
	Bootstrap(ctx context.Context, args map[string]any, reply *map[string]any) error

	// Set the service status
	SetStatus(status ServiceStatus)

	// Shutdown the service
	Shutdown()

	// Set the container
	SetContainer(c Container)
}

type ServiceManager interface {
	RegisterService(service Service)
	DeleteService(name string)
	GetService(name string) (Service, error)
	ListServices() []Service
}

type ServiceName string

// Core services names
const (
	ServiceStorageName  = ServiceName("storage")
	ServiceAuthName     = ServiceName("auth")
	ServiceDatabaseName = ServiceName("database")
	ServiceSettingsName = ServiceName("settings")
)
