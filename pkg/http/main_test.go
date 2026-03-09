/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package httpi

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup any global test configuration here

	// Run all tests
	code := m.Run()

	// Cleanup any global test resources here

	os.Exit(code)
}
