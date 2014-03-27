(ns funcgo.core-test
	(:require [midje.sweet :as test])
	(:require [funcgo.core :as fgo]))

(defn parse [expr]
  (fgo/funcgo-parse (str "package foo;import ()" expr)))

(defn parsed [expr]
  (str "(ns foo (:gen-class))(set! *warn-on-reflection* true)\n\n" expr "\n"))

(test/fact "can call function"
      (parse "f()")            => (parsed "(f)")
      (parse "f(x)")           => (parsed "(f x)")
      ;;(parse "x->f()")         => (parsed "(f x)")
      ;;(parse "x->f(y,z)")      => (parsed "(f x y z)")
      (parse "f(x,y,z)")       => (parsed "(f x y z)"))
(test/fact "can outside functions"
      (parse "o.f(x)")         => (parsed "(o/f x)"))
(test/fact "labels have no lower case"
      (parse "FOO")            => (parsed ":foo")
      (parse "FOO_BAR")        => (parsed ":foo-bar")
      (parse "OVER18")         => (parsed ":over18"))
(test/fact "dictionary literals"
      (parse "{}")             => (parsed "{}")
      (parse "{A:1}")          => (parsed "{:a 1 }")
      (parse "{A:1, B:2}")     => (parsed "{:a 1 :b 2 }")
      (parse "{A:1, B:2, C:3}")=> (parsed "{:a 1 :b 2 :c 3 }"))
(test/fact "named functions"
      (parse "func n(){d}")     => (parsed "(defn n [] d)")
      (parse "func n(a){d}")    => (parsed "(defn n [a] d)")
      (parse "func n(a,b){d}")  => (parsed "(defn n [a b] d)")
      (parse "func n(a,b,c){d}")=> (parsed "(defn n [a b c] d)"))
(test/fact "named functions space"
      (parse "func n(a,b) {c}") => (parsed "(defn n [a b] c)"))
(test/fact "named multifunctions"
      (parse "func n(a){b}(c){d}")=>(parsed "(defn n ([a] b) ([c] d))"))
(test/fact "named varadic"
      (parse "func n(a...){d}")  =>(parsed "(defn n [& a] d)")
      (parse "func n(a,b...){d}")=>(parsed "(defn n [a & b] d)")
      (parse "func n(a,b,c...){d}")=>(parsed "(defn n [a b & c] d)"))
(test/fact "anonymous functions"
      (parse "func(){c}")      => (parsed "(fn [] c)")
      (parse "func(a,b){c}")   => (parsed "(fn [a b] c)"))
(test/fact "anon multifunctions"
      (parse "func(a){b}(c){d}")=> (parsed "(fn ([a] b) ([c] d))"))
(test/fact "anon varadic"
      (parse "func(a...){d}")    => (parsed "(fn [& a] d)")
      (parse "func(a,b...){d}")  => (parsed "(fn [a & b] d)")
      (parse "func(a,b,c...){d}")=>(parsed "(fn [a b & c] d)"))
(test/fact "can have raw strings"
      (parse "`one two`")      => (parsed "\"one two\""))
(test/fact "can have strings"
      (parse "\"one two\"")    => (parsed "\"one two\""))
(test/fact "characters in raw"
      (parse "`\n'\"\b`")      => (parsed "\"\\n'\\\"\\b\""))
(test/fact "backslash in raw"
      (parse "`foo\\bar`")      => (parsed "\"foo\\\\bar\""))
(test/fact "characters in strings"
      (parse  "\"\\n\"")    =>
      (parsed "\"\\n\"")
      (parse  "\"\n'\b\"")  =>
      (parsed "\"\n'\b\""))
;; (test/fact "quotes in strings"
;; (parse "\"foo\"bar\"")   => (parsed "\"foo\"bar\""))  TODO implement
(test/fact "multiple expr "
      (parse "1;2;3")          => (parsed "1 2 3")
      (parse "1\n2\n3")        => (parsed "1 2 3"))
(test/fact "const"
      (parse "const(a = 2)a")=> (parsed "(let [a 2] a)")
      (parse " const(  a = 2 ) a")=> (parsed "(let [a 2] a)")
      (parse "const(\na = 2\n)\na")=> (parsed "(let [a 2] a)")
      (parse " const(\n  a = 2\n )\n a")=> (parsed "(let [a 2] a)")
      (parse "const(a = 2)f(a,b)")=> (parsed "(let [a 2] (f a b))"))
(test/fact "comment"
      (parse "//0 blah blah\naaa0")          => (parsed "aaa0")
      (parse " //1 blah blah \naaa1")        => (parsed "aaa1")
      (parse " //2 blah blah \naaa2")        => (parsed "aaa2")
      (parse " //3 blah blah\naaa3")         => (parsed "aaa3")
      (parse "\n //4 blah blah\n \naaa4")    => (parsed "aaa4")
      (parse "// comment\n     aaa5")        => (parsed "aaa5")
      (parse "// comment\n// another\naaa6") => (parsed "aaa6")
      (parse "// comment\n// another\naaa7") => (parsed "aaa7")
      (parse "\n\n//////\n// This file is part of the Funcgo compiler.\naaa8")  => (parsed "aaa8")
      (parse "///////\naaa9")                => (parsed "aaa9"))
(test/fact "regex"
      (parse "/aaa/")          => (parsed "#\"aaa\"")
      (parse "/[0-9]+/")       => (parsed "#\"[0-9]+\""))
;;   (parse "/aaa\/bbb/"       => (parsed "#\"aaa/bbb"")) TODO implement
(test/fact "if"
      (parse "if a {b}") => (parsed "(when a b)")
      (parse "if a {b;c}") => (parsed "(when a b c)")
      (parse "if a {b\nc}") => (parsed "(when a b c)")
      (parse "if a {b}else{c}") => (parsed "(if a b c)")
      (parse "if a {  b  }else{ c  }") => (parsed "(if a b c)")
      (parse "if a {b;c} else {d;e}") => (parsed "(if a (do b c) (do d e))"))
(test/fact "new"
      (parse "new Foo()") => (parsed "(Foo.)")
      (parse "new Foo(a)") => (parsed "(Foo. a)")
      (parse "new Foo(a,b,c)") => (parsed "(Foo. a b c)"))
(test/fact "try catch"
      (parse "try{a}catch T e{b}") => (parsed "(try a (catch T e b))")
      (parse "try{a}catch T1 e1{b} catch T2 e2{c}")
      => (parsed "(try a (catch T1 e1 b) (catch T2 e2 c))")
      (parse "try{a;b}catch T e{c;d}") => (parsed "(try a b (catch T e c d))")
      (parse "try{a}catch T e{b}finally{c}") => (parsed "(try a (catch T e b) (finally c))")
      (parse "try { a } catch T e{ b } ") => (parsed "(try a (catch T e b))"))
(test/fact "for"
      (parse "for x:=range xs{f(x)}")    => (parsed "(doseq [x xs] (f x))")
      (parse "for x := range xs {f(x)}") => (parsed "(doseq [x xs] (f x))")
      (parse "for x:= lazy xs{f(x)}") => (parsed "(for [x xs] (f x))")
      (parse "for x:= lazy xs if a{f(x)}") => (parsed "(for [x xs] :when a (f x))")
      (parse "for i:= times n {f(i)}") => (parsed "(dotimes [i n] (f i))"))
(test/fact "Camelcase is converted to dash-separated"
      (parse "foo") => (parsed "foo")
      (parse "fooBar") => (parsed "foo-bar")
      (parse "fooBarBaz") => (parsed "foo-bar-baz")
      (parse "foo_bar") => (parsed "foo_bar")
      (parse "Foo") => (parsed "Foo")
      (parse "FooBar") => (parsed "Foo-bar")
      (parse "FOO") => (parsed ":foo")
      (parse "FOO_BAR") => (parsed ":foo-bar")
      (parse "A") => (parsed ":a"))
(test/fact "leading underscore to dash"
      (parse "_main") => (parsed "-main"))
(test/fact "is to questionmark"
      (parse "isFoo") => (parsed "foo?"))
(test/fact "mutate to exclamation mark"
      (parse "mutateFoo") => (parsed "foo!"))
(test/fact "java method calls"
      (parse "foo->bar")                     => (parsed "(. foo bar)")
      (parse "foo->bar(a,b)")                => (parsed  "(. foo (bar a b))")
      (parse "foo->bar()")                   => (parsed  "(. foo (bar))")
      (parse "\"fred\"->toUpperCase()")      => (parsed "(. \"fred\" (toUpperCase))")
      (parse "println(a, e->getMessage())") => (parsed "(println a (. e (getMessage)))")
      (parse "System::getProperty(\"foo\")")  => (parsed "(System/getProperty \"foo\")")
      (parse "Math::PI")                      => (parsed "Math/PI")
      (parse "999 * f->foo()")                => (parsed  "(* 999 (. f (foo)))")
      (parse "f->foo() / b->bar()")           => (parsed  "(/ (. f (foo)) (. b (bar)))")
      (parse "999 * f->foo() / b->bar()")     => (parsed  "(/ (* 999 (. f (foo))) (. b (bar)))")
      (parse "999 * f->foo")                  => (parsed  "(* 999 (. f foo))")
      (parse "f->foo / b->bar")             => (parsed  "(/ (. f foo) (. b bar))")
      (parse "999 * f->foo / b->bar")         => (parsed  "(/ (* 999 (. f foo)) (. b bar))")
      )
(test/fact "there are some non-alphanumeric symbols"
           (parse "foo(a,=>,b)") => (parsed "(foo a => b)")
           (parse "test.fact(\"interesting\", parse(\"a\"), =>, parsed(\"a\"))")
           => (parsed "(test/fact \"interesting\" (parse \"a\") => (parsed \"a\"))")
           (parse "=>(a,b)") => (parsed "(=> a b)")
           )
(test/fact "infix"
           (parse "a b c") => (parsed "(b a c)")
           (parse "22 / 7") => (parsed "(/ 22 7)")
           )

(test/fact "equality"
           (parse "a == b") => (parsed "(= a b)")
           (parse "a isIdentical b") => (parsed "(identical? a b)"))

(test/fact "character literals",
           (parse "'a'") => (parsed "\\a")
           (parse "['a', 'b', 'c']") => (parsed "[\\a \\b \\c]")
           (parse "'\\n'") => (parsed "\\newline")
           (parse "' '") => (parsed "\\space")
           (parse "'\\t'") => (parsed "\\tab")
           (parse "'\\b'") => (parsed "\\backspace")
           (parse "'\\r'") => (parsed "\\return")
           (parse "'\\uDEAD'") => (parsed "\\uDEAD")
           (parse "'\\ubeef'") => (parsed "\\ubeef")
           (parse "'\\u1234'") => (parsed "\\u1234")
           (parse "'\\234'") => (parsed "\\o234"))

(test/fact "indexing"
           (parse "aaa(BBB)") => (parsed "(aaa :bbb)")
           (parse "aaa[bbb]") => (parsed "(nth aaa bbb)")
           (parse "v(6)") => (parsed "(v 6)")
           (parse "v[6]") => (parsed "(nth v 6)")
           (parse "aaa[6]") => (parsed "(nth aaa 6)"))

(test/fact "precedent"
           (parse "a || b < c")   => (parsed "(or a (< b c))")
           (parse "a || b && c") => (parsed "(or a (and b c))")
           (parse "a && b || c") => (parsed "(or (and a b) c)")
           (parse "a * b - c")    => (parsed "(- (* a b) c)")
           (parse "c + a * b")    => (parsed "(+ c (* a b))")
           (parse "a / b + c")    => (parsed "(+ (/ a b) c)")
           )

(test/fact "associativity"
           (parse "x / y * z") => (parsed "(* (/ x y) z)")
           (parse "x * y / z") => (parsed "(/ (* x y) z)")
           (parse "x + y - z") => (parsed "(- (+ x y) z)")
           (parse "x - y + z") => (parsed "(+ (- x y) z)"))

(test/fact "parentheses"
           (parse "(a or b) and c") => (parsed "(and (or a b) c)")
           (parse "a * b - c") => (parsed "(- (* a b) c)"))

(test/fact "unary"
           (parse "+a")  => (parsed "(+ a)")
           (parse "-a")  => (parsed "(- a)")
           (parse "!a")  => (parsed "(not a)")
           (parse "^a")  => (parsed "(bit-not a)")
           (parse "*a") => (parsed "@a"))

(test/fact "float litersls"
           (parse "2.000") => (parsed "2.000")
           (parse "0.") => (parsed "0.")
           (parse "72.40") => (parsed "72.40")
           (parse "072.40") => (parsed "072.40")
           (parse "2.71828") => (parsed "2.71828")
           (parse "1.e+0") => (parsed "1.e+0")
           (parse "6.67428e-11") => (parsed "6.67428e-11")
           (parse "1E6") => (parsed "1E6")
           (parse ".25") => (parsed ".25")
           (parse ".12345E+5") => (parsed ".12345E+5"))

(test/fact "symbols can start with keywords",
           (parse "format") => (parsed "format")
           (parse "ranged") => (parsed "ranged"))

(test/fact "full source file" (fgo/funcgo-parse "
package foo
import(
  b \"bar/baz\"
  ff \"foo/faz/fedudle\"
)

x := b.bbb(`blah blah`)
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
")  => "(ns foo (:gen-class)
  (:require [bar.baz :as b])
  (:require [foo.faz.fedudle :as ff]))(set! *warn-on-reflection* true)

(def x (b/bbb \"blah blah\")) (defn Foo-bar [iii jjj] (ff/fumanchu {:ooo (fn [m n] (str m n)) :ppp (fn [m n] (str m n)) :qqq qq }))
")
