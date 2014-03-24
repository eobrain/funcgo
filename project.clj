(defproject funcgo "0.1.19"
  :description "Compiler from Functional Go into Clojure"
  :url "http://funcgo.com"
  :license {:name "Eclipse Public License"
            :url "http://www.eclipse.org/legal/epl-v10.html"}
  :dependencies [[org.clojure/clojure "1.5.1"]
                 [instaparse "1.2.16"]
                 [org.clojure/tools.cli "0.3.1"]
                 [inflections "0.9.5" :scope "test"]
                 [midje "1.5.1" :scope "test"]]
  :profiles {:dev {:plugins [[lein-midje "3.1.1"]]}}

  :main funcgo.main)
