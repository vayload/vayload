/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Domain
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package domain

import (
	"context"
	"time"

	"github.com/vayload/vayload/internal/shared/snowflake"
)

// KPIsResponse representa los KPIs globales del dashboard
type KPIsResponse struct {
	TotalUsers  int64            `json:"total_users"`
	ActiveUsers int64            `json:"active_users"`
	ByPlatform  map[string]int64 `json:"by_platform"`
}

// WeeklyAnalyticsResponse representa la actividad semanal
type WeeklyAnalyticsResponse struct {
	Labels   []string `json:"labels"`
	NewUsers []int64  `json:"new_users"`
	Active   []int64  `json:"active"`
}

// CountryCount representa el conteo de usuarios por país
type CountryCount struct {
	Country string `json:"country"`
	Count   int64  `json:"count"`
}

// UserActivityResponse representa la actividad de un usuario
type UserActivityResponse struct {
	Labels       []string `json:"labels"`
	Interactions []int64  `json:"interactions"`
}

// UserLogEntry representa una entrada de log del usuario
type UserLogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Result    string    `json:"result"`
	IPAddress *string   `json:"ip_address,omitempty"`
	UserAgent *string   `json:"user_agent,omitempty"`
}

// RecentUser representa un usuario reciente con su último login
type RecentUser struct {
	ID        snowflake.ID `json:"id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Platform  string       `json:"platform"`
	Country   string       `json:"country"`
	City      string       `json:"city"`
	Status    string       `json:"status"`
	LastLogin time.Time    `json:"last_login"`
}

// AnalyticsFilter filtros para consultas de analytics
type AnalyticsFilter struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
}

// AnalyticsRepository interface para acceso a datos de analytics
type AnalyticsRepository interface {
	// GetKPIs obtiene los KPIs globales
	GetKPIs(ctx context.Context) (*KPIsResponse, error)

	// GetWeeklyAnalytics obtiene la actividad semanal
	GetWeeklyAnalytics(ctx context.Context) (*WeeklyAnalyticsResponse, error)

	// GetCountryRanking obtiene el ranking de usuarios por país
	GetCountryRanking(ctx context.Context) ([]CountryCount, error)

	// GetUserActivity obtiene la actividad de un usuario específico
	GetUserActivity(ctx context.Context, userID snowflake.ID, filter AnalyticsFilter) (*UserActivityResponse, error)

	// GetUserLogs obtiene los logs de un usuario específico
	GetUserLogs(ctx context.Context, userID snowflake.ID) ([]UserLogEntry, error)

	// GetRecentUsers obtiene los últimos usuarios logueados
	GetRecentUsers(ctx context.Context, limit int) ([]RecentUser, error)
}
