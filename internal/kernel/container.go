/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package kernel

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/vayload/vayload/internal/vayload"
)

var (
	ErrContainerNilTarget              = errors.New("cannot register a nil target")
	ErrContainerInvalidResolver        = errors.New("resolver must have signature func(Container) (any, error)")
	ErrContainerResolverCannotBeShared = errors.New("resolvers cannot be marked as shared (shared is only for direct instances)")
	ErrContainerEmptyName              = errors.New("service name cannot be empty")
	ErrContainerOverride               = errors.New("service already exists")
)

type container struct {
	instances map[string]any // used for shared instances, like singletons or static values
	resolvers map[string]func(vayload.Container) (any, error)
	shared    map[string]bool // used to track which instances are shared

	mutex sync.RWMutex
}

func NewContainer() *container {
	return &container{
		instances: make(map[string]any),
		resolvers: make(map[string]func(vayload.Container) (any, error)),
		shared:    make(map[string]bool),
	}
}

func (container *container) setTarget(name string, target any, shared bool, override bool) error {
	if name == "" {
		return ErrContainerEmptyName
	}

	if target == nil {
		return ErrContainerNilTarget
	}

	container.mutex.Lock()
	defer container.mutex.Unlock()

	if !override {
		if _, ok := container.instances[name]; ok {
			return ErrContainerOverride
		}
		if _, ok := container.resolvers[name]; ok {
			return ErrContainerOverride
		}
	}

	t := reflect.TypeOf(target)

	// Dereference *func(...)
	if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Func {
		t = t.Elem()
		target = reflect.ValueOf(target).Elem().Interface()
	}

	if t.Kind() == reflect.Func {
		// func(Container) (any, error)
		if t.NumIn() != 1 || t.NumOut() != 2 ||
			t.In(0) != reflect.TypeOf((*vayload.Container)(nil)).Elem() ||
			t.Out(0) != reflect.TypeOf((*any)(nil)).Elem() ||
			t.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			return ErrContainerInvalidResolver
		}

		resolver := target.(func(vayload.Container) (any, error))

		container.resolvers[name] = resolver

		if shared {
			container.shared[name] = true
		} else {
			delete(container.shared, name)
		}

		delete(container.instances, name)
		return nil
	}

	// Direct instance
	container.instances[name] = target

	if shared {
		container.shared[name] = true
	} else {
		delete(container.shared, name)
	}

	delete(container.resolvers, name)
	return nil
}

func (container *container) setResolver(name string, resolver func(vayload.Container) (any, error), shared bool) error {
	if name == "" {
		return ErrContainerEmptyName
	}

	if resolver == nil {
		return ErrContainerNilTarget
	}

	container.mutex.Lock()
	defer container.mutex.Unlock()

	// Prevent override
	if _, ok := container.instances[name]; ok {
		return ErrContainerOverride
	}
	if _, ok := container.resolvers[name]; ok {
		return ErrContainerOverride
	}

	container.resolvers[name] = resolver

	if shared {
		container.shared[name] = true
	} else {
		// Ensure resolver is not marked as shared
		delete(container.shared, name)
	}

	return nil
}

func (container *container) Set(name string, target any, shared bool) error {
	return container.setTarget(name, target, shared, false)
}

func (container *container) Override(name string, target any, shared bool) error {
	return container.setTarget(name, target, shared, true)
}

func (container *container) Singleton(name string, target any) error {
	return container.setTarget(name, target, true, false)
}

func (container *container) Get(name string) (any, error) {
	container.mutex.RLock()
	instance, found := container.instances[name]
	container.mutex.RUnlock()

	if found {
		return instance, nil
	}

	container.mutex.RLock()
	resolver, found := container.resolvers[name]
	shared := container.shared[name]
	container.mutex.RUnlock()

	if !found {
		return nil, fmt.Errorf("dependency '%s' not registered", name)
	}

	instance, err := resolver(container)
	if err != nil {
		return nil, err
	}

	if shared {
		// double-check (race-safe singleton)
		container.mutex.Lock()
		if existing, ok := container.instances[name]; ok {
			container.mutex.Unlock()
			return existing, nil
		}
		container.instances[name] = instance
		container.mutex.Unlock()
	}

	return instance, nil
}

func (container *container) ResolveInto(name string, target any) error {
	// Ensure target is a non-nil pointer
	valueOfTarget := reflect.ValueOf(target)
	if valueOfTarget.Kind() != reflect.Pointer || valueOfTarget.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}

	instance, err := container.Get(name)
	if err != nil {
		return err
	}

	valueOfTarget.Elem().Set(reflect.ValueOf(instance))
	return nil
}

func (container *container) Has(name string) bool {
	container.mutex.RLock()

	defer container.mutex.RUnlock()

	if _, found := container.instances[name]; found {
		return true
	}

	if _, found := container.resolvers[name]; found {
		return true
	}

	return false
}

func (container *container) Deffered(name string, target func(vayload.Container) (any, error), shared bool) error {
	return container.setResolver(name, target, shared)
}

var _ vayload.Container = (*container)(nil)
