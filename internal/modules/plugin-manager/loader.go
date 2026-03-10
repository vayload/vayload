/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Service Plugin Manager
 * Subpackage: Loader
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package plugin_manager

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	"github.com/vayload/vayload/internal/shared/ds"
	"github.com/vayload/vayload/pkg/logger"
	// External dependencies
)

type PluginSchema struct {
	Name         string                  `json:"name"`
	Version      string                  `json:"version"`
	Description  string                  `json:"description"`
	Author       PluginAuthor            `json:"author"`
	Runtime      PluginRuntimeConfig     `json:"runtime"`
	Host         PluginHostConfig        `json:"host"`
	Dependencies []string                `json:"dependencies"`
	Permissions  PluginPermissionsConfig `json:"permissions"`
	Exports      PluginExportsConfig     `json:"exports"`
	Config       map[string]any          `json:"config"`
}

type PluginAuthor struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Homepage string `json:"homepage"`
}

type PluginRuntimeConfig struct {
	Type   string `json:"type"`   // e.g., lua
	Engine string `json:"engine"` // e.g., 5.1, jit, 5.4
	Main   string `json:"main"`   // e.g., ./scripts/init.lua
}

type PluginHostConfig struct {
	APIVersion      string   `json:"api_version"`
	RequiredModules []string `json:"required_modules"`
}

type PluginPermissionsConfig struct {
	Filesystem         string   `json:"filesystem"` // e.g., read-only, none, read-write
	Network            []string `json:"network"`
	MaxMemoryMB        int      `json:"max_memory_mb"`
	MaxExecutionTimeMS int      `json:"max_execution_time_ms"`
}

type PluginExportsConfig struct {
	Hooks    []string `json:"hooks"`
	Commands []string `json:"commands"`
}

type pluginLoader struct {
	plugins *ds.HashMap[string, PluginSchema]
	path    string
}

func NewPluginLoader() *pluginLoader {
	return &pluginLoader{
		plugins: ds.NewHashMap[string, PluginSchema](),
		path:    "plugins",
	}
}

// Load scans the configured plugins directory for valid plugins.
// It walks through subdirectories and looks for a 'plugin.json' file.
func (pl *pluginLoader) Load() error {
	// !TODO: improve performance with best way to load plugins (dir scan or file watch)
	// case 1: load all plugins at once
	// case 2: load plugins on demand like nvim
	// case 3: parallel load plugins and use on demand with cache and compile to byte code
	return filepath.WalkDir(pl.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.E(err, logger.Fields{"path": path, "scope": "fs_walk"})
			return nil
		}

		if path == pl.path {
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		pluginDirName := d.Name()
		manifestPath := filepath.Join(path, "plugin.json")

		// Verification: Does plugin.json exist?
		if _, err := os.Stat(manifestPath); errors.Is(err, os.ErrNotExist) {
			// Log a warning if the directory structure is incorrect
			logger.W("Directory skipped: missing plugin.json", logger.Fields{
				"path":      path,
				"directory": pluginDirName,
			})
			return fs.SkipDir // Do not traverse deeper into this directory
		}

		// Prevent double registration
		if pl.plugins.Has(pluginDirName) {
			logger.W("Plugin skipped: already loaded", logger.Fields{"plugin": pluginDirName})
			return fs.SkipDir
		}

		plugin, err := pl.parsePlugin(path)
		if err != nil {
			// Log the error but do not stop the loader; just skip this bad plugin
			logger.E(err, logger.Fields{
				"path":   path,
				"plugin": pluginDirName,
				"scope":  "plugin_parse",
			})
			return fs.SkipDir
		}

		pl.plugins.Set(pluginDirName, *plugin)
		logger.I("Plugin loaded successfully", logger.Fields{
			"name":    plugin.Name,
			"version": plugin.Version,
		})

		// Skip deeper traversal as plugins are not nested
		return fs.SkipDir
	})
}

func (pl *pluginLoader) parsePlugin(path string) (*PluginSchema, error) {
	manifestPath := filepath.Join(path, "plugin.json")

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var plugin PluginSchema
	if err := json.Unmarshal(data, &plugin); err != nil {
		return nil, fmt.Errorf("invalid json syntax: %w", err)
	}

	if plugin.Name == "" {
		return nil, errors.New("field 'name' is required")
	}
	if plugin.Runtime.Main == "" {
		return nil, errors.New("field 'runtime.main' is required")
	}

	// Normalize the path to the main entry point
	plugin.Runtime.Main = filepath.Join(path, plugin.Runtime.Main)

	// Check entry point exists
	if _, err := os.Stat(plugin.Runtime.Main); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("entry point does not exist")
	}

	return &plugin, nil
}

// Dump prints the list of loaded plugins and their details to standard output.
// This is primarily used for debugging purposes.
func (pl *pluginLoader) Dump() {
	fmt.Println("=== Loaded Plugins ===")
	if pl.plugins.Size() == 0 {
		fmt.Println("No plugins loaded")
		return
	}
	pl.plugins.Range(func(id string, p PluginSchema) bool {
		fmt.Printf(
			" Name: %s\n  Version: %s\n  Runtime: %s (%s)\n  Entry: %s\n\n",
			p.Name, p.Version, p.Runtime.Type, p.Runtime.Engine, p.Runtime.Main,
		)
		return true
	})
}
