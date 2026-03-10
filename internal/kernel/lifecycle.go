package kernel

import (
	"github.com/vayload/vayload/internal/shared/ds"
	"github.com/vayload/vayload/internal/vayload"
)

type ServiceLifecycleDispatcher struct {
	registered *ds.ArrayList[vayload.ServiceRegisteredListener]
	starting   *ds.ArrayList[vayload.ServiceStartingListener]
	started    *ds.ArrayList[vayload.ServiceStartedListener]
	stopping   *ds.ArrayList[vayload.ServiceStoppingListener]
	stopped    *ds.ArrayList[vayload.ServiceStoppedListener]
	errors     *ds.ArrayList[vayload.ServiceErrorListener]
}

func NewServiceLifecycleDispatcher() *ServiceLifecycleDispatcher {
	return &ServiceLifecycleDispatcher{
		registered: ds.NewArrayList(func(a, b vayload.ServiceRegisteredListener) bool { return a == b }),
		starting:   ds.NewArrayList(func(a, b vayload.ServiceStartingListener) bool { return a == b }),
		started:    ds.NewArrayList(func(a, b vayload.ServiceStartedListener) bool { return a == b }),
		stopping:   ds.NewArrayList(func(a, b vayload.ServiceStoppingListener) bool { return a == b }),
		stopped:    ds.NewArrayList(func(a, b vayload.ServiceStoppedListener) bool { return a == b }),
		errors:     ds.NewArrayList(func(a, b vayload.ServiceErrorListener) bool { return a == b }),
	}
}

func (d *ServiceLifecycleDispatcher) AddRegisteredListener(l vayload.ServiceRegisteredListener) {
	d.registered.Add(l)
}

func (d *ServiceLifecycleDispatcher) AddStartingListener(l vayload.ServiceStartingListener) {
	d.starting.Add(l)
}

func (d *ServiceLifecycleDispatcher) AddStartedListener(l vayload.ServiceStartedListener) {
	d.started.Add(l)
}

func (d *ServiceLifecycleDispatcher) AddStoppingListener(l vayload.ServiceStoppingListener) {
	d.stopping.Add(l)
}

func (d *ServiceLifecycleDispatcher) AddStoppedListener(l vayload.ServiceStoppedListener) {
	d.stopped.Add(l)
}

func (d *ServiceLifecycleDispatcher) AddErrorListener(l vayload.ServiceErrorListener) {
	d.errors.Add(l)
}

func (d *ServiceLifecycleDispatcher) ServiceRegistered(event vayload.ServiceRegisteredEvent) {
	for _, l := range d.registered.ToSlice() {
		l.OnServiceRegistered(event)
	}
}

func (d *ServiceLifecycleDispatcher) ServiceStarting(event vayload.ServiceStartingEvent) {
	for _, l := range d.starting.ToSlice() {
		l.OnServiceStarting(event)
	}
}

func (d *ServiceLifecycleDispatcher) ServiceStarted(event vayload.ServiceStartedEvent) {
	for _, l := range d.started.ToSlice() {
		l.OnServiceStarted(event)
	}
}

func (d *ServiceLifecycleDispatcher) ServiceStopping(event vayload.ServiceStoppingEvent) {
	for _, l := range d.stopping.ToSlice() {
		l.OnServiceStopping(event)
	}
}

func (d *ServiceLifecycleDispatcher) ServiceStopped(event vayload.ServiceStoppedEvent) {
	for _, l := range d.stopped.ToSlice() {
		l.OnServiceStopped(event)
	}
}

func (d *ServiceLifecycleDispatcher) ServiceError(event vayload.ServiceErrorEvent) {
	for _, l := range d.errors.ToSlice() {
		l.OnServiceError(event)
	}
}

type KernelLifecycleDispatcher struct {
	bootstrap *ds.ArrayList[vayload.KernelBootstrapListener]
	started   *ds.ArrayList[vayload.KernelStartedListener]
	shutdown  *ds.ArrayList[vayload.KernelShutdownListener]
}

func NewKernelLifecycleDispatcher() *KernelLifecycleDispatcher {
	return &KernelLifecycleDispatcher{
		bootstrap: ds.NewArrayList(func(a, b vayload.KernelBootstrapListener) bool { return a == b }),
		started:   ds.NewArrayList(func(a, b vayload.KernelStartedListener) bool { return a == b }),
		shutdown:  ds.NewArrayList(func(a, b vayload.KernelShutdownListener) bool { return a == b }),
	}
}

func (d *KernelLifecycleDispatcher) RegisterBootstrapListener(l vayload.KernelBootstrapListener) {
	d.bootstrap.Add(l)
}

func (d *KernelLifecycleDispatcher) RegisterStartedListener(l vayload.KernelStartedListener) {
	d.started.Add(l)
}

func (d *KernelLifecycleDispatcher) RegisterShutdownListener(l vayload.KernelShutdownListener) {
	d.shutdown.Add(l)
}

func (d *KernelLifecycleDispatcher) KernelBootstrap(event vayload.KernelBootstrapEvent) {
	for _, l := range d.bootstrap.ToSlice() {
		l.OnKernelBootstrap(event)
	}
}

func (d *KernelLifecycleDispatcher) KernelStarted(event vayload.KernelStartedEvent) {
	for _, l := range d.started.ToSlice() {
		l.OnKernelStarted(event)
	}
}

func (d *KernelLifecycleDispatcher) KernelShutdown(event vayload.KernelShutdownEvent) {
	for _, l := range d.shutdown.ToSlice() {
		l.OnKernelShutdown(event)
	}
}
