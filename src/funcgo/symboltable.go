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

// A symbol table is a mutable state that keeps track of the symbols
// declared, so that the codegenerator can throw exceptions when it
// encounters an undefined symbol.

package symboltable
import "clojure/string"

// Return a new symbol table.
func New() {
	ref({
		"long": TYPE,
		"double": TYPE,
		"boolean": TYPE
	})
}

// Add a package symbol to the table.
func AddPackage(st, pkg) {
	dosync(alter(st, assoc, pkg, PACKAGE))
}

// Add a package symbol to the table.
func AddType(st, pkg) {
	dosync(alter(st, assoc, pkg, TYPE))
}

// Has this package been previously been added to the table?
func HasPackage(st, pkg) {
	(*st)(pkg) == PACKAGE
}

// Has this type been previously been added to the table?
func HasType(st, typ) {
	(*st)(typ) == TYPE
}

// Return a string representation of packages in the table.
func Packages(st) {
	const packages = for [symbol, key] := lazy *st if key == PACKAGE { symbol }
	str("[", ", " string.join packages, "]")
}

// Return a string representation of types in the table.
func Types(st) {
	const packages = for [symbol, key] := lazy *st if key == TYPE { symbol }
	str("[", ", " string.join packages, "]")
}
