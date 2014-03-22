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

package  funcgo.core
import (
        insta   instaparse.core
        pprint  clojure.pprint
        failure instaparse.failure
        string  clojure.string
	parser  funcgo.parser
	codegen funcgo.codegen
)


func funcgoParse(fgo) {
        const(
                parsed = parser.funcgoParser(string.replace(fgo, /\t/, "        "))
        )
        if insta.isFailure(parsed) {
                failure.pprintFailure(parsed)
                throw(new Exception(`"SYNTAX ERROR"`))
        } else {
	    // pprint.pprint(parsed)
            insta.transform(codegen.codeGenerator, parsed)
        }
}

