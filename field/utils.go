package field

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

// IndexToId converts the given field width and space index into a unique string
// (coordinate pairs).
func IndexToId(width, index int) string {
	return fmt.Sprintf("%d,%d", index%width+1, index/width+1)
}

// max returns the maximum of the provided values.
func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of the provided values.
func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// remove removes the given element from the given slice.
func remove[T comparable](s []T, a T) []T {
	for i, b := range s {
		if b == a {
			s[i] = s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	return s
}
