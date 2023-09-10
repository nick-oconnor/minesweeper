package field

import (
	"fmt"
)

type State int

const (
	Unknown  State = iota
	Flagged  State = iota
	Revealed State = iota
)

// String returns a text version of the State enum.
func (s State) String() string {
	switch s {
	case Unknown:
		return "unknown"
	case Flagged:
		return "flagged"
	case Revealed:
		return "revealed"
	}
	panic(fmt.Sprintf("invalid state %d", s))
}
