package dataselect

import (
	"strings"
	"time"
)

type StdComparableString string

func (s StdComparableString) Compare(otherV ComparableValue) int {
	other := otherV.(StdComparableString)
	return strings.Compare(string(s), string(other))
}

func (s StdComparableString) Contains(otherV ComparableValue) bool {
	other := otherV.(StdComparableString)
	return strings.Contains(string(s), string(other))
}

type StdComparableTime time.Time

func (s StdComparableTime) Compare(otherV ComparableValue) int {
	other := otherV.(StdComparableTime)
	return ints64Compare(time.Time(s).Unix(), time.Time(other).Unix())
}

func (s StdComparableTime) Contains(otherV ComparableValue) bool {
	return s.Compare(otherV) == 0
}

func ints64Compare(a, b int64) int {
	if a > b {
		return 1
	} else if a == b {
		return 0
	}
	return -1
}
