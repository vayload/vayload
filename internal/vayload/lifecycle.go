package vayload

import (
	"context"
)

/*
EVENTS
*/

type KernelBootstrapEvent struct {
	Context context.Context
	Args    []string
}

type KernelStartedEvent struct {
	Context context.Context
}

type KernelShutdownEvent struct {
	Context context.Context
}

type ServiceRegisteredEvent struct {
	Service Service
}

type ServiceStartingEvent struct {
	Context context.Context
	Service Service
}

type ServiceStartedEvent struct {
	Context context.Context
	Service Service
}

type ServiceStoppingEvent struct {
	Context context.Context
	Service Service
}

type ServiceStoppedEvent struct {
	Context context.Context
	Service Service
}

type ServiceErrorEvent struct {
	Context context.Context
	Service Service
	Error   error
}

/*
LISTENERS
*/

type KernelBootstrapListener interface {
	OnKernelBootstrap(KernelBootstrapEvent)
}

type KernelStartedListener interface {
	OnKernelStarted(KernelStartedEvent)
}

type KernelShutdownListener interface {
	OnKernelShutdown(KernelShutdownEvent)
}

type ServiceRegisteredListener interface {
	OnServiceRegistered(ServiceRegisteredEvent)
}

type ServiceStartingListener interface {
	OnServiceStarting(ServiceStartingEvent)
}

type ServiceStartedListener interface {
	OnServiceStarted(ServiceStartedEvent)
}

type ServiceStoppingListener interface {
	OnServiceStopping(ServiceStoppingEvent)
}

type ServiceStoppedListener interface {
	OnServiceStopped(ServiceStoppedEvent)
}

type ServiceErrorListener interface {
	OnServiceError(ServiceErrorEvent)
}

type ServiceLifecycleEvents interface {
	ServiceRegisteredListener
	ServiceStartingListener
	ServiceStartedListener
	ServiceStoppingListener
	ServiceStoppedListener
	ServiceErrorListener
}
