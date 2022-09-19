package main

import (
	"fmt"
	"time"
)

func listenToChan(ch <-chan int) {
	for {
		i := <-ch
		fmt.Printf("Got %d from channel\n", i)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	ch := make(chan int, 5)
	go listenToChan(ch)

	for i := 0; i < 10; i++ {
		ch <- i
		fmt.Printf("Sent %d to channel\n", i)
	}
}
