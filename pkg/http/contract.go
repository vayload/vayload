/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package httpi

import (
	"github.com/vayload/vayload/internal/vayload"
)

const (
	HTTP_AUTH_KEY  = "__auth__"
	HTTP_USER_KEY  = "__user__"
	HTTP_PERMS_KEY = "__permissions__"
)

type HttpMethod string

const (
	GET     HttpMethod = "GET"
	POST    HttpMethod = "POST"
	PUT     HttpMethod = "PUT"
	DELETE  HttpMethod = "DELETE"
	PATCH   HttpMethod = "PATCH"
	OPTIONS HttpMethod = "OPTIONS"
	HEAD    HttpMethod = "HEAD"
)

type HttpRoute struct {
	Path           string
	Method         HttpMethod
	Handler        vayload.HttpHandler
	Middleware     []vayload.HttpHandler
	PermissionRule string // Optional permission rule for authorization
	Public         bool   // Indicates if the route is public
}

type Body struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
	Meta   any    `json:"meta,omitempty"`
}

type HttpError struct {
	Code       string `json:"code"`
	Reason     string `json:"reason,omitempty"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
	StatusCode int    `json:"-,omitempty"` // HTTP status code, not included in JSON response
	Cause      error  `json:"-"`           // Original error, not included in JSON response
}

type Error struct {
	Status string    `json:"status"`
	Error  HttpError `json:"error"`
	Meta   any       `json:"meta,omitempty"`
}

// RequestBody is a generic structure for HTTP request bodies.
type RequestBody[T any] struct {
	Data     T        `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	RequestID string `json:"request_id"`
}

// ResponseBody is a generic structure for HTTP response bodies.
type ResponseBody[T any] struct {
	Status   string    `json:"status"` // always "success"
	Data     T         `json:"data"`
	Metadata *RespMeta `json:"metadata,omitempty"`
}

type RespMeta struct {
	RequestID string `json:"request_id"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
}

// ErrorResponse is a structure for HTTP error responses.
type ErrorResponse struct {
	Status string    `json:"status"` // always "error"
	Error  HttpError `json:"error"`
	Meta   any       `json:"meta,omitempty"`
}
