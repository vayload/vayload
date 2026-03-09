/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Event Bus
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package kernel

import (
	"context"
	"sync"

	"github.com/vayload/vayload/internal/vayload"
)

type EventListener struct {
	Handler vayload.EventHandler
	OneTime bool
}

type eventBus struct {
	listeners map[string][]*EventListener
	workers   int
	queue     chan func()

	mutex sync.RWMutex
}

func NewEventBus(workers int) *eventBus {
	eb := &eventBus{
		listeners: make(map[string][]*EventListener),
		workers:   workers,
		queue:     make(chan func(), 10000),
	}

	// Start the worker pool.
	for range eb.workers {
		go func() {
			for task := range eb.queue {
				task()
			}
		}()
	}

	return eb
}

func (e *eventBus) Publish(ctx context.Context, event vayload.Event) {
	eventType := event.Name()

	e.mutex.RLock()
	listeners, exists := e.listeners[eventType]
	e.mutex.RUnlock()

	// Check if there are listeners for this event.
	if !exists || len(listeners) == 0 {
		return
	}

	e.mutex.RLock()
	listenersCopy := make([]*EventListener, len(listeners))
	copy(listenersCopy, listeners)
	e.mutex.RUnlock()

	for _, listener := range listenersCopy {
		l := listener

		select {
		case e.queue <- func() {
			l.Handler(event)
			if l.OneTime {
				e.removeListener(ctx, eventType, l)
			}
		}:
		case <-ctx.Done():
			return
		}
	}
}

func (e *eventBus) removeListener(_ context.Context, event string, target *EventListener) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	listeners := e.listeners[event]
	for i, l := range listeners {
		// Compare pointers to find the exact listener instance
		if l == target {
			// Delete preserving order is not strictly necessary but safer;
			// using swap-delete is faster but changes order.
			// Here we use standard slice removal:
			e.listeners[event] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
}

func (e *eventBus) Subscribe(ctx context.Context, event string, handler vayload.EventHandler) {
	e.addListener(ctx, event, handler, false)
}

func (e *eventBus) SubscribeOnce(ctx context.Context, event string, handler vayload.EventHandler) {
	e.addListener(ctx, event, handler, true)
}

func (e *eventBus) addListener(_ context.Context, event string, handler vayload.EventHandler, once bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	listener := &EventListener{
		Handler: handler,
		OneTime: once,
	}

	e.listeners[event] = append(e.listeners[event], listener)
}

func (e *eventBus) Unsubscribe(ctx context.Context, event string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	delete(e.listeners, event)
}

var _ vayload.EventBus = (*eventBus)(nil)
