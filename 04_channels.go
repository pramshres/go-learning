package main

import (
	"fmt"
)

// The problem with this main function is that we are creating two goroutines that run together with main.
// If main reaches the end before goroutines have sent or received, the program would end and the goroutines would be killed
// In order for us to receive the value before ending, we can for example make the receivef() blocking by removing go.
func main() {
	chl := make(chan float64) //creates a new channel that carries float64 values (pain, bidirectional channel type)
	go sendf(chl)             //launches a goroutine (non-blocking because go)
	go receivef(chl)          //launches a goroutine (non-blocking because go)
	//receivef(chl)
}

// send-only channel - this function is only allowed to send on ch, not receive.
func sendf(ch chan<- float64) {
	ch <- 0.5 //0.5 is sent on the channel
}

// receive-only channel - this function can only receive from ch, not send
func receivef(ch <-chan float64) {
	v := <-ch //<-ch evaluates to the next value taken off the channel, assigned to v.
	fmt.Println(v)
}
