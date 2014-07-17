// Mathematical row vector

package row

exclude +, *
import (
	"clojure/core"
	"funcgo/reference/contract"
)


func +(a, b) {
	map(core.+, a, b)
}

// dot product
func *(v1, v2) {
	contract.Require(func{ count(v1) == count(v2) })
	core.+  reduce  map(core.*, v1, v2)
}
