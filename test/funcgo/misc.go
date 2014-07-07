package misc

import test "midje/sweet"

// func mostFrequentN(n, items) {
// 	->>(
// 		items,
// 		frequencies,
// 		sortBy(val),
// 		reverse,
// 		take(n),
// 		map(first)
// 	)
// }

func mostFrequentN(n, items) {
	first  map  (n  take  reverse(val  sortBy  frequencies(items)))
}

test.fact("can find 2 most common items in a sequence",
	mostFrequentN(2, ["a", "bb", "a", "x", "bb", "ccc", "dddd", "dddd", "bb", "dddd", "bb"]),
	=>, ["bb", "dddd"]
)
