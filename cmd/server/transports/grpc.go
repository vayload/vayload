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
	"net"
	"sync/atomic"

	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/vayload"
	"github.com/vayload/vayload/pkg/logger"

	// Transports dependencies
	"google.golang.org/grpc"
)

type GrpcTransport struct {
	listener net.Listener
	server   *grpc.Server

	config      *config.Config
	isListening atomic.Bool
}

func CreateGrpcServer(config *config.Config) *GrpcTransport {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GRPC.Host, config.GRPC.Port))
	if err != nil {
		logger.F(err, logger.Fields{"context": "grpc_listen"})
	}

	grpcServer := grpc.NewServer()

	return &GrpcTransport{
		listener: listen,
		server:   grpcServer,
		config:   config,
	}
}

func (t *GrpcTransport) Serve() error {
	go func() {
		logger.I("gRPC server started", logger.Fields{"port": t.config.GRPC.Port})
		if err := t.server.Serve(t.listener); err != nil {
			logger.F(err, logger.Fields{"context": "grpc_serve"})
		}

		t.isListening.Store(true)
	}()
	return nil
}

func (t *GrpcTransport) Shutdown(ctx context.Context) error {
	t.server.GracefulStop()
	t.isListening.Store(false)
	return nil
}

func (t *GrpcTransport) RegisterServices(services []vayload.GrpcServiceDescriptor, version string) {
	// for _, service := range services {
	// 	t.server.RegisterService(service.Service, service.Handler)
	// }
}

func (t *GrpcTransport) IsListening() bool {
	return t.isListening.Load()
}

var _ vayload.GrpcTransport = (*GrpcTransport)(nil)
