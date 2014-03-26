package leiningen.fgoc
import(
        sh clojure.java.shell
        io clojure.java.io
)

func fgoc(project, args...) {
        const(
		cmdLine = ["java", "-jar", "bin/funcgo-0.1.21-standalone.jar"] concat args
                result  = sh.sh apply cmdLine
        )
        println(result(ERR))
        println(result(OUT))
        if result(EXIT) == 0 {
                println("Compile finished")
        } else {
                println("ERROR")
        }
}
