package goscript_test
import (
        test "midje/sweet"
        fgo "funcgo/core"
)


test.fact("Clojure script import",
        fgo.Parse("p.gos", `
package p
import(
  ef "enfocus/core"
)
import macros(
  em "enfocus/macros"
)

ef.someFunction
em.someMacro
`),
        =>,
        str(
		`(ns p (:require [enfocus.core :as ef]) (:require-macros [enfocus.macros :as em]))`,
		` ef/some-function em/some-macro`
	)
)

func parse(expr) {
	fgo.Parse("foo.gos", "package foo;"  str  expr)
}

func parsed(expr) {
        str("(ns foo ) ", expr)
}


test.fact("symbol",
	parse(`apple`), =>, parsed(`apple`)
)

test.fact("enfocus",
	fgo.Parse("fgosite/client.gos", `
package client
import (
	ef "enfocus/core"
        "enfocus/effects"
        "enfocus/events"
	"clojure/browser/repl"
)
import macros(
	"enfocus/macros"
)
ef.a
effects.b
events.c
repl.d
macros.e
`
	), =>, str(
		`(ns fgosite.client`,
		` (:require [enfocus.core :as ef]`,
		` [enfocus.effects :as effects]`,
		` [enfocus.events :as events]`,
		` [clojure.browser.repl :as repl])`,
		` (:require-macros [enfocus.macros :as macros]))`,
		` ef/a effects/b events/c repl/d macros/e`
	)
)
