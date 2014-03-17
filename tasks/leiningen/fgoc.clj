(ns
 leiningen.fgoc
 (:gen-class)
 (:require [clojure.java.shell :as sh])
 (:require [clojure.java.io :as io]))

(defn
 fgoc
 [project & args]
 (let
  [result (sh/sh "java" "-jar" "bin/funcgo-0.1.9-standalone.jar")]
  (println (result :err))
  (println (result :out))
  (if
   (= (result :exit) 0)
   (println "Compile finished")
   (println "ERROR"))))

