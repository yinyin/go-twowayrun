package twowayrun

import "log"
import "context"

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
// Context object ctx will be pass into RunForward and RunBackward methods of runners.
func (r TwoWayRunners) Run(ctx context.Context) (err error) {
	var prevError error
	stoppedAt := (int)(-1)
	for idx, runner := range r {
		e := runner.RunForward(ctx)
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
		e := runner.RunBackward(ctx)
		if nil != e {
			log.Printf("TwoWayRunners.Run: failed at backward run: (index=%d, error=%v)", idx, e)
		}
	}
	return newTwoWayRunError(prevError, stoppedAt)
}

// RunForward runs forward runners only.
// Context object ctx will be pass into RunForward method of runners.
// Will stop on runner which result into error if stopOnError is set to true.
// If any error occurs the resulted error will be TwoWayRunErrors structure.
func (r TwoWayRunners) RunForward(ctx context.Context, stopOnError bool) (err error) {
	errInst := newTwoWayRunErrors()
	for idx, runner := range r {
		e := runner.RunForward(ctx)
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

// RunBackward runs backward runners only. Runners will be activate in reverse order.
// Context object ctx will be pass into RunBackward method of runners.
// Will stop on runner which result into error if stopOnError is set to true.
// If any error occurs the resulted error will be TwoWayRunErrors structure.
func (r TwoWayRunners) RunBackward(ctx context.Context, stopOnError bool) (err error) {
	errInst := newTwoWayRunErrors()
	for idx := len(r) - 1; idx >= 0; idx-- {
		runner := r[idx]
		e := runner.RunBackward(ctx)
		if nil != e {
			errInst.appendRunnerError(e, idx)
			log.Printf("TwoWayRunners.RunBackward: failed: (index=%d, error=%v)", idx, e)
			if stopOnError {
				break
			}
		}
	}
	return errInst.toError()
}
