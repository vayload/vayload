/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package datetime

import (
	"errors"
	"time"
)

const (
	DATETIME = "2006-01-02 15:04:05"
	DATE     = "2006-01-02"
	TIME     = "15:04:05"
	ISO8601  = "2006-01-02T15:04:05Z07:00"
	ISO8601N = "2006-01-02T15:04:05.999999999Z07:00"
	RFC3339  = "2006-01-02T15:04:05Z07:00"
	RFC3339N = "2006-01-02T15:04:05.999999999Z07:00"
	RFC1123  = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z = "Mon, 02 Jan 2006 15:04:05 -0700"
	RFC822   = "02 Jan 06 15:04 MST"
	RFC822Z  = "02 Jan 06 15:04 -0700"
	ANSIC    = "Mon Jan _2 15:04:05 2006"
	Kitchen  = "3:04PM"
)

const (
	CHILE_TIMEZONE = "America/Santiago"
	UTC            = "UTC"
)

type Carbon struct {
	time.Time
}

func NewCarbon(tz ...string) (Carbon, error) {
	if len(tz) == 0 {
		tz = append(tz, UTC)
	}
	loc, err := time.LoadLocation(tz[0])
	if err != nil {
		return Carbon{}, err
	}
	return Carbon{time.Now().In(loc)}, nil
}

func NowToTimezone(tz string) (Carbon, error) {
	utcNow, _ := NewCarbon(UTC)
	return utcNow.To(tz)
}

func Parse(dateStr string, format string, tz string) (Carbon, error) {
	if tz == "" {
		tz = UTC
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return Carbon{}, err
	}
	parsed, err := time.ParseInLocation(format, dateStr, loc)
	if err != nil {
		return Carbon{}, err
	}
	return Carbon{parsed}, nil
}

func ParseISO(dateStr string, tz string) (Carbon, error) {
	if tz == "" {
		tz = UTC
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return Carbon{}, err
	}
	parsed, err := time.Parse(time.RFC3339Nano, dateStr)
	if err != nil {
		return Carbon{}, err
	}
	return Carbon{parsed.In(loc)}, nil
}

func From(date any, format string, tz string) (Carbon, error) {
	switch v := date.(type) {
	case time.Time:
		return Carbon{v.In(getLocationOrUTC(tz))}, nil
	case string:
		return Parse(v, format, tz)
	default:
		return Carbon{}, errors.New("unsupported type")
	}
}

func (c Carbon) Copy() *Carbon {
	return &Carbon{c.Time}
}

func (c Carbon) To(tz string) (Carbon, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return Carbon{}, err
	}
	return Carbon{c.Time.In(loc)}, nil
}

func (c Carbon) ToUTC() Carbon {
	return Carbon{c.Time.UTC()}
}

func (c *Carbon) WithUserTimeOffset(offset int, read bool) Carbon {
	if offset > 0 {
		if read {
			c.Time = c.Add(time.Duration(offset) * time.Hour)
		} else {
			c.Time = c.Add(-time.Duration(offset) * time.Hour)
		}
	} else {
		if read {
			c.Time = c.Add(-time.Duration(-offset) * time.Hour)
		} else {
			c.Time = c.Add(time.Duration(-offset) * time.Hour)
		}
	}
	return *c
}

func (c *Carbon) StartOfDay() Carbon {
	y, m, d := c.Date()
	c.Time = time.Date(y, m, d, 0, 0, 0, 0, c.Location())
	return *c
}

func (c *Carbon) EndOfDay() Carbon {
	y, m, d := c.Date()
	c.Time = time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), c.Location())
	return *c
}

func (c *Carbon) StartOfWeek() Carbon {
	offset := int(c.Weekday())
	c.Time = c.AddDate(0, 0, -offset)
	return c.StartOfDay()
}

func (c *Carbon) EndOfWeek() Carbon {
	offset := int(6 - c.Weekday())
	c.Time = c.AddDate(0, 0, offset)
	return c.EndOfDay()
}

func (c *Carbon) StartOfMonth() Carbon {
	y, m, _ := c.Date()
	c.Time = time.Date(y, m, 1, 0, 0, 0, 0, c.Location())
	return *c
}

func (c *Carbon) EndOfMonth() Carbon {
	y, m, _ := c.Date()
	t := time.Date(y, m+1, 0, 23, 59, 59, int(time.Second-time.Nanosecond), c.Location())
	c.Time = t
	return *c
}

func (c Carbon) InThisWeek(ref Carbon) bool {
	start := ref.Copy().StartOfWeek()
	end := ref.Copy().EndOfWeek()
	return !c.Before(start.Time) && !c.After(end.Time)
}

func (c Carbon) LessThanTo(d time.Time) bool {
	return c.Before(d)
}

func (c Carbon) LessThanOrEqualTo(d time.Time) bool {
	return c.Before(d) || c.Equal(d)
}

func (c Carbon) GreaterThanTo(d time.Time) bool {
	return c.After(d)
}

func (c Carbon) GreaterThanOrEqualTo(d time.Time) bool {
	return c.After(d) || c.Equal(d)
}

func getLocationOrUTC(tz string) *time.Location {
	if tz == "" {
		loc, _ := time.LoadLocation(UTC)
		return loc
	}
	loc, _ := time.LoadLocation(tz)
	return loc
}
