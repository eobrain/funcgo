package main10

import "fmt"

func split(sum long) {
	const(
		x = sum * 4 / 9
		y = sum - x
	)
	[x, y]
}

func main() {
    fmt.Println(split(17))
}
