/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package migrator

import "time"

type Migration struct {
	Version   string     `json:"version"`
	Name      string     `json:"name"`
	UpSQL     string     `json:"up_sql"`
	DownSQL   string     `json:"down_sql"`
	LuaScript string     `json:"lua_script,omitempty"`
	Applied   bool       `json:"applied"`
	AppliedAt *time.Time `json:"applied_at"`
}

type MigrationSQLModel struct {
	Version   string     `db:"version"`
	Name      string     `db:"name"`
	AppliedAt *time.Time `db:"applied_at"`
}

type MigrationModel struct {
	Version   string     `json:"version"`
	Name      string     `json:"name"`
	UpSQL     string     `json:"up_sql"`
	DownSQL   string     `json:"down_sql"`
	LuaScript string     `json:"lua_script,omitempty"`
	Applied   bool       `json:"applied"`
	AppliedAt *time.Time `json:"applied_at"`
}

type MigrationStatus struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	Applied bool   `json:"applied"`
}
