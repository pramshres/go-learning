package main

import "fmt"

type Pair struct{ X, Y float64 }

// Pair appears inside Triple but without a field name. A normal struct field would look like somePair Pair (a name, then a type)
// Here, it's just Pair by itself. The comment explains: "no variable name -> implicit field." This is Go's embedding syntax
// — you're embedding the entire Pair struct inside Triple, and Go automatically gives it an implicit field name (which happens to be Pair, matching the type name).
type Triple struct {
	Pair //no variable name -> implicit field
	Z    float64
}

func main() {
	//constructs a Triple by providing a Pair value (for the embedded field) and a float64 value (for Z).
	triple := Triple{Pair{1, 2}, 3}
	fmt.Println(triple.X)
	//Notice: we access triple.X directly, even though X isn't literally a field of Triple. it's a field of the embedded Pair.
	//Go automatically "promotes" the embedded struct's fields (and, though not shown here, its methods too) up to the outer struct,
	//so you can access them as if they belonged directly to Triple — you don't have to write the more verbose triple.Pair.X (though that would also work).
}
