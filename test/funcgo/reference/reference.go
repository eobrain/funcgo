package reference

import (
	test "midje/sweet"
)
import type (
	java.util.logging.Logger
)

const (
	a = 55
	b = 66
	log = Logger::getLogger(str(\`*ns*`))
)

test.fact("Most things are Expression",

	{
		const smaller = if a < b {
			a
		} else {
			b
		}
		smaller
	}, =>, 55,

	{
		const (
			digits = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
			squares = for d := lazy digits {
				d * d
			}
		)
		squares
	}, =>, [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
)

test.fact("syntax",

	withOutStr(
		if a < b {
			println("Conclusion:")
			println(a, "is smaller than", b)
		}
	), =>, `Conclusion:
55 is smaller than 66
`,
	withOutStr(
		if a < b { println("Conclusion:"); println(a, "is smaller than", b) }
	), =>, `Conclusion:
55 is smaller than 66
`
)

test.fact("can destructure",

	{
		const (
			vec = [111, 222, 333, 444]
			[a, b, c, d] = vec
		)
		b
	}, =>, 222,

	{
		const vec = [111, 222, 333, 444]

		func theSecond([a, b, c, d]) {
			b
		}
		theSecond(vec)
	}, =>, 222,

	{
		const (
			vec = [111, 222, 333, 444]
			[first, rest...] = vec
		)
		rest
	}, =>, [222, 333, 444],

	{
		const (
			dict = {AAA: 11,  BBB: 22,  CCC: 33,  DDD: 44}
			{c: CCC, a: AAA} = dict
		)
		c
	}, =>, 33,

	{
		const dict = {AAA: 11,  BBB: 22,  CCC: 33,  DDD: 44}

		func extractCCC({c: CCC}) {
			c
		}
		extractCCC(dict)
	}, =>, 33,

	{
		const (
			planets = [
				{NAME: "Mercury", RADIUS_KM: 2440},
				{NAME: "Venus",   RADIUS_KM: 6052},
				{NAME: "Earth",   RADIUS_KM: 6371},
				{NAME: "Mars",    RADIUS_KM: 3390}
			]
			[_, _, {earthRadiusKm: RADIUS_KM}, _] = planets
		)
		earthRadiusKm
	}, =>, 6371
)

test.fact("Looping with tail recursion",

	{
		func sumSquares(vec) {
			if isEmpty(vec) {
				0
			} else {
				const x = first(vec)
				x * x + sumSquares(rest(vec))
			}
		}
		sumSquares([3, 4, 5, 10])
	}, =>, 150,


	{
		func sumSquares(vec) {
			func sumSq(accum, v) {
				if isEmpty(v) {
					accum
				} else {
					const x = first(v)
					recur(accum + x * x, rest(v))
				}
			}
			sumSq(0, vec)
		}
		sumSquares([3, 4, 5, 10])
	}, =>, 150,

	{
		func sumSquares(vec) {
			loop(accum=0, v=vec) {
				if isEmpty(v) {
					accum
				} else {
					const x int = first(v)
					recur(accum + x * x, rest(v))
				}
			}
		}
		sumSquares([3, 4, 5, 10])
	}, =>, 150,

	loop(vec=[], count = 0) {
		if count < 10 {
			const v = vec  conj  count
			recur(v, count + 1)
		} else {
			vec
		}
	}, =>, [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]

)

test.fact("Curly Brace Blocks",
	{
		const product = {
			log->info("doing the multiplication")
			100 * 100
		}
		product
	}, =>, 10000
)
