package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

var (
	seatsCapacity = 5
	arrivalRate   = 300
	cutDuration   = 1000 * time.Millisecond
	openDuration  = 10 * time.Second
)

func main() {
	rand.Seed(time.Now().UnixNano())

	color.Magenta("The Sleeping Barber Problem")
	color.Magenta("---------------------------")

	clientChan := make(chan string, seatsCapacity)
	doneChan := make(chan bool)

	shop := BarberShop{
		ShopCapacity:    seatsCapacity,
		CutDuration:     cutDuration,
		NumberOfBarbers: 0,
		Open:            true,
		ClientsChan:     clientChan,
		DoneChan:        doneChan,
	}
	color.Green("The shop is open")

	shop.AddBarber("Frank")
	shop.AddBarber("George")
	time.Sleep(1 * time.Second)

	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(openDuration)
		shopClosing <- true
		shop.Close()
		closed <- true
	}()

	i := 1
	go func() {
		for {
			randomMilliseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				// simulate clients coming to the closed shop
				for j := 0; j < 10; j++ {
					time.Sleep(time.Millisecond * time.Duration(randomMilliseconds))
					shop.AddClient(fmt.Sprintf("Client #%d", i))
					i++
				}
				return
			case <-time.After(time.Millisecond * time.Duration(randomMilliseconds)):
				shop.AddClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	<-closed
}
