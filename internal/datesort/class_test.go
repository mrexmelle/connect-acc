package datesort

import (
	"sort"
	"testing"
)

type SortTestCase struct {
	name     string
	input    []string
	expected []string
}

func TestIsCreated(t *testing.T) {
	tc := []SortTestCase{
		{
			name:     "Sorting without indeteriminate dates",
			input:    []string{"1985-04-24", "1994-04-07", "1000-09-09"},
			expected: []string{"1000-09-09", "1985-04-24", "1994-04-07"},
		},
		{
			name:     "Sorting with an indeterminate date",
			input:    []string{"1985-04-24", "", "1994-04-07", "1000-09-09"},
			expected: []string{"1000-09-09", "1985-04-24", "1994-04-07", ""},
		},
	}

	for _, c := range tc {
		sort.Sort(DateStringSlice(c.input))
		for i := 0; i < len(c.input); i++ {
			if c.input[i] != c.expected[i] {
				t.Errorf("[%s]\nresult: %v\nexpected: %v\n",
					c.name,
					c.input,
					c.expected,
				)
				break
			}
		}
	}
}
