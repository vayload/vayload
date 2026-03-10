/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Transports
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package vayload

import (
	"bufio"
	"context"
	"io"
	"mime/multipart"
	"time"

	// Only necessary for grpc transport
	"github.com/vayload/vayload/internal/shared/snowflake"
	"google.golang.org/grpc"
)

// ===========================================================
// HTTP Transport
// ===========================================================

type HttpAuth struct {
	UserId      snowflake.ID  `json:"user_id"`
	Role        string        `json:"role"`
	AccessToken string        `json:"access_token,omitempty"` // Optional access token for the user
	CountryId   *snowflake.ID `json:"country_id,omitempty"`
}

type HttpMethod string

const (
	HttpGet     HttpMethod = "GET"
	HttpPost    HttpMethod = "POST"
	HttpPut     HttpMethod = "PUT"
	HttpDelete  HttpMethod = "DELETE"
	HttpPatch   HttpMethod = "PATCH"
	HttpOptions HttpMethod = "OPTIONS"
	HttpHead    HttpMethod = "HEAD"
)

func (m HttpMethod) String() string {
	return string(m)
}

type HttpHandler func(req HttpRequest, res HttpResponse) error

type HttpRoute struct {
	Path           string
	Method         HttpMethod
	Handler        HttpHandler
	Middlewares    []HttpHandler
	PermissionRule string
	Public         bool
}

type HttpRoutesGroup struct {
	Prefix      string
	Public      bool
	Middlewares []HttpHandler
	Routes      []HttpRoute
}

type Cookie struct {
	Name        string
	Value       string
	Path        string
	Domain      string
	MaxAge      int
	Expires     time.Time
	Secure      bool
	HttpOnly    bool
	SameSite    string
	SessionOnly bool
}

type StreamWriter func(w *bufio.Writer) error

type HttpRequest interface {
	GetParam(key string, defaultValue ...string) string
	GetParamInt(key string, defaultValue ...int) (int, error)
	GetBody() []byte
	GetHeader(key string) string
	GetHeaders() map[string]string
	GetMethod() string
	GetPath() string
	GetQuery(key string, defaultValue ...string) string
	GetQueryInt(key string, defaultValue ...int) int
	Queries() map[string]string
	GetIP() string
	GetUserAgent() string
	GetHost() string
	ParseBody(any) error
	File(key string) (*multipart.FileHeader, error)
	FormData(key string) []string
	SaveFile(file *multipart.FileHeader, destination string) error
	GetCookie(name string) string
	Context() context.Context
	Locals(key string, value any) any
	Auth() *HttpAuth
	GetLocal(key string) any
	Next() error
	Validate(any) error     // Validate the request body using a validator
	ValidateBody(any) error // Parse and validate the request body
}

type HttpResponse interface {
	SetStatus(status int)
	SetHeader(key string, value string)
	Send(data []byte) error
	JSON(data any) error
	Json(data any) error
	File(path string) error
	Stream(stream io.Reader) error
	Status(status int) HttpResponse
	Redirect(path string, status int) error
	SetBodyStreamWriter(writer StreamWriter) error
	Cookie(cookie *Cookie) HttpResponse
	Cookies(cookies ...*Cookie) HttpResponse
}

// Http Exposer for services expose http routes
type HttpExposer interface {
	HttpRoutes() []HttpRoutesGroup
}

// ====================================================================================
// Grpc Transport
// ====================================================================================

type GrpcServiceDescriptor interface {
	ServiceName() string
	Register(server *grpc.Server)
}

type GrpcExposer interface {
	GrpcServices() []GrpcServiceDescriptor
}

// ====================================================================================
// Console Transport
// ====================================================================================

type ConsoleExposer interface {
	ConsoleCommands() []ConsoleCommand
}

// ====================================================================================
// Transports signature
// ====================================================================================

type Transport interface {
	Serve() error
	Shutdown(context.Context) error
	IsListening() bool
}

type HttpTransport interface {
	Transport
}

type GrpcTransport interface {
	Transport
	RegisterServices(services []GrpcServiceDescriptor, version string)
}

type McpTransport interface {
	Transport
	RegisterTools(tools []any)
}
