package main

import (
	"github.com/fatih/color"
	"time"
)

type BarberShop struct {
	CutDuration     time.Duration
	ShopCapacity    int
	NumberOfBarbers int
	BarberDoneChan  chan bool
	ClientsChan     chan string
}

func (b *BarberShop) AddBarber(barber string) {
	b.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("\t%s goes to the waiting room to check for clients", barber)

		for {
			if !isSleeping && len(b.ClientsChan) == 0 {
				isSleeping = true
				color.Yellow("\tThere is nothing to do, so %s takes a nap", barber)
			}

			client, shopOpen := <-b.ClientsChan

			if shopOpen {
				if isSleeping {
					isSleeping = false
					color.Yellow("\t%s wakes up", barber)
				}

				b.cutHair(barber, client)
			} else {
				b.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (b *BarberShop) cutHair(barber, client string) {
	color.Green("\t%s started cutting %s's hair", barber, client)
	time.Sleep(b.CutDuration)
	color.Green("\t%s is finished cutting %s's hair", barber, client)
}

func (b *BarberShop) sendBarberHome(barber string) {
	color.Blue("\t%s is going home", barber)
	b.BarberDoneChan <- true
}

func (b *BarberShop) AddClient(client string) {
	color.Magenta("*** %s arrives! ***", client)

	select {
	case b.ClientsChan <- client:
		color.Blue("%s takes a seat in waiting room", client)
	default:
		color.Red("The waiting room is full, so %s leaves", client)
	}
}

func (b *BarberShop) Start() (soonClosed, closed chan bool) {
	soonClosed = make(chan bool)
	closed = make(chan bool)

	go func() {
		<-time.After(openDuration)
		soonClosed <- true
		b.Close()
		closed <- true
	}()

	return soonClosed, closed
}

func (b *BarberShop) Close() {
	color.Blue("Closing shop for the day")
	close(b.ClientsChan)

	for a := 1; a <= b.NumberOfBarbers; a++ {
		<-b.BarberDoneChan
	}

	close(b.BarberDoneChan)
	color.Green("--------------------------------------------------------")
	color.Green("The barber shop is now closed and everyone has gone home")
}
