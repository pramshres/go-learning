package main

import (
	"fmt"
	"math"
)

type AbsI interface{ Abs() float64 }

type Pair struct{ X, Y float64 }
type Triple struct{ X, Y, Z float64 }

func (x *Pair) Abs() float64 {
	return math.Sqrt(x.X*x.X + x.Y*x.Y)
}
func (x *Triple) Abs() float64 {
	return math.Sqrt(x.X*x.X + x.Y*x.Y + x.Z*x.Z)
}

func main() {
	//var a AbsI - declares a variable a of the interface type AbsI.
	var a AbsI //must contain something that implements AbsI (Whatever gets assigned to a must have an Abs() float64 method)
	pair := Pair{1, 2}
	triple := Triple{3, 4, 5}

	a = &pair            //assigning a pointer to a Pair
	fmt.Println(a.Abs()) // Pair's Abs()

	a = &triple          //assigning a pointer to a Triple
	fmt.Println(a.Abs()) // Triple's Abs()

	//a = pair 	//a Pair is not ok, except pair := &Pair{1,2}
	//a = pair — this is not okay (a compile error)! Why? Because Abs() was defined specifically on Pair (a pointer to Pair), not on plain Pair itself.
	//A plain, non-pointer Pair value does not automatically have the Abs() method attached to it — only pointers to Pair do.
	//The comment clarifies the fix: pair := &Pair{1,2} would make pair itself a pointer from the start, which would work.
}

//Why this matters — this is a subtle but genuinely important Go-specific detail: whether a method is attached via a pointer receiver (*Pair) or a value receiver (Pair)
//actually changes which values satisfy a given interface! This isn't just a style choice — it has real consequences for interface satisfaction, as this example demonstrates directly.
