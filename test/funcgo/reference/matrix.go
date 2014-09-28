// Operations on matrices, stored as sequences of row vectors
package matrix
import "clojure/core"
exclude ( +, * )

// Begin private functions

func colCount(m) { count(first(m)) }
func dotProduct(v1, v2) {
	core.+  reduce  map(core.*, v1, v2)
}
func vecSum(a, b) { map(core.+, a, b) }

// Begin exported functions

func +(m1, m2) { map(vecSum, m1, m2) }

func Transpose(m) {
	firstColumnT := first  map  m
	if colCount(m) == 1 {
		 [firstColumnT]
	 } else {
		 firstColumnT  cons  Transpose(rest  map  m)
	 }
}

func *(m1, m2) {
	for m1row := lazy m1 {
		for m2col := lazy Transpose(m2) {
			m1row  dotProduct  m2col
		}
	}
}
