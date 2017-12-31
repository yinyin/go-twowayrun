package runtwoway

import (
	"testing"
)

type mockRunner1 struct {
	errOfForward  error
	errOfBackward error
}

func newMockRunner1(errOfForward, errOfBackward error) (r *mockRunner1) {
	return &mockRunner1{
		errOfForward:  errOfForward,
		errOfBackward: errOfBackward,
	}
}

func (r *mockRunner1) RunForward() (err error) {
	return r.errOfForward
}

func (r *mockRunner1) RunBackward() (err error) {
	return r.errOfBackward
}

func newMockRunners1(size int) TwoWayRunners {
	result := NewTwoWayRunners()
	for i := 0; i < size; i++ {
		result = result.Append(newMockRunner1(nil, nil))
	}
	return result
}

func TestTwoWayRunners_Run_e0(t *testing.T) {
	m1 := newMockRunners1(5)
	err := m1.Run()
	if nil != err {
		t.Errorf("Expect fully success")
	}
}
