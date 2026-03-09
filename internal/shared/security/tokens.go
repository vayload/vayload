/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Security
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type randomizer struct{}

func NewRandomizer() *randomizer {
	return &randomizer{}
}

func (r *randomizer) GenerateRandomString(length int, urlSafe bool) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		timestamp := time.Now().UnixNano()
		for i := range bytes {
			bytes[i] = byte((timestamp >> (i * 8)) & 0xFF)
		}
	}

	if urlSafe {
		return hex.EncodeToString(bytes)
	}

	return base64.StdEncoding.EncodeToString(bytes)
}

func (r *randomizer) GenerateRandomNumericCode(length int) string {
	if length <= 0 {
		return ""
	}

	max := big.NewInt(10)
	max.Exp(max, big.NewInt(int64(length)), nil)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		timestamp := time.Now().UnixNano()
		code := int(timestamp%int64(max.Int64())) + 1
		return formatNumericCode(code, length)
	}

	return formatNumericCode(int(n.Int64()), length)
}

func (r *randomizer) GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, nil
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (r *randomizer) GenerateUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return ""
	}
	// Set version to 4 (random)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// Set variant to 10xx
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func (r *randomizer) SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// utils
func formatNumericCode(code int, length int) string {
	codeStr := strconv.Itoa(code)
	if len(codeStr) < length {
		prefix := strings.Repeat("0", length-len(codeStr))
		return prefix + codeStr
	}
	return codeStr
}

type otpCode struct {
	generator *randomizer
}

func NewOtpCodeStrategy() *otpCode {
	return &otpCode{
		generator: NewRandomizer(),
	}
}

func (o *otpCode) GenerateOtpCode() string {
	return o.generator.GenerateRandomNumericCode(6)
}

func (o *otpCode) CompareOtpCode(inputCode string, storedCode string) bool {
	return subtle.ConstantTimeCompare([]byte(inputCode), []byte(storedCode)) == 1
}
