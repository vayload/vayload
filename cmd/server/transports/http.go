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
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	cfg "github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/shared/errors"
	"github.com/vayload/vayload/internal/vayload"
	httpi "github.com/vayload/vayload/pkg/http"
	"github.com/vayload/vayload/pkg/logger"
	"github.com/vayload/vayload/pkg/operator"

	// External dependencies
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	httplogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var startTime = time.Now()

type HttpTransport struct {
	server      *fiber.App
	isListening atomic.Bool

	config *cfg.Config
}

func CreateHttpServer(config *cfg.Config) *HttpTransport {
	server := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return httpi.HttpErrorHandler(httpi.NewHttpRequest(c), httpi.NewHttpResponse(c), err)
		},
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		EnableTrustedProxyCheck: true, // Enable trusted proxy check for security
		StreamRequestBody:       true,
	})

	server.Use(httplogger.New())
	server.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(config.Cors.Origins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowCredentials: true,
	}))
	server.Use(helmet.New(helmet.Config{
		CrossOriginEmbedderPolicy: "unsafe-none",
		CrossOriginResourcePolicy: "cross-origin",
		ContentSecurityPolicy: "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.tailwindcss.com; " +
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
			"font-src 'self' https://fonts.gstatic.com data:; " +
			"img-src 'self' data: https:; " +
			"connect-src 'self' https:; " +
			"frame-src 'self'; " +
			"media-src 'self' data:; " +
			"object-src 'none'; " +
			"frame-ancestors 'none'; " +
			"form-action 'self'; " +
			"base-uri 'self';",
	}))
	server.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	server.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        100,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			request := httpi.NewHttpRequest(c)
			if request.Auth() != nil && request.Auth().UserId != 0 {
				return request.Auth().UserId.String()
			}

			return operator.Coalesce(
				request.GetHeader("x-forwarded-for"),
				request.GetHeader("x-real-ip"),
				request.GetIP(),
			)
		},
		LimitReached: func(c *fiber.Ctx) error {
			return httpi.ErrTooManyRequests(errors.New("rate limit exceeded"))
		},
	}))

	server.Add("GET", "/api/_rest/health", func(c *fiber.Ctx) error {
		uptime := time.Since(startTime).Truncate(time.Second).String()

		if strings.Contains(c.Get("Accept"), "text/html") {
			c.Set("Content-Type", "text/html")

			return c.Status(http.StatusOK).SendString(fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Vayload Status</title>
<style>
body {
	background:#0b0b0b;
	color:#eaeaea;
	font-family:system-ui,-apple-system,Segoe UI,Roboto,sans-serif;
	display:flex;
	align-items:center;
	justify-content:center;
	height:100vh;
	margin:0;
}

.box {
	text-align:center;
}

.status {
	font-size:28px;
	font-weight:600;
	color:#ff7a18;
	margin-bottom:8px;
}

.meta {
	font-size:13px;
	color:#777;
}
</style>
</head>
<body>

<div class="box">
	<div class="status">Vayload is healthy</div>
	<div class="meta">uptime %s</div>
</div>

</body>
</html>
`, uptime))
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status": "healthy",
			"uptime": uptime,
		})
	})

	// If app mode is full, serve static files and SPA Admin
	if config.App.Mode == cfg.AppModeFull {
		server.Static("/", filepath.Join(config.WorkDir, "public/build"), fiber.Static{
			Compress: true,
		})
		server.Get("*", func(c *fiber.Ctx) error {
			if strings.HasPrefix(c.Path(), "/api") {
				return fiber.ErrNotFound
			}

			index := filepath.Join(config.WorkDir, "public/build/index.html")
			// check if file exists
			if _, err := os.Stat(index); os.IsNotExist(err) {
				return fiber.ErrNotFound
			}

			return c.SendFile(index)
		})
	}

	return &HttpTransport{server: server, config: config}
}

func (t *HttpTransport) Serve() error {
	// Run http server in background
	go func() {
		fmt.Printf("Http server start listening in http://localhost:%d\n", t.config.HTTP.Port)
		if err := t.server.Listen(fmt.Sprintf(":%d", t.config.HTTP.Port)); err != nil && err != http.ErrServerClosed {
			logger.F(err, logger.Fields{"context": "http_server"})
		}

		t.isListening.Store(true)
	}()

	return nil
}

func (t *HttpTransport) Shutdown(ctx context.Context) error {
	if err := t.server.ShutdownWithContext(ctx); err != nil {
		return err
	}

	t.isListening.Store(false)
	return nil
}

func (t *HttpTransport) OnServiceStarted(e vayload.ServiceStartedEvent) {
	if exposer, ok := e.Service.(vayload.HttpExposer); ok {
		logger.I("Discovered HTTP routes for service", logger.Fields{
			"service": e.Service.Name(),
		})

		t.RegisterRouteGroups(exposer.HttpRoutes(), "v1", e.Service.Name())
	}
}

func (t *HttpTransport) RegisterRouteGroups(groups []vayload.HttpRoutesGroup, version string, service string) {
	// Create a new group for the version if not exists
	base := t.server.Group(fmt.Sprintf("/api/%s/_rest", version))

	RegisterHttpRoutes(base, service, groups)
}

func (t *HttpTransport) Server() *fiber.App {
	return t.server
}

func (t *HttpTransport) IsListening() bool {
	return t.isListening.Load()
}

var _ vayload.HttpTransport = (*HttpTransport)(nil)
var _ vayload.ServiceStartedListener = (*HttpTransport)(nil)

func RegisterHttpRoutes(app fiber.Router, service string, handlers []vayload.HttpRoutesGroup) {
	fmt.Println()
	fmt.Println(cyan + "🚀 Discovering Http routes for service: " + service + reset)

	for _, fh := range handlers {

		publicGroup := app.Group("/public" + fh.Prefix)
		privateGroup := app.Group(fh.Prefix)

		// Global middleware
		if len(fh.Middlewares) > 0 {
			for _, mw := range fh.Middlewares {
				privateGroup.Use(httpi.FiberWrap(mw))
			}
		}

		for _, route := range fh.Routes {
			path := strings.TrimPrefix(route.Path, "/")
			fullPath := pathJoin(fh.Prefix, route.Path)

			var handlers []fiber.Handler
			var group fiber.Router

			if route.Public {
				group = publicGroup
				fmt.Printf("PUBLIC   - %s %s\n", route.Method, fullPath)
			} else {
				group = privateGroup

				for _, mw := range route.Middlewares {
					handlers = append(handlers, httpi.FiberWrap(mw))
				}
			}

			handlers = append(handlers, httpi.FiberWrap(route.Handler))

			group.Add(string(route.Method), path, handlers...)

			LogRegisteredRoute(string(route.Method), fullPath)
		}
	}

	fmt.Println()
}

func pathJoin(parts ...string) string {
	return "/" + strings.Trim(strings.Join(parts, "/"), "/")
}

func PrintRegisteredRoutes(app *fiber.App) {

	fmt.Println()
	fmt.Println(cyan + "📦 Registered routes:" + reset)

	for _, route := range app.GetRoutes() {
		methodColor := methodToColor(route.Method)
		fmt.Printf("  %s%-6s%s %s\n", methodColor, route.Method, reset, route.Path)
	}
}

func LogRegisteredRoute(method, path string) {

	methodColor := methodToColor(method)

	fmt.Printf("  %s%-6s%s %s\n", methodColor, method, reset, path)
}

func methodToColor(method string) string {

	switch method {
	case "GET":
		return green
	case "POST":
		return yellow
	case "PUT", "PATCH":
		return cyan
	case "DELETE":
		return "\033[31m"
	default:
		return reset
	}
}

const (
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	reset  = "\033[0m"
)
