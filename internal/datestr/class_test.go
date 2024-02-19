package datestr

import (
	"testing"
)

type IsCreatedTestCase struct {
	name      string
	s         string
	isCreated bool
}

type OffsetAndCloneTestCase struct {
	name     string
	s        string
	offset   int
	expected string
}

type CompareTestCase struct {
	name     string
	s1       string
	s2       string
	expected bool
}

func TestIsCreated(t *testing.T) {
	tc := []IsCreatedTestCase{
		{
			name:      "Valid date",
			s:         "1985-04-24",
			isCreated: true,
		},
		{
			name:      "Empty date",
			s:         "",
			isCreated: true,
		},
		{
			name:      "Invalid date",
			s:         "2004-02-30",
			isCreated: false,
		},
		{
			name:      "Alphabetical date",
			s:         "ABCD-EF-GH",
			isCreated: false,
		},
	}

	for _, c := range tc {
		_, err := NewFromString(c.s)
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

func TestOffsetAndClone(t *testing.T) {
	tc := []OffsetAndCloneTestCase{
		{
			name:     "Dates in non leap year",
			s:        "2023-02-28",
			offset:   1,
			expected: "2023-03-01",
		},
		{
			name:     "Dates in leap year",
			s:        "2024-02-28",
			offset:   1,
			expected: "2024-02-29",
		},
		{
			name:     "Dates in two different months",
			s:        "1985-05-04",
			offset:   -10,
			expected: "1985-04-24",
		},
	}

	for _, c := range tc {
		obj1, _ := NewFromString(c.s)
		obj2 := obj1.OffsetAndClone(c.offset)
		if obj2.AsString() != c.expected {
			t.Errorf("[%s]\nresult: %s\nexpected: %s\n",
				c.name,
				obj2.AsString(),
				c.expected,
			)
		}
	}
}

func TestEquals(t *testing.T) {
	tc := []CompareTestCase{
		{
			name:     "Equal dates",
			s1:       "1985-04-24",
			s2:       "1985-04-24",
			expected: true,
		},
		{
			name:     "Different dates",
			s1:       "1985-04-24",
			s2:       "1994-04-07",
			expected: false,
		},
		{
			name:     "One empty",
			s1:       "",
			s2:       "1985-04-24",
			expected: false,
		},
		{
			name:     "Both empty",
			s1:       "",
			s2:       "",
			expected: true,
		},
	}

	for _, c := range tc {
		obj1, _ := NewFromString(c.s1)
		obj2, _ := NewFromString(c.s2)
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

func TestIsBefore(t *testing.T) {
	tc := []CompareTestCase{
		{
			name:     "Dates are same",
			s1:       "1985-04-24",
			s2:       "1985-04-24",
			expected: false,
		},
		{
			name:     "Date is actually before",
			s1:       "1985-04-24",
			s2:       "1985-04-25",
			expected: true,
		},
		{
			name:     "Date is after",
			s1:       "1985-04-24",
			s2:       "1985-04-23",
			expected: false,
		},
		{
			name:     "Date is empty",
			s1:       "",
			s2:       "1985-04-24",
			expected: false,
		},
		{
			name:     "Compared date is empty",
			s1:       "1985-04-24",
			s2:       "",
			expected: true,
		},
		{
			name:     "Both empty",
			s1:       "",
			s2:       "",
			expected: false,
		},
	}

	for _, c := range tc {
		obj1, _ := NewFromString(c.s1)
		obj2, _ := NewFromString(c.s2)
		out := obj1.IsBefore(obj2)
		if out != c.expected {
			t.Errorf("[%s]\nresult: %t\nexpected: %t\n",
				c.name,
				out,
				c.expected,
			)
		}
	}
}

func TestIsBeforeOrEquals(t *testing.T) {
	tc := []CompareTestCase{
		{
			name:     "Dates are same",
			s1:       "1985-04-24",
			s2:       "1985-04-24",
			expected: true,
		},
		{
			name:     "Date is actually before",
			s1:       "1985-04-24",
			s2:       "1985-04-25",
			expected: true,
		},
		{
			name:     "Date is after",
			s1:       "1985-04-24",
			s2:       "1985-04-23",
			expected: false,
		},
		{
			name:     "Date is empty",
			s1:       "",
			s2:       "1985-04-24",
			expected: false,
		},
		{
			name:     "Compared date is empty",
			s1:       "1985-04-24",
			s2:       "",
			expected: true,
		},
		{
			name:     "Both empty",
			s1:       "",
			s2:       "",
			expected: true,
		},
	}

	for _, c := range tc {
		obj1, _ := NewFromString(c.s1)
		obj2, _ := NewFromString(c.s2)
		out := obj1.IsBeforeOrEquals(obj2)
		if out != c.expected {
			t.Errorf("[%s]\nresult: %t\nexpected: %t\n",
				c.name,
				out,
				c.expected,
			)
		}
	}
}

func TestIsAfter(t *testing.T) {
	tc := []CompareTestCase{
		{
			name:     "Dates are same",
			s1:       "1985-04-24",
			s2:       "1985-04-24",
			expected: false,
		},
		{
			name:     "Date is actually before",
			s1:       "1985-04-24",
			s2:       "1985-04-25",
			expected: false,
		},
		{
			name:     "Date is after",
			s1:       "1985-04-24",
			s2:       "1985-04-23",
			expected: true,
		},
		{
			name:     "Date is empty",
			s1:       "",
			s2:       "1985-04-24",
			expected: true,
		},
		{
			name:     "Compared date is empty",
			s1:       "1985-04-24",
			s2:       "",
			expected: false,
		},
		{
			name:     "Both empty",
			s1:       "",
			s2:       "",
			expected: false,
		},
	}

	for _, c := range tc {
		obj1, _ := NewFromString(c.s1)
		obj2, _ := NewFromString(c.s2)
		out := obj1.IsAfter(obj2)
		if out != c.expected {
			t.Errorf("[%s]\nresult: %t\nexpected: %t\n",
				c.name,
				out,
				c.expected,
			)
		}
	}
}

func TestIsAfterOrEquals(t *testing.T) {
	tc := []CompareTestCase{
		{
			name:     "Dates are same",
			s1:       "1985-04-24",
			s2:       "1985-04-24",
			expected: true,
		},
		{
			name:     "Date is actually before",
			s1:       "1985-04-24",
			s2:       "1985-04-25",
			expected: false,
		},
		{
			name:     "Date is after",
			s1:       "1985-04-24",
			s2:       "1985-04-23",
			expected: true,
		},
		{
			name:     "Date is empty",
			s1:       "",
			s2:       "1985-04-24",
			expected: true,
		},
		{
			name:     "Compared date is empty",
			s1:       "1985-04-24",
			s2:       "",
			expected: false,
		},
		{
			name:     "Both empty",
			s1:       "",
			s2:       "",
			expected: true,
		},
	}

	for _, c := range tc {
		obj1, _ := NewFromString(c.s1)
		obj2, _ := NewFromString(c.s2)
		out := obj1.IsAfterOrEquals(obj2)
		if out != c.expected {
			t.Errorf("[%s]\nresult: %t\nexpected: %t\n",
				c.name,
				out,
				c.expected,
			)
		}
	}
}
