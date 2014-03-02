(ns funcgo.core-test
  (:use midje.sweet)
  (:require [funcgo.core :refer :all]))

(fact "smallest complete program has no import and a single expression"
      (funcgo-parse "package foo;import ();12345") =>
      [:SourceFile
       "(ns foo)"
       [:Expression
        [:UnaryExpr
         [:PrimaryExpr [:Operand [:Literal [:BasicLit [:int_lit [:decimal_lit "12345"]]]]]]]]] )
