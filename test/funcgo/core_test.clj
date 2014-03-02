(ns funcgo.core-test
  (:require [clojure.test :refer :all]
            [funcgo.core :refer :all]))

(deftest replace-me ;; FIXME: write
  (is (funcgo-parser "
package foo
import (
aaa.bbb
aaa.ccc
)
12345
") [:SourceFile
   [:PackageClause]
   [:ImportDecl
    [:ImportSpec [:ImportPath [:identifier "aaa"] [:identifier "bbb"]]]
    [:ImportSpec [:ImportPath [:identifier "aaa"] [:identifier "ccc"]]]]
   [:Expression
    [:UnaryExpr
     [:PrimaryExpr [:Operand [:Literal [:BasicLit [:int_lit [:decimal_lit "12345"]]]]]]]]] ))
