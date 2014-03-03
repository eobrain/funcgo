(ns funcgo.core-test
  (:use midje.sweet)
  (:require [funcgo.core :refer :all]))

(fact "smallest complete program has no import and a single expression"
      (funcgo-parse "package foo;import ();12345") =>
      "(ns foo) 12345")

(fact "package can be dotted"
      (funcgo-parse "package foo.bar;import ();12345") =>
      "(ns foo.bar) 12345")
