package main

import (
	"fmt"
	"sync"
)

const producerCount int = 4
const consumerCount int = 3

var messages = [][]string{ //Each producer gets one nested list
	{"In", "brightest", "day", "in", "blackest", "night"},     //producer 0
	{"no", "evil", "shall", "escape", "my", "sight"},          //producer 1
	{"those", "who", "worship", "evil's", "might"},            //producer 2
	{"beware", "my", "power", "Green", "Lantern's", "light!"}, //producer 3
}

func main() {
	link := make(chan string) //no buffer size, so it is unbuffered and synchronous. Each send waits for a consumer to be ready

	//wg *sync.WaitGroup is passed as a pointer. if you passed it by value, each function would get
	//a copy with its own counter, and Done() would decrement the wrong one. All prties must share the same WaitGroup.
	wp := &sync.WaitGroup{}
	wc := &sync.WaitGroup{}

	wp.Add(producerCount)
	wc.Add(consumerCount)

	for i := 0; i < producerCount; i++ {
		go produce(link, i, wp)
	}
	for i := 0; i < consumerCount; i++ {
		go consume(link, i, wc)
	}

	wp.Wait()   //waits until all producers are finished, so nothing more will ever be sent
	close(link) //permanent "no more message"
	wc.Wait()   //waits until all consumer have exited, so every message has been processes before main ends.

	//ORDER MATTERS!
	//- Never closing would leave the consumer's range loops waiting forever, so wc.Wait() would never return. Deadlock.
	//- Closing before the producer finish would make a later producer send on a closed channel, which causes
	//  runtime panic. You drain a closed channel, but you cannot send on one.
	//
	//That is why you close only after wp.Wait(). This is also a Go convention: the sender side closes,
	//because only it knows when sending is finished-
	fmt.Println("DONE!")
}

// The producer gets a send-only view of the channel (compiler stop if we receive message here)
func produce(link chan<- string, id int, wg *sync.WaitGroup) {
	defer wg.Done() //Done() is registered at the top and fires when the function exits.
	for _, msg := range messages[id] {
		fmt.Printf("Producer %d put '%v' into the channel\n", id, msg)
		link <- msg
	}
}

// The consumer gets a receive-only view of the channel (compiler stops if we produce message here)
func consume(link <-chan string, id int, wg *sync.WaitGroup) {
	defer wg.Done()         //Done() is registered at the top and fires when the function exits.
	for msg := range link { //Consumer receives until this channel is closed and drained
		fmt.Printf("Message '%v' is consumed by consumer %v\n", msg, id)
	}
}

//Which consumer gets which message is non-deterministic, as with several senders on one channel.
//With four producers and three consumers on one shared channel, any consumer can pick up any message.
