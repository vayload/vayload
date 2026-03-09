/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package events

import (
	"reflect"
	"sync"
)

type EventBus interface {
	Publish(event any)
	Listen(handler any)
	ListenOnce(handler any)
	UnsubscribeType(eventType reflect.Type)
	UnsubscribeAll()
}

type Listeners interface {
	ListenOf(EventBus)
}

type EventListener struct {
	Handler   reflect.Value
	EventType reflect.Type
	OneTime   bool
}

type eventBus struct {
	listeners map[reflect.Type][]*EventListener
	workers   int
	queue     chan func()
	mutex     sync.RWMutex
}

func NewEventBus(workers int) *eventBus {
	eb := &eventBus{
		listeners: make(map[reflect.Type][]*EventListener),
		workers:   workers,
		queue:     make(chan func(), 10000),
	}

	for range eb.workers {
		go func() {
			for task := range eb.queue {
				task()
			}
		}()
	}

	return eb
}

func (e *eventBus) Publish(event any) {
	eventType := reflect.TypeOf(event)

	e.mutex.RLock()
	listeners := e.listeners[eventType]
	e.mutex.RUnlock()

	remaining := make([]*EventListener, 0, len(listeners))

	for _, listener := range listeners {
		l := listener
		eventValue := reflect.ValueOf(event)

		select {
		case e.queue <- func() {
			l.Handler.Call([]reflect.Value{eventValue})
		}:
			if !l.OneTime {
				remaining = append(remaining, l)
			}
		default:
			// Queue llena, descarta el evento
			if !l.OneTime {
				remaining = append(remaining, l)
			}
		}
	}

	// Actualizar listeners removiendo los OneTime
	if len(remaining) != len(listeners) {
		e.mutex.Lock()
		e.listeners[eventType] = remaining
		e.mutex.Unlock()
	}
}

func (e *eventBus) Listen(handler any) {
	e.addListener(handler, false)
}

func (e *eventBus) ListenOnce(handler any) {
	e.addListener(handler, true)
}

func (e *eventBus) addListener(handler any, once bool) {
	handlerType := reflect.TypeOf(handler)

	// Validar que sea una función
	if handlerType.Kind() != reflect.Func {
		panic("handler debe ser una función")
	}

	// Validar que tenga exactamente 1 parámetro
	if handlerType.NumIn() != 1 {
		panic("handler debe tener exactamente 1 parámetro")
	}

	// Obtener el tipo del evento (primer parámetro)
	eventType := handlerType.In(0)
	handlerValue := reflect.ValueOf(handler)

	e.mutex.Lock()
	defer e.mutex.Unlock()

	listener := &EventListener{
		Handler:   handlerValue,
		EventType: eventType,
		OneTime:   once,
	}

	e.listeners[eventType] = append(e.listeners[eventType], listener)
}

func (e *eventBus) UnsubscribeType(eventType reflect.Type) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	delete(e.listeners, eventType)
}

func (e *eventBus) UnsubscribeAll() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.listeners = make(map[reflect.Type][]*EventListener)
}

// Función de conveniencia para obtener el tipo de un evento
func TypeOf[T any]() reflect.Type {
	var zero T
	return reflect.TypeOf(zero)
}
