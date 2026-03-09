/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Config
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

type AppMode string

const (
	// App serves both API endpoints and frontend static assets (SPA).
	AppModeFull AppMode = "full"

	// App serves API endpoints only. Frontend is hosted externally.
	AppModeAPI AppMode = "api"
)

type Config struct {
	App       AppConfig       `toml:"app"`
	HTTP      HTTPConfig      `toml:"http"`
	GRPC      GRPCConfig      `toml:"grpc"`
	MCP       MCPConfig       `toml:"mcp"`
	Cors      CorsConfig      `toml:"cors"`
	Database  DatabaseConfig  `toml:"database"`
	RateLimit RateLimitConfig `toml:"rate_limit"`
	Logging   LoggingConfig   `toml:"logging"`
	Security  SecurityConfig  `toml:"security"`
	MagicLink MagicLinkConfig `toml:"magic_link"`
	OAuth     OAuthConfig     `toml:"oauth"`
	Queue     QueueConfig     `toml:"queue"`
	Plugins   PluginsConfig   `toml:"plugins"`

	ConfigPath string `toml:"-"`
	WorkDir    string `toml:"-"`
	DataDir    string `toml:"-"`
}

type AppConfig struct {
	Env        string  `toml:"env"`
	WorkingDir string  `toml:"working_dir"`
	Domain     string  `toml:"domain"`
	SecretKey  string  `toml:"secret_key"`
	APIKey     string  `toml:"api_key"`
	Mode       AppMode `toml:"mode"`
	LogLevel   string  `toml:"log_level"`
}

type CorsConfig struct {
	Origins []string `toml:"origins"`
}

type HTTPConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type GRPCConfig struct {
	Enabled bool   `toml:"enabled"`
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
}

type MCPConfig struct {
	Enabled bool `toml:"enabled"`
}

type DatabaseConfig struct {
	Driver         string `toml:"driver"`
	User           string `toml:"user"`
	Password       string `toml:"password"`
	Host           string `toml:"host"`
	Port           int    `toml:"port"`
	Schema         string `toml:"schema"`
	MigrationsPath string `toml:"migrations_path"`
}

type RateLimitConfig struct {
	WindowMs int `toml:"window_ms"`
	Max      int `toml:"max"`
}

type LoggingConfig struct {
	Level string `toml:"level"`
}

type SecurityConfig struct {
	JwtPublicKeyBase64       string `toml:"jwt_public_key"`
	JwtPrivateKeyBase64      string `toml:"jwt_private_key"`
	JwtExpirationTime        int    `toml:"jwt_expiration_time"`
	JwtRefreshExpirationDays int    `toml:"jwt_refresh_expiration_days"`

	JwtPublicKey  []byte `toml:"-"`
	JwtPrivateKey []byte `toml:"-"`
}

type MagicLinkConfig struct {
	SecretKey        string `toml:"secret_key"`
	ExpiresInMinutes int    `toml:"expires_in_minutes"`
}

type OAuthProviderConfig struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	RedirectURL  string `toml:"redirect_url"`
}

type OAuthConfig struct {
	RedirectBase string              `toml:"redirect_base"`
	Google       OAuthProviderConfig `toml:"google"`
	Facebook     OAuthProviderConfig `toml:"facebook"`
	Apple        OAuthProviderConfig `toml:"apple"`
}

type QueueConfig struct {
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
	APIKey string `toml:"api_key"`
}

type PluginsConfig struct {
	RegistryURI string `toml:"registry_uri"`
}

func (c *Config) SetConfigPath(path string) {
	c.ConfigPath = path
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	cfg.SetConfigPath(path)

	if cfg.Security.JwtPublicKeyBase64 != "" {
		keyBytes, err := base64.StdEncoding.DecodeString(cfg.Security.JwtPublicKeyBase64)
		if err != nil {
			return nil, fmt.Errorf("failed to decode jwt_public_key: %w", err)
		}
		cfg.Security.JwtPublicKey = keyBytes
	}

	if cfg.Security.JwtPrivateKeyBase64 != "" {
		keyBytes, err := base64.StdEncoding.DecodeString(cfg.Security.JwtPrivateKeyBase64)
		if err != nil {
			return nil, fmt.Errorf("failed to decode jwt_private_key: %w", err)
		}
		cfg.Security.JwtPrivateKey = keyBytes
	}

	return &cfg, nil
}

var (
	configOnce sync.Once
	config     *Config
)

func GetConfig(path string) (*Config, error) {
	var err error
	configOnce.Do(func() {
		workdir := os.Getenv("WORKDIR")
		if workdir == "" {
			workdir = "."
		}
		config, err = LoadConfig(filepath.Join(workdir, path))
		if err != nil {
			return
		}

		config.WorkDir = workdir
		config.DataDir = filepath.Join(workdir, "data")
	})

	return config, err
}

func MustConfig(path string) *Config {
	config, err := GetConfig(path)
	if err != nil {
		panic(err)
	}

	return config
}
