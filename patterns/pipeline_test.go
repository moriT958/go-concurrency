package patterns

import (
	"context"
	"fmt"
	"testing"
)

func TestPipe(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pipe := Double(ctx, Offset(ctx, Double(ctx, Generator(ctx, nums...))))
	for v := range pipe {
		fmt.Println(v)
	}

	fmt.Println("finished")
}
