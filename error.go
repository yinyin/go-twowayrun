package runtwoway

import "fmt"

// TwoWayRunError represents error on running runners.
// Error returns by runner will be keep in PrevError field.
// Index of runner in TwoWayRunners instance will be keep in StopIndex field.
type TwoWayRunError struct {
	PrevError error
	StopIndex int
}

func newTwoWayRunError(prevError error, stopIndex int) (e *TwoWayRunError) {
	return &TwoWayRunError{
		PrevError: prevError,
		StopIndex: stopIndex,
	}
}

func (e *TwoWayRunError) Error() string {
	return fmt.Sprintf("TwoWayRunError(PrevError: %s, StopIndex: %d)", e.PrevError.Error(), e.StopIndex)
}
