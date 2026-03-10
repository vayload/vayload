/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/shared/cache"
	"github.com/vayload/vayload/internal/shared/snowflake"
)

const (
	AUDIT_LOG_TABLE = "audit_log_entries"
	USERS_TABLE     = "users"

	// Cache keys
	CACHE_KEY_KPIS      = "analytics:kpis"
	CACHE_KEY_WEEKLY    = "analytics:weekly"
	CACHE_KEY_COUNTRIES = "analytics:countries"

	// Cache TTL
	CACHE_TTL = 1 * time.Minute
)

type AnalyticsRepository struct {
	database connection.DatabaseConnection
	cache    cache.Cache
}

func NewAnalyticsRepository(database connection.DatabaseConnection) *AnalyticsRepository {
	return &AnalyticsRepository{
		database: database,
		cache:    cache.NewLRUCache[string, any](100),
	}
}

// GetKPIs obtiene los KPIs globales
func (repo *AnalyticsRepository) GetKPIs(ctx context.Context) (*domain.KPIsResponse, error) {
	// Check cache
	if cached, ok := repo.cache.Get(ctx, CACHE_KEY_KPIS); ok {
		return cached.(*domain.KPIsResponse), nil
	}

	// Single query to get all KPIs at once
	query := `
		SELECT
			(SELECT COUNT(*) FROM ` + USERS_TABLE + `) as total_users,
			(SELECT COUNT(DISTINCT actor_id) FROM ` + AUDIT_LOG_TABLE + `
			 WHERE action = 'sign-in' AND created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)) as active_users`

	var kpiResult struct {
		TotalUsers  int64 `db:"total_users"`
		ActiveUsers int64 `db:"active_users"`
	}
	if err := repo.database.SelectOne(ctx, &kpiResult, query); err != nil {
		kpiResult.TotalUsers = 0
		kpiResult.ActiveUsers = 0
	}

	// By platform - single query
	platformQuery := `
		SELECT
			COALESCE(JSON_UNQUOTE(JSON_EXTRACT(payload, '$.method')), 'unknown') as platform,
			COUNT(DISTINCT actor_id) as count
		FROM ` + AUDIT_LOG_TABLE + `
		WHERE action = 'sign-in'
		GROUP BY platform`

	var platformResults []struct {
		Platform string `db:"platform"`
		Count    int64  `db:"count"`
	}
	byPlatform := make(map[string]int64)
	if err := repo.database.Select(ctx, &platformResults, platformQuery); err == nil {
		for _, p := range platformResults {
			byPlatform[p.Platform] = p.Count
		}
	}

	result := &domain.KPIsResponse{
		TotalUsers:  kpiResult.TotalUsers,
		ActiveUsers: kpiResult.ActiveUsers,
		ByPlatform:  byPlatform,
	}

	// Store in cache
	repo.cache.Set(ctx, CACHE_KEY_KPIS, result, CACHE_TTL)

	return result, nil
}

func (repo *AnalyticsRepository) GetWeeklyAnalytics(ctx context.Context) (*domain.WeeklyAnalyticsResponse, error) {
	if cached, ok := repo.cache.Get(ctx, CACHE_KEY_WEEKLY); ok {
		return cached.(*domain.WeeklyAnalyticsResponse), nil
	}

	dayNames := []string{"Dom", "Lun", "Mar", "Mié", "Jue", "Vie", "Sáb"}

	labels := make([]string, 7)
	newUsers := make([]int64, 7)
	active := make([]int64, 7)

	dateToIndex := make(map[string]int)
	for i := 0; i < 7; i++ {
		date := time.Now().AddDate(0, 0, -(6 - i))
		dateStr := date.Format("2006-01-02")
		dateToIndex[dateStr] = i
		labels[i] = dayNames[int(date.Weekday())]
	}

	newUsersQuery := `
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM ` + USERS_TABLE + `
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)
		GROUP BY DATE(created_at)`

	var newUserResults []struct {
		Date  time.Time `db:"date"`
		Count int64     `db:"count"`
	}
	if err := repo.database.Select(ctx, &newUserResults, newUsersQuery); err == nil {
		for _, r := range newUserResults {
			dateStr := r.Date.Format("2006-01-02")
			if idx, ok := dateToIndex[dateStr]; ok {
				newUsers[idx] = r.Count
			}
		}
	}

	activeQuery := `
		SELECT DATE(created_at) as date, COUNT(DISTINCT actor_id) as count
		FROM ` + AUDIT_LOG_TABLE + `
		WHERE action = 'sign-in' AND created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)
		GROUP BY DATE(created_at)`

	var activeResults []struct {
		Date  time.Time `db:"date"`
		Count int64     `db:"count"`
	}
	if err := repo.database.Select(ctx, &activeResults, activeQuery); err == nil {
		for _, r := range activeResults {
			dateStr := r.Date.Format("2006-01-02")
			if idx, ok := dateToIndex[dateStr]; ok {
				active[idx] = r.Count
			}
		}
	}

	result := &domain.WeeklyAnalyticsResponse{
		Labels:   labels,
		NewUsers: newUsers,
		Active:   active,
	}

	repo.cache.Set(ctx, CACHE_KEY_WEEKLY, result, CACHE_TTL)

	return result, nil
}

func (repo *AnalyticsRepository) GetCountryRanking(ctx context.Context) ([]domain.CountryCount, error) {
	if cached, ok := repo.cache.Get(ctx, CACHE_KEY_COUNTRIES); ok {
		return cached.([]domain.CountryCount), nil
	}

	query := `select c.name country, COUNT(u.id) as count FROM users u JOIN countries c ON u.country_id = c.id GROUP BY c.name ORDER BY count DESC limit 6`
	var results []struct {
		Country string `db:"country"`
		Count   int64  `db:"count"`
	}

	if err := repo.database.Select(ctx, &results, query); err != nil {
		return []domain.CountryCount{}, nil
	}

	countries := make([]domain.CountryCount, len(results))
	for i, r := range results {
		countries[i] = domain.CountryCount{
			Country: r.Country,
			Count:   r.Count,
		}
	}

	repo.cache.Set(ctx, CACHE_KEY_COUNTRIES, countries, CACHE_TTL)

	return countries, nil
}

func (repo *AnalyticsRepository) GetUserActivity(ctx context.Context, userID snowflake.ID, filter domain.AnalyticsFilter) (*domain.UserActivityResponse, error) {
	cacheKey := fmt.Sprintf("analytics:user:%d:activity", userID)

	if cached, ok := repo.cache.Get(ctx, cacheKey); ok {
		return cached.(*domain.UserActivityResponse), nil
	}

	from := time.Now().AddDate(0, 0, -30)
	to := time.Now()

	if filter.From != nil {
		from = *filter.From
	}
	if filter.To != nil {
		to = *filter.To
	}

	query := `
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM ` + AUDIT_LOG_TABLE + `
		WHERE actor_id = ? AND created_at BETWEEN ? AND ?
		GROUP BY DATE(created_at)
		ORDER BY date ASC`

	var results []struct {
		Date  time.Time `db:"date"`
		Count int64     `db:"count"`
	}

	if err := repo.database.Select(ctx, &results, query, userID, from, to); err != nil {
		return &domain.UserActivityResponse{
			Labels:       []string{},
			Interactions: []int64{},
		}, nil
	}

	labels := make([]string, len(results))
	interactions := make([]int64, len(results))

	for i, r := range results {
		labels[i] = r.Date.Format("2 Jan")
		interactions[i] = r.Count
	}

	result := &domain.UserActivityResponse{
		Labels:       labels,
		Interactions: interactions,
	}

	repo.cache.Set(ctx, cacheKey, result, CACHE_TTL)

	return result, nil
}

func (repo *AnalyticsRepository) GetUserLogs(ctx context.Context, userID snowflake.ID) ([]domain.UserLogEntry, error) {
	cacheKey := fmt.Sprintf("analytics:user:%d:logs", userID)

	if cached, ok := repo.cache.Get(ctx, cacheKey); ok {
		return cached.([]domain.UserLogEntry), nil
	}

	query := `
		SELECT
			created_at,
			action,
			ip_address,
			user_agent
		FROM ` + AUDIT_LOG_TABLE + `
		WHERE actor_id = ?
		ORDER BY created_at DESC
		LIMIT 100`

	var results []struct {
		CreatedAt time.Time `db:"created_at"`
		Action    string    `db:"action"`
		IPAddress *string   `db:"ip_address"`
		UserAgent *string   `db:"user_agent"`
	}

	if err := repo.database.Select(ctx, &results, query, userID); err != nil {
		return []domain.UserLogEntry{}, nil
	}

	logs := make([]domain.UserLogEntry, len(results))
	for i, r := range results {
		logs[i] = domain.UserLogEntry{
			Timestamp: r.CreatedAt,
			Action:    r.Action,
			Result:    "Success",
			IPAddress: r.IPAddress,
			UserAgent: r.UserAgent,
		}
	}

	repo.cache.Set(ctx, cacheKey, logs, CACHE_TTL)

	return logs, nil
}

func (repo *AnalyticsRepository) GetRecentUsers(ctx context.Context, limit int) ([]domain.RecentUser, error) {
	cacheKey := fmt.Sprintf("analytics:recent_users:%d", limit)

	if cached, ok := repo.cache.Get(ctx, cacheKey); ok {
		return cached.([]domain.RecentUser), nil
	}

	query := `SELECT
        u.id,
        u.first_name as name,
        u.email,
        if(u.auth_type = "" or u.auth_type is null, "vayload", u.auth_type) platform,
        COALESCE(c.name, "Sin asignar") country,
        COALESCE(p.name, "Sin asignar") city,
        -- CASE WHEN u.email_confirmed_at IS NOT NULL THEN 'verified' ELSE 'pending' END as status,
        IF(u.email_confirmed_at IS NOT NULL OR u.phone_confirmed_at IS NOT NULL, 'verified', 'pending') as status,
        COALESCE(al.last_login, u.created_at) as last_login
    FROM users u
        left join user_addresses ua on ua.user_id = u.id
        LEFT JOIN countries c ON ua.country_id = c.id
        LEFT JOIN provinces p ON ua.province_id = p.id
        left join (
            select actor_id, max(created_at) last_login from audit_log_entries where action = 'sign-in' group by 1
        ) al on al.actor_id = u.id
    where u.deleted_at is null
		LIMIT ?`

	var results []struct {
		ID        snowflake.ID `db:"id"`
		Name      string       `db:"name"`
		Email     string       `db:"email"`
		Platform  string       `db:"platform"`
		Country   string       `db:"country"`
		City      string       `db:"city"`
		Status    string       `db:"status"`
		LastLogin time.Time    `db:"last_login"`
	}

	if err := repo.database.Select(ctx, &results, query, limit); err != nil {
		return []domain.RecentUser{}, nil
	}

	users := make([]domain.RecentUser, len(results))

	for i, r := range results {
		users[i] = domain.RecentUser{
			ID:        r.ID,
			Name:      r.Name,
			Email:     r.Email,
			Platform:  r.Platform,
			Country:   r.Country,
			City:      r.City,
			Status:    r.Status,
			LastLogin: r.LastLogin,
		}
	}

	repo.cache.Set(ctx, cacheKey, users, CACHE_TTL)

	return users, nil
}

var _ domain.AnalyticsRepository = (*AnalyticsRepository)(nil)
