package main

import (
	"fmt"
	"net/http"
)

// This is the concurrent fix to the sequential problem from Slide 14: all four checklink goroutines start up essentially at the same time and make their HTTP requests in parallel,
// rather than one after another. The total time this program takes is now roughly "however long the slowest single request takes" (since they're all happening simultaneously),
// rather than "the sum of all four requests' times" (as it would be sequentially).
func main() {
	links := []string{
		"https://google.com",
		"https://facebook.com",
		"https://golang.org",
		"https://stackoverflow.com",
	}
	c := make(chan string)
	//for loop that iterates over every link in the list, and for each one, launches a separate goroutine running checklink(link, c)
	//It fires off all four goroutines immediately, one after another with minimal delay between each launch.
	for _, link := range links { //The _ in for _, link is Go's "I don't care about this value" placeholder. Here we could get an index.
		go checklink(link, c)
	}
	//Second loop, which runs exactly 4 times (once per link), each time receiving a value from channel c (<-c) and printing it
	//The second loop receives results in whatever order they happen to arrive on the channel - not necessarily in the original order of links!
	//If facebook.com happens to be faster than goodle.com, its result could print first, even though google.com was listed first. Order of arrival is non-deterministic
	for i := 0; i < len(links); i++ {
		fmt.Println(<-c)
	}
	fmt.Println("DONE!")
}

func checklink(link string, c chan string) {
	_, err := http.Get(link) //pairs as language builtins - here the _ is discarded, as we don't need the HTTP response body for this simple check.
	// Returning (result, error) pairs is a core, idiomatic Go pattern for error handling
	if err != nil {
		c <- link + " is down!" //sending to the channel
		return
	}
	c <- link + " is up!" //sending to the channel
}

//Good to know: We have written in second loop that we are only supposed to be receive exactly 4 values from the channel (len(links))
//After we receive all values, the program ends - we are telling the main goroutine to wait for exactly 4 values before ending.
//If we are to put len(links)+1 in the second loop, telling it to wait for 5 values, we end up in a deadlock.
