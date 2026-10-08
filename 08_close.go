package main

import "fmt"

func main() {
	ch := make(chan int)
	go send(ch)

	//inifite loop, it keeps receiving forever, until ok comes back as false.
	for {
		//Receiving from a channel can optionally return a second value, ok (a boolean), alongside the received value i.
		//This is analogous to the (result, error) pattern we have seen.
		//Receiving from a channel can tell us not just what was received, but also whether the channel is still open.
		i, ok := <-ch
		//ok == true (the channel is still open, or this is the last value sent before closing)
		//ok == false (the channel has been closed, and there are no more values left to receive)
		//- i in this case would be just 0 (because channels content here is int)
		if !ok {
			break
		}
		fmt.Println(i)
	}
}

// send() function sends 1, then 2, then calls close(ch) - this marks the channel
// as closed, signaling "no more values will ever be sent on this channel"
func send(ch chan<- int) {
	ch <- 1
	ch <- 2
	close(ch)
}

//This is a critical new concept: closing is different from just "nothing left to send right now." close(ch) is an explicit signal of "I am permanently done sending"
//— it lets a receiver distinguish between "the channel is temporarily empty, but more might come later" (in which case <-ch would simply block, waiting) versus
//"the channel is closed, and it's safe to stop trying to receive more" (in which case <-ch returns immediately with ok=false).
//This maps directly onto the EOS (End-of-Stream) concept from last lecture's Merge filter example — close is Go's built-in, language-level mechanism for exactly
//that "end of stream" signal, rather than needing to invent a special sentinel value like EOS yourself.
