/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package auth_http

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/auth/services/analytics"
	"github.com/vayload/vayload/internal/modules/auth/services/login"
	"github.com/vayload/vayload/internal/modules/auth/services/recovery"
	"github.com/vayload/vayload/internal/modules/auth/services/registration"
	"github.com/vayload/vayload/internal/shared/snowflake"
	"github.com/vayload/vayload/internal/vayload"
	httpi "github.com/vayload/vayload/pkg/http"
	"github.com/vayload/vayload/pkg/logger"

	// Shared errors
	"github.com/vayload/vayload/internal/shared/errors"
)

type AuthHttpHandler struct {
	loginService     *login.LoginService
	registerService  *registration.RegisterService
	recoveryService  *recovery.RecoveryService
	analyticsService *analytics.AnalyticsService

	config   *config.Config
	registry vayload.Container
}

func (h *AuthHttpHandler) Identify(req vayload.HttpRequest, res vayload.HttpResponse) error {
	input := login.IdentifyInput{}
	if err := req.ValidateBody(&input); err != nil {
		return httpi.ErrWrapping(err)
	}

	authStep, err := h.loginService.IdentifyUserForFactor(req.Context(), input)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[*login.AuthStepCredentials]{
		Status: "success",
		Data:   authStep,
	})
}

func (h *AuthHttpHandler) LoginWithPassword(req vayload.HttpRequest, res vayload.HttpResponse) error {
	input := login.LoginInput{}
	if err := req.ValidateBody(&input); err != nil {
		return httpi.ErrWrapping(err)
	}

	authContext := domain.AuthContext{
		UserAgent: req.GetUserAgent(),
		IP:        req.GetIP(),
	}
	session, err := h.loginService.LoginWithPassword(req.Context(), input, authContext)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return h.withAuthCookies(session, res)
}

func (h *AuthHttpHandler) LoginWithOtpCode(req vayload.HttpRequest, res vayload.HttpResponse) error {
	input := login.OtpCodeInput{}
	if err := req.ParseBody(&input); err != nil {
		return httpi.ErrBadRequest(err)
	}

	authContext := domain.AuthContext{
		IP:        req.GetIP(),
		UserAgent: req.GetUserAgent(),
	}

	session, err := h.loginService.LoginWithOtpCode(req.Context(), input, authContext)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return h.withAuthCookies(session, res)
}

func (h *AuthHttpHandler) GenerateOtpCode(req vayload.HttpRequest, res vayload.HttpResponse) error {
	dto := login.OtpCodeGenInput{}
	if err := req.ParseBody(&dto); err != nil {
		return httpi.ErrBadRequest(err)
	}

	err := h.loginService.GenerateOtpCode(req.Context(), dto)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[string]{
		Status: "success",
		Data:   "OTP code sent",
	})
}

func (h *AuthHttpHandler) RegisterUser(req vayload.HttpRequest, res vayload.HttpResponse) error {
	input := registration.RegisterInput{}
	if err := req.ParseBody(&input); err != nil {
		return httpi.ErrBadRequest(err)
	}

	response, err := h.registerService.RegisterUser(req.Context(), input)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[*domain.User]{
		Status: "success",
		Data:   response,
	})
}

func (h *AuthHttpHandler) ValidateRegister(req vayload.HttpRequest, res vayload.HttpResponse) error {
	input := registration.RegisterValidationInput{}
	if err := req.ParseBody(&input); err != nil {
		return httpi.ErrBadRequest(err)
	}

	session, err := h.registerService.ValidateRegister(req.Context(), input)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return h.withAuthCookies(session, res)
}

func (h *AuthHttpHandler) SendEmailVerificationCode(req vayload.HttpRequest, res vayload.HttpResponse) error {
	input := registration.RegisterValidationInput{}
	if err := req.ParseBody(&input); err != nil {
		return httpi.ErrBadRequest(err)
	}

	err := h.registerService.SendEmailVerificationCode(req.Context(), input)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[string]{
		Status: "success",
		Data:   "Verification code sent to email",
	})
}

func (h *AuthHttpHandler) Logout(req vayload.HttpRequest, res vayload.HttpResponse) error {
	cookies := FlushAuthCookies(h.config)

	return res.Status(200).Cookies(cookies...).Json(&httpi.ResponseBody[string]{
		Status: "success",
		Data:   "Logged out successfully",
	})
}

func (h *AuthHttpHandler) RefreshToken(req vayload.HttpRequest, res vayload.HttpResponse) error {
	refreshToken := req.GetCookie(COOKIE_REFRESH_TOKEN)
	if refreshToken == "" {
		return httpi.ErrUnauthorized(errors.New("missing refresh token"))
	}

	session, err := h.loginService.RefreshToken(req.Context(), refreshToken)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	cookies := CreateAuthCookies(session, h.config, false)

	return res.Status(200).Cookies(cookies...).Json(session.ToJson())
}

func (h *AuthHttpHandler) GetUserPermissions(req vayload.HttpRequest, res vayload.HttpResponse) error {
	permissions, err := h.loginService.GetUserPermissions(req.Context(), req.Auth().UserId)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[any]{
		Status: "success",
		Data:   permissions,
	})
}

func (h *AuthHttpHandler) GetCurrentUser(req vayload.HttpRequest, res vayload.HttpResponse) error {
	session, err := h.loginService.GetCurrentSession(req.Context(), req.Auth().UserId.String(), req.Auth().AccessToken)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(session.ToJson())
}

func (h *AuthHttpHandler) LoginWithOAuth2(req vayload.HttpRequest, res vayload.HttpResponse) error {
	var input struct {
		Provider domain.OAuth2Provider `json:"provider"`
		Origin   string                `json:"origin"`
	}

	if err := req.ParseBody(&input); err != nil {
		return httpi.ErrBadRequest(err)
	}

	authURL, err := h.loginService.GetOAuth2URL(req.Context(), input.Provider, domain.OAuth2State{
		Origin: input.Origin,
		Nonce:  fmt.Sprintf("%d", time.Now().UnixNano()),
	})
	if err != nil {
		return httpi.ErrBadRequest(err)
	}

	response := map[string]string{
		"redirect": authURL,
	}

	return res.Status(200).Json(&httpi.ResponseBody[map[string]string]{
		Status: "success",
		Data:   response,
	})
}

func (h *AuthHttpHandler) HandleOAuth2Callback(req vayload.HttpRequest, res vayload.HttpResponse) error {
	provider := domain.OAuth2Provider(req.GetParam("provider"))
	if !provider.IsSupported() {
		// Only occurs when the provider is not found
		return httpi.ErrNotFound(nil)
	}

	code := req.GetQuery("code")
	if code == "" {
		codeData := req.FormData("code")
		if len(codeData) > 0 {
			code = codeData[0]
		}
	}

	state := req.GetQuery("state")
	if state == "" {
		stateData := req.FormData("state")
		if len(stateData) > 0 {
			state = stateData[0]
		}
	}

	if code == "" || state == "" {
		return httpi.ErrBadRequest(fmt.Errorf("code and state parameters are required"))
	}

	var data struct {
		Origin string `json:"origin"`
		Nonce  string `json:"nonce"`
	}

	rawState, _ := base64.URLEncoding.DecodeString(state)
	if err := json.Unmarshal(rawState, &data); err != nil {
		return httpi.ErrBadRequest(err)
	}

	authContext := domain.AuthContext{
		UserAgent: req.GetUserAgent(),
		IP:        req.GetIP(),
	}

	session, err := h.loginService.LoginWithOAuth2(req.Context(), domain.OAuth2Provider(provider), code, authContext)
	if err != nil {
		logger.E(err, logger.Fields{"context": "HandleOAuth2Callback", "provider": provider})

		return res.Redirect(data.Origin, 302)
	}

	cookies := CreateAuthCookies(session, h.config, true)
	return res.Status(200).Cookies(cookies...).Redirect(data.Origin, 302)
}

func (h *AuthHttpHandler) ForgotPasswordRequest(req vayload.HttpRequest, res vayload.HttpResponse) error {
	var input struct {
		Email string `json:"email"`
	}
	if err := req.ParseBody(&input); err != nil {
		return httpi.ErrBadRequest(err)
	}

	if err := h.recoveryService.RequestPasswordRecovery(req.Context(), input.Email); err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[string]{
		Status: "success",
		Data:   "recovery email sent",
	})
}

func (h *AuthHttpHandler) ForgotPasswordReset(req vayload.HttpRequest, res vayload.HttpResponse) error {
	var input struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := req.ParseBody(&input); err != nil {
		return httpi.ErrBadRequest(err)
	}

	if err := h.recoveryService.ResetPassword(req.Context(), input.Token, input.NewPassword); err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[string]{
		Status: "success",
		Data:   "password reset successfully",
	})
}

func (h *AuthHttpHandler) RegisterChangeEmailRequest(req vayload.HttpRequest, res vayload.HttpResponse) error {
	var input struct {
		Identifier string `json:"identifier"`
		Type       string `json:"type"`
		NewEmail   string `json:"new_email"`
	}
	if err := req.ParseBody(&input); err != nil {
		return httpi.ErrBadRequest(err)
	}

	changeInput := registration.RegisterEmailChangeInput{
		Identifier: input.Identifier,
		Type:       input.Type,
		NewEmail:   input.NewEmail,
	}

	if err := h.registerService.RequestEmailVerificationChange(req.Context(), changeInput); err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[string]{
		Status: "success",
		Data:   "email change requested",
	})
}

func (h *AuthHttpHandler) RegisterChangeEmailConfirm(req vayload.HttpRequest, res vayload.HttpResponse) error {
	var input struct {
		CurrentToken string `json:"current_token"`
		NewToken     string `json:"new_token"`
	}
	if err := req.ParseBody(&input); err != nil {
		return httpi.ErrBadRequest(err)
	}

	confirmationInput := registration.RegisterEmailChangeConfirmInput{
		CurrentToken: input.CurrentToken,
		NewToken:     input.NewToken,
	}

	confirmationErr := h.registerService.ConfirmEmailVerificationChange(req.Context(), confirmationInput)
	if confirmationErr != nil {
		return errors.MappingErrToHttp(confirmationErr)
	}

	return res.Status(200).Json(&httpi.ResponseBody[string]{
		Status: "success",
		Data:   "email changed successfully",
	})
}

// ============================= Internal Methods ==============================

// This method is used for internal services, for setup authentication (create or signing the current user)
func (h *AuthHttpHandler) InternalSetup(req vayload.HttpRequest, res vayload.HttpResponse) error {
	input := login.SetupUserInput{}
	if err := req.ValidateBody(&input); err != nil {
		return httpi.ErrWrapping(err)
	}

	authContext := domain.AuthContext{
		IP:        req.GetIP(),
		UserAgent: req.GetUserAgent(),
	}
	session, err := h.loginService.SetupUser(req.Context(), input, authContext)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(session.ToJson())
}

// ============================= Analytics Methods ==============================

func (h *AuthHttpHandler) GetAnalyticsKPIs(req vayload.HttpRequest, res vayload.HttpResponse) error {
	kpis, err := h.analyticsService.GetKPIs(req.Context())
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[*domain.KPIsResponse]{
		Status: "success",
		Data:   kpis,
	})
}

func (h *AuthHttpHandler) GetAnalyticsWeekly(req vayload.HttpRequest, res vayload.HttpResponse) error {
	weekly, err := h.analyticsService.GetWeeklyAnalytics(req.Context())
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[*domain.WeeklyAnalyticsResponse]{
		Status: "success",
		Data:   weekly,
	})
}

func (h *AuthHttpHandler) GetAnalyticsCountries(req vayload.HttpRequest, res vayload.HttpResponse) error {
	countries, err := h.analyticsService.GetCountryRanking(req.Context())
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[[]domain.CountryCount]{
		Status: "success",
		Data:   countries,
	})
}

func (h *AuthHttpHandler) GetRecentUsers(req vayload.HttpRequest, res vayload.HttpResponse) error {
	users, err := h.analyticsService.GetRecentUsers(req.Context(), 10)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[[]domain.RecentUser]{
		Status: "success",
		Data:   users,
	})
}

func (h *AuthHttpHandler) GetUserActivity(req vayload.HttpRequest, res vayload.HttpResponse) error {
	userIdStr := req.GetParam("id")
	if userIdStr == "" {
		return httpi.ErrBadRequest(fmt.Errorf("user id is required"))
	}

	userId, err := snowflake.FromString(userIdStr)
	if err != nil {
		return httpi.ErrBadRequest(fmt.Errorf("invalid user id"))
	}

	input := analytics.UserActivityInput{
		UserID: userId,
	}

	if from := req.GetQuery("from"); from != "" {
		input.From = &from
	}
	if to := req.GetQuery("to"); to != "" {
		input.To = &to
	}

	activity, err := h.analyticsService.GetUserActivity(req.Context(), input)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[*domain.UserActivityResponse]{
		Status: "success",
		Data:   activity,
	})
}

func (h *AuthHttpHandler) GetUserLogs(req vayload.HttpRequest, res vayload.HttpResponse) error {
	userIdStr := req.GetParam("id")
	if userIdStr == "" {
		return httpi.ErrBadRequest(fmt.Errorf("user id is required"))
	}

	userId, err := snowflake.FromString(userIdStr)
	if err != nil {
		return httpi.ErrBadRequest(fmt.Errorf("invalid user id"))
	}

	logs, err := h.analyticsService.GetUserLogs(req.Context(), userId)
	if err != nil {
		return errors.MappingErrToHttp(err)
	}

	return res.Status(200).Json(&httpi.ResponseBody[[]domain.UserLogEntry]{
		Status: "success",
		Data:   logs,
	})
}

func (h *AuthHttpHandler) withAuthCookies(session *domain.OAuthSession, res vayload.HttpResponse) error {
	cookies := CreateAuthCookies(session, h.config, true)
	return res.Status(200).Cookies(cookies...).Json(session.ToJson())
}

type HttpServices struct {
	LoginService     *login.LoginService
	RegisterService  *registration.RegisterService
	RecoveryService  *recovery.RecoveryService
	AnalyticsService *analytics.AnalyticsService
}

func NewHttpHandler(config *config.Config, registry vayload.Container, services HttpServices) *AuthHttpHandler {
	handler := &AuthHttpHandler{
		loginService:     services.LoginService,
		registerService:  services.RegisterService,
		recoveryService:  services.RecoveryService,
		analyticsService: services.AnalyticsService,
		config:           config,
		registry:         registry,
	}

	return handler
}

func (handler *AuthHttpHandler) HttpRoutes() []vayload.HttpRoutesGroup {
	authGuard := AuthHttpMiddleware(handler.registry, handler.config)

	routes := []vayload.HttpRoute{
		// Identify route
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/identify",
			Handler: handler.Identify,
		},

		// Authentication routes
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/password",
			Handler: handler.LoginWithPassword,
		},
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/oauth2",
			Handler: handler.LoginWithOAuth2,
		},
		{
			Method:  vayload.HttpMethod(httpi.GET),
			Path:    "/auth/:provider/callback",
			Handler: handler.HandleOAuth2Callback,
		},
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/:provider/callback", // Apple uses POST form_post
			Handler: handler.HandleOAuth2Callback,
		},
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/otp",
			Handler: handler.LoginWithOtpCode,
		},

		// User permissions
		{
			Method:      vayload.HttpMethod(httpi.GET),
			Path:        "/auth/permissions",
			Handler:     handler.GetUserPermissions,
			Middlewares: []vayload.HttpHandler{authGuard},
		},

		// Generate methods for login
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/otp/generate",
			Handler: handler.GenerateOtpCode,
		},

		// Registration routes
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/register",
			Handler: handler.RegisterUser,
		},
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/register/validate-account", // This method is used to validate the account after registration
			Handler: handler.ValidateRegister,
		},

		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/register/send-verification-code",
			Handler: handler.SendEmailVerificationCode,
		},

		// Forgot password routes (public)
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/forgot-password/request",
			Handler: handler.ForgotPasswordRequest,
		},
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/forgot-password/reset",
			Handler: handler.ForgotPasswordReset,
		},

		// Registration email change routes (pre-verificación)
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/register/change-email/request",
			Handler: handler.RegisterChangeEmailRequest,
		},
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/register/change-email/confirm",
			Handler: handler.RegisterChangeEmailConfirm,
		},

		// Logout route
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/logout",
			Handler: handler.Logout,
		},

		// Refresh token route
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/refresh-token",
			Handler: handler.RefreshToken,
		},

		{
			Method:      vayload.HttpMethod(httpi.GET),
			Path:        "/auth/me",
			Handler:     handler.GetCurrentUser,
			Middlewares: []vayload.HttpHandler{authGuard},
		},

		// ============================== INTERNAL ROUTES ==============================
		// These routes are for internal services use only and should not be exposed to the public
		{
			Method:  vayload.HttpMethod(httpi.POST),
			Path:    "/auth/__internal__/setup",
			Handler: handler.InternalSetup,
		},

		// ============================== ANALYTICS ROUTES ==============================
		// These routes require admin authentication via X-Auth-Token header
		{
			Method:  vayload.HttpMethod(httpi.GET),
			Path:    "/auth/analytics/kpis",
			Handler: handler.GetAnalyticsKPIs,
		},
		{
			Method:  vayload.HttpMethod(httpi.GET),
			Path:    "/auth/analytics/weekly",
			Handler: handler.GetAnalyticsWeekly,
		},
		{
			Method:  vayload.HttpMethod(httpi.GET),
			Path:    "/auth/analytics/countries",
			Handler: handler.GetAnalyticsCountries,
		},
		{
			Method:  vayload.HttpMethod(httpi.GET),
			Path:    "/auth/analytics/users/recent",
			Handler: handler.GetRecentUsers,
		},
		{
			Method:  vayload.HttpMethod(httpi.GET),
			Path:    "/auth/analytics/users/:id/activity",
			Handler: handler.GetUserActivity,
		},
		{
			Method:  vayload.HttpMethod(httpi.GET),
			Path:    "/auth/analytics/users/:id/logs",
			Handler: handler.GetUserLogs,
		},
	}

	return []vayload.HttpRoutesGroup{
		{
			Prefix: "auth",
			Routes: routes,
		},
	}
}
