package main

import "fmt"

// A buffered Go channel, with its fixed capacity, behaves almost exactly like the asnchronous channel model.
// A buffered channel with capacity N is conceptually a bounded buffer of size N.
func main() {
	ch := make(chan int, 2) //Creates a buffered channel with a capacity of 2
	//The channel can hold up to 2 values "in flight" without requiring and immediate receiver
	ch <- 1 //does not block
	ch <- 2 //does not block
	//send only blocks once the buffer is ful, since the buffer here has capacity 2 and we only send 2 values, neither send blocks.
	//Both values simply sit in the channel's internal buffer, waiting to be picked up later, with no receiver needed to be
	//ready at the moment of sending.
	fmt.Println(<-ch) //1 blocks (FIFO order, as always) b
	fmt.Println(<-ch) //2 blocks

	//Once the buffer already holds 2 values, its full. The third send, ch <- 3, must now block.
	//The program is stuck waiting for a receive that will never happen (because the code that
	//would do the receiving is unreachable, blocked behind the stuck send). We are running a single goroutine (main itself)
	ch <- 1
	ch <- 2
	ch <- 3 //deadlock
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

//Buffered channels = finite asycn like behavior
//Unbuffered channels = fully synchronous (rendezvous) behavior

//Why this matters — this is a clean, minimal illustration of buffered-channel deadlock. This is structurally identical to the Bounded Buffer "full buffer, no consumer available" scenario
//— except here, made worse by the fact that it's single-threaded (no separate consumer goroutine exists at all to eventually drain the buffer), so there's no possible way for this
//program to ever make progress once it hits that third send.
