/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - ds
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package ds

import "sync"

type HashMap[K comparable, V any] struct {
	mu   sync.RWMutex // for concurrent access and safe multireaders
	data map[K]V
}

func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		data: make(map[K]V),
	}
}

func (m *HashMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *HashMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

func (m *HashMap[K, V]) Has(key K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exits := m.data[key]
	return exits
}

func (m *HashMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
}

func (m *HashMap[K, V]) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

func (m *HashMap[K, V]) Range(fn func(K, V) bool) {
	m.mu.RLock()
	copyData := make(map[K]V, len(m.data))
	for k, v := range m.data {
		copyData[k] = v
	}
	m.mu.RUnlock()

	for k, v := range copyData {
		if !fn(k, v) {
			return
		}
	}
}
