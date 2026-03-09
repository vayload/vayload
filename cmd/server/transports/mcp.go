/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Server Transports
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package server_transports

import (
	"context"
	"sync/atomic"

	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/vayload"

	// Transports dependencies
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/mark3labs/mcp-go/server"
)

type McpTransport struct {
	server *server.StreamableHTTPServer

	listener *fiber.App
	IsListen atomic.Bool
}

func CreateMCPServer(config *config.Config, httpListener *fiber.App) *McpTransport {
	mcpServer := server.NewMCPServer("vayload-mcp", "0.0.1")

	httpServer := server.NewStreamableHTTPServer(mcpServer, server.WithEndpointPath("/"))

	return &McpTransport{
		server:   httpServer,
		listener: httpListener,
	}
}

func (t *McpTransport) Serve() error {
	// Run server in background
	t.listener.Use("/mcp", adaptor.HTTPHandler(t.server))
	// go func() {
	// 	fmt.Printf("Http server start listening in http://localhost:%d\n", t.config.HTTP.Port)
	// 	if err := t.server.Listen(fmt.Sprintf(":%d", t.config.HTTP.Port)); err != nil && err != http.ErrServerClosed {
	// 		logger.F(err, logger.Fields{"context": "http_server"})
	// 	}

	// 	t.IsListen.Store(true)
	// }()
	t.IsListen.Store(true)

	return nil
}

func (t *McpTransport) Shutdown(ctx context.Context) error {
	if err := t.server.Shutdown(ctx); err != nil {
		return err
	}

	t.IsListen.Store(false)
	return nil
}

func (t *McpTransport) RegisterTools(tools []any) {
	// Create a new group for the version
	// v1 := t.server.Group(fmt.Sprintf("/api/%s/_rest", version))

	// // Add each route to the group
	// for _, route := range routes {
	// 	v1.Add(string(route.Method()), route.Path(), func(c *fiber.Ctx) error {
	// 		return route.Handler()(httpi.NewHttpRequest(c), httpi.NewHttpResponse(c))
	// 	})
	// }
}

func (t *McpTransport) IsListening() bool {
	return t.IsListen.Load()
}

var _ vayload.McpTransport = (*McpTransport)(nil)
