package fgoc
import(
        sh "clojure/java/shell"
)

func fgoc(project, args...) {
        const(
		cmdLine = [
			"java", "-jar", "bin/funcgo-compiler-0.2.3-standalone.jar",
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
