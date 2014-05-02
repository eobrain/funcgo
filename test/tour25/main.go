package main

import (
	"fmt"
)

func newton(n int, x, z float64) float64 {
	if n == 0 {
		z
	} else {
		newton(n-1, x, z-(z*z-x)/(2*x))
	}
}

func Sqrt(x float64) float64 {
	return newton(500, x, x/2)
}

func main() {
	fmt.Println(Sqrt(100))
}
