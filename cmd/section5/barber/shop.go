package main

import (
	"github.com/fatih/color"
	"sync"
	"time"
)

type BarberShop struct {
	ShopCapacity    int
	CutDuration     time.Duration
	NumberOfBarbers int
	DoneChan        chan bool
	ClientsChan     chan string
	ClientsDoneChan chan bool
	closed          bool
	mutex           sync.RWMutex
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
	b.DoneChan <- true
}

func (b *BarberShop) AddClient(client string) {
	color.Magenta("*** %s arrives! ***", client)

	b.mutex.RLock()
	if b.closed {
		color.Red("The shop is already closed, so %s leaves", client)
	} else {
		select {
		case b.ClientsChan <- client:
			color.Blue("%s takes a seat in waiting room", client)
		default:
			color.Red("The waiting room is full, so %s leaves", client)
		}
	}
	b.mutex.RUnlock()
}

func (b *BarberShop) Start() chan bool {
	shopClosedChan := make(chan bool)

	go func() {
		<-time.After(openDuration)

		b.mutex.Lock()
		b.closed = true
		b.mutex.Unlock()

		b.Close()
		shopClosedChan <- true
	}()

	return shopClosedChan
}

func (b *BarberShop) Close() {
	color.Blue("Closing shop for the day")
	close(b.ClientsChan)

	for a := 1; a <= b.NumberOfBarbers; a++ {
		<-b.DoneChan
	}

	close(b.DoneChan)
	color.Green("--------------------------------------------------------")
	color.Green("The barber shop is now closed and everyone has gone home")
}
