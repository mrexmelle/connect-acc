package datesort

import (
	"github.com/mrexmelle/connect-emp/internal/datestr"
)

type DateStringSlice []string

func (s DateStringSlice) Len() int {
	return len(s)
}
func (s DateStringSlice) Less(i, j int) bool {
	if s[i] == datestr.Indeterminate {
		return false
	} else if s[j] == datestr.Indeterminate {
		return true
	} else {
		return s[i] < s[j]
	}
}

func (s DateStringSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
