package datestr

import (
	"database/sql"
	"errors"
	"time"
)

const (
	FormatDefault = "2006-01-02"
	Indeterminate = ""
)

type Class struct {
	date string
}

func NewFromString(s string) (*Class, error) {
	if s != "" {
		_, err := time.Parse(FormatDefault, s)
		if err != nil {
			return nil, err
		}
	}

	return &Class{
		date: s,
	}, nil
}

func NewFromTime(t time.Time) *Class {
	return &Class{
		date: t.Format(FormatDefault),
	}
}

func (c *Class) AsString() string {
	return c.date
}

func (c *Class) AsTime() time.Time {
	t, _ := time.Parse(FormatDefault, c.date)
	return t
}

func (c *Class) AsSqlNullTime() (sql.NullTime, error) {
	var t sql.NullTime
	var err error

	if c.date == "" {
		t.Valid = false
		return t, nil
	}

	t.Time, err = time.Parse("2006-01-02", c.date)
	if err != nil {
		t.Valid = false
		return t, errors.New("string must be in " + FormatDefault + " format")
	}

	t.Valid = true
	return t, nil
}

func (c *Class) OffsetAndClone(offset int) *Class {
	t := c.AsTime()
	afterOffset := t.Add(time.Duration(offset) * 24 * time.Hour)
	return NewFromTime(afterOffset)
}

func (c *Class) IsIndeterminate() bool {
	return (c.date == Indeterminate)
}

func (c *Class) IsValidStartDate() bool {
	return !c.IsIndeterminate()
}

func (c *Class) IsValidEndDate() bool {
	return true
}

func (c *Class) Equals(other *Class) bool {
	return (c.IsIndeterminate() && other.IsIndeterminate()) ||
		(c.date == other.date)
}

func (c *Class) IsBefore(other *Class) bool {
	if c.IsIndeterminate() && !other.IsIndeterminate() {
		return false
	}

	if !c.IsIndeterminate() && other.IsIndeterminate() {
		return true
	}

	return (c.date < other.date)
}

func (c *Class) IsBeforeOrEquals(other *Class) bool {
	if c.IsIndeterminate() && !other.IsIndeterminate() {
		return false
	}

	if !c.IsIndeterminate() && other.IsIndeterminate() {
		return true
	}

	return (c.date <= other.date)
}

func (c *Class) IsAfter(other *Class) bool {
	if c.IsIndeterminate() && !other.IsIndeterminate() {
		return true
	}

	if !c.IsIndeterminate() && other.IsIndeterminate() {
		return false
	}

	return (c.date > other.date)
}

func (c *Class) IsAfterOrEquals(other *Class) bool {
	if c.IsIndeterminate() && !other.IsIndeterminate() {
		return true
	}

	if !c.IsIndeterminate() && other.IsIndeterminate() {
		return false
	}

	return (c.date >= other.date)
}
