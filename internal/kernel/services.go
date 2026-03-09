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
	kernel            vayload.Kernel
	shouldFailOnError bool
	dependencies      []vayload.ServiceName
	container         vayload.Container
}

func NewBaseService(name string, shouldFailOnError bool, dependencies ...vayload.ServiceName) BaseService {
	return BaseService{
		name:              name,
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

func (service *BaseService) Kernel() vayload.Kernel {
	return service.kernel
}

func (service *BaseService) SetKernel(k vayload.Kernel) {
	service.kernel = k
}

func (service *BaseService) Container() vayload.Container {
	return service.container
}

func (service *BaseService) SetContainer(c vayload.Container) {
	service.container = c
}

type manager struct {
	services map[vayload.ServiceName]vayload.Service
	mu       sync.RWMutex
}

func NewServiceManager() *manager {
	return &manager{
		services: make(map[vayload.ServiceName]vayload.Service),
	}
}

func (m *manager) RegisterService(service vayload.Service) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.services[vayload.ServiceName(service.Name())] = service
}

func (m *manager) DeleteService(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.services, vayload.ServiceName(name))
}

func (m *manager) GetService(name string) (vayload.Service, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.services[vayload.ServiceName(name)]
	if !ok {
		return nil, fmt.Errorf("service %s not found", name)
	}
	return s, nil
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

	for _, s := range ordered {
		s.SetStatus(vayload.ServiceStarting)
		err := s.Bootstrap(ctx, nil, nil)
		if err != nil {
			s.SetStatus(vayload.ServiceError)
			if s.ShouldFailOnError() {
				return fmt.Errorf("critical service failed: %s -> %w", s.Name(), err)
			}
			continue
		}
		s.SetStatus(vayload.ServiceRunning)
	}

	return nil
}

func (m *manager) resolveDependencies() ([]vayload.Service, error) {
	var sorted []vayload.Service
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

	for name := range m.services {
		if !visited[name] {
			if err := visit(name); err != nil {
				return nil, err
			}
		}
	}

	return sorted, nil
}

var _ vayload.ServiceManager = (*manager)(nil)
