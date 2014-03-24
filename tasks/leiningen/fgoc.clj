;; Compiled from ./tasks/leiningen/fgoc.go
(ns
 leiningen.fgoc
 (:gen-class)
 (:require [clojure.java.shell :as sh])
 (:require [clojure.java.io :as io]))

(defn
 fgoc
 [project & args]
 (let
  [cmd-line
   (concat ["java" "-jar" "bin/funcgo-0.1.18-standalone.jar"] args)
   result
   (apply sh/sh cmd-line)]
  (println (result :err))
  (println (result :out))
  (if
   (= (result :exit) 0)
   (println "Compile finished")
   (println "ERROR"))))

