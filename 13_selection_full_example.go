package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	//time.After is a Go standard-library function that returns a channel,
	//which will automatically receive a value after the specified duration has elapsed
	//(30 seconds here). This is a clever trick: Go represents "a timer going off" as just another channel event,
	//letting you use the exact same select mechanism to wait for a timeout alongside waiting for real data.
	done := time.After(30 * time.Second)
	echo := make(chan []byte)
	go readStdin(echo)

	for {
		select {
		case buf := <-echo:
			if _, err := os.Stdout.Write(buf); err != nil {
				fmt.Println("write error:", err)
				os.Exit(1)
			}
		case <-done:
			fmt.Println("Timer ran out")
			os.Exit(0)
		}
	}
	fmt.Println("DONE!")
}

func readStdin(out chan<- []byte) {
	for {
		data := make([]byte, 1024)
		l, err := os.Stdin.Read(data)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("read error:", err)
			return
		}
		if l > 0 {
			out <- data
		}
	}
}

//Why this matters — this demonstrates select's real power: racing a "real" operation against a timeout, cleanly.
//Without select, implementing "wait for input, but give up after 30 seconds if nothing comes" would be awkward
//— you'd need some separate, manual timer-checking mechanism. With select, you simply treat the timeout as just another channel (thanks to time.After)
//and let select naturally race it against your actual data channel, reacting to whichever happens first.
