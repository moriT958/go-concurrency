package main

import (
	"fmt"
	"math/rand"
	"time"
)

func getLuckyNum(c chan<- int) {
	fmt.Println("...")

	rand.New(rand.NewSource(time.Now().Unix()))
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)

	num := rand.Intn(10)

	c <- num
}

func main() {
	fmt.Println("what is today's lucky number?")

	// チャネルを使った例
	c := make(chan int)
	go getLuckyNum(c)

	num := <-c

	fmt.Printf("Today's your lucky number is %d!\n", num)

	close(c)

	// WaitGroupを使った実行同期の例
	// var wg sync.WaitGroup
	// wg.Add(1)

	// go func() {
	// 	defer wg.Done()

	// 	getLuckyNum()
	// }()

	// wg.Wait()
}
