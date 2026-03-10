/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Services
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package kernel

import (
	"context"
	"fmt"
	"sync"

	"github.com/vayload/vayload/internal/vayload"
)

type BaseService struct {
	name              string
	status            vayload.ServiceStatus
	shouldFailOnError bool
	dependencies      []vayload.ServiceName
	container         vayload.Container
	eventBus          vayload.EventBus
}

func NewBaseService(name vayload.ServiceName, shouldFailOnError bool, dependencies ...vayload.ServiceName) BaseService {
	return BaseService{
		name:              string(name),
		status:            vayload.ServiceStopped,
		shouldFailOnError: shouldFailOnError,
		dependencies:      dependencies,
	}
}

func (service *BaseService) SetStatus(status vayload.ServiceStatus) {
	service.status = status
}

func (service *BaseService) Name() string {
	return service.name
}

func (service *BaseService) Status() vayload.ServiceStatus {
	return service.status
}

func (service *BaseService) RequiredServices() []vayload.ServiceName {
	return service.dependencies
}

func (service *BaseService) IsRunning() bool {
	return service.status == vayload.ServiceRunning
}

func (service *BaseService) ShouldFailOnError() bool {
	return service.shouldFailOnError
}

func (service *BaseService) Container() vayload.Container {
	return service.container
}

func (service *BaseService) SetContainer(c vayload.Container) {
	service.container = c
}

func (service *BaseService) EventBus() vayload.EventBus {
	return service.eventBus
}

func (service *BaseService) SetEventBus(c vayload.EventBus) {
	service.eventBus = c
}

type manager struct {
	*ServiceLifecycleDispatcher
	services map[vayload.ServiceName]vayload.Service
	ordered  []vayload.Service
	registry vayload.Container
	events   vayload.EventBus

	mu sync.RWMutex
}

func NewServiceManager(registry vayload.Container, events vayload.EventBus) *manager {
	return &manager{
		services:                   make(map[vayload.ServiceName]vayload.Service),
		ordered:                    make([]vayload.Service, 0),
		ServiceLifecycleDispatcher: NewServiceLifecycleDispatcher(),
		registry:                   registry,
		events:                     events,
	}
}

func (m *manager) RegisterService(service vayload.Service) {
	m.mu.Lock()
	defer m.mu.Unlock()

	service.SetContainer(m.registry)
	service.SetEventBus(m.events)
	m.services[vayload.ServiceName(service.Name())] = service

	m.ServiceRegistered(vayload.ServiceRegisteredEvent{
		Service: service,
	})
}

func (m *manager) DeleteService(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.services, vayload.ServiceName(name))
}

func (m *manager) GetService(name string) (vayload.Service, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.services[vayload.ServiceName(name)]
	return s, ok
}

func (m *manager) ListServices() []vayload.Service {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := make([]vayload.Service, 0, len(m.services))
	for _, s := range m.services {
		list = append(list, s)
	}
	return list
}

func (m *manager) StartAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ordered, err := m.resolveDependencies()
	if err != nil {
		return err
	}

	m.ordered = ordered

	for _, s := range ordered {
		s.SetStatus(vayload.ServiceStarting)
		m.ServiceStarting(vayload.ServiceStartingEvent{
			Context: ctx,
			Service: s,
		})

		err := s.Bootstrap(ctx, nil, nil)
		if err != nil {
			s.SetStatus(vayload.ServiceError)
			m.ServiceError(vayload.ServiceErrorEvent{
				Context: ctx,
				Service: s,
				Error:   err,
			})

			if s.ShouldFailOnError() {
				return fmt.Errorf("critical service failed: %s -> %w", s.Name(), err)
			}
			continue
		}

		s.SetStatus(vayload.ServiceRunning)
		m.ServiceStarted(vayload.ServiceStartedEvent{
			Context: ctx,
			Service: s,
		})
	}

	return nil
}

func (m *manager) resolveDependencies() ([]vayload.Service, error) {
	sorted := make([]vayload.Service, 0, len(m.services))
	visited := make(map[vayload.ServiceName]bool)
	temporary := make(map[vayload.ServiceName]bool)

	var visit func(name vayload.ServiceName) error
	visit = func(name vayload.ServiceName) error {
		if temporary[name] {
			return fmt.Errorf("circular dependency detected at service: %s", name)
		}
		if !visited[name] {
			temporary[name] = true
			s, ok := m.services[name]
			if !ok {
				return fmt.Errorf("required service %s is not registered", name)
			}

			for _, dep := range s.RequiredServices() {
				if err := visit(dep); err != nil {
					return err
				}
			}

			visited[name] = true
			delete(temporary, name)
			sorted = append(sorted, s)
		}
		return nil
	}

	// We need to ensure we visit all services
	// The order here doesn't strictly matter as long as dependencies are visited first
	for name := range m.services {
		if !visited[name] {
			if err := visit(name); err != nil {
				return nil, err
			}
		}
	}

	return sorted, nil
}

func (m *manager) StopAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Stop in reverse order of startup
	for i := len(m.ordered) - 1; i >= 0; i-- {
		s := m.ordered[i]

		if s == nil {
			continue
		}

		m.ServiceStopping(vayload.ServiceStoppingEvent{
			Context: ctx,
			Service: s,
		})

		s.SetStatus(vayload.ServiceStopped)
		err := s.Shutdown(ctx)
		if err != nil {
			s.SetStatus(vayload.ServiceError)
			m.ServiceError(vayload.ServiceErrorEvent{
				Context: ctx,
				Service: s,
				Error:   err,
			})
			return fmt.Errorf("failed to stop service %s: %w", s.Name(), err)
		}

		m.ServiceStopped(vayload.ServiceStoppedEvent{
			Context: ctx,
			Service: s,
		})
	}

	return nil
}

func (m *manager) OnServiceRegistered(l vayload.ServiceRegisteredListener) {
	m.ServiceLifecycleDispatcher.AddRegisteredListener(l)
}

func (m *manager) OnServiceStarting(l vayload.ServiceStartingListener) {
	m.ServiceLifecycleDispatcher.AddStartingListener(l)
}

func (m *manager) OnServiceStarted(l vayload.ServiceStartedListener) {
	m.ServiceLifecycleDispatcher.AddStartedListener(l)
}

func (m *manager) OnServiceStopping(l vayload.ServiceStoppingListener) {
	m.ServiceLifecycleDispatcher.AddStoppingListener(l)
}

func (m *manager) OnServiceStopped(l vayload.ServiceStoppedListener) {
	m.ServiceLifecycleDispatcher.AddStoppedListener(l)
}

func (m *manager) OnServiceError(l vayload.ServiceErrorListener) {
	m.ServiceLifecycleDispatcher.AddErrorListener(l)
}

var _ vayload.ServiceManager = (*manager)(nil)
