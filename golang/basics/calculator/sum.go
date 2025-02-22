package calculator

import "fmt"

var offset float64 = 1
var Offset float64 = 1

func Sum(a, b float64) float64 {
	fmt.Println("multiply: ", multiply(a, b))
	fmt.Println("multiply: ", Multiply(a, b))
	return a + b + offset
}
