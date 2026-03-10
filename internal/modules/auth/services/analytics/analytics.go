/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package analytics

import (
	"context"
	"time"

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/shared/snowflake"
)

const SERVICE_KEY = "auth.analytics.service"

type AnalyticsService struct {
	repository domain.AnalyticsRepository
}

func NewAnalyticsService(repository domain.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{
		repository: repository,
	}
}

// GetKPIs obtiene los KPIs globales
func (s *AnalyticsService) GetKPIs(ctx context.Context) (*domain.KPIsResponse, error) {
	return s.repository.GetKPIs(ctx)
}

// GetWeeklyAnalytics obtiene la actividad semanal
func (s *AnalyticsService) GetWeeklyAnalytics(ctx context.Context) (*domain.WeeklyAnalyticsResponse, error) {
	return s.repository.GetWeeklyAnalytics(ctx)
}

// GetCountryRanking obtiene el ranking de usuarios por país
func (s *AnalyticsService) GetCountryRanking(ctx context.Context) ([]domain.CountryCount, error) {
	return s.repository.GetCountryRanking(ctx)
}

// UserActivityInput input para obtener actividad de usuario
type UserActivityInput struct {
	UserID snowflake.ID
	From   *string
	To     *string
}

// GetUserActivity obtiene la actividad de un usuario específico
func (s *AnalyticsService) GetUserActivity(ctx context.Context, input UserActivityInput) (*domain.UserActivityResponse, error) {
	filter := domain.AnalyticsFilter{}

	if input.From != nil {
		if t, err := time.Parse("2006-01-02", *input.From); err == nil {
			filter.From = &t
		}
	}

	if input.To != nil {
		if t, err := time.Parse("2006-01-02", *input.To); err == nil {
			// Agregar fin del día
			t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filter.To = &t
		}
	}

	return s.repository.GetUserActivity(ctx, input.UserID, filter)
}

// GetUserLogs obtiene los logs de un usuario específico
func (s *AnalyticsService) GetUserLogs(ctx context.Context, userID snowflake.ID) ([]domain.UserLogEntry, error) {
	return s.repository.GetUserLogs(ctx, userID)
}

// GetRecentUsers obtiene los últimos usuarios logueados
func (s *AnalyticsService) GetRecentUsers(ctx context.Context, limit int) ([]domain.RecentUser, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repository.GetRecentUsers(ctx, limit)
}
