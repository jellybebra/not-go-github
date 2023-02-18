package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"time"
)

const DickSize = 13 // untyped const automatically takes the type.

var i, j int = 1, 2
var c, python, java bool

// var c, python, java = true, false, "no!"
// If an initializer is present, the type can be omitted

var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)

/*
int, uint, uintptr
	are 32 bits wide (32-bit system)
	and 64 bits wide (64-bit system).

Use "int" for integers
unless you have a reason
to use other int types.
*/

/*
bool

string

int  int8  int16  int32  int64
uint uint8 uint16 uint32 uint64 uintptr

byte // alias for uint8

rune // alias for int32
     // represents a Unicode code point

float32 float64

complex64 complex128
*/

func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}

func add(x int, y int) int {
	return x + y
}

func swap(x, y string) (string, string) {
	return y, x
}

func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}

func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim { // "v" is only in scope until the end of the if/else
		return v
	} else {
		fmt.Printf("%g >= %g\n", v, lim)
	}
	// can't use v here
	return lim
}

func main() {
	fmt.Println(
		pow(3, 2, 10),
		pow(3, 3, 20),
	)

	fmt.Println(sqrt(2), sqrt(-4))

	// intead of "var k int = 3"
	k := 3            // int
	m := 3.142        // float64
	l := 0.867 + 0.5i // complex128

	fmt.Printf("v is of type %T\n", k)
	fmt.Printf("v is of type %T\n", m)
	fmt.Printf("v is of type %T\n", l)

	fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)
	fmt.Printf("Type: %T Value: %v\n", MaxInt, MaxInt)
	fmt.Printf("Type: %T Value: %v\n", z, z)

	fmt.Println(i, c, python, java, k)

	fmt.Println("The time is", time.Now())
	rand.Seed(time.Now().UnixNano())
	fmt.Println(rand.Intn(10))

	fmt.Printf("Now you have %g problems.\n", math.Sqrt(7))
	fmt.Println(add(42, 13))
	fmt.Println(split(17))

	a, b := swap("hello", "world")
	fmt.Println(a, b)

	var i int
	var f float64
	var d bool
	var s string
	fmt.Printf("%v %v %v %q\n", i, f, d, s)

	var x, y int = 3, 4
	var r float64 = math.Sqrt(float64(x*x + y*y))
	var z uint = uint(r)
	fmt.Println(x, y, z)

	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	sum = 1
	for sum < 1000 {
		sum += sum
	}

	// infinite loop
	for {
		break
	}

}
