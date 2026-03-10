/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package auth_http

import (
	"strings"

	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/kernel"
	"github.com/vayload/vayload/internal/shared/security"
	"github.com/vayload/vayload/internal/vayload"
	httpi "github.com/vayload/vayload/pkg/http"
	"github.com/vayload/vayload/pkg/logger"
)

func AuthHttpMiddleware(registry vayload.Container, config *config.Config) vayload.HttpHandler {
	return func(req vayload.HttpRequest, res vayload.HttpResponse) error {
		jwtManager, err := kernel.MapTo[*security.JwtManager](registry, security.JWT_MANAGER_KEY)
		if err != nil {
			logger.E(err, logger.Fields{"context": "AuthMiddleware", "action": "retrieve JwtManager"})

			return res.Status(500).Json(httpi.ErrorResponse{
				Status: "error",
				Error: httpi.HttpError{
					Code:    "INTERNAL_ERROR",
					Message: "Unable to retrieve authentication service",
				},
			})
		}

		// Get token from Authorization header or cookie
		tokenRaw := req.GetHeader("Authorization")
		if len(tokenRaw) == 0 {
			tokenRaw = req.GetCookie(COOKIE_ACCESS_TOKEN)
		}

		tokenStr := strings.TrimPrefix(tokenRaw, "Bearer ")
		if tokenStr == "" {
			return res.Status(401).Json(httpi.ErrorResponse{
				Status: "error",
				Error: httpi.HttpError{
					Code:    "UNAUTHORIZED",
					Message: "Expected authorization token",
				},
			})
		}

		token, err := jwtManager.ValidateToken(tokenStr)
		if err != nil {
			safeToken := maskToken(tokenStr)
			logger.E(err, logger.Fields{"context": "AuthMiddleware", "token": safeToken})

			return res.Status(401).Json(httpi.ErrorResponse{
				Status: "error",
				Error: httpi.HttpError{
					Code:    "UNAUTHORIZED",
					Message: "Invalid or expired token",
				},
			})
		}

		req.Locals(httpi.HTTP_AUTH_KEY, &vayload.HttpAuth{
			UserId:      token.ID,
			Role:        string(token.Role),
			AccessToken: tokenStr,
			CountryId:   token.CountryId,
		})

		return req.Next()
	}
}

func maskToken(token string) string {
	if len(token) <= 10 {
		return token
	}
	return token[:6] + "..." + token[len(token)-4:]
}
