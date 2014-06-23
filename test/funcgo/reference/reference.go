package reference

import (
	test "midje/sweet"
)

const (
	a = 55
	b = 66
)

test.fact("Everything is an Expression",

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
