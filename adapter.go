package twowayrun

import (
	"context"
)

// TwoWayRunner defines required methods for a two-way runner.
type TwoWayRunner interface {
	// RunForward performs operation forward to commit.
	// Nil should be returned if operation is successful.
	// Otherwise an error should be returned.
	RunForward(ctx context.Context) (err error)

	// RunBackward performs operation rollbacks forward operation.
	// Nil should be returned if operation is successful.
	// Otherwise an error should be returned.
	RunBackward(ctx context.Context) (err error)
}
