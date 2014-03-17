(ns leiningen.fgoc
  (:require [clojure.java.shell :as sh])
  (:require [clojure.java.io :as io]))

(defn fgoc [project & args]
  (let
      [result (sh/sh "java" "-jar" "bin/funcgo-0.1.9-standalone.jar")]
    (println (:err result))
    (println (:out result))
    (if (= (:exit result) 0)
      (println "Compile finished")
      (println "ERROR"))))
