package patterns

import (
	"context"
)

func Generator(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case <-ctx.Done():
				// cancelコンテキストを受信したら終了
				return
			case out <- n:
			}
		}
	}()
	return out
}

func Double(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case <-ctx.Done():
				return
			case out <- n * 2:
			}
		}
	}()
	return out
}

func Offset(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case <-ctx.Done():
				return
			case out <- n + 2:
			}
		}
	}()
	return out
}
