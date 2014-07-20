package reference

import (
	test "midje/sweet"
	"funcgo/reference/matrix"
	"clojure/tools/logging"
	"cljLoggingConfig/log4j"
)
import type (
	java.util.ArrayList
	java.lang.Iterable
)
log4j.mutateSetLogger(LEVEL, WARN)

var a = 55
var b = 66

test.fact("Most things are Expression",

	{
		smaller := if a < b {
			a
		} else {
			b
		}
		smaller
	}, =>, 55,

	{
		digits  := [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
		squares := for d := lazy digits {
			d * d
		}
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
		vec          := [111, 222, 333, 444]
		[a, b, c, d] := vec
		b
	}, =>, 222,

	{
		vec := [111, 222, 333, 444]

		func theSecond([a, b, c, d]) {
			b
		}
		theSecond(vec)
	}, =>, 222,

	{
		vec              := [111, 222, 333, 444]
		[first, rest...] := vec
		rest
	}, =>, [222, 333, 444],

	{
		dict             := {AAA: 11,  BBB: 22,  CCC: 33,  DDD: 44}
		{c: CCC, a: AAA} := dict
		c
	}, =>, 33,

	{
		dict := {AAA: 11,  BBB: 22,  CCC: 33,  DDD: 44}

		func extractCCC({c: CCC}) {
			c
		}
		extractCCC(dict)
	}, =>, 33,

	{
		planets := [
			{NAME: "Mercury", RADIUS_KM: 2440},
			{NAME: "Venus",   RADIUS_KM: 6052},
			{NAME: "Earth",   RADIUS_KM: 6371},
			{NAME: "Mars",    RADIUS_KM: 3390}
		]
		[_, _, {earthRadiusKm: RADIUS_KM}, _] := planets
		earthRadiusKm
	}, =>, 6371
)

test.fact("Looping with tail recursion",

	{
		func sumSquares(vec) {
			if isEmpty(vec) {
				0
			} else {
				x := first(vec)
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
					x := first(v)
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
					x int := first(v)
					recur(accum + x * x, rest(v))
				}
			}
		}
		sumSquares([3, 4, 5, 10])
	}, =>, 150,

	loop(vec=[], count = 0) {
		if count < 10 {
			v := vec  conj  count
			recur(v, count + 1)
		} else {
			vec
		}
	}, =>, [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]

)

test.fact("Curly Brace Blocks",
	{
		product := {
			logging.info("doing the multiplication")
			100 * 100
		}
		product
	}, =>, 10000
)

test.fact("Type switch",

	{
		func plus(a, b) {
			switch a.(type) {
			case Number:   a + b
			case String:   a  str  b
			case Iterable: vec(a  concat  b)
			default:       str("Unknown types for ", a, " and ", b)
			}
		}

		[
			2       plus  3,
			0.5     plus  0.75,
			[P, Q]  plus  [R, S, T],
			"foo"   plus  "bar",
                        FOO     plus  BAR
		]
	}, =>, [
		5,
		1.25,
		[P, Q, R, S, T],
		"foobar",
		"Unknown types for :foo and :bar"
	]
)

test.fact("select (1)",

	{
		c1 := make(chan, 1)
		c2 := make(chan, 1)
		thread {
			Thread::sleep(10)
			c1 <- 111
		}
		c2 <- 222
		select {
		case x = <-c1:
			x * 100
		case x = <-c2:
			x * 100
		}
	}, =>, 22200,

	{
		c1 := make(chan, 1)
		c2 := make(chan, 1)
		go func(){
			Thread::sleep(10)
			c1 <- 111
		}()
		c2 <- 222
		select {
		case x = <-c1:
			x * 100
		case x = <-c2:
			x * 100
		}
	}, =>, 22200
)

test.fact("select (2)",
	{
		c1 := make(chan, 1)
		c2 := make(chan, 1)
		go {
			for i := times(10000) { var x = i }
			c1 <: 111
		}
		go {
			c2 <: 222
		}
		<-go {
			select {
			case x = <:c1:
				x * 100
			case x = <:c2:
				x * 100
			}
		}
	}, =>, 22200,

	{
		c1 := make(chan)
		c2 := make(chan)
		go func(){
			Thread::sleep(10)
			<-c1
		}()
		go func(){
			<-c2
		}()
		select {
		case c1 <- 111:
			"wrote to c1"
		case c2 <- 222:
			"wrote to c2"
		}
	}, =>, "wrote to c2",

)

test.fact("select (3)",
	{
		c1 := make(chan)
		c2 := make(chan)
		go {
			<:c2
		}
		go {
			for i := times(10000) { var x = i }
			<:c1
		}
		<-go {
			select {
			case c1 <: 111:
				"wrote to c1"
			case c2 <: 222:
				"wrote to c2"
			}
		}
	}, =>, "wrote to c2",

	{
		c1 := make(chan, 1)
		c2 := make(chan)
		thread {
			Thread::sleep(20)
			c1 <- 111
		}
		thread {
			Thread::sleep(10)
			<-c2
		}
		select {
		case x = <-c1:
			x * 100
		case c2 <- 222:
			"wrote to c2"
		default:
			"nothing ready"
		}
	}, =>, "nothing ready",

)

test.fact("inline",

	str("foo", "bar"),
	=>, "foobar",

	"foo"  str  "bar",
	=>, "foobar"

)

func truthTable(op) {
	[
		false op false,
		false op true,
		true  op false,
		true  op true
	]
}

test.fact("operators",

        3 * 4                   , =>, 12,
        16.0 / 2.0              , =>, 8.0,
	12 % 5                  , =>, 2,
        0xCAFE << 4             , =>, 0xCAFE0,
        0xCAFE >> 4             , =>, 0xCAF,
        0xFACADE &  0xFFF000    , =>, 0xFAC000,
        0xFACADE &^ 0x000FFF    , =>, 0xFAC000,
	3 + 4                   , =>, 7,
        3 - 4                   , =>, -1,
        0xFACADE |  0xFFF000    , =>, 0xFFFADE,
        0xFACADE ^  0x000FFF    , =>, 0xFAC521,
        5 == 5                  , =>, true,
        5 == 4                  , =>, false,
        5 == "5"                , =>, false,
        "5" == "5"              , =>, true,
	[A, B, C] == [A, B, C]  , =>, true,
	[A, B, C] == [A, B, DD] , =>, false,
	{A:1,B:2} == {A:1,B:2}  , =>, true,
	{A:1,B:2} == {A:1,B:9}  , =>, false,
        5 > 5                   , =>, false,
        5 > 4                   , =>, true,
        5 < 5                   , =>, false,
        5 < 4                   , =>, false,
        5 >= 5                  , =>, true,
        5 >= 4                  , =>, true,
        5 <= 5                  , =>, true,
        5 <= 4                  , =>, false,
	truthTable(func{$1 && $2}), =>, [false, false, false, true],
	truthTable(func{$1 || $2}), =>, [false,  true,  true, true]
)

var (
	a = randInt(100)
	b = randInt(100)
	c = randInt(100)
	p = randInt(2) == 0
	q = randInt(2) == 0
	r = randInt(2) == 0
)

test.fact("precedence",
	^a * b         , =>, (^a) * b,
	a * b - c      , =>, (a * b) - c,
	a + b < c      , =>, (a + b) < c,
	a < b && b < c , =>, (a < b) && (b < c),
	p && q || r    , =>, (p && q) || r,
	p || q  str  r , =>, (p || q)  str  r
)

test.fact("vars",
	{
		var (
			pp = 111
			qq = 222
		)
		pp + qq
	}, =>, 333,

	{
		var rr = 111
		var ss = 222
		pp + qq
	}, =>, 333,

	{
		var tt int    = 111
		var uu string = "foo"
		uu  str  tt
	}, =>, "foo111"
)

test.fact("for",

	{
		fib        := [1, 1, 2, 3, 5, 8]
		fibSquared := for x := lazy fib {
			x * x
		}
		fibSquared
	}, =>, [1, 1, 4, 9, 25, 64],

	{
		fib        := [1, 1, 2, 3, 5, 8]
		fibSquared := func(x){ x * x }  map  fib
		fibSquared
	}, =>, [1, 1, 4, 9, 25, 64],

	withOutStr({
		fib := [1, 1, 2, 3, 5, 8]
		for x := lazy fib {
			print(" ", x)
		}
	}), =>, "",

	withOutStr({
		fib := [1, 1, 2, 3, 5, 8]
		func(x){ print(" ", x) }  map  fib
	}), =>, "",

	withOutStr({
		fib := [1, 1, 2, 3, 5, 8]
		for x := range fib {
			print(" ", x)
		}
	}), =>, "  1  1  2  3  5  8",

	withOutStr({
		for x := times 10 {
			print(" ", x)
		}
	}), =>, "  0  1  2  3  4  5  6  7  8  9"

)


test.fact("exceptions",
	{
		try {
			throw(new AssertionError("foo"))
		} catch OutOfMemoryError e {
			"out of memory"
		} catch AssertionError e {
			"assertion failed: "  str  e->getMessage()
		} finally {
			"useless"
		}
	}, =>, "assertion failed: foo",

	{
		mutex     := make(chan, 1)
		dangerous := new ArrayList()
		dangerous->add(0)
		mutex <- true  // initialize mutex
		for _ := times 1000 {
			thread {
				<-mutex   // grab mutex
				try {
					i := dangerous->get(0)
					dangerous->set(0, i + 1)
				} finally {
					mutex <- true   // release mutex
				}
			}
		}
		Thread::sleep(100)
		dangerous->get(0)
	}, =>, 1000
)

test.fact("you can transpose matrices",
	{
		m := [
			[1, 2, 3],
			[4, 5, 6]
		]

		matrix.Transpose(m)

	}, =>, [
		[1, 4],
		[2, 5],
		[3, 6]
	]
)


test.fact("you can multiply matrices together",
	{
		a := [[3, 4]]
		b := [
			[5],
			[6]
		]

		a  matrix.*  b

	}, =>, [[39]],

	{
		m := [
			[1, 2, 3],
			[4, 5, 6]
		]
		mT := [
			[1, 4],
			[2, 5],
			[3, 6]
		]

		m  matrix.*  mT

	}, =>, [
		[14, 32],
		[32, 77]
	]
)




