package exercise

import (
	"golang.org/x/tour/tree"
)

// Exercise: Equivalent Binary Trees
// https://go-tour-jp.appspot.com/concurrency/7

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	defer close(ch)
	walk(t, ch)
}

func walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}

	if t.Left != nil {
		walk(t.Left, ch)
	}
	ch <- t.Value

	if t.Right != nil {
		walk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {

	ch1, ch2 := make(chan int), make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2
		if (ok1 != ok2) || (v1 != v2) {
			return false
		}
		if !ok1 || !ok2 {
			break
		}
	}
	return true
}
