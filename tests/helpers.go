/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package tests

import (
	"reflect"
	"testing"
)

type TestCase struct {
	Name     string
	TestFunc func(t *testing.T)
	Parallel bool // Flag for parallel execution
}

func RunTestCases(t *testing.T, testCases []TestCase) {
	for _, testCase := range testCases {
		tc := testCase // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Parallel {
				t.Parallel()
			}
			tc.TestFunc(t)
		})
	}
}

type Assert struct {
	t *testing.T
}

func NewAssert(t *testing.T) *Assert {
	return &Assert{t: t}
}

func (a *Assert) Equal(expected, actual any, msg ...string) {
	a.t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		message := getMessage(msg, "values are not equal")
		a.t.Fatalf("%s\nExpected: %#v\nActual:   %#v", message, expected, actual)
	}
}

func (a *Assert) NotEqual(expected, actual any, msg ...string) {
	a.t.Helper()
	if reflect.DeepEqual(expected, actual) {
		message := getMessage(msg, "values should not be equal")
		a.t.Fatalf("%s\nUnexpected: %#v", message, actual)
	}
}

func (a *Assert) True(actual bool, msg ...string) {
	a.t.Helper()
	if !actual {
		a.t.Fatalf("%s", getMessage(msg, "expected true, got false"))
	}
}

func (a *Assert) False(actual bool, msg ...string) {
	a.t.Helper()
	if actual {
		a.t.Fatalf("%s", getMessage(msg, "expected false, got true"))
	}
}

func (a *Assert) Nil(actual any, msg ...string) {
	a.t.Helper()
	if !isNil(actual) {
		a.t.Fatalf("%s\nExpected nil, got: %#v", getMessage(msg, ""), actual)
	}
}

func (a *Assert) NotNil(actual any, msg ...string) {
	a.t.Helper()
	if isNil(actual) {
		a.t.Fatalf("%s\nExpected not nil, got: nil", getMessage(msg, ""))
	}
}

func (a *Assert) Error(err error, msg ...string) {
	a.t.Helper()
	if err == nil {
		a.t.Fatalf("%s\nExpected an error, got nil", getMessage(msg, ""))
	}
}

func (a *Assert) NoError(err error, msg ...string) {
	a.t.Helper()
	if err != nil {
		a.t.Fatalf("%s\nUnexpected error: %v", getMessage(msg, ""), err)
	}
}

func (a *Assert) Panics(f func(), msg ...string) {
	a.t.Helper()
	defer func() {
		if r := recover(); r == nil {
			a.t.Fatalf("%s\nExpected panic, but none occurred", getMessage(msg, ""))
		}
	}()
	f()
}

func (a *Assert) NotPanics(f func(), msg ...string) {
	a.t.Helper()
	defer func() {
		if r := recover(); r != nil {
			a.t.Fatalf("%s\nUnexpected panic: %v", getMessage(msg, ""), r)
		}
	}()
	f()
}

func (a *Assert) ApproxEqual(expected, actual, tolerance float64, msg ...string) {
	a.t.Helper()
	diff := expected - actual
	if diff < 0 {
		diff = -diff
	}
	if diff > tolerance {
		a.t.Fatalf("%s\nExpected: %f ±%f\nActual: %f", getMessage(msg, "values not within tolerance"), expected, tolerance, actual)
	}
}

func getMessage(msg []string, defaultMsg string) string {
	if len(msg) > 0 {
		return msg[0]
	}
	return defaultMsg
}

func isNil(v any) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	kind := val.Kind()
	return (kind >= reflect.Chan && kind <= reflect.Slice && val.IsNil())
}
