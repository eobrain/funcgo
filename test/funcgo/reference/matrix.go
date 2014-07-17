// Operations on matrices, stored as sequences of row vectors

package matrix
exclude ( +, * )
import (
	"funcgo/reference/contract"
	r "funcgo/reference/row"
)

rowCount := count
func colCount(m) { count(first(m)) }

func +(m1, m2) {
	map(r.+, m1, m2)
}

func Transpose(m) {
	firstColumnT := first map m
	if colCount(m) == 1 {
		 [firstColumnT]
	 } else {
		 firstColumnT cons Transpose(rest map m)
	 }
}




func *(m1, m2) {
	contract.Require(func{ colCount(m1) == rowCount(m2) })
	contract.Ensure(
		func{ rowCount($1)==rowCount(m1) && colCount($1)==colCount(m2) },

		for m1row := lazy m1 {
			for m2col := lazy Transpose(m2) {
				m1row  r.*  m2col
			}
		}
	)
}
