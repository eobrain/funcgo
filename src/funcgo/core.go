//////
// This file is part of the Funcgo compiler.
//
// Copyright (c) 2014 Eamonn O'Brien-Strain All rights
// reserved. This program and the accompanying materials are made
// available under the terms of the Eclipse Public License v1.0 which
// accompanies this distribution, and is available at
// http://www.eclipse.org/legal/epl-v10.html
//
// Contributors:
// Eamonn O'Brien-Strain e@obrain.com - initial author
//////

package  core
import (
        insta "instaparse/core"
        "clojure/pprint"
        "instaparse/failure"
        "clojure/string"
	"funcgo/parser"
	"funcgo/codegen"
)
import type java.io.IOException


func escapeRawString(s){
	string.replace(
		s,
		/\x60([^\x60]*)\x60/,
		func(matched){
			raw := matched[1]
			newlines := string.replace(raw, /[^\n]/, "")
			escaped := str(`"`, raw  string.escape  charEscapeString, '"')
			escaped str newlines
		}
	)
	// println("escaped=", str("<<",escaped,">>"), "newlines=",str("<<",newlines,">>"))
}
// Strip trailing comments and whitespace
func stripTrailing(s){ string.replace(s, /([ \t]*\/\/[^\n]*|[ \t]+)\n/, "\n") }
func untabify(s){      string.replace(s, /\t/,           "        ") }


func Parse(path, fgo) {
	Parse(path, fgo, SOURCEFILE)
} (path, fgo, startRule) {
	Parse(path, fgo, startRule, false, false)
} (path, fgo, startRule, isNodes, isSync) {
	preprocessed := untabify(stripTrailing(escapeRawString(fgo)))
        parsed := parser.Parse(preprocessed, START, startRule)
        if insta.isFailure(parsed) {
		"__preprocessed.go"  spit  preprocessed
		throw(new IOException(str(withOutStr(failure.pprintFailure(parsed)))))
        } else {
		if isNodes {
			pprint.pprint(parsed)
		}
		codegen.Generate(path, parsed, isSync)
        }
}
