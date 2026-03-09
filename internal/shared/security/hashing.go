/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Security
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package security

import (
	"github.com/vayload/vayload/pkg/crypto"
	"github.com/vayload/vayload/pkg/logger"
)

type hasher struct {
	// The core hasher interface
	scrypt crypto.Hasher
}

func NewHasher() *hasher {
	return &hasher{
		scrypt: crypto.NewScryptHasher(),
	}
}

func (h *hasher) HashPassword(password string) string {
	hash, err := h.scrypt.Generate([]byte(password))
	if err != nil {
		logger.E(err, logger.Fields{"context": "HashPassword", "action": "generate scrypt hash"})
		return ""
	}

	return string(hash)
}

func (h *hasher) VerifyPassword(password string, hash string) bool {
	result, err := h.scrypt.Compare([]byte(hash), []byte(password))
	if err != nil {
		logger.E(err, logger.Fields{"context": "VerifyPassword", "action": "compare scrypt hash"})
		return false
	}

	return result

}

func (h *hasher) VerboseVerifyPassword(password, hash string) (valid bool, algo string) {
	result, err := h.scrypt.Compare([]byte(hash), []byte(password))
	if err != nil {
		logger.E(err, logger.Fields{"context": "VerboseVerifyPassword", "action": "compare scrypt hash"})
		return false, "scrypt"
	}

	return result, "scrypt"
}
