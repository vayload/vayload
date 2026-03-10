/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package kernel

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/vayload/vayload/internal/vayload"
)

var (
	ErrServiceExists                   = errors.New("the service not exists")
	ErrContainerNilTarget              = errors.New("cannot register a nil target")
	ErrContainerEmptyName              = errors.New("service name cannot be empty")
	ErrContainerOverride               = errors.New("service already exists")
	ErrContainerInvalidResolver        = errors.New("resolver must have signature func(Container) (any, error)")
	ErrContainerResolverCannotBeShared = errors.New("resolvers cannot be marked as shared (shared is only for direct instances)")
)

// Closer allows a service to release resources when deleted
type Closer interface {
	Close() error
}

type ContextCloser interface {
	Close(ctx context.Context) error
}

// Provider creates a service instance
type Provider func(vayload.Container) (any, error)

// entry stores service metadata
type entry struct {
	provider func(vayload.Container) (any, error)
	shared   bool // Indicating is singleton or transient
	instance any
	typ      reflect.Type // cached type for fast GetInto
}

type container struct {
	mu sync.RWMutex

	services map[string]*entry
	ctx      context.Context
}

func NewContainer(ctx context.Context) *container {
	return &container{
		services: make(map[string]*entry, 10),
		ctx:      ctx,
	}
}

func (c *container) validate(name string, target any) error {
	if name == "" {
		return ErrContainerEmptyName
	}

	if target == nil {
		return ErrContainerNilTarget
	}

	return nil
}

func (c *container) closeInstance(instance any) {
	if instance == nil {
		return
	}
	if d, ok := instance.(ContextCloser); ok {
		_ = d.Close(c.ctx)
	} else if d, ok := instance.(Closer); ok {
		_ = d.Close()
	}
}

func (c *container) SetInstance(name string, instance any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.services[name]; exists {
		return ErrServiceExists
	}

	c.services[name] = &entry{
		instance: instance,
		shared:   true,
	}

	return nil
}

func (c *container) SetTarget(name string, target any, shared bool, override bool) error {
	if err := c.validate(name, target); err != nil {
		return err
	}

	// Skip and return error when service exists
	if !override {
		c.mu.RLock()
		if _, exists := c.services[name]; exists {
			c.mu.RUnlock()
			return ErrContainerOverride
		}

		c.mu.RUnlock()
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
			t.In(0) != reflect.TypeFor[vayload.Container]() ||
			t.Out(0) != reflect.TypeFor[any]() ||
			t.Out(1) != reflect.TypeFor[error]() {
			return ErrContainerInvalidResolver
		}

		resolver := target.(func(vayload.Container) (any, error))

		c.mu.Lock()
		c.services[name] = &entry{
			provider: resolver,
			shared:   shared,
		}
		c.mu.Unlock()
		return nil
	}

	// Direct instance
	c.mu.Lock()
	c.services[name] = &entry{
		instance: target,
		shared:   true,
		typ:      reflect.TypeOf(target),
	}
	c.mu.Unlock()

	return nil
}

// Has checks if a service is registered
func (c *container) Has(name string) bool {
	c.mu.RLock()
	_, ok := c.services[name]
	c.mu.RUnlock()

	return ok
}

// Set registers a direct instance. Fails if it already exists.
func (c *container) Set(name string, target any, shared bool) error {
	if err := c.validate(name, target); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.services[name]; exists {
		return ErrContainerOverride
	}

	c.services[name] = &entry{
		instance: target,
		shared:   shared,
		typ:      reflect.TypeOf(target),
	}

	return nil
}

// Override forcibly replaces an existing service, closing the old one if needed
func (c *container) Override(name string, target any, shared bool) error {
	if err := c.validate(name, target); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Release resources of the overridden instance if it exists
	if old, ok := c.services[name]; ok && old.instance != nil {
		if d, ok := old.instance.(Closer); ok {
			_ = d.Close()
		}
		if d, ok := old.instance.(ContextCloser); ok {
			_ = d.Close(c.ctx)
		}
	}

	c.services[name] = &entry{
		instance: target,
		shared:   shared,
		typ:      reflect.TypeOf(target),
	}

	return nil
}

// Deferred registers a provider to be resolved later (lazy loading)
func (c *container) Deferred(name string, target func(vayload.Container) (any, error), shared bool) error {
	if err := c.validate(name, target); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.services[name]; exists {
		return ErrContainerOverride
	}

	c.services[name] = &entry{
		provider: target,
		shared:   shared,
	}

	return nil
}

// Singleton is a syntactic sugar for a Deferred Singleton provider
func (c *container) Singleton(name string, target any) error {
	return c.Set(name, target, true)
}

// Get returns the service instance, creating it if necessary
func (c *container) Get(name string) (any, error) {
	c.mu.RLock()
	item, ok := c.services[name]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("container: service %s not found", name)
	}

	// Fast path for Singletons already instantiated or direct instances
	if item.shared && item.instance != nil {
		return item.instance, nil
	}

	// Slow path for initialization
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check locking
	if item.shared && item.instance != nil {
		return item.instance, nil
	}

	if item.provider == nil {
		return nil, fmt.Errorf("container: no provider for %s", name)
	}

	instance, err := item.provider(c)
	if err != nil {
		return nil, err
	}

	if item.shared {
		item.instance = instance
		item.typ = reflect.TypeOf(instance)
	}

	return instance, nil
}

// ResolveInto injects a service into a pointer. Validates type using cached entry.typ
func (c *container) GetInto(name string, target any) error {
	c.mu.RLock()
	item, ok := c.services[name]
	c.mu.RUnlock()

	if !ok {
		return fmt.Errorf("container: service %s not found", name)
	}

	tVal := reflect.ValueOf(target)
	if tVal.Kind() != reflect.Ptr || tVal.IsNil() {
		return errors.New("container: target must be a non-nil pointer")
	}

	// Type validation using the cached type
	if item.typ != nil && tVal.Elem().Type() != item.typ {
		return fmt.Errorf("container: target type mismatch for %s (expected *%v)", name, item.typ)
	}

	val, err := c.Get(name)
	if err != nil {
		return err
	}

	tVal.Elem().Set(reflect.ValueOf(val))
	return nil
}

// Delete removes a service and calls Closer interfaces if implemented
func (c *container) Delete(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.services[name]; ok {
		c.closeInstance(item.instance)
		delete(c.services, name)
	}
}

// Flush clears all services, dropping resources for singletons
func (c *container) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, item := range c.services {
		c.closeInstance(item.instance)
		delete(c.services, name)
	}
}

func (c *container) Context() context.Context {
	return c.ctx
}

// MapTo is a generic function that maps a service to a specific type
func MapTo[T any](c vayload.Container, name string) (T, error) {
	val, err := c.Get(name)
	if err != nil {
		var zero T
		return zero, err
	}
	typed, ok := val.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("container: service %s is not of type %T", name, zero)
	}
	return typed, nil
}

var _ vayload.Container = (*container)(nil)
