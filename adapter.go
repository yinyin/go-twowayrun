package twowayrun

// TwoWayRunner defines required methods for a two-way runner.
type TwoWayRunner interface {
	RunForward() (err error)
	RunBackward() (err error)
}
