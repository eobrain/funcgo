(ns funcgo.core-test
  (:use midje.sweet)
  (:require [funcgo.core :refer :all]))

(fact "smallest complete program has no import and a single expression"
      (funcgo-parser "package foo;import ();12345") => [:SourceFile
   [:PackageClause]
   [:ImportDecl]
   [:Expression
    [:UnaryExpr
     [:PrimaryExpr [:Operand [:Literal [:BasicLit [:int_lit [:decimal_lit "12345"]]]]]]]]] )
