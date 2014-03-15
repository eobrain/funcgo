(ns leiningen.fgoc
  (:require [clojure.java.shell :as sh])
  (:require [clojure.java.io :as io]))

(defn fgoc [project & args]
  (let
      [result (sh/sh "java" "-jar" "bin/funcgo-0.1.0-SNAPSHOT-standalone.jar" "src/funcgo/core.fgo")
       clj  (:out result)]
    (println (:err result))
    (if (= (:exit result) 0)
      (io/copy clj (io/file "src/funcgo/core.clj"))
      (println "ERROR"))))
    
