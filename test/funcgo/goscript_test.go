package funcgo/goscript_test
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

12345
`),
        =>,
        `(ns p (:require [enfocus.core :as ef]) (:require-macros [enfocus.macros :as em])) 12345`)

func parse(expr) {
	fgo.Parse("foo.gos", "package foo;" str expr)
}

func parsed(expr) {
        str("(ns foo ) ", expr)
}


test.fact("symbol",
	parse(`apple`), =>, parsed(`apple`)
)

test.fact("enfocus",
	fgo.Parse("client.gos", `
package fgosite/client
import (
	ef "enfocus/core"
        "enfocus/effects"
        "enfocus/events"
	"clojure/browser/repl"
)
import macros(
	"enfocus/macros"
)
aaa
`
	), =>,
	`(ns fgosite.client (:require [enfocus.core :as ef] [enfocus.effects :as effects] [enfocus.events :as events] [clojure.browser.repl :as repl]) (:require-macros [enfocus.macros :as macros])) aaa`
)
