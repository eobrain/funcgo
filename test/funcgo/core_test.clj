(ns funcgo.core-test
  (:use midje.sweet)
  (:require [funcgo.core :refer :all]))

(fact "smallest complete program has no import and a single expression"
      (funcgo-parse "package foo;import ()12345")  => "(ns foo (:gen-class))

12345
")

(fact "Can use newlines instead of semicolons"  (funcgo-parse "
package foo
import (
)
12345
")  =>  "(ns foo (:gen-class))

12345
")

(fact "package can be dotted"
      (funcgo-parse "package foo.bar;import ()12345")  =>  "(ns foo.bar (:gen-class))

12345
")

(fact "can import other packages" (funcgo-parse "
package foo
import(
  b bar
)
12345
")  => "(ns foo (:gen-class)
  (:require [bar :as b]))

12345
")

(defn parse [expr]
  (funcgo-parse (str "package foo;import ()" expr)))

(defn parsed [expr]
  (str "(ns foo (:gen-class))\n\n" expr "\n"))

(fact "can refer to symbols"
      (parse "a")              => (parsed "a"))
(fact "outside symbols"       (parse "other.foo")
      => (parsed "other/foo"))
(fact "can define things"     (parse "a := 12345")
      => (parsed "(def a 12345)"))
(fact "can call function"
      (parse "f()")            => (parsed "(f)")
      (parse "f(x)")           => (parsed "(f x)")
      ;;(parse "x->f()")         => (parsed "(f x)")
      ;;(parse "x->f(y,z)")      => (parsed "(f x y z)")
      (parse "f(x,y,z)")       => (parsed "(f x y z)"))
(fact "can outside functions"
      (parse "o.f(x)")         => (parsed "(o/f x)"))
(fact "labels have no lower case"
      (parse "FOO")            => (parsed ":foo")
      (parse "FOO_BAR")        => (parsed ":foo-bar")
      (parse "OVER18")         => (parsed ":over18"))
(fact "dictionary literals"
      (parse "{}")             => (parsed "{}")
      (parse "{A:1}")          => (parsed "{:a 1 }")
      (parse "{A:1, B:2}")     => (parsed "{:a 1 :b 2 }")
      (parse "{A:1, B:2, C:3}")=> (parsed "{:a 1 :b 2 :c 3 }"))
(fact "named functions"
      (parse "func n(){d}")     => (parsed "(defn n [] d)")
      (parse "func n(a){d}")    => (parsed "(defn n [a] d)")
      (parse "func n(a,b){d}")  => (parsed "(defn n [a b] d)")
      (parse "func n(a,b,c){d}")=> (parsed "(defn n [a b c] d)"))
(fact "named functions space"
      (parse "func n(a,b) {c}") => (parsed "(defn n [a b] c)"))
(fact "named multifunctions"
      (parse "func n(a){b}(c){d}")=>(parsed "(defn n ([a] b) ([c] d))"))
(fact "named varadic"
      (parse "func n(&a){d}")  =>(parsed "(defn n [& a] d)")
      (parse "func n(a,&b){d}")=>(parsed "(defn n [a & b] d)")
      (parse "func n(a,b,&c){d}")=>(parsed "(defn n [a b & c] d)"))
(fact "anonymous functions"
      (parse "func(){c}")      => (parsed "(fn [] c)")
      (parse "func(a,b){c}")   => (parsed "(fn [a b] c)"))
(fact "anon multifunctions"
      (parse "func(a){b}(c){d}")=> (parsed "(fn ([a] b) ([c] d))"))
(fact "anon varadic"
      (parse "func(&a){d}")    => (parsed "(fn [& a] d)")
      (parse "func(a,&b){d}")  => (parsed "(fn [a & b] d)")
      (parse "func(a,b,&c){d}")=>(parsed "(fn [a b & c] d)"))
(fact "can have raw strings"
      (parse "`one two`")      => (parsed "\"one two\""))
(fact "can have strings"
      (parse "\"one two\"")    => (parsed "\"one two\""))
(fact "characters in raw"
      (parse "`\n'\"\b`")      => (parsed "\"\\n'\\\"\\b\""))
(fact "backslash in raw"
      (parse "`foo\\bar`")      => (parsed "\"foo\\\\bar\""))
(fact "characters in strings"
      (parse "\"\n'\b\"")      => (parsed "\"\n'\b\""))
;; (fact "quotes in strings"
;; (parse "\"foo\"bar\"")   => (parsed "\"foo\"bar\""))  TODO implement
(fact "multiple expr "
      (parse "1;2;3")          => (parsed "1 2 3")
      (parse "1\n2\n3")        => (parsed "1 2 3"))
(fact "const"
      (parse "const(a=2)a")=> (parsed "(let [a 2] a)")
      (parse " const(  a=2 ) a")=> (parsed "(let [a 2] a)")
      (parse "const(\na=2\n)\na")=> (parsed "(let [a 2] a)")
      (parse " const(\n  a=2\n )\n a")=> (parsed "(let [a 2] a)")
      (parse "const(a=2)f(a,b)")=> (parsed "(let [a 2] (f a b))"))
(fact "comment"
      (parse "//blah blah\naaa")=> (parsed "aaa")
      (parse " //blah blah \naaa")=> (parsed "aaa")
      (parse "\n //blah blah\n \naaa")=> (parsed "aaa"))
(fact "regex"
      (parse "/aaa/")          => (parsed "#\"aaa\"")
      (parse "/[0-9]+/")       => (parsed "#\"[0-9]+\""))
;;   (parse "/aaa\/bbb/"       => (parsed "#\"aaa/bbb"")) TODO implement
(fact "if"
      (parse "if a {b}") => (parsed "(when a b)")
      (parse "if a {b;c}") => (parsed "(when a b c)")
      (parse "if a {b\nc}") => (parsed "(when a b c)")
      (parse "if a {b}else{c}") => (parsed "(if a b c)")
      (parse "if a {  b  }else{ c  }") => (parsed "(if a b c)")
      (parse "if a {b;c} else {d;e}") => (parsed "(if a (do b c) (do d e))"))
(fact "new"
      (parse "new Foo()") => (parsed "(Foo.)")
      (parse "new Foo(a)") => (parsed "(Foo. a)")
      (parse "new Foo(a,b,c)") => (parsed "(Foo. a b c)"))
(fact "try catch"
      (parse "try{a}catch T e{b}") => (parsed "(try a (catch T e b))")
      (parse "try{a}catch T1 e1{b} catch T2 e2{c}")
      => (parsed "(try a (catch T1 e1 b) (catch T2 e2 c))")
      (parse "try{a;b}catch T e{c;d}") => (parsed "(try a b (catch T e c d))")
      (parse "try{a}catch T e{b}finally{c}") => (parsed "(try a (catch T e b) (finally c))")
      (parse "try { a } catch T e{ b } ") => (parsed "(try a (catch T e b))"))
(fact "for"
      (parse "for x:=range xs{f(x)}")    => (parsed "(doseq [x xs] (f x))")
      (parse "for x := range xs {f(x)}") => (parsed "(doseq [x xs] (f x))")
      (parse "for x:= lazy xs{f(x)}") => (parsed "(for [x xs] (f x))")
      (parse "for x:= lazy xs if a{f(x)}") => (parsed "(for [x xs] :when a (f x))")
      (parse "for i:= times n {f(i)}") => (parsed "(dotimes [i n] (f i))"))
(fact "Camelcase is converted to dash-separated"
      (parse "foo") => (parsed "foo")
      (parse "fooBar") => (parsed "foo-bar")
      (parse "fooBarBaz") => (parsed "foo-bar-baz")
      (parse "foo_bar") => (parsed "foo_bar")
      (parse "Foo") => (parsed "Foo")
      (parse "FooBar") => (parsed "Foo-bar")
      (parse "FOO") => (parsed ":foo")
      (parse "FOO_BAR") => (parsed ":foo-bar")
      (parse "A") => (parsed ":a"))
(fact "leading underscore to dash"
      (parse "_main") => (parsed "-main"))
(fact "is to questionmark"
      (parse "isFoo") => (parsed "foo?"))
(fact "mutate to exclamation mark"
      (parse "mutateFoo") => (parsed "foo!"))
(fact "java method calls"
      (parse "foo->bar")                     => (parsed "(. foo bar)")
      (parse "foo->bar(a,b)")                => (parsed  "(. foo (bar a b))")
      (parse "foo->bar()")                   => (parsed  "(. foo (bar))")
      (parse "\"fred\"->toUpperCase()")      => (parsed "(. \"fred\" (toUpperCase))")
      (parse "println(a, e->getMessage())") => (parsed "(println a (. e (getMessage)))")
      ;;(parse "System::getProperty(\"foo\")")  => (parsed "(System/getProperty \"foo\")")
      ;;(parse "Math::PI")                      => (parsed "Math/PI")
      )



(fact "full source file" (funcgo-parse "
package foo
import(
  b bar.baz
  ff foo.faz.fedudle
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
  (:require [foo.faz.fedudle :as ff]))

(def x (b/bbb \"blah blah\")) (defn Foo-bar [iii jjj] (ff/fumanchu {:ooo (fn [m n] (str m n)) :ppp (fn [m n] (str m n)) :qqq qq }))
")
