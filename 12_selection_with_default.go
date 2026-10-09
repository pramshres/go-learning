package main

import "fmt"

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go send(ch1) //goroutine sending value on channel 1
	go send(ch2) //goroutine sending value on channel 2

	//When a default case is present, select NO longer blocks at all.
	//If none of the channel cases are ready right at this exact moment (nobody has a value available to receive yet, on either channel),
	//select immediately falls through to the default case instead of waiting around.
	select {
	case i1 := <-ch1:
		fmt.Printf("first call received: %d \n", i1)
	case i2 := <-ch2:
		fmt.Printf("second call received: %d \n", i2)
	default:
		fmt.Println("I just went to default, didn't wait")
	}
	fmt.Println("DONE!")
}

func send(ch chan<- int) {
	ch <- 1
}

//default makes select non-blocking; without it, select blocks until at least one case is ready.
