(ns leiningen.fgoc
  (:require [clojure.java.shell :as sh])
  (:require [clojure.java.io :as io]))

(defn fgoc [project & args]
  (let
      [result (sh/sh "java" "-jar" "bin/funcgo-0.1.0-standalone.jar" "src/funcgo/core.go")
       clj  (:out result)]
    (println (:err result))
    (if (= (:exit result) 0)
      (if (.startsWith clj "Parse error at line")
        (println clj)
        (io/copy clj (io/file "src/funcgo/core.clj")))
      (println "ERROR"))))
    
