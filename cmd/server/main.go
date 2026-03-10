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
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	server_transports "github.com/vayload/vayload/cmd/server/transports"
	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/kernel"
	"github.com/vayload/vayload/internal/modules/auth"
	"github.com/vayload/vayload/internal/modules/database"
	"github.com/vayload/vayload/internal/modules/storage"
	"github.com/vayload/vayload/internal/vayload"
	"github.com/vayload/vayload/pkg/logger"
)

const (
	shutdownTimeout = 10 * time.Second
	logMaxSize      = 100
	logMaxBackups   = 3
	logMaxAge       = 28
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}

func run() error {
	cfg, err := config.GetConfig("config.toml")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger.Init(logger.Config{
		Level:      logger.ParseLevel(cfg.App.LogLevel),
		FilePath:   filepath.Join(cfg.App.WorkingDir, "logs", "app.log"),
		MaxSize:    logMaxSize,
		MaxBackups: logMaxBackups,
		MaxAge:     logMaxAge,
		Compress:   true,
		Console:    true,
		TimeFormat: time.RFC3339,
	})

	// Single context for the whole app lifetime.
	appContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	container := kernel.NewContainer(appContext)
	eventsBus := kernel.NewEventBus(10)
	services := kernel.NewServiceManager(container, eventsBus)

	services.RegisterService(database.NewDatabaseService(cfg))
	services.RegisterService(storage.NewStorageService(cfg))
	services.RegisterService(auth.NewAuthService(cfg))

	vkernel := kernel.NewKernel("vayload-server", container, eventsBus, services)

	httpServer := server_transports.CreateHttpServer(cfg)
	vkernel.OnServiceStarted(httpServer)

	active := []vayload.Transport{httpServer}

	if cfg.GRPC.Enabled {
		grpcServer := server_transports.CreateGrpcServer(cfg)
		active = append(active, grpcServer)
		vkernel.OnServiceStarted(grpcServer)
	}

	if cfg.MCP.Enabled {
		// MCP piggybacks on the HTTP server — pass it after HTTP is built.
		active = append(active, server_transports.CreateMCPServer(cfg, httpServer.Server()))
	}

	// ============================ RUN KERNEL BOOTSTRAP ==============================
	if err := vkernel.Bootstrap(appContext, nil); err != nil {
		return fmt.Errorf("kernel bootstrap: %w", err)
	}

	// ============================ START TRANSPORTS ==============================
	// Start all active transports.
	for _, srv := range active {
		// Capture loop variable for the goroutine.
		srv := srv
		if err := srv.Serve(); err != nil {
			return fmt.Errorf("serve: %w", err)
		}
	}

	// ============================ WAIT FOR SHUTDOWN SIGNAL ==============================
	<-appContext.Done()
	logger.I("shutdown signal received, draining…")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	var errs []error

	for _, srv := range active {
		if !srv.IsListening() {
			continue
		}
		if err := srv.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown transport: %w", err))
		}
	}

	container.Flush()
	vkernel.Shutdown(appContext)

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	logger.I("shutdown complete")
	return nil
}
