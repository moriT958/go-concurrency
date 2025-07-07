package exercise

import (
	"testing"

	"golang.org/x/tour/tree"
)

func TestSame(t *testing.T) {
	t.Run("should return true, if trees are same.", func(t *testing.T) {
		t1, t2 := tree.New(1), tree.New(1)
		if !Same(t1, t2) {
			t.Error("t1 and t2 must be same, but return false")
		}
	})

	t.Run("should return false, if trees are false.", func(t *testing.T) {
		t1, t2 := tree.New(1), tree.New(2)
		if Same(t1, t2) {
			t.Error("t1 and t2 must not be same, but return true")
		}
	})
}
