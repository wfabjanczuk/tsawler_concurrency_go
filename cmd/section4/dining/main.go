package main

import (
	"github.com/fatih/color"
	"sync"
	"time"
)

type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

var philosophers = []Philosopher{
	{name: "Plato", leftFork: 4, rightFork: 0},
	{name: "Socrates", leftFork: 0, rightFork: 1},
	{name: "Aristotle", leftFork: 1, rightFork: 2},
	{name: "Pascal", leftFork: 2, rightFork: 3},
	{name: "Locke", leftFork: 3, rightFork: 4},
}

var (
	maxMeals     = 3
	eatingTime   = 200 * time.Millisecond
	thinkingTime = 300 * time.Millisecond
	results      = NewStats()
)

type stats struct {
	m        *sync.Mutex
	finished []string
}

func NewStats() *stats {
	return &stats{
		m:        &sync.Mutex{},
		finished: make([]string, 0, len(philosophers)),
	}
}

func (s *stats) appendFinished(name string) {
	s.m.Lock()
	color.Green("%s has finished", name)
	color.Green("%s has left the table", name)
	s.finished = append(s.finished, name)
	s.m.Unlock()
}

func (s *stats) String() string {
	msg := "The order: "
	for _, name := range s.finished {
		msg += name + ", "
	}
	return msg
}

func main() {
	color.Magenta("Dining Philosophers Problem")
	color.Magenta("---------------------------")
	color.Yellow("The table is empty.")

	dine()

	color.Green("-------------------")
	color.Green("The table is empty.")
}

func dine() {
	eatingWg := &sync.WaitGroup{}
	eatingWg.Add(len(philosophers))

	seatingWg := &sync.WaitGroup{}
	seatingWg.Add(len(philosophers))

	forks := make(map[int]*sync.Mutex)
	for i := range philosophers {
		forks[i] = &sync.Mutex{}
	}

	for _, p := range philosophers {
		go diningProblem(p, forks, eatingWg, seatingWg)
	}

	eatingWg.Wait()
	color.Magenta(results.String())
}

func diningProblem(philosopher Philosopher, forks map[int]*sync.Mutex, eatingWg, seatingWg *sync.WaitGroup) {
	defer eatingWg.Done()

	color.Green("%s is seated at the table.", philosopher.name)
	seatingWg.Done()
	seatingWg.Wait()

	for i := 0; i < maxMeals; i++ {
		if philosopher.rightFork == 0 {
			forks[philosopher.rightFork].Lock()
			color.Yellow("\t%s takes the right fork", philosopher.name)
			forks[philosopher.leftFork].Lock()
			color.Yellow("\t%s takes the left fork", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			color.Yellow("\t%s takes the left fork", philosopher.name)
			forks[philosopher.rightFork].Lock()
			color.Yellow("\t%s takes the right fork", philosopher.name)
		}

		color.Yellow("\t%s has both forks and is eating", philosopher.name)
		time.Sleep(eatingTime)

		color.Yellow("\t%s has both forks and is thinking", philosopher.name)
		time.Sleep(thinkingTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()
		color.Yellow("\t%s has put down the forks", philosopher.name)
	}

	results.appendFinished(philosopher.name)
}
