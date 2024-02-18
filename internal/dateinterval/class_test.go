package dateinterval

import (
	"testing"

	"github.com/mrexmelle/connect-emp/internal/datestr"
)

type IsCreatedTestCase struct {
	name      string
	startDate string
	endDate   string
	isCreated bool
}

type EqualsTestCase struct {
	name       string
	startDate1 string
	endDate1   string
	startDate2 string
	endDate2   string
	expected   bool
}

type IsEncompassingDateTestCase struct {
	name      string
	startDate string
	endDate   string
	d         string
	expected  bool
}

type IsEncompassingIntervalTestCase struct {
	name       string
	startDate1 string
	endDate1   string
	startDate2 string
	endDate2   string
	expected   bool
}

type CollideWithTestCase struct {
	name       string
	startDate1 string
	endDate1   string
	startDate2 string
	endDate2   string
	expected   []DateIntervalPair
}

type DateIntervalPair struct {
	startDate string
	endDate   string
}

func NewIntervalPair(sd string, ed string) DateIntervalPair {
	return DateIntervalPair{
		startDate: sd,
		endDate:   ed,
	}
}

func TestIsCreated(t *testing.T) {
	tc := []IsCreatedTestCase{
		{
			name:      "Correct sequence",
			startDate: "1985-04-01",
			endDate:   "1985-04-30",
			isCreated: true,
		},
		{
			name:      "Same dates",
			startDate: "1985-04-24",
			endDate:   "1985-04-24",
			isCreated: true,
		},
		{
			name:      "Wrong sequence",
			startDate: "1985-04-30",
			endDate:   "1985-04-01",
			isCreated: false,
		},
		{
			name:      "Empty end date",
			startDate: "1985-04-24",
			endDate:   "",
			isCreated: true,
		},
		{
			name:      "Empty start date",
			startDate: "",
			endDate:   "1985-04-24",
			isCreated: false,
		},
	}

	for _, c := range tc {
		_, err := NewFromStrings(c.startDate, c.endDate)
		created := (err == nil)
		if created != c.isCreated {
			t.Errorf("[%s]\nresult: %t\nexpected: %t\n",
				c.name,
				created,
				c.isCreated,
			)
		}
	}
}

func TestEquals(t *testing.T) {
	tc := []EqualsTestCase{
		{
			name:       "Equal structs",
			startDate1: "1985-04-24",
			endDate1:   "1985-04-24",
			startDate2: "1985-04-24",
			endDate2:   "1985-04-24",
			expected:   true,
		},
		{
			name:       "Different start dates",
			startDate1: "1985-04-24",
			endDate1:   "1985-04-24",
			startDate2: "1985-04-24",
			endDate2:   "1994-04-07",
			expected:   false,
		},
		{
			name:       "Different end dates",
			startDate1: "1985-04-24",
			endDate1:   "1985-04-24",
			startDate2: "1985-05-24",
			endDate2:   "1994-04-07",
			expected:   false,
		},
	}

	for _, c := range tc {
		obj1, _ := NewFromStrings(c.startDate1, c.endDate1)
		obj2, _ := NewFromStrings(c.startDate2, c.endDate2)
		out := obj1.Equals(obj2)
		if out != c.expected {
			t.Errorf("[%s]\nresult: %t\nexpected: %t\n",
				c.name,
				out,
				c.expected,
			)
		}
	}
}

func TestIsEncompassingDate(t *testing.T) {
	tc := []IsEncompassingDateTestCase{
		{
			name:      "Date is encompassed",
			startDate: "1985-04-01",
			endDate:   "1985-04-30",
			d:         "1985-04-24",
			expected:  true,
		},
		{
			name:      "Date is before interval",
			startDate: "1985-04-01",
			endDate:   "1985-04-30",
			d:         "1985-03-03",
			expected:  false,
		},
		{
			name:      "Date is after interval",
			startDate: "1985-04-01",
			endDate:   "1985-04-30",
			d:         "1985-05-03",
			expected:  false,
		},
		{
			name:      "Date is start date",
			startDate: "1985-04-01",
			endDate:   "1985-04-30",
			d:         "1985-04-01",
			expected:  true,
		},
		{
			name:      "Date is end date",
			startDate: "1985-04-01",
			endDate:   "1985-04-30",
			d:         "1985-04-30",
			expected:  true,
		},
		{
			name:      "Date is before start date but end date is empty",
			startDate: "1985-04-01",
			endDate:   "",
			d:         "1985-03-01",
			expected:  false,
		},
		{
			name:      "Date is empty and end date is empty",
			startDate: "1985-04-01",
			endDate:   "",
			d:         "",
			expected:  true,
		},
		{
			name:      "Date is after start date but end date is empty",
			startDate: "1985-04-01",
			endDate:   "",
			d:         "2004-03-01",
			expected:  true,
		},
		{
			name:      "Date is empty",
			startDate: "1985-04-01",
			endDate:   "1985-04-30",
			d:         "",
			expected:  false,
		},
	}

	for _, c := range tc {
		obj, _ := NewFromStrings(c.startDate, c.endDate)
		d, _ := datestr.NewFromString(c.d)
		out := obj.IsEncompassingDate(d)
		if out != c.expected {
			t.Errorf("[%s]\nresult: %t\nexpected: %t\n",
				c.name,
				out,
				c.expected,
			)
		}
	}
}

func TestIsEncompassingInterval(t *testing.T) {
	tc := []IsEncompassingIntervalTestCase{
		{
			name:       "Other is before interval",
			startDate1: "1985-04-01",
			endDate1:   "1985-04-30",
			startDate2: "1985-02-01",
			endDate2:   "1985-03-01",
			expected:   false,
		},
		{
			name:       "Other is after interval",
			startDate1: "1985-04-01",
			endDate1:   "1985-04-30",
			startDate2: "1985-05-01",
			endDate2:   "1985-06-01",
			expected:   false,
		},
		{
			name:       "Same start dates and within interval",
			startDate1: "1985-04-01",
			endDate1:   "1985-04-30",
			startDate2: "1985-04-01",
			endDate2:   "1985-04-02",
			expected:   true,
		},
		{
			name:       "Same end dates and within interval",
			startDate1: "1985-04-01",
			endDate1:   "1985-04-30",
			startDate2: "1985-04-24",
			endDate2:   "1985-04-30",
			expected:   true,
		},
		{
			name:       "Empty end dates and within interval",
			startDate1: "1985-04-01",
			endDate1:   "",
			startDate2: "1985-04-24",
			endDate2:   "",
			expected:   true,
		},
		{
			name:       "Empty end dates and before interval",
			startDate1: "1985-04-01",
			endDate1:   "",
			startDate2: "1985-03-01",
			endDate2:   "",
			expected:   false,
		},
		{
			name:       "Empty other's end date",
			startDate1: "1985-04-01",
			endDate1:   "1985-04-30",
			startDate2: "1985-04-24",
			endDate2:   "",
			expected:   false,
		},
	}

	for _, c := range tc {
		obj1, _ := NewFromStrings(c.startDate1, c.endDate1)
		obj2, _ := NewFromStrings(c.startDate2, c.endDate2)
		out := obj1.IsEncompassingInterval(obj2)
		if out != c.expected {
			t.Errorf("[%s]\nresult: %t\nexpected: %t\n",
				c.name,
				out,
				c.expected,
			)
		}
	}
}

func TestCollideWith(t *testing.T) {
	tc := []CollideWithTestCase{
		{
			name:       "Other intersects with first part",
			startDate1: "1985-04-01",
			endDate1:   "1985-04-30",
			startDate2: "1985-03-01",
			endDate2:   "1985-04-10",
			expected: []DateIntervalPair{
				NewIntervalPair("1985-04-11", "1985-04-30"),
			},
		},
		{
			name:       "Other's end date is class' start date",
			startDate1: "1985-04-01",
			endDate1:   "1985-04-30",
			startDate2: "1985-03-01",
			endDate2:   "1985-04-01",
			expected: []DateIntervalPair{
				NewIntervalPair("1985-04-02", "1985-04-30"),
			},
		},
		{
			name:       "Other intersects with last part",
			startDate1: "1985-04-01",
			endDate1:   "1985-04-30",
			startDate2: "1985-04-24",
			endDate2:   "1985-05-10",
			expected: []DateIntervalPair{
				NewIntervalPair("1985-04-01", "1985-04-23"),
			},
		},
		{
			name:       "Other's start date is class' end date",
			startDate1: "1985-04-01",
			endDate1:   "1985-04-30",
			startDate2: "1985-04-30",
			endDate2:   "1985-05-01",
			expected: []DateIntervalPair{
				NewIntervalPair("1985-04-01", "1985-04-29"),
			},
		},
		{
			name:       "Other is equal",
			startDate1: "1985-04-01",
			endDate1:   "1985-04-30",
			startDate2: "1985-04-01",
			endDate2:   "1985-04-30",
			expected:   []DateIntervalPair{},
		},
		{
			name:       "Other is bigger",
			startDate1: "1985-04-01",
			endDate1:   "1985-04-30",
			startDate2: "1985-03-01",
			endDate2:   "1985-05-01",
			expected:   []DateIntervalPair{},
		},
		{
			name:       "Other is smaller",
			startDate1: "1985-03-01",
			endDate1:   "1985-05-30",
			startDate2: "1985-04-01",
			endDate2:   "1985-04-30",
			expected: []DateIntervalPair{
				NewIntervalPair("1985-03-01", "1985-03-31"),
				NewIntervalPair("1985-05-01", "1985-05-30"),
			},
		},
		{
			name:       "Other doesn't collide",
			startDate1: "1985-04-01",
			endDate1:   "1985-04-30",
			startDate2: "1986-03-01",
			endDate2:   "1986-05-01",
			expected: []DateIntervalPair{
				NewIntervalPair("1985-04-01", "1985-04-30"),
			},
		},
	}

	for _, c := range tc {
		obj1, _ := NewFromStrings(c.startDate1, c.endDate1)
		obj2, _ := NewFromStrings(c.startDate2, c.endDate2)

		out := obj1.CollideWith(obj2)

		exp := []*Class{}
		for i := 0; i < len(c.expected); i++ {
			pair, _ := NewFromStrings(c.expected[i].startDate, c.expected[i].endDate)
			exp = append(exp, pair)
		}

		for i := 0; i < len(out); i++ {
			if out[i].Equals(exp[i]) == false {
				t.Errorf("[%s]\nresult: %v\nexpected: %v\n",
					c.name,
					*out[i],
					*exp[i],
				)
			}
		}
	}
}
