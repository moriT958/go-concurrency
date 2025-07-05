package patterns

import (
	"context"
	"fmt"
	"runtime"
	"testing"
)

func TestFanOutFanin(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nums := []int{1, 2, 3, 4, 5, 6, 7, 8}

	cores := runtime.NumCPU()

	outChs := make([]<-chan string, cores)
	inData := Generator(ctx, nums...)
	for i := 0; i < cores; i++ {
		outChs[i] = FanOut(ctx, inData, i+1)
	}

	for v := range FanIn(ctx, outChs...) {
		fmt.Println(v)
	}

	fmt.Println("finished")
}
