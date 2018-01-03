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
	errInst, ok := err.(*TwoWayRunError)
	if !ok {
		t.Errorf("expecting error type as TwoWayRunError: %v", err)
	}
	if mockErr != errInst.PrevError {
		t.Errorf("unexpected mock error instance: %v", errInst)
	}
	if 2 != errInst.StopIndex {
		t.Errorf("unexpected stop index: %v", errInst)
	}
	checkMockRunners1BothRunned(m1x, t, 2)
}

func TestTwoWayRunners_RunForward_e0(t *testing.T) {
	m1x := newMockRunners1(5)
	err := m1x.RunForward(false)
	if nil != err {
		t.Errorf("expecting fully success: %v", err)
	}
	checkMockRunners1ForwardRunned(m1x, t, 4)
}
