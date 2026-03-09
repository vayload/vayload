/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Runtime
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package vayload

import lua "github.com/yuin/gopher-lua"

type ScriptRuntime interface {
	Init()
	EvalString(code string) error
	LoadFile(path string) error
	Bind(fn func(*lua.LState))
	Close() error
}

type ScriptModule interface {
	Register(L *lua.LState) int
}
