package patterns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Fan-out/Fan-inは複数のタスクを並行で処理し、最終結果を1つにまとめるパターン

func HeavyTask(i int, id int) string {
	time.Sleep(200 * time.Millisecond)
	return fmt.Sprintf("result:%v (id:%v)", i*i, id)
}

func FanOut(ctx context.Context, in <-chan int, id int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-ctx.Done():
				return
			case out <- HeavyTask(v, id):
			}
		}
	}()
	return out
}

func FanIn(ctx context.Context, chs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	multiplex := func(ch <-chan string) {
		defer wg.Done()
		for text := range ch {
			select {
			case <-ctx.Done():
				return
			case out <- text:
			}
		}
	}

	wg.Add(len(chs))
	for _, ch := range chs {
		go multiplex(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
