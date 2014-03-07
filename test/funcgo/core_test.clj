(ns funcgo.core-test
  (:use midje.sweet)
  (:require [funcgo.core :refer :all]))

(fact "smallest complete program has no import and a single expression"
      (funcgo-parse "package foo;import ();12345")
      => "(ns foo)

12345
")

(fact "Can use newlines instead of semicolons"  (funcgo-parse "
package foo
import ()
12345
")  =>  "(ns foo)

12345
")

(fact "package can be dotted" (funcgo-parse "
package foo.bar
import ()
12345
")  =>  "(ns foo.bar)

12345
")

(fact "can import other packages" (funcgo-parse "
package foo
import(
  b bar
)
12345
")  => "(ns foo
  (:require [bar :as b]))

12345
")

(defn parse [expr]
  (funcgo-parse (str "package foo;import (;)" expr)))

(defn parsed [expr]
  (str "(ns foo)\n\n" expr "\n"))

(fact "can refer to symbols"  (parse "a")              => (parsed "a"))
(fact "outside symbols"       (parse "other.foo")      => (parsed "other/foo"))
(fact "can define things"     (parse "a := 12345")     => (parsed "(def a 12345)"))
(fact "can call functions"    (parse "f(x)")           => (parsed "(f x)"))
(fact "can outside functions" (parse "o.f(x)")         => (parsed "(o/f x)"))
(fact "labels are all-caps"   (parse "FOO")            => (parsed ":foo"))
(fact "dictionary literals"   (parse "{A:1, B:2}")     => (parsed "{:a 1 :b 2 }"))
(fact "named functions"       (parse "func n(a,b){c}") => (parsed "(defn n [a b]\n  c)"))
(fact "named functions 2"     (parse "func n(a,b) {c}")=> (parsed "(defn n [a b]\n  c)"))
(fact "anonymous functions"   (parse "func(a,b){c}")   => (parsed "(fn [a b]\n  c)"))
(fact "can have raw strings"  (parse "`one two`")      => (parsed "\"one two\""))
(fact "can have strings"      (parse "\"one two\"")    => (parsed "\"one two\""))
(fact "characters in raw"     (parse "`\n'\"\b`")      => (parsed "\"\\n'\\\"\\b\""))
(fact "characters in strings" (parse "\"\n'\b\"")      => (parsed "\"\n'\b\""))
;; (fact "quotes in strings"     (parse "\"foo\"bar\"")   => (parsed "\"foo\"bar\""))  TODO implement
(fact "multiple expr"         (parse "1;2;3")          => (parsed "1\n\n2\n\n3"))
(fact "multiple expr 2"       (parse "1\n2\n3")        => (parsed "1\n\n2\n\n3"))


(fact "full source file" (funcgo-parse "
package foo
import(
  b bar.baz
  ff foo.faz.fedudle
)

x := b.bbb(`blah blah`)

func FooBar(iii, jjj) {
  ff.fumanchu(
    333
  )
}
")  => "(ns foo
  (:require [bar.baz :as b])
  (:require [foo.faz.fedudle :as ff]))

(def x (b/bbb \"blah blah\"))

(defn FooBar [iii jjj]
  (ff/fumanchu 333))
")
