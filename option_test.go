package query

import (
	"testing"
)

func TestOpenTrace(t *testing.T) {
	NewAction[any](
		WithDB[any](nil),
		WithContext[any](nil),
		WithTable[any](nil),
		WithTracer[any](nil),
		WithIOperation[any](nil),
		WithIOperationX[any](nil),
		WithIAssociation[any](nil),
	)
}
