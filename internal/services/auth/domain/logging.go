/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Domain
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package domain

import "time"

type ExternalRequestLogging struct {
	RequestID      string
	URL            string
	Method         string
	RequestBody    any
	RequestAt      time.Time
	ResponseBody   any
	ResponseStatus int
	RequestElapsed time.Duration
	Error          *string // Optional error message
}
