(defproject funcgo "0.1.0"
  :description "Compiler from Functional Go into Clojure"
  :url "http://funcgo.com"
  :license {:name "Eclipse Public License"
            :url "http://www.eclipse.org/legal/epl-v10.html"}
  :dependencies [[org.clojure/clojure "1.5.1"]
                 [instaparse "1.2.16"]
                 [midje "1.5.1" :scope "test"]]
  :main funcgo.core)
