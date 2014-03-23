package funcgo.compiler_test
import (
	test midje.sweet
	fgo funcgo.core
)

test.fact("smallest complete program has no import and a single expression",
	fgo.funcgoParse("package foo;import ()12345"),
	=>,
	"(ns foo (:gen-class))

12345
")

test.fact("Can use newlines instead of semicolons",
	fgo.funcgoParse(`
package foo
import (
)
12345
`),
	=>,
	`(ns foo (:gen-class))

12345
`)

test.fact("package can be dotted",
	fgo.funcgoParse("package foo.bar;import ()12345"),
	=> ,
	`(ns foo.bar (:gen-class))

12345
`)

test.fact("can import other packages",
	fgo.funcgoParse(`
package foo
import(
  b bar
)
12345
`),
	=>,
	`(ns foo (:gen-class)
  (:require [bar :as b]))

12345
`)

func parse(expr) {
  fgo.funcgoParse("package foo;import ()" str expr)
}

func parsed(expr) {
	str("(ns foo (:gen-class))\n\n", expr, "\n")
}

test.fact("can refer to symbols",
	parse("a"), =>, parsed("a"),
	parse("foo"), =>, parsed("foo")
)

test.fact("vector",
	parse("[]"), =>, parsed("[]"),
	parse("[a]"), =>, parsed("[a]"),
	parse("[a,b]"), =>, parsed("[a b]"),
	parse("[a,b,c]"), =>, parsed("[a b c]"),
	parse("[foo,bar]"), =>, parsed("[foo bar]"),
	parse(" [ a, b, c ]"), =>, parsed("[a b c]"),
	parse(" [ a , b , c  ]"), =>, parsed("[a b c]"),
	parse(" [   a  , b,     c ]  "), =>, parsed("[a b c]")
)

test.fact("escaped identifier",
	parse(`\range`), =>, parsed("range"),
	parse(`\for`), =>, parsed("for")
)

test.fact("multiple expressions inside func",
	parse(`func(){if c {d}}`),    =>, parsed(`(fn [] (when c d))`),
	parse(`func(){b;c}`),         =>, parsed(`(fn [] b c)`),
	parse(`func(){b;if c {d}}`), =>, parsed(`(fn [] b (when c d))`)
)

test.fact("subsequent const nests",
	parse(`const(a=1){x;const(b=2)y}`), =>, parsed(`(let [a 1] x (let [b 2] y))`))


//	parse(``), =>, parsed(``),
