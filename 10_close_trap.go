package main

import "fmt"

// Receiving from a closed, empty channel does not block and does not panic/crash
// — it returns immediately, giving you the zero value for the channel's type (here, 0 for int).
// Since this particular receive expression (fmt.Println(<-ch)) doesn't check the ok boolean at all
// (it's just using the single-value form, <-ch, not i, ok := <-ch), you'd never actually know,
// just from this code, that the 0 you got back is a "fake" zero-value rather than a genuinely sent 0.
func main() {
	ch := make(chan int)
	go send(ch)
	fmt.Println(<-ch) //1
	fmt.Println(<-ch) //0 (not real sent value)
}

func send(ch chan<- int) {
	ch <- 1
	close(ch)
}

//Trap: If you only ever use the single-value receive form (<-ch) and never check ok, you cannot distinguish between
//"I received a genuinely sent zero" and "the channel is closed and I'm just getting the type's zero value by default."
//This is a real source of subtle bugs in Go code — the lesson is: if you need to know whether a channel is closed,
//you must use the two-value form (i, ok := <-ch) or range — never rely on the plain <-ch form alone if closing matters to your program's correctness.

//Every type in Go has a defined zero value — the default value a variable of that type gets when it's declared without explicit initialization
//(or, as we just saw, when you receive from a closed, empty channel). Here's the pattern across common types:
//
// Type                    Zero value
// ----                    ----------
// int, float64             0
// string                   "" (empty string)
// bool                     false
// pointer (*Pair, etc.)    nil
// error                    nil
// struct (e.g., Pair)      All fields set to their zero values.
//                          Example: Pair{} == Pair{X: 0, Y: 0}
