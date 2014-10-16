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


func untabify(s){      string.replace(s, /\t/,           "        ") }

func Ambiguity(fgo) {
  preprocessed := untabify(fgo)
  insta.parses(parser.Parse, preprocessed)
}

func parse(preprocessed, startRule, isAmbiguity) {
	if isAmbiguity {

		parsedList := insta.parses(parser.Parse, preprocessed, START, startRule)
		ambiguity := count(parsedList)
		switch ambiguity {
		case 0: {
			"__preprocessed.go"  spit  preprocessed
			throw(new IOException(
				"Parsing failure.  Turn off ambiguity flag to see details."))
		}
		case 1:
			parsedList[0]
		default: {
			print(" WARNING, ambiguity=", ambiguity)
			parsedList[0]
		}
		}

	} else {

		parsed := parser.Parse(preprocessed, START, startRule)
		if insta.isFailure(parsed) {
			"__preprocessed.go"  spit  preprocessed
			throw(new IOException(str(withOutStr(failure.pprintFailure(parsed)))))
		} else {
			parsed
		}

	}
}

func Parse(path, fgo) {
	Parse(path, fgo, SOURCEFILE)
} (path, fgo, startRule) {
	Parse(path, fgo, startRule, false, false, false)
} (path, fgo, startRule, isNodes, isSync, isAmbiguity) {
	preprocessed := untabify(fgo)
	parsed := parse(preprocessed, startRule, isAmbiguity)
	if isNodes {
		pprint.pprint(parsed)
	}
	codegen.Generate(path, parsed, isSync)
}
