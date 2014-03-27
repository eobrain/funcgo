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

package  funcgo/main
import (
        "clojure/java/io"
        "clojure/pprint"
        "clojure/string"
        "clojure/tools/cli"
        "funcgo/core"
)

cliOptions := [
        ["-n", "--nodes", "print out the parse tree that the parser produces"],
        ["-f", "--force", "Force compiling even if not out-of-date"],
        ["-h", "--help", "print help"]
]

// A version of pprint that preserves type hints.
// See https://groups.google.com/forum/#!topic/clojure/5LRmPXutah8
func prettyPrint(obj, writer) {
        const(
                origDispatch = \`pprint/*print-pprint-dispatch*`
        )
        pprint.withPprintDispatch(
                func(o) {
                        if meta(o) {
                                print("^")
                                origDispatch(meta(o)(TAG))
                                print(" ")
                                pprint.pprintNewline(FILL)
                        }
                        origDispatch(o)
                },
                pprint.pprint(obj, writer)
        )
}

func compileFile(inFile java.io.File, opts) {
        const(
                inPath = inFile->getPath()
                outFile = io.file(string.replace(inPath, /\.go$/, ".clj"))
        )
        if opts(FORCE) || outFile->lastModified() < inFile->lastModified() {
                println(inPath)
                const(
                        clj = core.funcgoParse(slurp(inFile), opts(NODES))
                        // TODO(eob) open using with-open
                        writer java.io.BufferedWriter = io.writer(outFile)
                )
                writer->write(str(";; Compiled from ", inFile, "\n"))
                for expr := range readString( str("[", clj, "]")) {
                        prettyPrint(expr, writer)
                        writer->newLine()
                }
                writer->close()
                println("  -->", outFile->getPath())
                if (outFile->length) / (inFile->length) < 0.5 {
                        println("WARNING: Output file is only",
                                int(100 * (outFile->length) / (inFile->length)),
                                "% the size of the input file")
                }
        }
}

 // Convert funcgo to clojure
func _main(args...) {
	const(
		cmdLine   = args cli.parseOpts cliOptions
		otherArgs = cmdLine(ARGUMENTS)
		opts      = cmdLine(OPTIONS)
	) {
		if cmdLine(ERRORS) || opts(HELP){
			println("ERROR: ", cmdLine(ERRORS))
			println(cmdLine(SUMMARY))
		}else{
			if not(seq(otherArgs)) {
				for f := range fileSeq(io.file(".")) {
					const(
						ff java.io.File = f
					)
					try {
						if ff->getName()->endsWith(".go") { 
							compileFile(ff, opts)
						}
					} catch Exception e {
						println("\n", e->getMessage())
					}
				}
			} else {
				for arg := range otherArgs {
					try {
						compileFile(io.file(arg), opts)
					} catch Exception e {
						println("\n", e->getMessage())
					}
				}
			}
		}
	}
}
