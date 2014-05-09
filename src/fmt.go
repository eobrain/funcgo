package fmt
import "clojure/string"
import type clojure.lang.PersistentVector

//var Println = println
var Printf = printf

func stripZeroAfterDecimal(s) { string.replace(str(s), /^([1-9][0-9]*)\.0$/, "$1") }


// Create a string in the same way Go does.
func toStringLikeGo(x) {
	switch x.(type) {
	case String:           x
	case Number:           stripZeroAfterDecimal(x)
	case PersistentVector: str(
		"[", 
		" " string.join (toStringLikeGo map x),
		"]"
	)
	default:               str(x)
	}
}

func Println(args...) {
	const as = func{toStringLikeGo(..)} map args
	println(" " string.join as)
}
