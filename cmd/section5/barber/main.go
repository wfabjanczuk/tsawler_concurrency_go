package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

var (
	seatsCapacity = 2
	arrivalRate   = 300
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
		ShopClosingChan: make(chan bool, 1),
		ClientsChan:     clientChan,
		DoneChan:        doneChan,
	}
	color.Green("The shop is open")

	shop.AddBarber("Frank")
	shop.AddBarber("George")
	time.Sleep(1 * time.Second)

	shopClosedChan := shop.Start()

	clientsClosing := make(chan bool)
	go func() {
		<-time.After(2 * openDuration)
		clientsClosing <- true
	}()

	i := 1
	go func() {
		for {
			randomMilliseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-clientsClosing:
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
	time.Sleep(2 * openDuration)
}
