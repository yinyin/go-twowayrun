package runtwoway

import "log"

// TwoWayRunners is a slice construct by two-way runners
type TwoWayRunners []TwoWayRunner

// NewTwoWayRunners creates an instance of TwoWayRunners
func NewTwoWayRunners() (r TwoWayRunners) {
	return make([]TwoWayRunner, 0)
}

// Append given runner instance to the end of TwoWayRunners instance and return
// appended instance.
func (r TwoWayRunners) Append(runner TwoWayRunner) TwoWayRunners {
	return append(r, runner)
}

// Run invokes RunForward() method of runners in TwoWayRunners instance one by one.
// The invocation will turn backward if any runner returns error.
func (r TwoWayRunners) Run() (err error) {
	var prevError error
	stoppedAt := (int)(-1)
	for idx, runner := range r {
		e := runner.RunForward()
		if nil != e {
			prevError = e
			stoppedAt = idx
			log.Printf("TwoWayRunners.Run: failed at forward run: (index=%d, error=%v)", idx, e)
			break
		}
	}
	if -1 == stoppedAt {
		return nil
	}
	for idx := stoppedAt; idx >= 0; idx-- {
		runner := r[idx]
		e := runner.RunBackward()
		if nil != e {
			log.Printf("TwoWayRunners.Run: failed at backward run: (index=%d, error=%v)", idx, e)
		}
	}
	return newTwoWayRunError(prevError, stoppedAt)
}

// RunForward runs forward runners only.
// Will stop on runner which result into error if stopOnError is set to true.
// If any error occurs the resulted error will be TwoWayRunErrors structure.
func (r TwoWayRunners) RunForward(stopOnError bool) (err error) {
	var errInst *TwoWayRunErrors
	for idx, runner := range r {
		e := runner.RunForward()
		if nil != e {
			errInst.appendRunnerError(e, idx)
			log.Printf("TwoWayRunners.RunForward: failed: (index=%d, error=%v)", idx, e)
			if stopOnError {
				break
			}
		}
	}
	return errInst.toError()
}
