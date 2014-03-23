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
        // An option with a required argument
        ["-p", "--port PORT", "Port number",
                DEFAULT, 80,
                PARSE_FN, func(x){Integer::parseInt(x)},
                VALIDATE, [
                        func(x){ func(x){0 < x && x < 65536}}, //0x10000}},
                        "Must be a number between 0 and 65536"
                ]
        ],
        // A non-idempotent option
        ["-v", nil, "Verbosity level",
                ID, VERBOSITY,
                DEFAULT, 0,
                ASSOC_FN, func(m, k, _) { updateIn( m, [k], inc)}
        ],
        // A boolean option defaulting to nil
        ["-h", "--help"]
]


func compileFile(inFile) {
        const(
                inPath = inFile->getPath()
                outFile = io.file(string.replace(inPath, /\.go$/, ".clj"))
        )
        if outFile->lastModified() < inFile->lastModified() {
                const(
                        clj = core.funcgoParse(slurp(inFile))
                        // TODO(eob) open using with-open
                        writer = io.writer(outFile)
                )
                println(inPath)
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
        //args cli.parseOpts cliOptions
	//println(cli.parseOpts(args, cliOptions))
        if not(seq(args)) {
                for f := range fileSeq(io.file(".")) {
                        try {
                                if f->getName()->endsWith(".go") { 
                                        compileFile(f)
                                }
                        } catch Exception e {
                                println("\n", e->getMessage())
                        }
                }
        }else{
                for arg := range args {
                        compileFile(io.file(arg))
                }
        }
}
