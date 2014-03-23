package leiningen.fgoc
import(
        sh clojure.java.shell
        io clojure.java.io
)

func fgoc(project, &args) {
        const( // TODO(eob) pass args to command line
                result = sh.sh("java", "-jar", "bin/funcgo-0.1.17-standalone.jar")
        )
        println(result[ERR])
        println(result[OUT])
        if result[EXIT] == 0 {
                println("Compile finished")
        } else {
                println("ERROR")
        }
}

