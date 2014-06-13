(defproject org.eamonn.funcgo/funcgo-compiler "0.2.5"
  :description "Compile Functional Go into Clojure"
  :url "http://funcgo.org"
  :license {:name "Eclipse Public License"
            :url "http://www.eclipse.org/legal/epl-v10.html"}
  :dependencies [[org.clojure/clojure "1.5.1"]
                 [org.clojure/core.async "0.1.303.0-886421-alpha"]
                 [instaparse "1.3.2"]
                 [jline "2.11"]
                 [org.clojure/tools.cli "0.3.1"]
                 [commons-lang/commons-lang "2.6"]
                 [inflections "0.9.5" :scope "test"]
                 [midje "1.5.1" :scope "test"]]
  :profiles {
             :dev {:plugins [[lein-midje "3.1.1"]]}
             :uberjar {:aot :all}}
  :main funcgo.main)
