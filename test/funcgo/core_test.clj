(ns funcgo.core-test
  (:use midje.sweet)
  (:require [funcgo.core :refer :all]))

(fact "smallest complete program has no import and a single expression"
      (funcgo-parse "package foo;import (;);12345;")
      => "(ns foo)

12345
")

(fact "Can use newlines instead of semicolons"  (funcgo-parse "
package foo
import (
)
12345
")  =>  "(ns foo)

12345
")

(fact "package can be dotted" (funcgo-parse "
package foo.bar
import (
)
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
  (funcgo-parse (str "package foo;import (;);" expr ";")))

(defn parsed [expr]
  (str "(ns foo)\n\n" expr "\n"))

(fact "can refer to symbols"  (parse "a")              => (parsed "a"))
(fact "outside symbols"       (parse "other.foo")      => (parsed "other/foo"))
(fact "can define things"     (parse "a := 12345")     => (parsed "(def a 12345)"))
(fact "can call functions 0"  (parse "f()")            => (parsed "(f)"))
(fact "can call functions 1"  (parse "f(x)")           => (parsed "(f x)"))
(fact "can call functions 3"  (parse "f(x,y,z)")       => (parsed "(f x y z)"))
(fact "can outside functions" (parse "o.f(x)")         => (parsed "(o/f x)"))
(fact "labels are all-caps"   (parse "FOO")            => (parsed ":foo"))
(fact "dictionary literals 0" (parse "{}")             => (parsed "{}"))
(fact "dictionary literals 1" (parse "{A:1}")          => (parsed "{:a 1 }"))
(fact "dictionary literals 2" (parse "{A:1, B:2}")     => (parsed "{:a 1 :b 2 }"))
(fact "dictionary literals 3" (parse "{A:1, B:2, C:3}")=> (parsed "{:a 1 :b 2 :c 3 }"))
(fact "named functions 0"     (parse "func n(){d}")    => (parsed "(defn n [] d)"))
(fact "named functions 1"     (parse "func n(a){d}")   => (parsed "(defn n [a] d)"))
(fact "named functions 2"     (parse "func n(a,b){d}") => (parsed "(defn n [a b] d)"))
(fact "named functions 3"     (parse "func n(a,b,c){d}")=> (parsed "(defn n [a b c] d)"))
(fact "named functions space" (parse "func n(a,b) {c}") => (parsed "(defn n [a b] c)"))
(fact "named multifunctions"  (parse "func n(a){b}(c){d}")=>(parsed "(defn n ([a] b) ([c] d))"))
(fact "named varadic 1"       (parse "func n(&a){d}")  =>(parsed "(defn n [& a] d)"))
(fact "named varadic 2"       (parse "func n(a,&b){d}")=>(parsed "(defn n [a & b] d)"))
(fact "named varadic 3"       (parse "func n(a,b,&c){d}")=>(parsed "(defn n [a b & c] d)"))
(fact "anonymous functions 0" (parse "func(){c}")      => (parsed "(fn [] c)"))
(fact "anonymous functions"   (parse "func(a,b){c}")   => (parsed "(fn [a b] c)"))
(fact "anon multifunctions"   (parse "func(a){b}(c){d}")=> (parsed "(fn ([a] b) ([c] d))"))
(fact "anon varadic 1"        (parse "func(&a){d}")    => (parsed "(fn [& a] d)"))
(fact "anon varadic 2"        (parse "func(a,&b){d}")  => (parsed "(fn [a & b] d)"))
(fact "anon varadic 3"        (parse "func(a,b,&c){d}")=>(parsed "(fn [a b & c] d)"))
(fact "can have raw strings"  (parse "`one two`")      => (parsed "\"one two\""))
(fact "can have strings"      (parse "\"one two\"")    => (parsed "\"one two\""))
(fact "characters in raw"     (parse "`\n'\"\b`")      => (parsed "\"\\n'\\\"\\b\""))
(fact "characters in strings" (parse "\"\n'\b\"")      => (parsed "\"\n'\b\""))
;; (fact "quotes in strings"     (parse "\"foo\"bar\"")   => (parsed "\"foo\"bar\""))  TODO implement
(fact "multiple expr"         (parse "1;2;3")          => (parsed "1\n\n2\n\n3"))
(fact "multiple expr 2"       (parse "1\n2\n3")        => (parsed "1\n\n2\n\n3"))
(fact "const"                 (parse "const(\na=2\n)\na")=> (parsed "(let [a 2] a)"))
(fact "const indent"          (parse " const(\n  a=2\n )\n a")=> (parsed "(let [a 2] a)"))
(fact "comment"               (parse "//blah blah\naaa")=> (parsed "aaa"))
(fact "comment 1"             (parse " //blah blah \naaa")=> (parsed "aaa"))
(fact "comment 2"             (parse "\n //blah blah\n \naaa")=> (parsed "aaa"))


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
")  => "(ns foo
  (:require [bar.baz :as b])
  (:require [foo.faz.fedudle :as ff]))

(def x (b/bbb \"blah blah\"))

(defn FooBar [iii jjj] (ff/fumanchu {:ooo (fn [m n] (str m n)) :ppp (fn [m n] (str m n)) :qqq qq }))
")
