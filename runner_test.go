package runtwoway

import (
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
		t.Errorf("Expect fully success")
	}
	checkMockRunners1ForwardRunned(m1x, t, 4)
}
