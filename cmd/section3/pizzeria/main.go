package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, pizzasTotal int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we are making
	i := 0

	// run forever or until we receive a quit notification
	// try to make pizzas
	for {
		currentPizzaOrder := makePizza(i)
		if currentPizzaOrder != nil {
			i = currentPizzaOrder.pizzaNumber
			pizzaMaker.data <- *currentPizzaOrder
		} else {
			quitChan := <-pizzaMaker.quit
			close(pizzaMaker.data)
			close(quitChan)
			return
		}
	}
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber > NumberOfPizzas {
		return nil
	}

	delay := rand.Intn(5) + 1
	color.Magenta("Received order number %d!\n", pizzaNumber)

	rnd := rand.Intn(12) + 1
	msg := ""
	success := false

	if rnd < 5 {
		pizzasFailed++
	} else {
		pizzasMade++
	}
	pizzasTotal++

	color.Yellow("Making pizza %d. It will take %d seconds...\n", pizzaNumber, delay)
	time.Sleep(time.Duration(delay) * time.Second)

	if rnd <= 2 {
		msg = fmt.Sprintf("We ran out of ingredients for pizza %d", pizzaNumber)
	} else if rnd <= 4 {
		msg = fmt.Sprintf("The cook quit while making pizza %d", pizzaNumber)
	} else {
		success = true
		msg = fmt.Sprintf("Pizza order %d is ready", pizzaNumber)
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
		message:     msg,
		success:     success,
	}
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func main() {
	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// print out a message
	color.Magenta("The Producer is open for business!")
	color.Magenta("----------------------------------")

	// create and run a producer
	producer := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}
	go pizzeria(producer)

	// create and run a consumer
	for o := range producer.data {
		if o.success {
			color.Green(o.message)
			color.Green("Order %d is out for delivery!", o.pizzaNumber)
		} else {
			color.Red(o.message)
			color.Red("The customer is really mad!")
		}

		if o.pizzaNumber == NumberOfPizzas {
			color.Magenta("Done making pizzas.")
			err := producer.Close()
			if err != nil {
				color.Red("Error closing channel!")
			}
		}
	}

	color.Green("-----------------")
	color.Green("Done for the day.")
	color.Green("We made %d pizzas, but failed to make %d, with %d attempts in total.", pizzasMade, pizzasFailed, pizzasTotal)

	switch {
	case pizzasFailed > 3:
		color.Red("It was an awful day")
	case pizzasFailed > 2:
		color.Red("It was not a very good day")
	case pizzasFailed > 0:
		color.Yellow("It was OK day")
	default:
		color.Green("It was great day")
	}
}
