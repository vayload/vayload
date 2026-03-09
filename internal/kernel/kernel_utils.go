/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Kernel Utils
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package kernel

import "github.com/spf13/cobra"

func toFlagMap(root *cobra.Command, flags map[string]any) map[string]any {
	flagPtrs := make(map[string]any)

	for flagName, defaultValue := range flags {
		switch v := defaultValue.(type) {
		case bool:
			ptr := new(bool)
			*ptr = v
			root.Flags().BoolVar(ptr, flagName, v, "")
			flagPtrs[flagName] = ptr
		case string:
			ptr := new(string)
			*ptr = v
			root.Flags().StringVar(ptr, flagName, v, "")
			flagPtrs[flagName] = ptr
		case int:
			ptr := new(int)
			*ptr = v
			root.Flags().IntVar(ptr, flagName, v, "")
			flagPtrs[flagName] = ptr
		}
	}

	return flagPtrs
}

func getFlagValues(flagPtrs map[string]any) map[string]any {
	flagValues := make(map[string]any)

	for flagName, ptr := range flagPtrs {
		switch p := ptr.(type) {
		case *bool:
			flagValues[flagName] = *p
		case *string:
			flagValues[flagName] = *p
		case *int:
			flagValues[flagName] = *p
		}
	}

	return flagValues
}

type cobraFlags struct {
	cmd *cobra.Command
}

func NewCobraFlags(cmd *cobra.Command) *cobraFlags {
	return &cobraFlags{cmd: cmd}
}

func (f *cobraFlags) GetString(name string, def string) string {
	val, err := f.cmd.Flags().GetString(name)
	if err != nil || val == "" {
		return def
	}
	return val
}

func (f *cobraFlags) GetBool(name string, def bool) bool {
	val, err := f.cmd.Flags().GetBool(name)
	if err != nil {
		return def
	}
	return val
}

func (f *cobraFlags) GetInt(name string, def int) int {
	val, err := f.cmd.Flags().GetInt(name)
	if err != nil {
		return def
	}
	return val
}
