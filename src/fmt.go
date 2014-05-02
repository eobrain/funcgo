package fmt
import "clojure/string"

//var Println = println
var Printf = printf

func stripZeroAfterDecimal(s) { string.replace(s, /^([1-9][0-9]*)\.0$/, "$1") }

func Println(args...) {
	const as = func{stripZeroAfterDecimal(str(..))} map args
	println(" " string.join as)
}
