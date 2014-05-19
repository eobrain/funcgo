package wiki
import(
        test "midje/sweet"
)

test.fact("Anonymous Functions",
	func(x){x * x}(3),
	=>, 9,

	map(func{list($1, inc($2))}, [1, 2, 3], [1, 2, 3]),
	=>, [[1, 2], [2, 3], [3, 4]],

	map(func(x, y){list(x, inc(y))}, [1, 2, 3], [1, 2, 3]),
	[[1, 2], [2, 3], [3, 4]],

	func{list($1, inc($1))} map [1, 2, 3],
	=>, [[1, 2], [2, 3], [3, 4]],

	func(x){list(x, inc(x))} map [1, 2, 3],
	=>, [[1, 2], [2, 3], [3, 4]],

	func{str apply $*}("Hello"),
	=>, "Hello",

	func{str apply $*}("Hello", ", ", "World!"),
	=>, "Hello, World!"
)
