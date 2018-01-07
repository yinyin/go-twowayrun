package runtwoway

import "fmt"

// TwoWayRunError represents error on running runners.
// Error returns by runner will be keep in RunnerError field.
// Index of runner in TwoWayRunners instance will be keep in StopIndex field.
type TwoWayRunError struct {
	RunnerError error
	StopIndex   int
}

func newTwoWayRunError(runnerError error, stopIndex int) (e *TwoWayRunError) {
	return &TwoWayRunError{
		RunnerError: runnerError,
		StopIndex:   stopIndex,
	}
}

func (e *TwoWayRunError) toError() error {
	if nil == e {
		return nil
	}
	return e
}

func (e *TwoWayRunError) Error() string {
	return fmt.Sprintf("TwoWayRunError(RunnerError: %s, StopIndex: %d)", e.RunnerError.Error(), e.StopIndex)
}

// TwoWayRunErrors represents errors on running runners where runner errors could be ignored.
// Errors from runner will be pack into TwoWayRunError and append into RunErrors field.
type TwoWayRunErrors struct {
	RunErrors []*TwoWayRunError
}

func newTwoWayRunErrors() (e *TwoWayRunErrors) {
	return &TwoWayRunErrors{
		RunErrors: make([]*TwoWayRunError, 0),
	}
}

func (e *TwoWayRunErrors) appendRunnerError(runnerError error, stopIndex int) {
	errInst := newTwoWayRunError(runnerError, stopIndex)
	e.RunErrors = append(e.RunErrors, errInst)
}

func (e *TwoWayRunErrors) toError() error {
	if (nil == e) || (0 == len(e.RunErrors)) {
		return nil
	}
	return e
}

func (e *TwoWayRunErrors) Error() string {
	var firstError *TwoWayRunError
	if len(e.RunErrors) > 0 {
		firstError = e.RunErrors[0]
	}
	return fmt.Sprintf("TwoWayRunErrors(len(RunErrors): %d, RunErrors[0]: %v)", len(e.RunErrors), firstError)
}
