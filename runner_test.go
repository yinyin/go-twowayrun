package runtwoway

import (
	"fmt"
	"testing"
)

type mockRunner1 struct {
	countForward  int
	errOfForward  error
	countBackward int
	errOfBackward error
}

func newMockRunner1(errOfForward, errOfBackward error) (r *mockRunner1) {
	return &mockRunner1{
		countForward:  0,
		errOfForward:  errOfForward,
		countBackward: 0,
		errOfBackward: errOfBackward,
	}
}

func (r *mockRunner1) RunForward() (err error) {
	r.countForward++
	return r.errOfForward
}

func (r *mockRunner1) RunBackward() (err error) {
	r.countBackward++
	return r.errOfBackward
}

func newMockRunners1(size int) TwoWayRunners {
	result := NewTwoWayRunners()
	for i := 0; i < size; i++ {
		result = result.Append(newMockRunner1(nil, nil))
	}
	return result
}

func checkMockRunners1BothRunned(runners TwoWayRunners, t *testing.T, toIdx int) {
	for idx, runner := range runners {
		m1 := runner.(*mockRunner1)
		if idx <= toIdx {
			if 1 != m1.countForward {
				t.Errorf("forward no run (index=%d)", idx)
			}
			if 1 != m1.countBackward {
				t.Errorf("backward no run (index=%d)", idx)
			}
		} else {
			if 0 != m1.countForward {
				t.Errorf("forward runned (index=%d)", idx)
			}
			if 0 != m1.countBackward {
				t.Errorf("backward runned (index=%d)", idx)
			}
		}
	}
}

func checkMockRunners1ForwardRunned(runners TwoWayRunners, t *testing.T, toIdx int) {
	for idx, runner := range runners {
		m1 := runner.(*mockRunner1)
		if idx <= toIdx {
			if 1 != m1.countForward {
				t.Errorf("forward no run (index=%d)", idx)
			}
		} else {
			if 0 != m1.countForward {
				t.Errorf("forward runned (index=%d)", idx)
			}
		}
		if 0 != m1.countBackward {
			t.Errorf("backward runned (index=%d)", idx)
		}
	}
}

func checkMockRunners1BackwardRunned(runners TwoWayRunners, t *testing.T, toIdx int) {
	for idx, runner := range runners {
		m1 := runner.(*mockRunner1)
		if 0 != m1.countForward {
			t.Errorf("forward runned (index=%d)", idx)
		}
		if idx <= toIdx {
			if 1 != m1.countBackward {
				t.Errorf("backward no run (index=%d)", idx)
			}
		} else {
			if 0 != m1.countBackward {
				t.Errorf("backward runned (index=%d)", idx)
			}
		}
	}
}

func castToTwoWayRunErrors(t *testing.T, err error, expErrorCount int) (errInst *TwoWayRunErrors) {
	errInst, ok := err.(*TwoWayRunErrors)
	if !ok {
		t.Errorf("expecting TwoWayRunErrors: %#v", err)
		return nil
	}
	if expErrorCount != len(errInst.RunErrors) {
		t.Errorf("expecting %d errors in TwoWayRunErrors but got %d", expErrorCount, len(errInst.RunErrors))
	}
	return errInst
}

func checkTwoWayRunError(t *testing.T, err error, expRunnerError error, expStopIndex int) {
	errInst, ok := err.(*TwoWayRunError)
	if !ok {
		t.Errorf("expecting error type as TwoWayRunError: %v", err)
		return
	}
	if expRunnerError != errInst.RunnerError {
		t.Errorf("unexpected error instance: %v", errInst)
	}
	if expStopIndex != errInst.StopIndex {
		t.Errorf("unexpected stop index: %v", errInst)
	}
}

func checkTwoWayRunErrorsElement(t *testing.T, err *TwoWayRunErrors, errIndex int, expRunnerError error, expStopIndex int) {
	if nil == err {
		return
	}
	if len(err.RunErrors) <= errIndex {
		t.Errorf("length of RunErrors not that much: %d <= %d", len(err.RunErrors), errIndex)
		return
	}
	e := err.RunErrors[errIndex]
	checkTwoWayRunError(t, e, expRunnerError, expStopIndex)
}

func TestTwoWayRunners_Run_e0(t *testing.T) {
	m1x := newMockRunners1(5)
	err := m1x.Run()
	if nil != err {
		t.Errorf("expecting fully success: %v", err)
	}
	checkMockRunners1ForwardRunned(m1x, t, 4)
}

func TestTwoWayRunners_Run_e1(t *testing.T) {
	m1x := newMockRunners1(5)
	mockErr := fmt.Errorf("mock error")
	m1x[2].(*mockRunner1).errOfForward = mockErr
	err := m1x.Run()
	if nil == err {
		t.Errorf("unexpected fully success")
	}
	checkTwoWayRunError(t, err, mockErr, 2)
	checkMockRunners1BothRunned(m1x, t, 2)
}

func TestTwoWayRunners_RunForward_e0(t *testing.T) {
	m1x := newMockRunners1(5)
	err := m1x.RunForward(false)
	if nil != err {
		t.Errorf("expecting fully success: %#v", err)
	}
	checkMockRunners1ForwardRunned(m1x, t, 4)
}

func TestTwoWayRunners_RunForward_e1a(t *testing.T) {
	m1x := newMockRunners1(5)
	mockErr1 := fmt.Errorf("mock error 1")
	m1x[2].(*mockRunner1).errOfForward = mockErr1
	mockErr2 := fmt.Errorf("mock error 2")
	m1x[4].(*mockRunner1).errOfForward = mockErr2
	err := m1x.RunForward(false)
	if nil == err {
		t.Errorf("expecting some errors: %#v", err)
	}
	checkMockRunners1ForwardRunned(m1x, t, 4)
	errInst := castToTwoWayRunErrors(t, err, 2)
	checkTwoWayRunErrorsElement(t, errInst, 0, mockErr1, 2)
	checkTwoWayRunErrorsElement(t, errInst, 1, mockErr2, 4)
}

func TestTwoWayRunners_RunForward_e1b(t *testing.T) {
	m1x := newMockRunners1(5)
	mockErr1 := fmt.Errorf("mock error 1")
	m1x[2].(*mockRunner1).errOfForward = mockErr1
	mockErr2 := fmt.Errorf("mock error 2")
	m1x[4].(*mockRunner1).errOfForward = mockErr2
	err := m1x.RunForward(true)
	if nil == err {
		t.Errorf("expecting some errors: %#v", err)
	}
	checkMockRunners1ForwardRunned(m1x, t, 2)
	errInst := castToTwoWayRunErrors(t, err, 1)
	checkTwoWayRunErrorsElement(t, errInst, 0, mockErr1, 2)
}
