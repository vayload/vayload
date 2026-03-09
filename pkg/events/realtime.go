/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package events

import (
	"fmt"
	"sync"
)

type Payload struct {
	Id        string `json:"id"`
	EventName string `json:"event"`
	Data      any    `json:"data"` // Data can be any type, typically a struct or map with json tags
}

type RealtimeHub interface {
	Subscribe(userId string) <-chan Payload
	Publish(userId string, payload Payload)
	Unsubscribe(userId string)
	// returns all subscribers ids.
	Subscribers() []string
}

type notifier struct {
	subscribers map[string]chan Payload
	mutex       sync.RWMutex
	bufferSize  int
}

func NewRealtimeEvents(bufferSize ...int) *notifier {
	size := 1
	if len(bufferSize) > 0 {
		size = bufferSize[0]
	}

	return &notifier{
		subscribers: make(map[string]chan Payload),
		bufferSize:  size,
	}
}

func (n *notifier) Subscribe(userId string) <-chan Payload {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if ch, exists := n.subscribers[userId]; exists {
		return ch
	}

	ch := make(chan Payload, n.bufferSize)
	n.subscribers[userId] = ch
	return ch
}

func (n *notifier) Publish(userId string, payload Payload) {
	n.mutex.RLock()
	ch, ok := n.subscribers[userId]
	n.mutex.RUnlock()

	if !ok {
		return
	}

	select {
	case ch <- payload:
	default:
		fmt.Printf("notifier channel full or no reader for key %v\n", userId)
	}
}

func (n *notifier) Unsubscribe(userId string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if ch, ok := n.subscribers[userId]; ok {
		fmt.Printf("Unsubscribing channel for key %v\n", userId)
		close(ch)
		delete(n.subscribers, userId)
	}
}

func (n *notifier) Subscribers() []string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	subscribers := make([]string, 0, len(n.subscribers))
	for userId := range n.subscribers {
		subscribers = append(subscribers, userId)
	}
	return subscribers
}
