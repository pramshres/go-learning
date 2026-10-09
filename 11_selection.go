package main

import "fmt"

// Waiting on several channels in parallel. The select statement can watch multiple channels (zero or more).
// Until something happens, it will wait (or execute a default statement, if supplied). When a channel has an event,
// the select statement will execute that event.
func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go send(ch1) //goroutine sending value on channel 1
	go send(ch2) //goroutine sending value on channel 2

	//select blocks until at least one of its cases becomes ready on either ch1 or ch2
	//and then it executes whichever case became ready, running that case's associated code
	select {
	case i1 := <-ch1:
		fmt.Printf("first call received: %d \n", i1)
	case i2 := <-ch2:
		fmt.Printf("second call received: %d \n", i2)
	}
	fmt.Println("DONE!")
}

func send(ch chan<- int) {
	ch <- 1
}

//If both channels happen to become ready around the same time (a genuine race between them), select picks one at random among the ready cases
//— this isn't specified as deterministic in any particular order; it's a fair, non-deterministic choice among whichever cases are currently ready.

//A useful analogy: select is like waiting at a two-window bank counter, ready to walk up to whichever window becomes free first
//— you're not committed in advance to window 1 or window 2; you react to whichever becomes available.
