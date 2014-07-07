;; Compiled from tasks/leiningen/fgoc.go
(ns leiningen.fgoc (:gen-class) (:require [clojure.java.shell :as sh]))

(set! *warn-on-reflection* true)

(defn-
 fgoc
 [project & args]
 (let
  [cmd-line
   (concat
    ["java"
     "-jar"
     "bin/funcgo-compiler-0.3.0-standalone.jar"
     "src"
     "test"
     "tasks"]
    args)
   result
   (apply sh/sh cmd-line)]
  (println (result :err))
  (println (result :out))
  (if
   (= (result :exit) 0)
   (println "Compile finished")
   (println "ERROR"))))

