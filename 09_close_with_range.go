package main

import "fmt"

// This program is identical to 08_close.go, but written more concisely using for i := range ch
func main() {
	ch := make(chan int)
	go send(ch) //start another goroutine that sends values to the channel

	//for i := range ch automatically receives values from ch, one at a time, assigning each to i and running the loop body
	//and it automatically stops (exits the loop) the moment the channel is closed and drained, without you needing to manually check an ok boolean yourself.
	//Go handles the "check ok, break if false" logic internally, behind the scenes, as part of what range means when applied to a channel.
	for i := range ch {
		fmt.Println(i)
	}

	fmt.Println("DONE!")
}

func send(ch chan<- int) {
	ch <- 1
	ch <- 2
	close(ch)
}
