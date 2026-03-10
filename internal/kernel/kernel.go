/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Kernel
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package kernel

import (
	"context"
	"sync"

	"github.com/vayload/vayload/internal/vayload"
	"github.com/vayload/vayload/pkg/logger"

	// External dependencies
	"github.com/spf13/cobra"
)

type mapFlags struct {
	bucket *sync.Map
}

func NewCommandFlags(bucket map[string]any) *mapFlags {
	store := &sync.Map{}
	for k, v := range bucket {
		store.Store(k, v)
	}

	return &mapFlags{bucket: store}
}

func (f *mapFlags) GetString(name string, def string) string {
	if val, ok := f.bucket.Load(name); ok {
		if str, ok := val.(string); ok && str != "" {
			return str
		}
	}
	return def
}

func (f *mapFlags) GetBool(name string, def bool) bool {
	if val, ok := f.bucket.Load(name); ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return def
}

func (f *mapFlags) GetInt(name string, def int) int {
	if val, ok := f.bucket.Load(name); ok {
		if i, ok := val.(int); ok {
			return i
		}
	}
	return def
}

// CommandMeta is a struct that contains metadata about a command.
type CommandMeta struct {
	name        string
	description string
	flags       map[string]any
}

func NewCommandMeta(name, description string, flags map[string]any) *CommandMeta {
	return &CommandMeta{
		name:        name,
		description: description,
		flags:       flags,
	}
}

func (m *CommandMeta) Name() string          { return m.name }
func (m *CommandMeta) Description() string   { return m.description }
func (m *CommandMeta) Flags() map[string]any { return m.flags }

type CommandRunner struct{}

func NewCommandRunner() *CommandRunner {
	return &CommandRunner{}
}

func (r *CommandRunner) Run(ctx context.Context, cmd vayload.ConsoleCommand, args []string, flags vayload.CommandFlags) error {
	if v, ok := cmd.(vayload.ConsoleValidator); ok {
		if err := v.Validate(ctx, args); err != nil {
			return err
		}
	}

	if p, ok := cmd.(vayload.ConsolePreRunner); ok {
		if err := p.PreRun(ctx, args); err != nil {
			return err
		}
	}

	if err := cmd.Execute(ctx, args, flags); err != nil {
		return err
	}

	if p, ok := cmd.(vayload.ConsolePostRunner); ok {
		if err := p.PostRun(ctx, args); err != nil {
			return err
		}
	}

	return nil
}

type consoleKernel struct {
	mu       sync.Mutex
	registry vayload.Container
	events   vayload.EventBus
	commands map[string]vayload.ConsoleCommand
	runner   vayload.ConsoleRunner

	// Internal dependencies
	cobraRoot *cobra.Command
}

func NewConsoleKernel(registry vayload.Container, events vayload.EventBus) *consoleKernel {
	root := &cobra.Command{
		Use:   "vayload",
		Short: "Vayload command line interface",
		Long:  "Vayload CLI is a tool for managing vayload resources and services.",
	}

	root.CompletionOptions.DisableDefaultCmd = true

	return &consoleKernel{
		registry:  registry,
		events:    events,
		commands:  make(map[string]vayload.ConsoleCommand),
		runner:    NewCommandRunner(),
		cobraRoot: root,
	}
}

func (k *consoleKernel) Bootstrap(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return nil
	}

	cmd, ok := k.commands[args[0]]
	if !ok {
		return nil
	}

	return k.runner.Run(ctx, cmd, args[1:], NewCommandFlags(nil))
}

func (k *consoleKernel) RegisterCommand(command vayload.ConsoleCommand) {
	k.mu.Lock()
	defer k.mu.Unlock()

	root := &cobra.Command{
		Use:   command.Name(),
		Short: command.Description(),
	}

	flags := command.Flags()
	flagPtrs := toFlagMap(root, flags)

	if len(command.SubCommands()) > 0 {
		subcommands := command.SubCommands()
		for _, subcommand := range subcommands {
			if len(subcommand.SubCommands()) > 0 {
				logger.F(nil, logger.Fields{"context": "RegisterCommand", "command": subcommand.Name(), "error": "subcommands not supported"})
			}

			flags := subcommand.Flags()

			subroot := &cobra.Command{
				Use:   subcommand.Name(),
				Short: subcommand.Description(),
			}

			sflagPtrs := toFlagMap(subroot, flags)
			subroot.Run = func(cmd *cobra.Command, args []string) {
				flagValues := getFlagValues(sflagPtrs)
				err := subcommand.Execute(cmd.Context(), args, NewCommandFlags(flagValues))
				if err != nil {
					logger.F(err, logger.Fields{"context": "subcommand.Execute", "command": subcommand.Name()})
				}
			}

			root.AddCommand(subroot)

		}
	} else {
		root.Run = func(cmd *cobra.Command, args []string) {
			flagValues := getFlagValues(flagPtrs)
			err := command.Execute(cmd.Context(), args, NewCommandFlags(flagValues))
			if err != nil {
				logger.F(err, logger.Fields{"context": "command.Execute", "command": command.Name()})
			}
		}
	}

	k.cobraRoot.AddCommand(root)
}

func (k *consoleKernel) Name() string {
	return k.cobraRoot.Name()
}

func (k *consoleKernel) ShortDescription() string {
	return k.cobraRoot.Short
}

func (k *consoleKernel) LongDescription() string {
	return k.cobraRoot.Long
}

func (k *consoleKernel) Events() vayload.EventBus {
	return k.events
}

func (k *consoleKernel) Container() vayload.Container {
	return k.registry
}

type kernel struct {
	name     string
	registry vayload.Container
	events   vayload.EventBus
	services vayload.ServiceManager

	*KernelLifecycleDispatcher
}

func NewKernel(name string, registry vayload.Container, events vayload.EventBus, services vayload.ServiceManager) *kernel {
	return &kernel{
		name:                      name,
		registry:                  registry,
		events:                    events,
		services:                  services,
		KernelLifecycleDispatcher: NewKernelLifecycleDispatcher(),
	}
}

func (k *kernel) Name() string {
	return k.name
}

func (k *kernel) PublishEvent(event vayload.KernelEvent) {
	k.events.Publish(context.Background(), event)
}

func (k *kernel) SubscribeEvent(topic string, handler func(vayload.KernelEvent)) {
	k.events.Subscribe(context.Background(), topic, func(e vayload.Event) {
		if ke, ok := e.(vayload.KernelEvent); ok {
			handler(ke)
		}
	})
}

func (k *kernel) Events() vayload.EventBus {
	return k.events
}

func (k *kernel) Container() vayload.Container {
	return k.registry
}

func (k *kernel) Services() vayload.ServiceManager {
	return k.services
}

func (k *kernel) Bootstrap(ctx context.Context, args []string) error {
	k.KernelBootstrap(vayload.KernelBootstrapEvent{
		Context: ctx,
		Args:    args,
	})

	if err := k.services.StartAll(ctx); err != nil {
		return err
	}

	k.KernelStarted(vayload.KernelStartedEvent{
		Context: ctx,
	})

	return nil
}

func (k *kernel) Shutdown(ctx context.Context) error {
	k.KernelShutdown(vayload.KernelShutdownEvent{
		Context: ctx,
	})

	return k.services.StopAll(ctx)
}

func (k *kernel) OnServiceRegistered(l vayload.ServiceRegisteredListener) {
	k.services.OnServiceRegistered(l)
}

func (k *kernel) OnServiceStarting(l vayload.ServiceStartingListener) {
	k.services.OnServiceStarting(l)
}

func (k *kernel) OnServiceStarted(l vayload.ServiceStartedListener) {
	k.services.OnServiceStarted(l)
}

func (k *kernel) OnServiceStopping(l vayload.ServiceStoppingListener) {
	k.services.OnServiceStopping(l)
}

func (k *kernel) OnServiceStopped(l vayload.ServiceStoppedListener) {
	k.services.OnServiceStopped(l)
}

func (k *kernel) OnServiceError(l vayload.ServiceErrorListener) {
	k.services.OnServiceError(l)
}

var _ vayload.Kernel = (*kernel)(nil)
