package query

import (
	"testing"
)

func TestOpenTrace(t *testing.T) {
	NewAction[any](OpenTrace[any](), WithDB[any](nil), WithContext[any](nil), WithTable[any](nil))
}
