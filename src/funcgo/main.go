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

package  funcgo.main
import (
        io     clojure.java.io
        pprint clojure.pprint
        string clojure.string
        cli    clojure.tools.cli
        core   funcgo.core
)

cliOptions := [
        ["-n", "--nodes", "print out the parse tree that the parser produces"],
        ["-f", "--force", "Force compiling even if not out-of-date"],
        ["-h", "--help", "print help"]
]


func compileFile(inFile, opts) {
        const(
                inPath = inFile->getPath()
                outFile = io.file(string.replace(inPath, /\.go$/, ".clj"))
        )
        if opts[FORCE] || outFile->lastModified() < inFile->lastModified() {
                println(inPath)
                const(
                        clj = core.funcgoParse(slurp(inFile), opts[NODES])
                        // TODO(eob) open using with-open
                        writer = io.writer(outFile)
                )
                writer->write(str(";; Compiled from ", inFile, "\n"))
                for expr := range readString( str("[", clj, "]")) {
                        pprint.pprint(expr, writer)
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
func _main(&args) {
	const(
		cmdLine   = args cli.parseOpts cliOptions
		otherArgs = cmdLine[ARGUMENTS]
		opts      = cmdLine[OPTIONS]
	) {
		if cmdLine[ERRORS] || opts[HELP]{
			println("ERROR: ", cmdLine[ERRORS])
			println(cmdLine[SUMMARY])
		}else{
			if not(seq(otherArgs)) {
				for f := range fileSeq(io.file(".")) {
					try {
						if f->getName()->endsWith(".go") { 
							compileFile(f, opts)
						}
					} catch Exception e {
						println("\n", e->getMessage())
					}
				}
			} else {
				for arg := range otherArgs {
					compileFile(io.file(arg), opts)
				}
			}
		}
	}
}
