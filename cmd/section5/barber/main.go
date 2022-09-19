package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

var (
	seatsCapacity = 2
	arrivalRate   = 0
	cutDuration   = 1000 * time.Millisecond
	openDuration  = 5 * time.Second
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
		ClientsChan:     clientChan,
		BarberDoneChan:  doneChan,
	}
	color.Green("The shop is open")

	shop.AddBarber("Frank")
	shop.AddBarber("George")
	time.Sleep(1 * time.Second)

	shopSoonClosedChan, shopClosedChan := shop.Start()

	i := 1
	go func() {
		for {
			randomMilliseconds := arrivalRate
			select {
			case <-shopSoonClosedChan:
				color.Cyan("<><><><><><><><><><><><><><><><><")
				color.Cyan("Pissed off clients stopped coming")
				color.Cyan("<><><><><><><><><><><><><><><><><")
				return
			case <-time.After(time.Millisecond * time.Duration(randomMilliseconds)):
				shop.AddClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	<-shopClosedChan
}
