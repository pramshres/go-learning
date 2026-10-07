package main

import (
	"fmt"
	"math"
)

type Pair struct{ X, Y float64 } //declares a new type Pair, which is a struct (a plain record) with two fields (both of type float64)

func main() {
	var pair1 Pair      //declares a variable pair1 of type Pair (since no value is give, it gets the zero value for a Pair {0,0}.
	pair1 = Pair{3, 4}  //assigns a new Pair value to pair1
	pair2 := Pair{1, 2} //:= operator is Go's short variable declaration.
	//It declares and initializes a variable in one step.

	var result = pair1.Abs() + pair2.Abs()
	fmt.Println(result)
}

// This defines a method Abs that's attached to the type Pair (a pointer to Pair).
// The (pair *Pair) part before the function name is called the receiver. Here we are
// attaching methods to a type without needing a class definition.
func (pair *Pair) Abs() float64 {
	return math.Sqrt(pair.X*pair.X + pair.Y*pair.Y)
}

//Good to know: `*Pair` (pointer receiver) vs. just `Pair` (value receiver) is a meaningful distinction in Go
//(pointer receivers can modify the original struct; value receivers get a copy)
