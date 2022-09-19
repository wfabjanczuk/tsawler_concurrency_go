package main

import (
	"fmt"
	"strings"
)

func main() {
	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)
	pingPong(ping, pong)

	fmt.Println("--------------------------")
	fmt.Println("All done. Closing channels")
	close(ping)
	close(pong)
}

func pingPong(ping chan<- string, pong <-chan string) {
	fmt.Println("Type something and press ENTER (or Q to quit)")

	for {
		fmt.Print("-> ")
		var userInput string
		_, _ = fmt.Scanln(&userInput)

		if strings.ToLower(userInput) == "q" {
			return
		}

		ping <- userInput
		response, ok := <-pong

		if ok {
			fmt.Println("\tResponse:", response)
		}
	}
}

func shout(ping <-chan string, pong chan<- string) {
	for {
		s, ok := <-ping
		if ok {
			pong <- fmt.Sprintf("%s!!!", strings.ToUpper(s))
		}
	}
}
