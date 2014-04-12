//////
// This file is part of the Funcgo compiler.
//
// Copyright (c) 2012,2013 Eamonn O'Brien-Strain All rights
// reserved. This program and the accompanying materials are made
// available under the terms of the Eclipse Public License v1.0 which
// accompanies this distribution, and is available at
// http://www.eclipse.org/legal/epl-v10.html
//
// Contributors:
// Eamonn O'Brien-Strain e@obrain.com - initial author
//////

// This file contains the entry point for the standalone version of
// the Funcgo compile and is called by the Leiningen plugin.

package  main
import (
        "clojure/java/io"
        "clojure/pprint"
        "clojure/string"
        "clojure/tools/cli"
        "funcgo/core"
)

var commandLineOptions = [
        ["-n", "--nodes", "print out the parse tree that the parser produces"],
        ["-f", "--force", "Force compiling even if not out-of-date"],
        ["-h", "--help", "print help"]
]

// A version of pprint that preserves type hints.
// See https://groups.google.com/forum/#!topic/clojure/5LRmPXutah8
func prettyPrint(obj, writer) {
        const origDispatch = \`pprint/*print-pprint-dispatch*`
        pprint.withPprintDispatch(
                func(o) {
			const met = meta(o)
                        if met {
				print("^")
				if count(met) == 1 {
					if met(TAG) {
						origDispatch(met(TAG))
					} else {
						if met(PRIVATE) == true {
							origDispatch(PRIVATE)
						} else {
							origDispatch(met)
						}
					}
				} else {
					origDispatch(met)
				}
				print(" ")
				pprint.pprintNewline(FILL)
			}
                        origDispatch(o)
                },
                pprint.pprint(obj, writer)
        )
}

func writePrettyTo(cljText, writer java.io.BufferedWriter) {
	for expr := range readString( str("[", cljText, "]")) {
		prettyPrint(expr, writer)
		writer->newLine()
	}
	writer->close()
}


func compileExpression(inPath, fgoText) {
	const (
		cljText = core.Parse(inPath, fgoText, EXPR)
		strWriter = new java.io.StringWriter()
		writer = new java.io.BufferedWriter(strWriter)
	){
		cljText writePrettyTo writer
		strWriter->toString()
	}
}

// func repl(consoleReader jline.console.ConsoleReader) {
// 	const (
// 		//cljText = compileExpression("repl.go", readLine())
// 		cljText = first(core.Parse("repl.go", consoleReader->readLine("fgo=>     "), EXPR))
// 	){
// 		println("Clojure: ", cljText)
// 		println("Result:  ", eval(readString(cljText)))
// 	}
// 	recur(consoleReader)
// }

// func Repl() {
// 	repl(new jline.console.ConsoleReader())
// }

func Repl(){
	const(
		consoleReader = new jline.console.ConsoleReader()
	){
		loop(){
			const (
				cljText = first(core.Parse(
					"repl.go",
					consoleReader->readLine("fgo=>     "),
					EXPR
				))
			){
				println("Clojure: ", cljText)
				println("Result:  ", eval(readString(cljText)))
			}
			recur()
		}
	}
}

func CompileString(inPath, fgoText) {
	const (
		cljText = core.Parse(inPath, fgoText)
		strWriter = new java.io.StringWriter()
		writer = new java.io.BufferedWriter(strWriter)
	){
		cljText writePrettyTo writer
		strWriter->toString()
	}
}

func compileFile(inFile java.io.File, root java.io.File, opts) {
        const(
                inPath = inFile->getPath()
                outFile = io.file(string.replace(inPath, /\.go(s?)$/, ".clj$1"))
        )
        if opts(FORCE) || outFile->lastModified() < inFile->lastModified() {
                println("  ", inPath)
                const(
			prefixLen = root->getAbsolutePath()->length()
			relative = subs(inFile->getAbsolutePath(), prefixLen + 1)
                        cljText = core.Parse(relative, slurp(inFile), SOURCEFILE, opts(NODES))
                        // TODO(eob) open using with-open
                        writer = io.writer(outFile)
                )
		writer->write(str(";; Compiled from ", inFile, "\n"))
		cljText writePrettyTo writer
                println("    -->", outFile->getPath())
                if (outFile->length) / (inFile->length) < 0.5 {
                        println("WARNING: Output file is only",
                                int(100 * (outFile->length) / (inFile->length)),
                                "% the size of the input file")
                }
        }
}

func compileTree(root java.io.File, opts) {
	println("Compiling under root ", root->getName())
	for f := range fileSeq(root) {
		const (
			ff java.io.File = f
			name = ff->getName
		)
		try {
			if name->endsWith(".go") || name->endsWith(".gos") { 
				compileFile(ff, root, opts)
			}
		} catch Exception e {
			println("\n", e->getMessage())
		}
	}
}

func printError(cmdLine) {
	println()
	if cmdLine(ERRORS) {
		println(cmdLine(ERRORS))
	}
	println("USAGE:  fgoc [options] path ...")
	println("options:")
	println(cmdLine(SUMMARY))
}

// Convert Funcgo files to clojure files, using the commandLineOptions
// to parse the arguments.  By default compiles all modified files
// under the current directory.
func Compile(args...) {
	const(
		cmdLine   = args cli.parseOpts commandLineOptions
		otherArgs = cmdLine(ARGUMENTS)
		opts      = cmdLine(OPTIONS)
		here      = io.file(".")
	) {
		if cmdLine(ERRORS) || opts(HELP){
			println(cmdLine(SUMMARY))
		}else{
			if not(seq(otherArgs)) {
				println("Missing directory or file argument.")
				printError(cmdLine)
			} else {
				// file argumentd
				for arg := range otherArgs {
					const (
						file = io.file(arg)
					) {
						if file->isDirectory {
							compileTree(file, opts)
						} else {
							try {
								compileFile(file, here, opts)
							} catch Exception e {
								println("\n", e->getMessage())
							}
						}
					}
				}
			}
		}
	}
}

// Entry point for stand-alone compiler. Usage is the same as for the
// Compile function.
func _main(args...) {
	Compile apply args
}
