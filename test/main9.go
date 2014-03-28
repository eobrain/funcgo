package main9

import "fmt"

func swap(x, y String) {
    [y, x]
}

func main() {
	const [a, b] = swap("hello", "world")
	fmt.Println(a, b)
}
