/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Server
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	server_transports "github.com/vayload/vayload/cmd/server/transports"
	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/pkg/logger"
)

func main() {
	config, err := config.GetConfig("config.toml")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	logger.Init(logger.Config{
		Level:      logger.ParseLevel(config.App.LogLevel),
		FilePath:   "./logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
		Console:    true,
		TimeFormat: time.RFC3339,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create servers
	httpServer := server_transports.CreateHttpServer(config)
	grpcServer := server_transports.CreateGrpcServer(config)
	mcpServer := server_transports.CreateMCPServer(config, httpServer.Server())

	if err := httpServer.Serve(); err != nil {
		logger.F(err, logger.Fields{"context": "http_serve"})
	}

	// Serve grpc server only is enabled by enviroment
	if config.GRPC.Enabled {
		if err := grpcServer.Serve(); err != nil {
			logger.F(err, logger.Fields{"context": "grpc_serve"})
		}
	}

	// Server mcp server only is enabled by enviroment
	if config.MCP.Enabled {
		if err := mcpServer.Serve(); err != nil {
			logger.F(err, logger.Fields{"context": "mcp_serve"})
		}
	}

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	logger.I("got shutdown signal, shutting down server...")

	localCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	servers := []interface {
		Shutdown(context.Context) error
		IsListening() bool
	}{
		httpServer,
		grpcServer,
		mcpServer,
	}

	// Close if servers is listening
	for _, srv := range servers {
		if srv.IsListening() {
			if err := srv.Shutdown(localCtx); err != nil {
				logger.F(err, logger.Fields{"context": "http_shutdown"})
			}
		}
	}

	logger.I("server shutdown complete")
}
