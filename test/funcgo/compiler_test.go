package funcgo/compiler_test
import (
        test "midje/sweet"
        fgo "funcgo/core"
)

test.fact("smallest complete program has no import and a single expression",
        fgo.funcgoParse("package foo;import ()12345"),
        =>,
        "(ns foo (:gen-class))(set! *warn-on-reflection* true)

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
        `(ns foo (:gen-class))(set! *warn-on-reflection* true)

12345
`)

test.fact("package can be dotted",
        fgo.funcgoParse("package foo/bar;import ()12345"),
        => ,
        `(ns foo.bar (:gen-class))(set! *warn-on-reflection* true)

12345
`)

test.fact("can import other packages",
        fgo.funcgoParse(`
package foo
import(
  b "bar"
)
12345
`),
        =>,
        `(ns foo (:gen-class)
  (:require [bar :as b]))(set! *warn-on-reflection* true)

12345
`)

func parse(expr) {
  fgo.funcgoParse("package foo;import ()" str expr)
}

func parsed(expr) {
        str("(ns foo (:gen-class))(set! *warn-on-reflection* true)\n\n", expr, "\n")
}

test.fact("can refer to symbols",
        parse("a"), =>, parsed("a"),
        parse("foo"), =>, parsed("foo")
)

test.fact("can refer to numbers",
	parse("99"),  =>, parsed("99"),
	parse("9"),   =>, parsed("9"),
	parse("0"),   =>, parsed("0")
)

test.fact("outside symbols",
	parse("other.foo"), =>, parsed("other/foo")
)

test.fact("can define things",
	parse("a := 12345"), =>, parsed("(def a 12345)"),
	parse("var a = 12345"), =>, parsed("(def a 12345)"),
	parse("var a Foo = 12345"), =>, parsed("(def ^Foo a 12345)")
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

// See http://blog.jayfields.com/2010/07/clojure-destructuring.html
test.fact("Vector Destructuring",
        parse(`const([a,b]=ab) f(a,b)`),
        =>,
        parsed(`(let [[a b] ab] (f a b))`),

        parse(`const([x, more...] = indexes) f(x, more)`), 
        =>,
        parsed(`(let [[x & more] indexes] (f x more))`),

        parse(`const([x, more..., AS, full] = indexes) f(x, more, full)`),
        =>,
        parsed(`(let [[x & more :as full] indexes] (f x more full))`),

        // TODO(eob) implement KEYS:
        //parse(`const({KEYS: [x, y]} = point) f(x, y)`),
        //=>,
        //parsed(`(let [{:keys [x y]} point] (f x y))`),

        parse(`const([[a,b],[c,d]] = numbers) f(a, b, c, d)`),
        =>,
        parsed(`(let [[[a b] [c d]] numbers] (f a b c d))`)
)

test.fact("Map Destructuring",
        parse(`const({theX: X, theY: Y} = point) f(theX, theY)`),
        =>,
        parsed(`(let [{the-x :x the-y :y} point] (f the-x the-y))`),

        parse(`const({name: NAME, {KEYS: [pages, \isbn10]}: DETAILS} = book) f(name,pages,\isbn10)`),
        =>,
        parsed(`(let [{name :name {:keys [pages isbn10]} :details} book] (f name pages isbn10))`),
        
        parse(`const({name: NAME, [hole1, hole2]: SCORES} = golfer) f(name, hole1, hole2)`),
        =>,
        parsed(`(let [{name :name [hole1 hole2] :scores} golfer] (f name hole1 hole2))`),

        parse(`func printStatus({name: NAME, [hole1, hole2]: SCORES}) { f(name, hole1, hole2) }`),
        =>,
        parsed(`(defn print-status [{name :name [hole1 hole2] :scores}] (f name hole1 hole2))`),

        parse(`printStatus( {NAME: "Jim", SCORES: [3, 5, 4, 5]} )`),
        =>,
        parsed(`(print-status {:name "Jim" :scores [3 5 4 5] })`)
)

test.fact("type hints",
        parse(`const(a Foo = 3) f(a)`), =>, parsed(`(let [^Foo a 3] (f a))`),
        parse(`func g(a Foo) { f(a) }`),  =>, parsed(`(defn g [^Foo a] (f a))`),
        parse(`func(a Foo) { f(a) }`),  =>, parsed(`(fn [^Foo a] (f a))`)
)

//      parse(``), =>, parsed(``),
