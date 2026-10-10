package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	lock := make(chan bool, 1) //buffered channel with capacity 1, carrying bool values
	for i := 0; i < 7; i++ {
		wg.Add(1)
		go worker(i, lock, &wg)
	}
	wg.Wait()
	fmt.Println("DONE!")
}

func worker(id int, lock chan bool, group *sync.WaitGroup) {
	defer group.Done()
	fmt.Printf("Worker %d wants the lock \n", id)
	lock <- true //a worker sends a value into the lock channel if the channel is currently empty, else the send blocks.

	fmt.Printf("Worker %d has the lock \n", id)
	time.Sleep(500 * time.Millisecond)

	fmt.Printf("Worker %d releasing the lock \n", id)
	<-lock //the worker receives (empties) the value of the channel, making room for the next worker's lock <- true to succeed (release/unlock step)
}

//A buffered channel with capacity 1 is a mutex/binary semaphore — send (filling the buffer) ≃ acquiring the lock (like a P/Psem operation), receive (emptying the buffer) ≃ releasing the lock (like a V/Vsem operation).
//You could, if you wanted, directly rewrite this using the FIFO Semaphore monitor code from way back in the Monitors lecture, and it would behave identically in spirit — this is the same underlying synchronization mechanism,
//just expressed using Go's channel syntax instead of monitor syntax.
