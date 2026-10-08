package main

import (
	"fmt"
	"os"
	"sync"
)

//NEW PROBLEM: in 05_checklink.go we had a convenient way to know when all goroutine were done - we just received exactly len(links) values from the channel
//But what if our goroutines don't naturally produce a result value to receive (e.g., they just perform som side-effecting work, like compressing files - nothing to send back)
//SOLUTION: sync.WaitGroup - a synchronization primitive from Go's sync package, used to avoid needing channels for this specific "wait until a group of goroutines are all done" pattern.

//Three methods:
//* Add  - adds a counter to the WaitGroup
//* Done - decrements the counter by one
//* Wait - called from the main goroutine, blocks execution until the counter reaches zero.

//USAGE: go run 06_waitgroup.go file1.txt file2.txt file3.txt

func main() {
	var wg sync.WaitGroup //a WaitGroup is declared like an ordinary variable - its "zero value"
	var i int = -1
	var file string
	for i, file = range os.Args[1:] { //os.Args[1:], Go's way of accessing command-line arguments (os.Args[0] is always the program's own name)
		fmt.Println("File: " + file)
		wg.Add(1)              //add before asynchronous call! CRITICAL!
		go func(file string) { //anonymous function (a function with no name, defined right here on the spot)
			compress(file)
			wg.Done()
		}(file)
	}
	wg.Wait() //the main goroutine blocks here until every launched goroutine has called wg.Done()
	fmt.Printf("compressed %d files \n", i+1)
}

func compress(filename string) error {
	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close() //defer schedules in.Close() to run automatically, right when compress is about to return (cleanup)

	out, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer out.Close() ////defer schedules out.Close() to run automatically, right when compress is about to return (cleanup)
	//gzout := gzip.NewWriter(out)
	//read file and write compressed content
	return nil
}

//Why wrap compress(file) in an anonymous function instead of calling go compress(file) directly? Because we need to call wg.Done() right after compress finishes
//— but compress itself doesn't know anything about wg. By wrapping both calls (compress(file) then wg.Done()) inside one anonymous function, we create a small,
//self-contained unit of work that does the compression and reports completion, and go launches this entire bundle as one goroutine.

//Why does the anonymous function take file as its own parameter, rather than just directly referencing the outer file variable from the loop? This is a subtle,
//genuinely important Go gotcha worth understanding carefully: the loop variable file is reused across iterations (it's the same variable, just reassigned each time around the loop)
//— recall it was declared outside the loop. If the goroutine's anonymous function directly captured (closed over) this outer file variable without taking it as its own parameter,
//then by the time the goroutine actually runs (which might be slightly delayed, since goroutines don't necessarily start executing immediately), the loop could have already moved on
//and changed file to the next filename! This would cause multiple goroutines to accidentally compress the same (wrong, most-recently-assigned) file, rather than each one handling its
//own distinct file. By passing file as an explicit parameter to the anonymous function — and immediately calling it with (file) at the point of the go statement
//— Go evaluates file's current value right then and there, and gives each goroutine its own private, frozen copy, safe from being overwritten by later loop iterations.

//Why wg.Add(1) must happen before go func(...): Think about what could go wrong if you added after launching the goroutine instead. If the main goroutine reaches wg.Wait()
//before some slower-to-start goroutine has gotten around to calling its own wg.Add(1), the WaitGroup's counter might still read 0 at that moment (since that particular Add hasn't executed yet)
//— so Wait() would see "counter is 0" and incorrectly conclude everything is already done, returning immediately, even though that goroutine hasn't even started its work yet! Adding before
//launching the goroutine guarantees the counter correctly reflects "how many goroutines are currently outstanding" at every point in time, with no race condition on the counter itself.

//defer is a general-purpose Go keyword for delaying any function call until the surrounding function is about to return. It has nothing to do with concurrency or WaitGroups at all
//— you could use defer in a program that has no goroutines, no channels, no WaitGroups whatsoever. Its only job is: "run this call later, right before this function exits, no matter how it exits."
