//go:build !cache_redis
// +build !cache_redis

/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Cache
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package cache

func NewCache(config CacheConfig) (*LRUCache[string, any], error) {
	return NewLRUCache[string, any](100), nil
}
