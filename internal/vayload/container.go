/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package vayload

type Container interface {
	// Set registers a new dependency in the container.
	Set(name string, target any, shared bool) error

	// Override overrides an existing dependency in the container.
	Override(name string, target any, shared bool) error

	// Singleton registers a new singleton dependency in the container.
	Singleton(name string, target any) error

	// Deffered registers a new deferred dependency in the container.
	Deffered(name string, target func(Container) (any, error), shared bool) error

	// Get retrieves a dependency from the container.
	Get(name string) (any, error)

	// ResolveInto resolves a dependency and injects it into the provided target.
	// The target must be a non-nil pointer.
	ResolveInto(name string, target any) error

	// Has checks if a dependency is registered in the container.
	Has(name string) bool
}
