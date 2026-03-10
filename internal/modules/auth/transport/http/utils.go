/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package auth_http

import (
	"fmt"
	"time"

	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/vayload"
)

const (
	COOKIE_REFRESH_TOKEN  = "izrefre"
	COOKIE_ACCESS_TOKEN   = "iztoken"
	COOKIE_FLAG_LOGGED_IN = "izloggedin"
	COOKIE_AVATAR_URL     = "izavatar"
)

func CreateAuthCookies(session *domain.OAuthSession, config *config.Config, withRefresh bool) []*vayload.Cookie {
	cookies := []*vayload.Cookie{
		{
			Name:        COOKIE_ACCESS_TOKEN,
			Value:       fmt.Sprintf("%s %s", session.TokenType, session.AccessToken),
			HttpOnly:    true,
			Secure:      true,
			Expires:     session.ExpiresAt,
			MaxAge:      int(session.ExpiresIn),
			SameSite:    "None",
			Domain:      config.App.Domain,
			SessionOnly: false,
		},
		{
			Name:        COOKIE_FLAG_LOGGED_IN,
			Value:       "true",
			HttpOnly:    false,
			Secure:      true,
			Expires:     session.ExpiresAt,
			MaxAge:      int(session.ExpiresIn),
			SameSite:    "None",
			Domain:      config.App.Domain,
			SessionOnly: false,
		},
	}

	if withRefresh {
		cookies = append(cookies, &vayload.Cookie{
			Name:        COOKIE_REFRESH_TOKEN,
			Value:       session.RefreshToken,
			HttpOnly:    true,
			Secure:      true,
			Expires:     session.ExpiresRefreshAt,
			MaxAge:      int(session.ExpiresRefreshIn),
			SameSite:    "None",
			Domain:      config.App.Domain,
			SessionOnly: false,
		})
	}

	if session.Meta != nil {
		if avatar, ok := session.Meta["avatar_url"].(string); ok && avatar != "" {
			cookies = append(cookies, &vayload.Cookie{
				Name:        COOKIE_AVATAR_URL,
				Value:       avatar,
				HttpOnly:    false,
				Secure:      true,
				Expires:     session.ExpiresAt,
				MaxAge:      int(session.ExpiresIn),
				SameSite:    "None",
				Domain:      config.App.Domain,
				SessionOnly: false,
			})
		}
	}

	return cookies
}

func FlushAuthCookies(config *config.Config) []*vayload.Cookie {
	expired := time.Unix(0, 0)

	return []*vayload.Cookie{
		{
			Name:        COOKIE_REFRESH_TOKEN,
			Value:       "",
			Path:        "/",
			Domain:      config.App.Domain,
			HttpOnly:    true,
			Secure:      true,
			SameSite:    "None",
			Expires:     expired,
			MaxAge:      -1,
			SessionOnly: false,
		},
		{
			Name:        COOKIE_ACCESS_TOKEN,
			Value:       "",
			Path:        "/",
			Domain:      config.App.Domain,
			HttpOnly:    true,
			Secure:      true,
			SameSite:    "None",
			Expires:     expired,
			MaxAge:      -1,
			SessionOnly: false,
		},
		{
			Name:        COOKIE_FLAG_LOGGED_IN,
			Value:       "",
			Path:        "/",
			Domain:      config.App.Domain,
			HttpOnly:    false,
			Secure:      true,
			SameSite:    "None",
			Expires:     expired,
			MaxAge:      -1,
			SessionOnly: false,
		},
		{
			Name:        COOKIE_AVATAR_URL,
			Value:       "",
			Path:        "/",
			Domain:      config.App.Domain,
			HttpOnly:    false,
			Secure:      true,
			SameSite:    "None",
			Expires:     expired,
			MaxAge:      -1,
			SessionOnly: false,
		},
	}
}
