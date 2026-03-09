/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Events
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package vayload

import "context"

type Event interface {
	Name() string
}

type EventHandler func(Event)

type EventBus interface {
	// Publish an event with event interface its contains name of event
	Publish(ctx context.Context, event Event)

	// Subscribe to an event
	Subscribe(ctx context.Context, eventName string, handler EventHandler)

	// Subscribe to an event once
	SubscribeOnce(ctx context.Context, eventName string, handler EventHandler)

	// Unsubscribe from an event
	Unsubscribe(ctx context.Context, eventName string)
}
