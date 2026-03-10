/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package database

import (
	"database/sql"
	"time"

	"github.com/vayload/vayload/internal/shared/snowflake"
)

func NilIfInvalidString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func NilIfInvalidInt64(ni sql.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}

func NilIfInvalidInt32(ni sql.NullInt32) *int32 {
	if ni.Valid {
		return &ni.Int32
	}
	return nil
}

func NilIfInvalidFloat64(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}

func NilIfInvalidBool(nb sql.NullBool) *bool {
	if nb.Valid {
		return &nb.Bool
	}
	return nil
}

func NilIfInvalidTime(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

func NilIfInvalidID(ni snowflake.ID) *snowflake.ID {
	if ni > 0 {
		return &ni
	}
	return nil
}

func NilIfInvalidFlakeId(ni sql.NullInt64) *snowflake.ID {
	if ni.Valid {
		return snowflake.FromInt64Ptr(&ni.Int64)
	}

	return nil
}
