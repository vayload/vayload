/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package login

import (
	"strings"
)

func MaskIdentity(value, typ string) string {
	switch typ {
	case "phone":
		if len(value) <= 5 {
			return value
		}
		prefix := value[:2]
		suffix := value[len(value)-3:]
		maskedLength := len(value) - 5
		return prefix + strings.Repeat("*", maskedLength) + suffix

	case "email":
		atIndex := strings.Index(value, "@")
		if atIndex == -1 {
			return value // Invalid email, return as is
		}

		username := value[:atIndex]
		domain := value[atIndex+1:]

		dotParts := strings.Split(domain, ".")
		if len(dotParts) == 0 {
			return username + "@****"
		}

		tld := dotParts[len(dotParts)-1]
		return username + "@****." + tld

	default:
		return value // Unknown type, return as is
	}
}

func detectIdentifierType(identifier string) string {
	if strings.Contains(identifier, "@") {
		return "email"
	}

	return "phone"
}
