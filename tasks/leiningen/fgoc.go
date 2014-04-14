package fgoc
import(
        sh "clojure/java/shell"
        "clojure/java/io"
)

func fgoc(project, args...) {
        const(
		cmdLine = [
			"java", "-jar", "bin/funcgo-0.1.30-standalone.jar",
			"src", "test", "tasks"
		] concat args
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
