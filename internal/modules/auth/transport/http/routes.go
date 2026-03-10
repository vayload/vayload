/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package auth_http

import (
	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/vayload"
	httpi "github.com/vayload/vayload/pkg/http"
)

type Handler interface {
	// Start with Identify and Login methods
	Identify(req vayload.HttpRequest, res vayload.HttpResponse) error

	// Generating methods for login
	GenerateOtpCode(req vayload.HttpRequest, res vayload.HttpResponse) error

	// Authentication methods
	LoginWithPassword(req vayload.HttpRequest, res vayload.HttpResponse) error
	LoginWithOtpCode(req vayload.HttpRequest, res vayload.HttpResponse) error
	LoginWithOAuth2(req vayload.HttpRequest, res vayload.HttpResponse) error

	// Handle OAuth2 callback
	HandleOAuth2Callback(req vayload.HttpRequest, res vayload.HttpResponse) error

	GetUserPermissions(req vayload.HttpRequest, res vayload.HttpResponse) error

	RegisterUser(req vayload.HttpRequest, res vayload.HttpResponse) error
	ValidateRegister(req vayload.HttpRequest, res vayload.HttpResponse) error
	SendEmailVerificationCode(req vayload.HttpRequest, res vayload.HttpResponse) error

	// Password recovery
	ForgotPasswordRequest(req vayload.HttpRequest, res vayload.HttpResponse) error
	ForgotPasswordReset(req vayload.HttpRequest, res vayload.HttpResponse) error

	// Registration email change
	RegisterChangeEmailRequest(req vayload.HttpRequest, res vayload.HttpResponse) error
	RegisterChangeEmailConfirm(req vayload.HttpRequest, res vayload.HttpResponse) error

	Logout(req vayload.HttpRequest, res vayload.HttpResponse) error
	RefreshToken(req vayload.HttpRequest, res vayload.HttpResponse) error
	GetCurrentUser(req vayload.HttpRequest, res vayload.HttpResponse) error

	// Internal setup method for initializing the auth module (used for vayload services)
	InternalSetup(req vayload.HttpRequest, res vayload.HttpResponse) error

	// Analytics methods
	GetAnalyticsKPIs(req vayload.HttpRequest, res vayload.HttpResponse) error
	GetAnalyticsWeekly(req vayload.HttpRequest, res vayload.HttpResponse) error
	GetAnalyticsCountries(req vayload.HttpRequest, res vayload.HttpResponse) error
	GetRecentUsers(req vayload.HttpRequest, res vayload.HttpResponse) error
	GetUserActivity(req vayload.HttpRequest, res vayload.HttpResponse) error
	GetUserLogs(req vayload.HttpRequest, res vayload.HttpResponse) error
}

func RegisterRoutes(handler Handler, _container vayload.Container, config *config.Config) []httpi.HttpRoute {
	routes := []httpi.HttpRoute{
		// Identify route
		{
			Method:  httpi.POST,
			Path:    "/auth/identify",
			Handler: handler.Identify,
		},

		// Authentication routes
		{
			Method:  httpi.POST,
			Path:    "/auth/password",
			Handler: handler.LoginWithPassword,
		},
		{
			Method:  httpi.POST,
			Path:    "/auth/oauth2",
			Handler: handler.LoginWithOAuth2,
		},
		{
			Method:  httpi.GET,
			Path:    "/auth/:provider/callback",
			Handler: handler.HandleOAuth2Callback,
		},
		{
			Method:  httpi.POST,
			Path:    "/auth/:provider/callback", // Apple uses POST form_post
			Handler: handler.HandleOAuth2Callback,
		},
		{
			Method:  httpi.POST,
			Path:    "/auth/otp",
			Handler: handler.LoginWithOtpCode,
		},

		// User permissions
		{
			Method:  httpi.GET,
			Path:    "/auth/permissions",
			Handler: handler.GetUserPermissions,
			Middleware: []vayload.HttpHandler{
				AuthHttpMiddleware(_container, config),
			},
		},

		// Generate methods for login
		{
			Method:  httpi.POST,
			Path:    "/auth/otp/generate",
			Handler: handler.GenerateOtpCode,
		},

		// Registration routes
		{
			Method:  httpi.POST,
			Path:    "/auth/register",
			Handler: handler.RegisterUser,
		},
		{
			Method:  httpi.POST,
			Path:    "/auth/register/validate-account", // This method is used to validate the account after registration
			Handler: handler.ValidateRegister,
		},

		{
			Method:  httpi.POST,
			Path:    "/auth/register/send-verification-code",
			Handler: handler.SendEmailVerificationCode,
		},

		// Forgot password routes (public)
		{
			Method:  httpi.POST,
			Path:    "/auth/forgot-password/request",
			Handler: handler.ForgotPasswordRequest,
		},
		{
			Method:  httpi.POST,
			Path:    "/auth/forgot-password/reset",
			Handler: handler.ForgotPasswordReset,
		},

		// Registration email change routes (pre-verificación)
		{
			Method:  httpi.POST,
			Path:    "/auth/register/change-email/request",
			Handler: handler.RegisterChangeEmailRequest,
		},
		{
			Method:  httpi.POST,
			Path:    "/auth/register/change-email/confirm",
			Handler: handler.RegisterChangeEmailConfirm,
		},

		// Logout route
		{
			Method:  httpi.POST,
			Path:    "/auth/logout",
			Handler: handler.Logout,
		},

		// Refresh token route
		{
			Method:  httpi.POST,
			Path:    "/auth/refresh-token",
			Handler: handler.RefreshToken,
		},

		{
			Method:  httpi.GET,
			Path:    "/auth/me",
			Handler: handler.GetCurrentUser,
			Middleware: []vayload.HttpHandler{
				AuthHttpMiddleware(_container, config),
			},
		},

		// ============================== INTERNAL ROUTES ==============================
		// These routes are for internal services use only and should not be exposed to the public
		{
			Method:  httpi.POST,
			Path:    "/auth/__internal__/setup",
			Handler: handler.InternalSetup,
		},

		// ============================== ANALYTICS ROUTES ==============================
		// These routes require admin authentication via X-Auth-Token header
		{
			Method:  httpi.GET,
			Path:    "/auth/analytics/kpis",
			Handler: handler.GetAnalyticsKPIs,
		},
		{
			Method:  httpi.GET,
			Path:    "/auth/analytics/weekly",
			Handler: handler.GetAnalyticsWeekly,
		},
		{
			Method:  httpi.GET,
			Path:    "/auth/analytics/countries",
			Handler: handler.GetAnalyticsCountries,
		},
		{
			Method:  httpi.GET,
			Path:    "/auth/analytics/users/recent",
			Handler: handler.GetRecentUsers,
		},
		{
			Method:  httpi.GET,
			Path:    "/auth/analytics/users/:id/activity",
			Handler: handler.GetUserActivity,
		},
		{
			Method:  httpi.GET,
			Path:    "/auth/analytics/users/:id/logs",
			Handler: handler.GetUserLogs,
		},
	}

	return routes
}
