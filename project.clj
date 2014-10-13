(defproject org.eamonn.funcgo/funcgo-compiler "0.5.0-SNAPSHOT"
  :description "Compile Functional Go into Clojure"
  :url "http://funcgo.org"
  :license {:name "Eclipse Public License"
            :url "http://www.eclipse.org/legal/epl-v10.html"}
  :dependencies [[org.clojure/clojure "1.6.0"]
                 [org.clojure/core.async "0.1.303.0-886421-alpha"]
                 [instaparse "1.3.3"]
                 [jline "2.11"]
                 [org.clojure/tools.cli "0.3.1"]
                 [commons-lang/commons-lang "2.6"]
                 [inflections "0.9.5"               :scope "test"]
                 [org.clojure/tools.logging "0.3.0" :scope "test"]
                 [clj-logging-config "1.9.12"       :scope "test"]
                 [midje "1.6.3"                     :scope "test"]]
  :profiles {
             :dev {:plugins [[lein-midje "3.1.1"]]}
             :uberjar {:aot :all}}
  :main funcgo.main)
