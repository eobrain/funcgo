package compiler_test
import (
        test "midje/sweet"
        fgo "funcgo/core"
        fgoc "funcgo/main"
        "clojure/string"
)

func compileString(path, fgoText) {
	string.trim(
		string.replace(
			fgoc.CompileString(path, fgoText),
			/\s+/,
			" "
		)
	)
}


test.fact("smallest complete program has no import and a single expression",
        compileString("foo.go", "package foo;12345"),
        =>,
        `(ns foo (:gen-class)) (set! *warn-on-reflection* true) 12345`)

test.fact("Can use newlines instead of semicolons",
        compileString("foo.go", `
package foo
12345
`),
        =>,
        `(ns foo (:gen-class)) (set! *warn-on-reflection* true) 12345`)

test.fact("package can be dotted",
        compileString("foo/bar.go", "package bar;12345"),
        => ,
        `(ns foo.bar (:gen-class)) (set! *warn-on-reflection* true) 12345`)

test.fact("package can be dotted",
        compileString("yippee/yaday/yahoo/boys.go", "package boys;12345"),
        => ,
        `(ns yippee.yaday.yahoo.boys (:gen-class)) (set! *warn-on-reflection* true) 12345`)

test.fact("can import other packages",
        compileString("foo.go", `
package foo
import(
  b "bar"
)
b.xxx
`),
        =>,
        `(ns foo (:gen-class) (:require [bar :as b])) (set! *warn-on-reflection* true) b/xxx`
)


func parse(expr) {
	parse(expr, [], [])
} (expr, pkgs) {
	parse(expr, list(pkgs), [])
} (expr, pkgs, types) {
	const(
		imports = if count(pkgs) == 0 {
			""
		} else {
			const lines = for p := lazy pkgs {str(`"`, p, `"`)}
			str("import(\n", "\n" string.join lines, "\n)\n")
		}
		importtypes = if count(types) == 0 {
			""
		} else {
			str("import type(\n", "\n" string.join types, "\n)\n")
		}
	)
	compileString("foo.go", 
		str("package foo\n", imports, importtypes, expr)
	)
}

func parseJs(expr) {
	parseJs(expr, [], [])
} (expr, pkgs) {
	parseJs(expr, list(pkgs), [])
} (expr, pkgs, types) {
	const(
		imports = if count(pkgs) == 0 {
			""
		} else {
			const lines = for p := lazy pkgs {str(`"`, p, `"`)}
			str("import(\n", "\n" string.join lines, "\n)\n")
		}
		importtypes = if count(types) == 0 {
			""
		} else {
			str("import type(\n", "\n" string.join types, "\n)\n")
		}
	)
	compileString("foo.gos", 
		str("package foo\n", imports, importtypes, expr)
	)
}
//func parseJs(expr, pkgs...) {
//	compileString("foo.gos", 
//		if count(pkgs) == 0 {
//			str("package foo;", expr)
//		}else {
//			const imports = func{str(`"`, .., `"`)} map pkgs
//			str("package foo; import(", "\n" string.join imports, ");", expr)
//		}
//	)
//}

func parsed(expr) {
	parsed(expr, [], [])
} (expr, pkgs) {
	parsed(expr, list(pkgs), [])
} (expr, pkgs, types) {
	const(
		imports = if count(pkgs) == 0 {
			""
		} else {
			const lines = for p := lazy pkgs {str("[", p, " :as ", p, "]")}
			str(" (:require ", " " string.join lines, ")")
		}
		importtypes = if count(types) == 0 {
			""
		} else {
			const lines = for t := lazy types {str("(", t, ")")}
			str(" (:import ", " " string.join lines, ")")
		}
	)
	str("(ns foo (:gen-class)",
		imports,
		importtypes,
		") (set! *warn-on-reflection* true) ",
		expr
	)
}

func parsedJs(expr) {
	parsedJs(expr, [], [])
} (expr, pkgs) {
	parsedJs(expr, list(pkgs), [])
} (expr, pkgs, types) {
	const(
		imports = if count(pkgs) == 0 {
			""
		} else {
			const lines = for p := lazy pkgs {str("[", p, " :as ", p, "]")}
			str(" (:require ", " " string.join lines, ")")
		}
		importtypes = if count(types) == 0 {
			""
		} else {
			const lines = for t := lazy types {str("(", t, ")")}
			str(" (:import ", " " string.join lines, ")")
		}
	)
	str("(ns foo",
		imports,
		importtypes,
		") ",
		expr
	)
}

//func parsedJs(expr, pkgs...) {
//	if count(pkgs) == 0 {
//		str("(ns foo) ", expr)
//	}else{
//		const imports = func{str(" [", .., " :as ", .., "]")} map pkgs
//		str("(ns foo (:require", (str apply imports), ")) ", expr)
//	}
//}

func parseNoPretty(expr) {
	fgo.Parse("foo.go", "package foo;" str expr)
}

func parsedNoPretty(expr) {
        str("(ns foo (:gen-class) ) (set! *warn-on-reflection* true) ", expr)
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
	parse("other.foo", "other"), =>, parsed("other/foo", "other")
)

test.fact("can define things",
	parse("a := 12345"), =>, parsed("(def ^:private a 12345)"),
	parse("var a = 12345"), =>, parsed("(def ^:private a 12345)"),
	parse("var a FooType = 12345", [], ["foo.FooType"]),
		=>, parsed("(def ^{:private true, :tag FooType} a 12345)", [], ["foo FooType"]),
	parse("Foo := 12345"), =>, parsed("(def Foo 12345)"),
	parse("var Foo = 12345"), =>, parsed("(def Foo 12345)"),
	parse("var Foo FooType = 12345", [], ["foo.FooType"]),
	=>, parsed("(def ^FooType Foo 12345)", [], ["foo FooType"])
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
        parse(`func(){b;c}`),         =>, parsed(`(fn [] (do b c))`),
        parse(`func(){b;if c {d}}`), =>, parsed(`(fn [] (do b (when c d)))`)
)

test.fact("subsequent const nests",
        parse(`{const(a=1)x;{const(b=2)y}}`), =>, parsed(`(let [a 1] x (let [b 2] y))`),
        parse(`{const a=1; {const b=2; y}}`), =>, parsed(`(let [a 1] (let [b 2] y))`)
)

// See http://blog.jayfields.com/2010/07/clojure-destructuring.html
test.fact("Vector Destructuring",
        parse(`{const([a,b]=ab) f(a,b)}`),
        =>,
        parsed(`(let [[a b] ab] (f a b))`),

        parse(`{const [a,b]=ab; f(a,b)}`),
        =>,
        parsed(`(let [[a b] ab] (f a b))`),

        parse(`{const([x, more...] = indexes) f(x, more)}`), 
        =>,
        parsed(`(let [[x & more] indexes] (f x more))`),

        parse(`{const [x, more...] = indexes; f(x, more)}`), 
        =>,
        parsed(`(let [[x & more] indexes] (f x more))`),

        parse(`{const([x, more..., AS, full] = indexes) f(x, more, full)}`),
        =>,
        parsed(`(let [[x & more :as full] indexes] (f x more full))`),

        // TODO(eob) implement KEYS:
        //parse(`const({KEYS: [x, y]} = point) f(x, y)`),
        //=>,
        //parsed(`(let [{:keys [x y]} point] (f x y))`),

        parse(`{const([[a,b],[c,d]] = numbers) f(a, b, c, d)}`),
        =>,
        parsed(`(let [[[a b] [c d]] numbers] (f a b c d))`)
)

test.fact("Map Destructuring",
        parse(`{const({theX: X, theY: Y} = point) f(theX, theY)}`),
        =>,
        parsed(`(let [{the-x :x, the-y :y} point] (f the-x the-y))`),

        parse(`{const {theX: X, theY: Y} = point; f(theX, theY)}`),
        =>,
        parsed(`(let [{the-x :x, the-y :y} point] (f the-x the-y))`),

        parse(`{const({name: NAME, {[pages, \isbn10]: KEYS}: DETAILS} = book) f(name,pages,\isbn10)}`),
        =>,
        parsed(`(let [{name :name, {[pages isbn10] :keys} :details} book] (f name pages isbn10))`),
        
        parse(`{const({name: NAME, [hole1, hole2]: SCORES} = golfer) f(name, hole1, hole2)}`),
        =>,
        parsed(`(let [{name :name, [hole1 hole2] :scores} golfer] (f name hole1 hole2))`),

        parse(`func printStatus({name: NAME, [hole1, hole2]: SCORES}) { f(name, hole1, hole2) }`),
        =>,
        parsed(`(defn- print-status [{name :name, [hole1 hole2] :scores}] (f name hole1 hole2))`),

        parse(`func PrintStatus({name: NAME, [hole1, hole2]: SCORES}) { f(name, hole1, hole2) }`),
        =>,
        parsed(`(defn Print-status [{name :name, [hole1 hole2] :scores}] (f name hole1 hole2))`),

        parse(`printStatus( {NAME: "Jim", SCORES: [3, 5, 4, 5]} )`),
        =>,
        parsed(`(print-status {:name "Jim", :scores [3 5 4 5]})`)
)

test.fact("type hints",
        parse(`{const(a FooType = 3) f(a)}`, [], ["foo.FooType"]),
	=>, parsed(`(let [^FooType a 3] (f a))`, [], ["foo FooType"]),

        parse(`{const a FooType = 3; f(a)}`, [], ["foo.FooType"]),
	=>, parsed(`(let [^FooType a 3] (f a))`, [], ["foo FooType"]),

        parse(`func g(a FooType) { f(a) }`, [], ["foo.FooType"]),
	=>, parsed(`(defn- g [^FooType a] (f a))`, [], ["foo FooType"]),

        parse(`func(a FooType) { f(a) }`, [], ["foo.FooType"]),
	=>, parsed(`(fn [^FooType a] (f a))`, [], ["foo FooType"]),

        parse(`func g(a) FooType { f(a) }`, [], ["foo.FooType"]),
	=>, parsed(`(defn- g ^FooType [a] (f a))`, [], ["foo FooType"]),

        parse(`func(a) FooType { f(a) }`, [], ["foo.FooType"]),
	=>, parsed(`(fn ^FooType [a] (f a))`, [], ["foo FooType"]),

        parse(`func f(a) long {a/3} (a, b) double {a+b}`),
	=>,
	parsed(`(defn- f (^long [a] (/ a 3)) (^double [a b] (+ a b)))`),

        parse(`func(a) long {a/3} (a, b) double {a+b}`),
	=>,
	parsed(`(fn (^long [a] (/ a 3)) (^double [a b] (+ a b)))`)
)

test.fact("expression",
	parse(`1<<64 - 1`),         =>, parsed(`(- (bit-shift-left 1 64) 1)`),
	parse(`var a = 1<<64 - 1`), =>, parsed(`(def ^:private a (- (bit-shift-left 1 64) 1))`)
)

test.fact("quoting",
	parse("quote(foo(a))"),           =>, parsed("'(foo a)"),
	parseNoPretty("syntax foo(a)"),     =>, parsedNoPretty("`(foo a)"),
	parseNoPretty("syntax \\`(foo a)`"), =>, parsedNoPretty("`(foo a)"),

	parseNoPretty(`syntax fred(x, unquote x, lst, unquotes lst, 7, 8, NINE)`),
	=>,
	parsedNoPretty("`(fred x ~x lst ~@lst 7 8 :nine)")
)

test.fact("symbol beginning with underscore",
	parse(`_main`), =>, parsed(`-main`),
	parse(`_foo`),  =>, parsed(`-foo`),
	parseJs(`mutateSet( js.window->_onload, start)`), 
	=>, parsedJs(`(set! (. js/window -onload) start)`)
)

test.fact("Javascript",
	parseJs(`new js.Date()->toISOString`), 
	=>, parsedJs(`(. (js.Date.) toISOString)`)
)


test.fact("can call function",
	parse("f()")            ,=>, parsed("(f)"),
	parse("f(x)")           ,=>, parsed("(f x)"),
	//parse("x->f()")         ,=>, parsed("(f x)"),
	//parse("x->f(y,z)")      ,=>, parsed("(f x y z)"),
	parse("f(x,y,z)")       ,=>, parsed("(f x y z)"))

test.fact("can call outside functions",
	parse("o.f(x)", "o")         ,=>, parsed("(o/f x)", "o")
)
test.fact("labels have no lower case",
	parse("FOO")            ,=>, parsed(":foo"),
	parse("FOO_BAR")        ,=>, parsed(":foo-bar"),
	parse("OVER18")         ,=>, parsed(":over18")
)
test.fact("dictionary literals",
	parse("{}")             ,=>, parsed("{}"),
	parse("{A:1}")          ,=>, parsed("{:a 1}"),
	parse("{A:1, B:2}")     ,=>, parsed("{:a 1, :b 2}"),
	parse("{A:1, B:2, C:3}"),=>, parsed("{:a 1, :b 2, :c 3}")
)
test.fact("dictionary literals with trailing comma",
	parse("{}")             ,=>, parsed("{}"),
	parse("{A:1,}")          ,=>, parsed("{:a 1}"),
	parse("{A:1, B:2,}")     ,=>, parsed("{:a 1, :b 2}"),
	parse("{A:1, B:2, C:3,}"),=>, parsed("{:a 1, :b 2, :c 3}")
)
test.fact("set literals",
	parse("set{}")                 ,=>, parsed("#{}"),
	parse("set{A}")                ,=>, parsed("#{:a}"),
	parse("set{A, B}")             ,=>, parsed("#{:a :b}"),
	parse("set{A, B, C}")          ,=>, parsed("#{:a :c :b}"),
	parse(`set{"A", "B", "C"}`)    ,=>, parsed(`#{"A" "B" "C"}`),
	parse(`set{'A', 'B', 'C'}`)    ,=>, parsed(`#{\A \B \C}`),
	parse(`set{A, "B", 'C', 999}`) ,=>, parsed(`#{"B" \C 999 :a}`)
)
test.fact("private named functions",
	parse("func foo(){d}")     ,=>, parsed("(defn- foo [] d)"),
	parse("func foo(a){d}")    ,=>, parsed("(defn- foo [a] d)"),
	parse("func foo(a,b){d}")  ,=>, parsed("(defn- foo [a b] d)"),
	parse("func foo(a,b,c){d}"),=>, parsed("(defn- foo [a b c] d)")
)
test.fact("named functions",
	parse("func Foo(){d}")     ,=>, parsed("(defn Foo [] d)"),
	parse("func Foo(a){d}")    ,=>, parsed("(defn Foo [a] d)"),
	parse("func Foo(a,b){d}")  ,=>, parsed("(defn Foo [a b] d)"),
	parse("func Foo(a,b,c){d}"),=>, parsed("(defn Foo [a b c] d)")
)
test.fact("named functions space",
      parse("func n(a,b) {c}") ,=>, parsed("(defn- n [a b] c)")
)
test.fact("named multifunctions",
      parse("func n(a){b}(c){d}"),=>,parsed("(defn- n ([a] b) ([c] d))")
)
test.fact("named varadic",
	parse("func n(a...){d}")  ,=>,parsed("(defn- n [& a] d)"),
	parse("func n(a,b...){d}"),=>,parsed("(defn- n [a & b] d)"),
	parse("func n(a,b,c...){d}"),=>,parsed("(defn- n [a b & c] d)")
)
test.fact("anonymous functions",
	parse("func(){c}")      ,=>, parsed("(fn [] c)"),
	parse("func(a,b){c}")   ,=>, parsed("(fn [a b] c)")
)
test.fact("anon multifunctions",
      parse("func(a){b}(c){d}"),=>, parsed("(fn ([a] b) ([c] d))")
)
test.fact("anon varadic",
	parse("func(a...){d}")    ,=>, parsed("(fn [& a] d)"),
	parse("func(a,b...){d}")  ,=>, parsed("(fn [a & b] d)"),
	parse("func(a,b,c...){d}"),=>,parsed("(fn [a b & c] d)")
)
test.fact("can have raw strings",
      parse("`one two`")      ,=>, parsed(`"one two"`)
)
test.fact("can have strings",
      parse(`"one two"`)    ,=>, parsed(`"one two"`)
)
test.fact("characters in raw",
	parse("`\n'\b`")   ,=>, parsed(`"\n'\b"`),
	parse(str("`", `"`, "`"))   ,=>, parsed(`"\""`)
)
test.fact("backslash in raw",
      parse("`foo\\bar`")      ,=>, parsed(`"foo\\bar"`)
)
test.fact("characters in strings",
	parse( `"\n"`)    ,=>,  parsed(`"\n"`)
)
// test.fact("quotes in strings",
// parse("\"foo\"bar\"")   ,=>, parsed("\"foo\"bar\""))  TODO implement
test.fact("multiple expr ",
	parse("1;2;3")          ,=>, parsed("1 2 3"),
	parse("1\n2\n3")        ,=>, parsed("1 2 3")
)
test.fact("const",
	parse("{const(a = 2)a}"),=>, parsed("(let [a 2] a)"),
	parse("{ const(  a = 2 ) a}"),=>, parsed("(let [a 2] a)"),
	parse("{const(\nb = 2\n)\na}"),=>, parsed("(let [b 2] a)"),
	parse("{ const(\n  c = 2\n )\n a}"),=>, parsed("(let [c 2] a)"),
	parse("{const(a = 2)f(a,b)}"),=>, parsed("(let [a 2] (f a b))")
)
test.fact("comment",
	parse("//0 blah blah\naaa0")          ,=>, parsed("aaa0"),
	parse(" //1 blah blah \naaa1")        ,=>, parsed("aaa1"),
	parse(" //2 blah blah \naaa2")        ,=>, parsed("aaa2"),
	parse(" //3 blah blah\naaa3")         ,=>, parsed("aaa3"),
	parse("\n //4 blah blah\n \naaa4")    ,=>, parsed("aaa4"),
	parse("// comment\n     aaa5")        ,=>, parsed("aaa5"),
	parse("// comment\n// another\naaa6") ,=>, parsed("aaa6"),
	parse("// comment\n// another\naaa7") ,=>, parsed("aaa7"),
	parse("\n\n//////\n// This file is part of the Funcgo compiler.\naaa8")  ,=>, parsed("aaa8"),
	parse("///////\naaa9")                ,=>, parsed("aaa9")
)
test.fact("regex",
	parse("/aaa/")          ,=>, parsed(`#"aaa"`),
	parse("/[0-9]+/")       ,=>, parsed(`#"[0-9]+"`)
)
//   parse("/aaa\/bbb/"       ,=>, parsed("#\"aaa/bbb"")) TODO implement
test.fact("if",
	parse("if a {b}") ,=>, parsed("(when a b)"),
	parse("if a {b;c}") ,=>, parsed("(when a (do b c))"),
	parse("if a {b\nc}") ,=>, parsed("(when a (do b c))"),
	parse("if a {b}else{c}") ,=>, parsed("(if a b c)"),
	parse("if a {  b  }else{ c  }") ,=>, parsed("(if a b c)"),
	parse("if a {b;c} else {d;e}") ,=>, parsed("(if a (do b c) (do d e))")
)
test.fact("new",
	parse("new Foo()", [], ["foo.Foo"]) ,=>, parsed("(Foo.)", [], ["foo Foo"]),
	parse("new Foo(a)", [], ["foo.Foo"]) ,=>, parsed("(Foo. a)", [], ["foo Foo"]),
	parse("new Foo(a,b,c)", [], ["foo.Foo"]) ,=>, parsed("(Foo. a b c)", [], ["foo Foo"])
)
test.fact("try catch",
	parse("try{a}catch T e{b}", [], ["a.T"]),
	=>, parsed("(try a (catch T e b))", [], ["a T"]),

	parse("try{a}catch T1 e1{b} catch T2 e2{c}", [], ["a.{T1,T2}"]),
	=>, parsed("(try a (catch T1 e1 b) (catch T2 e2 c))", [], ["a T1 T2"]),

	parse("try{a;b}catch T e{c;d}", [], ["a.T"]),
	=>, parsed("(try a b (catch T e c d))", [], ["a T"]),

	parse("try{a}catch T e{b}finally{c}", [], ["a.T"]),
	=>, parsed("(try a (catch T e b) (finally c))", [], ["a T"]),

	parse("try { a } catch T e{ b }", [], ["a.T"]),
	=>, parsed("(try a (catch T e b))", [], ["a T"])
)
test.fact("for",
	parse("for x:=range xs{f(x)}")    ,=>, parsed("(doseq [x xs] (f x))"),
	parse("for x := range xs {f(x)}") ,=>, parsed("(doseq [x xs] (f x))"),
	parse("for x:= lazy xs{f(x)}") ,=>, parsed("(for [x xs] (f x))"),
	parse("for x:= lazy xs if a{f(x)}") ,=>, parsed("(for [x xs :when a] (f x))"),
	parse("for i:= times n {f(i)}") ,=>, parsed("(dotimes [i n] (f i))"),
	parse("for [a,b]:= lazy xs{f(a,b)}") ,=>, parsed("(for [[a b] xs] (f a b))"),
	parse("for x:=lazy xs if x<0 {f(x)}") ,=>, parsed("(for [x xs :when (< x 0)] (f x))")
)
test.fact("Camelcase is converted to dash-separated",
	parse("foo") ,=>, parsed("foo"),
	parse("fooBar") ,=>, parsed("foo-bar"),
	parse("fooBarBaz") ,=>, parsed("foo-bar-baz"),
	parse("foo_bar") ,=>, parsed("foo_bar"),
	parse("Foo") ,=>, parsed("Foo"),
	parse("FooBar") ,=>, parsed("Foo-bar"),
	parse("FOO") ,=>, parsed(":foo"),
	parse("FOO_BAR") ,=>, parsed(":foo-bar"),
	parse("A") ,=>, parsed(":a")
)
test.fact("leading underscore to dash",
	parse("_main") ,=>, parsed("-main")
)
test.fact("is to questionmark",
	parse("isFoo") ,=>, parsed("foo?")
)
test.fact("mutate to exclamation mark",
	parse("mutateFoo") ,=>, parsed("foo!")
)
test.fact("java method calls",
	parse("foo->bar")                     ,=>, parsed("(. foo bar)"),
	parse("foo->bar(a,b)")                ,=>, parsed( "(. foo (bar a b))"),
	parse("foo->bar()")                   ,=>, parsed( "(. foo (bar))"),
	parse(`"fred"->toUpperCase()`)      ,=>, parsed(`(. "fred" (toUpperCase))`),
	parse("println(a, e->getMessage())") ,=>, parsed("(println a (. e (getMessage)))"),
	parse(`System::getProperty("foo")`)  ,=>, parsed(`(System/getProperty "foo")`),
	parse("Math::PI")                      ,=>, parsed("Math/PI"),
	parse("999 * f->foo()")                ,=>, parsed( "(* 999 (. f (foo)))"),
	parse("f->foo() / b->bar()")           ,=>, parsed( "(/ (. f (foo)) (. b (bar)))"),
	parse("999 * f->foo() / b->bar()")     ,=>, parsed( "(/ (* 999 (. f (foo))) (. b (bar)))"),
	parse("999 * f->foo")                  ,=>, parsed( "(* 999 (. f foo))"),
	parse("f->foo / b->bar")             ,=>, parsed( "(/ (. f foo) (. b bar))"),
	parse("999 * f->foo / b->bar")         ,=>, parsed( "(/ (* 999 (. f foo)) (. b bar))")
)
test.fact("there are some non-alphanumeric symbols",
	parse("foo(a,=>,b)") ,=>, parsed("(foo a => b)"),
	parse(`test.fact("interesting", parse("a"), =>, parsed("a"))`, "test"),
	=>,
	parsed(`(test/fact "interesting" (parse "a") => (parsed "a"))`, "test")
)
test.fact("infix",
	parse("a b c")        ,=>, parsed("(b a c)"),
	parse("22 / 7")       ,=>, parsed("(/ 22 7)"),
	parse("22 / (7 + 4)") ,=>, parsed("(/ 22 (+ 7 4))")
)

test.fact("equality",
	parse("a == b") ,=>, parsed("(= a b)"),
	parse("a isIdentical b") ,=>, parsed("(identical? a b)")
)

test.fact("character literals",
	parse("'a'") ,=>, parsed("\\a"),
	parse("['a', 'b', 'c']") ,=>, parsed("[\\a \\b \\c]"),
	parse("'\\n'") ,=>, parsed("\\newline"),
	parse("' '") ,=>, parsed("\\space"),
	parse("'\\t'") ,=>, parsed("\\tab"),
	parse("'\\b'") ,=>, parsed("\\backspace"),
	parse("'\\r'") ,=>, parsed("\\return"),
	parseNoPretty("'\\uDEAD'") ,=>, parsedNoPretty("\\uDEAD"),
	parseNoPretty("'\\ubeef'") ,=>, parsedNoPretty("\\ubeef"),
	parseNoPretty("'\\u1234'") ,=>, parsedNoPretty("\\u1234"),
	parseNoPretty("'\\234'") ,=>, parsedNoPretty("\\o234")
)

test.fact("indexing",
	parse("aaa(BBB)") ,=>, parsed("(aaa :bbb)"),
	parse("aaa[bbb]") ,=>, parsed("(nth aaa bbb)"),
	parse("v(6)") ,=>, parsed("(v 6)"),
	parse("v[6]") ,=>, parsed("(nth v 6)"),
	parse("aaa[6]") ,=>, parsed("(nth aaa 6)")
)

test.fact("precedent",
	parse("a || b < c")   ,=>, parsed("(or a (< b c))"),
	parse("a || b && c") ,=>, parsed("(or a (and b c))"),
	parse("a && b || c") ,=>, parsed("(or (and a b) c)"),
	parse("a * b - c")    ,=>, parsed("(- (* a b) c)"),
	parse("c + a * b")    ,=>, parsed("(+ c (* a b))"),
	parse("a / b + c")    ,=>, parsed("(+ (/ a b) c)")
)

test.fact("associativity",
	parse("x / y * z") ,=>, parsed("(* (/ x y) z)"),
	parse("x * y / z") ,=>, parsed("(/ (* x y) z)"),
	parse("x + y - z") ,=>, parsed("(- (+ x y) z)"),
	parse("x - y + z") ,=>, parsed("(+ (- x y) z)")
)

test.fact("parentheses",
	parse("(a or b) and c") ,=>, parsed("(and (or a b) c)"),
	parse("a * b - c") ,=>, parsed("(- (* a b) c)")
)

test.fact("unary",
	parse("+a")  ,=>, parsed("(+ a)"),
	parse("-a")  ,=>, parsed("(- a)"),
	parse("!a")  ,=>, parsed("(not a)"),
	parse("^a")  ,=>, parsed("(bit-not a)"),
	parse("*a") ,=>, parsed("@a")
)

test.fact("float literals",
	parse("2.0") ,=>, parsed("2.0"),
	parse("2.000") ,=>, parsed("2.0"),
	parse("0.0") ,=>, parsed("0.0"),
	parse("0.") ,=>, parsed("0.0"),
	parse("72.40") ,=>, parsed("72.4"),
	parse("072.40") ,=>, parsed("72.4"),
	parse("2.71828") ,=>, parsed("2.71828"),
	parse("1.0") ,=>, parsed("1.0"),
	parse("1.e+0") ,=>, parsed("1.0"),
	parse("6.67428E-11") ,=>, parsed("6.67428E-11"),
	parse("6.67428e-11") ,=>, parsed("6.67428E-11"),
	parse("1000000.0") ,=>, parsed("1000000.0"),
	parse("1E6") ,=>, parsed("1000000.0"),
	parse(".25") ,=>, parsed(".25"),
	parse(".12345E+5") ,=>, parsed(".12345E+5")
)

test.fact("symbols can start with keywords",
	parse("format") ,=>, parsed("format"),
	parse("ranged") ,=>, parsed("ranged")
)

test.fact("full source file", fgo.Parse("foo.go", `
package foo
import(
  b "bar/baz"
  ff "foo/faz/fedudle"
)

x := b.bbb("blah blah")
func FooBar(iii, jjj) {
  ff.fumanchu(
    {
      OOO: func(m,n) {str(m,n)},
      PPP: func(m,n) {
        str(m,n)
      },
      QQQ: qq
    }
  )
}
`)  ,=>, `(ns foo (:gen-class) (:require [bar.baz :as b] [foo.faz.fedudle :as ff])) (set! *warn-on-reflection* true) (def ^:private x (b/bbb "blah blah")) (defn Foo-bar [iii jjj] (ff/fumanchu {:ooo (fn [m n] (str m n)) :ppp (fn [m n] (str m n)) :qqq qq }))`)


test.fact("Escaped string terminater",
      parse(`"aaa\"aaa"`), =>, parsed(`"aaa\"aaa"`)
)

test.fact("Escaped regex terminater",
	parse(`/aaa\/bbb/`)          ,=>, parsed(`#"aaa/bbb"`)
)

test.fact("tail recursion",
	parse(`loop(){a;recur()}`),           =>, parsed(`(loop [] a (recur))`),
	parse(`loop(a=b){c;recur(d)}`),       =>, parsed(`(loop [a b] c (recur d))`),
	parse(`loop(a=b;c=d){e;recur(f,g)}`), =>, parsed(`(loop [a b c d] e (recur f g))`)
)

test.fact("short anonymous functions",
	parseNoPretty(`func{s}`),             =>, parsedNoPretty(`(fn [] s)`),
	parseNoPretty(`func{a+1}`),           =>, parsedNoPretty(`#(+ a 1)`),
	parseNoPretty(`func{..+..}`),         =>, parsedNoPretty(`#(+ % %)`),
	parseNoPretty(`func{..1+..2}`),       =>, parsedNoPretty(`#(+ %1 %2)`),
	parseNoPretty(`func{.. + ..}`),       =>, parsedNoPretty(`#(+ % %)`),
	parseNoPretty(`func{..1 + ..2}`),     =>, parsedNoPretty(`#(+ %1 %2)`),
	parseNoPretty(`func{str apply ...}`), =>, parsedNoPretty(`#(apply str %&)`)

)

test.fact("Effective Go",
	parse(`if a := b; f(a) {c}`),
	=>,
	parsed("(let [a b] (when (f a) c))"),

	parse(`if err := file.Chmod(0664); err != nil {
    log.Print(err)
    err
}`, ["file", "log"], []),
	=>,
	parsed("(let [err (file/Chmod 436)] (when (not= err nil) (do (log/Print err) err)))",
		["file", "log"],[])
)

test.fact("interface",
	parse(`type Ia interface{
		f(a, b)
		g()
	}`),
	=>, parsed(`(defprotocol Ia (f [this a b]) (g [this]))`)
)

test.fact("An interface defining a sliceable object",
	parse(`type ISliceable interface{
		slice(s int, e int)
		sliceCount() int
	}`),
	=>, parsed(`(defprotocol ISliceable (slice [this ^int s ^int e]) (^long sliceCount [this]))`)
)

test.fact("interface with three methods",
      parse(`

type Interface interface {
        // Len is the number of elements in the collection.
        Len() int
        // Less reports whether the element with
        // index i should sort before the element with index j.
        Less(i, j int) boolean
        // Swap swaps the elements with indexes i and j.
        Swap(i, j int)
}
`), =>, parsed(str(
	`(defprotocol Interface`,
	` (^long Len [this])`,
	` (^boolean Less [this i ^int j])`,
	` (Swap [this i ^int j]))`))
)

//test.fact("",
//      parse(`
//type Sequence []int
//`), =>, parsed(`(defrecord Sequence [??]`)
//)

test.fact("implements",

	parse(`implements Ia func (Ty) f(a) {b}`, [], ["a.Ia"]),
	=>, parsed(`(extend-type Ty Ia (f [this a] b))`, [], ["a Ia"]),
		
	parse(`implements Ia func(Ty)(f(a) {b}; g() {c})`, [], ["a.Ia"]),
	=>, parsed(`(extend-type Ty Ia (f [this a] b) (g [this] c))`, [], ["a Ia"])
)

test.fact("Methods required by sort.Interface",
      parse(`
implements Interface 
func (Sequence) (
  Len() int {
      len(this)
  }
  Less(i, j int) boolean {
      this[i] < this[j]
  }
  Swap(i, j int) {
      this += {i: this[j], j: this[i]}
  }
)
`, [], ["sort.Interface"]),
	=>, parsed(str(
		`(extend-type Sequence Interface`,
		` (^long Len [this] (len this))`,
		` (^boolean Less [this i ^long j] (< (nth this i) (nth this j)))`,
		` (Swap [this i ^long j] (assoc this i (nth this j) j (nth this i))))`),
		[], ["sort Interface"])
)

test.fact("Method for printing - sorts the elements before printing.",
	parse(`
implements Stringer
func (Sequence) String() String {
    str("[",  " " join sort.Sort(this),  "]")
}

`, ["sort"], ["fmt.Stringer"]), =>, parsed(str(
	`(extend-type Sequence Stringer`,
	` (^String String [this]`,
	` (str "[" (join " " (sort/Sort this)) "]")))`
), ["sort"], ["fmt Stringer"]))

test.fact("struct",
	parse(`type TreeNode struct{val; l; r}`), 
	=>,
	parsed(str(`(defrecord TreeNode [val l r]`,
		` Object (toString [this] (str "{" val " " l " " r "}")))`))
)

test.fact("typed struct",
	parse(`type TreeNode struct{}; type TreeNode struct{val; l TreeNode; r TreeNode}`), 
	=>,
	parsed(str(`(defrecord TreeNode []) (defrecord TreeNode [val ^TreeNode l ^TreeNode r]`,
		` Object (toString [this] (str "{" val " " l " " r "}")))`))
)

test.fact("switch",
	parse(`switch {case a: b; case c: d; default: e}`),
	=>, parsed(`(cond a b c d :else e)`),

	parse(`switch x.(type) {case String: x; case Integer: str(x*x); default: str(x)}`),
	=>, parsed(`(cond (instance? String x) x (instance? Integer x) (str (* x x)) :else (str x))`),

	parse(`switch x {case A: b; case C: d; default: e}`),
	=>, parsed(`(case x :a b :c d e)`),

	parse(`switch x {case P, Q, R: b; case S, T, U: d; default: e}`),
	=>, parsed(`(case x (:p :q :r) b (:s :t :u) d e)`)
)
test.fact("Error if external package not imported",
	parse("huh.bar"),
	=>, test.throws(Exception, `package "huh" in huh.bar does not appear in imports []`),

	parse("huh.bar", ["aaa", "bbb"], []),
	=>, test.throws(Exception, `package "huh" in huh.bar does not appear in imports [bbb, aaa]`)
)

test.fact("Error if import not used",
	parse("1234", "aaa"),
	=>, test.throws(Exception, `Packages imported but never used: [aaa]`),

	parse("1234", ["aaa", "bbb"], []),
	=>, test.throws(Exception, `Packages imported but never used: [aaa, bbb]`),

	parse("aaa.xxx", ["aaa", "bbb"], []),
	=>, test.throws(Exception, `Packages imported but never used: [bbb]`),

	parse("1234", [], ["a.Aaa", "b.Bbb"]),
	=>, test.throws(Exception, `Types imported but never used: [Aaa, Bbb]`),

	parse("Aaa::xxx", [], ["a.Aaa", "b.Bbb"]),
	=>, test.throws(Exception, `Types imported but never used: [Bbb]`)
)

test.fact("import type",
	compileString("joy/java.go", `
package java
import type (
  java.util.{HashMap, List}
  java.util.concurrent.atomic.AtomicLong
)

new HashMap({"happy?": true})
new AtomicLong(42)
new List()
`),
	=>,
	`(ns joy.java (:gen-class) (:import (java.util HashMap List) (java.util.concurrent.atomic AtomicLong))) (set! *warn-on-reflection* true) (HashMap. {"happy?" true}) (AtomicLong. 42) (List.)`
)
	
test.fact("assoc",
	parse(`x += {AA: aaa, BB: bbb}`), =>, parsed(`(assoc x :aa aaa :bb bbb)`)
)

test.fact("dissoc",
	parse(`x -= {AA: aaa, BB: bbb}`), =>, parsed(`(dissoc x :aa aaa :bb bbb)`)
)

test.fact("assoc-in",
	parse(`x += {4 AAA 6 8: aaa }`), =>, parsed(`(assoc-in x [4 :aaa 6 8] aaa)`)
)

test.fact("Vertex struct",
	parse(`type Vertex struct {
	Lat, Long float64
}
`), =>, parsed(str(`(defrecord Vertex [^double Lat ^double Long]`,
		` Object (toString [this] (str "{" Lat " " Long "}")))`))
)

test.fact("struct literal",
      parse(`Vertex{
		40.68433, -74.39967
	}`,[],["a.Vertex"]), =>, parsed(`(Vertex. 40.68433 (- 74.39967))`,[],["a Vertex"])
)

test.fact("struct literal with trailing comma",
      parse(`Vertex{
		40.68433, -74.39967,
	}`,[],["a.Vertex"]), =>, parsed(`(Vertex. 40.68433 (- 74.39967))`,[],["a Vertex"])
)

//test.fact("type assertion",
//	parse(`a.(string)`), =>, parsed(`[x (instance? String x)]`)
//)

//test.fact("",
//      parse(``), =>, parsed(``),
//)
