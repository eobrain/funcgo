package funcgo/goscript_test
import (
        test "midje/sweet"
        fgo "funcgo/core"
)


test.fact("Clojure script import",
        fgo.funcgoParse("p.gos", `
package p
import(
  em "macros//enfocus/macros"
  ef "enfocus/core"
)

12345
`),
        =>,
        `(ns p 
  (:require-macros [enfocus.macros :as em])
  (:require [enfocus.core :as ef]))

12345
`)

func parse(expr) {
	fgo.funcgoParse("foo.gos", "package foo;import ()" str expr)
}

func parsed(expr) {
        str("(ns foo )\n\n", expr, "\n")
}


test.fact("symbol",
	parse(`apple`), =>, parsed(`apple`)
)
