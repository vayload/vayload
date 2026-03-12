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

type KPIsResponse struct {
	TotalUsers  int64            `json:"total_users"`
	ActiveUsers int64            `json:"active_users"`
	ByPlatform  map[string]int64 `json:"by_platform"`
}

type WeeklyAnalyticsResponse struct {
	Labels   []string `json:"labels"`
	NewUsers []int64  `json:"new_users"`
	Active   []int64  `json:"active"`
}

type CountryCount struct {
	Country string `json:"country"`
	Count   int64  `json:"count"`
}

type UserActivityResponse struct {
	Labels       []string `json:"labels"`
	Interactions []int64  `json:"interactions"`
}

type UserLogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Result    string    `json:"result"`
	IPAddress *string   `json:"ip_address,omitempty"`
	UserAgent *string   `json:"user_agent,omitempty"`
}

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

type AnalyticsFilter struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
}

type AnalyticsRepository interface {
	GetKPIs(ctx context.Context) (*KPIsResponse, error)
	GetWeeklyAnalytics(ctx context.Context) (*WeeklyAnalyticsResponse, error)
	GetCountryRanking(ctx context.Context) ([]CountryCount, error)
	GetUserActivity(ctx context.Context, userID snowflake.ID, filter AnalyticsFilter) (*UserActivityResponse, error)
	GetUserLogs(ctx context.Context, userID snowflake.ID) ([]UserLogEntry, error)
	GetRecentUsers(ctx context.Context, limit int) ([]RecentUser, error)
}
