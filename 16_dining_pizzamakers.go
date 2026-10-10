package main

import (
	"fmt"
	"sync"
	"time"
)

const count = 5

var pWgroup sync.WaitGroup

// A PizzaTool is just a mutex. Embedding sync.Mutex means you can call .Lock() and
// .Unlock() directly on it. Only one goroutine can hold a tool at a time
type PizzaTool struct{ sync.Mutex }

// Each Pizzamaker has an id and pointers to two tools: one on their left, one on their right
type PizzaMaker struct {
	id             int
	leftPizzaTool  *PizzaTool
	rightPizzaTool *PizzaTool
}

func main() {
	// Create pizza tools
	pizzaTools := make([]*PizzaTool, count) //5 tools
	for i := 0; i < count; i++ {
		pizzaTools[i] = new(PizzaTool)
	}

	// Create pizza makers, assign them 2 tools and send them to the dining table
	pizzaMakers := make([]*PizzaMaker, count)
	for i := 0; i < count; i++ {
		pizzaMakers[i] = &PizzaMaker{
			id:             i,
			leftPizzaTool:  pizzaTools[i],
			rightPizzaTool: pizzaTools[(i+1)%count]} //The % count wraps around so maker 4 uses tools 4 and 0
		pWgroup.Add(1) //Add before making new goroutine
		go pizzaMakers[i].eat()
	}
	pWgroup.Wait() //makes main wait until all five finish, so it doesn't exit early
	fmt.Println("All pizzas done!")
}

//Maker i uses tool i and tool i+1. The % count wraps around so maker 4 uses tools 4 and 0. That means neighboring makers share a tool, and it’s a circle:
//Maker0: tools 0,1
//Maker1: tools 1,2
//Maker2: tools 2,3
//Maker3: tools 3,4
//Maker4: tools 4,0

func (p PizzaMaker) eat() {
	defer pWgroup.Done()

	// Fix for deadlock: the last maker picks up tools in reverse order, (breaking the symmetry)
	// which breaks the circular wait.
	left, right := p.leftPizzaTool, p.rightPizzaTool
	if p.id == count-1 {
		left, right = right, left //The last maker will pick up right tool first.
		//Maker 4 now picks 0 and then 4, instead of 4 and then 0.
	}
	for j := 0; j < 3; j++ {
		left.Lock()
		right.Lock()

		p.say("preparing pizza")
		time.Sleep(time.Second)

		right.Unlock()
		left.Unlock()

		p.say("finished preparing pizza")
		time.Sleep(time.Second)
	}
}

func (p *PizzaMaker) say(action string) {
	fmt.Printf("Pizza maker %d is %s\n", p.id, action)
}
