package dateinterval

import (
	"fmt"

	"github.com/mrexmelle/connect-emp/internal/datestr"
	"github.com/mrexmelle/connect-emp/internal/localerror"
)

type Class struct {
	StartDate *datestr.Class
	EndDate   *datestr.Class
}

func NewFromStrings(startDate string, endDate string) (*Class, error) {
	sd, err := datestr.NewFromString(startDate)
	if err != nil {
		return nil, err
	}

	ed, err := datestr.NewFromString(endDate)
	if err != nil {
		return nil, err
	}

	if sd.IsAfter(ed) {
		fmt.Println("sd is after ed")
		return nil, localerror.ErrBadDateSequence
	}

	return &Class{
		StartDate: sd,
		EndDate:   ed,
	}, nil
}

func NewFromDateStrings(startDate *datestr.Class, endDate *datestr.Class) (*Class, error) {
	if startDate.IsAfter(endDate) {
		return nil, localerror.ErrBadDateSequence
	}

	return &Class{
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

func (c *Class) Equals(other *Class) bool {
	return (c.StartDate.AsString() == other.StartDate.AsString() &&
		c.EndDate.AsString() == other.EndDate.AsString())
}

func (c *Class) IsEncompassingDate(d *datestr.Class) bool {
	return (c.StartDate.IsBeforeOrEquals(d) && c.EndDate.IsAfterOrEquals(d))
}

func (c *Class) IsEncompassingInterval(other *Class) bool {
	return c.StartDate.IsBeforeOrEquals(other.StartDate) &&
		(c.EndDate.IsAfterOrEquals(other.EndDate) || (c.EndDate.IsIndeterminate() && other.EndDate.IsIndeterminate()))
}

func (c *Class) CollideWith(other *Class) []*Class {
	if other.IsEncompassingInterval(c) {
		return []*Class{}
	} else if c.StartDate.IsBefore(other.StartDate) && other.IsEncompassingDate(c.EndDate) {
		leftOverEndDate := other.StartDate.OffsetAndClone(-1)
		leftOver, _ := NewFromDateStrings(c.StartDate, leftOverEndDate)
		return []*Class{
			leftOver,
		}
	} else if other.IsEncompassingDate(c.StartDate) && c.EndDate.IsAfter(other.EndDate) {
		leftOverStartDate := other.EndDate.OffsetAndClone(+1)
		leftOver, _ := NewFromDateStrings(leftOverStartDate, c.EndDate)
		return []*Class{
			leftOver,
		}
	} else if c.IsEncompassingInterval(other) {
		leftOver1EndDate := other.StartDate.OffsetAndClone(-1)
		leftOver2StartDate := other.EndDate.OffsetAndClone(+1)
		leftOver1, _ := NewFromDateStrings(c.StartDate, leftOver1EndDate)
		leftOver2, _ := NewFromDateStrings(leftOver2StartDate, c.EndDate)
		return []*Class{
			leftOver1,
			leftOver2,
		}
	}
	return []*Class{c}
}
