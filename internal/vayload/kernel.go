/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Kernel
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package vayload

import "context"

type KernelEvent interface {
	Name() string
	Service() string
}

type PublicKernel interface {
	PublishEvent(event KernelEvent)
	SubscribeEvent(topic string, handler func(KernelEvent))

	Events() EventBus
	Container() Container
	Services() ServiceManager
}

// Base kernel
type Kernel interface {
	Name() string
	PublicKernel

	// Bostrap this kernel and return error when corrured
	Bootstrap(ctx context.Context, args []string) error
}

// For work with command flags
type CommandFlags interface {
	GetString(name string, def string) string
	GetBool(name string, def bool) bool
	GetInt(name string, def int) int
}

// Command defines the interface for CLI commands
type ConsoleCommand interface {
	// Name returns the command name (used by CLI)
	Name() string

	// Description returns a short description of the command
	Description() string

	// Execute runs the command logic
	Execute(ctx context.Context, args []string, flags CommandFlags) error

	// Get Subcommands
	SubCommands() []ConsoleCommand

	// Flags returns a map of flag names to default values or types
	// This can be used to auto-register flags with Cobra or another flag parser
	Flags() map[string]any
}

type ConsoleValidator interface {
	// Validate checks preconditions or arguments before Execute
	Validate(ctx context.Context, args []string) error
}

type ConsolePreRunner interface {
	// PreRun executes logic before the main Execute (optional)
	PreRun(ctx context.Context, args []string) error
}

type ConsolePostRunner interface {
	// PostRun executes logic after the main Execute (optional)
	PostRun(ctx context.Context, args []string) error
}

type ConsoleRunner interface {
	Run(ctx context.Context, cmd ConsoleCommand, args []string, flags CommandFlags) error
}

// Console kernel for cli tool -> vayload [args]
type ConsoleKernel interface {
	Kernel
	LongDescription() string
	RegisterCommand(command ConsoleCommand)
	ShortDescription() string
}

// Http kernel uses for register http routes (rest, graphql)
type HttpKernel interface {
	Kernel
	RegisterRoute()
}
