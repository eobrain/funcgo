package main5

import (
    "fmt"
    "math"
)

func main() {
    fmt.Printf("Now you have %g problems.",
        math.Nextafter(2, 3))
}
