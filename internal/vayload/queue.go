/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - JobQueue
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package vayload

import "context"

type Job interface {
	Topic() string
	Payload() any
}

type JobHandler func(Job) error

type JobQueue interface {
	Enqueue(ctx context.Context, job Job) error
	Consume(ctx context.Context, topic string, handler JobHandler)
}
