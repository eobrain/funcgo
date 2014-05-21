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
		"boolean": TYPE,
		UNUSED_PACKAGES: set{},
		UNUSED_TYPES: set{}
	})
}

// Add a package symbol to the table.
func PackageImported(st, pkg) {
	dosync(st alter func{
		$1 += {
			pkg: PACKAGE,
			UNUSED_PACKAGES: (*st)(UNUSED_PACKAGES) conj pkg
		}
	})
}

// Add a package symbol to the table, but don't require it to be used.
func PackageCreated(st, pkg) {
	dosync(st alter func{$1 += {
		pkg: PACKAGE
	}})
}

// Add a package symbol to the table.
func TypeImported(st, typ) {
	dosync(st alter func{$1 += {
		typ: TYPE,
		UNUSED_TYPES: (*st)(UNUSED_TYPES) conj typ
	}})
	//dosync{
	//	st := $1 += {
	//		typ: TYPE,
	//		UNUSED_TYPES: (*st)(UNUSED_TYPES) conj typ
	//	}
	//}
}

// Add a package symbol to the table.
func TypeCreated(st, typ) {
	dosync(st alter func{$1 += {
		typ: TYPE
	}})
	//dosync{
	//	st := $1 += {typ: TYPE}
	//}
}

// Has this package been previously been added to the table?
func HasPackage(st, pkg) {
	dosync(st alter func{$1 += {
		UNUSED_PACKAGES: (*st)(UNUSED_PACKAGES) disj pkg
	}})
	(*st)(pkg) == PACKAGE
}

// Has this type been previously been added to the table?
func HasType(st, typ) {
	dosync(st alter func{$1 += {
		UNUSED_TYPES: (*st)(UNUSED_TYPES) disj typ
	}})
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

func CheckAllUsed(st) {
	const (
		pkgs = (*st)(UNUSED_PACKAGES)
		typs = (*st)(UNUSED_TYPES)
	)
	if notEmpty(pkgs) {
		const pkgsS = ", " string.join pkgs
		throw(new Exception(str("Packages imported but never used: [", pkgsS, "]")))
	}
	if notEmpty(typs) {
		const typsS = ", " string.join typs
		throw(new Exception(str("Types imported but never used: [", typsS, "]")))
	}

}
