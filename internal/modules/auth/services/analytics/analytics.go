/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Services/Analytics
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package analytics

import (
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

func (s *AnalyticsService) GetKPIs(ctx interface{}) (*domain.KPIsResponse, error) {
	return nil, nil
}

func (s *AnalyticsService) GetWeeklyAnalytics(ctx interface{}) (*domain.WeeklyAnalyticsResponse, error) {
	return nil, nil
}

func (s *AnalyticsService) GetCountryRanking(ctx interface{}) ([]domain.CountryCount, error) {
	return nil, nil
}

type UserActivityInput struct {
	UserID snowflake.ID
	From   *string
	To     *string
}

func (s *AnalyticsService) GetUserActivity(ctx interface{}, input UserActivityInput) (*domain.UserActivityResponse, error) {
	return nil, nil
}

func (s *AnalyticsService) GetUserLogs(ctx interface{}, userID snowflake.ID) ([]domain.UserLogEntry, error) {
	return nil, nil
}

func (s *AnalyticsService) GetRecentUsers(ctx interface{}, limit int) ([]domain.RecentUser, error) {
	return nil, nil
}

var _ interface{} = (*AnalyticsService)(nil)
