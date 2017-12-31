package runtwoway

type TwoWayRunner interface {
	RunForward() (success bool)
	RunBackward() (success bool)
}
