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

package  codegen
import (
        s     "clojure/string"
        insta "instaparse/core"
	symbols "funcgo/symboltable"
)

// Returns a map of parser targets to functions that generate the
// corresponding Clojure code.
func codeGenerator(symbolTable) {

	func listStr(item...) {
		str("(", s.join(" ", item), ")")
	}
	
	func blankJoin (xs...){
		" " s.join xs
	}
	
	func vecStr(item...) {
		str("[", s.join(" ", item), "]")
	}
	
	func infix(expression) {
		expression
	} (left, operator, right) {
		listStr(operator, left, right)
	}
	
	// Capitalized
	func isPublic(identifier) {
		// TODO(eob) handle general Unicode
		/^[A-Z]/ reFind identifier
	}
	
	// Return a function that always returns the given constant string.
	func constantFunc(s) {
		func{s}
	}
	
	func splitPath(path String) {
		const(
			slash = path->lastIndexOf(int('/'))
			beforeSlash = subs(path, 0, slash + 1)
			afterSlash = subs(path, slash + 1)
		)
		[
			s.replace(beforeSlash, '/', '.'),
			s.replace(afterSlash, /\.gos?$/, "")
		]
	}
	
	func declBlockFunc(typ) {
		func(xs...){
			const(
				consts = butlast(xs)
				expressions = last(xs)
			)
			str("(", typ, " [",
				" " s.join consts,
				"] ",
				expressions,
				")")
		}
	}

	func importSpec(imported) {
		importSpec(last(imported s.split /\./), imported)
	} (identifier, imported) {
		symbolTable symbols.AddPackage identifier
		str("[", imported, " :as ", identifier, "]")
	}
	func externImportSpec(identifier) {
		symbolTable symbols.AddPackage identifier
		""
	}

	// Mapping from parse tree to generators of CLJ code.
	{
		SOURCEFILE:     func(header, body) {str(header, " ", body)},
		IMPORTDECLS: blankJoin,
		IMPORTDECL:     func() {
			""
		} (importSpecs...) {
			listStr apply (":require" cons importSpecs)
		},
		MACROIMPORTDECL:     func() {
			""
		} (importSpecs...) {
			listStr apply (":require-macros" cons importSpecs)
		},
		IMPORTSPEC: importSpec,
		EXTERNIMPORTSPEC: externImportSpec,
		PRECEDENCE0: infix,
		PRECEDENCE1: infix,
		PRECEDENCE2: infix,
		PRECEDENCE3: infix,
		PRECEDENCE4: infix,
		PRECEDENCE5: infix,
		IFELSEEXPR: func(condition, exprs) {
			listStr("when", condition, exprs)
		} (condition, block1, block2) {
			listStr("if", condition, block1, block2)
		},
		LETIFELSEEXPR: func(lhs, rhs, condition, exprs) {
			str("(let [", lhs, " ", rhs, "] ",
				listStr("when", condition, exprs),
				")"
			)
		} (lhs, rhs, condition, block1, block2) {
			str("(let [", lhs, " ", rhs, "] ",
				listStr("if", condition, block1, block2),
			")"
			)
		},
		BOOLSWITCH: func(clauses...) {
			listStr apply ("cond" cons clauses)
		},
		BOOLCASECLAUSE: blankJoin,
		BOOLSWITCHCASE: func(){
			":else"
		} (cond) {
			cond
		},
		CONSTSWITCH: func(expr, clauses...) {
			listStr apply ("case" cons (expr cons clauses))
		},
		CONSTCASECLAUSE: blankJoin,
		CONSTANTLIST: func(c) {
			c
		}(c0, c...){
			listStr apply (c0 cons c)
		},
		CONSTSWITCHCASE: func(){
			""
		} (cond) {
			cond
		},
		FORRANGE: func(identifier, seq, expressions) {
			str("(doseq [", identifier, " ", seq, "] ", expressions, ")")
		},
		FORLAZY: func(identifier, seq, expressions) {
			str("(for [", identifier, " ", seq, "] ", expressions, ")")
		} (identifier, seq, condition, expressions) {
			str("(for [", identifier, " ", seq, " :when ", condition, "] ", expressions, ")")
		},
		FORTIMES: func(identifier, count, expressions) {
			str("(dotimes [", identifier, " ", count, "] ", expressions, ")")
		},
		TRYEXPR: func(expressions, catches) {
			listStr("try", expressions, catches)
		} (expressions, catches, finally) {
			listStr("try", expressions, catches, finally)
		},
		CATCHES: blankJoin,
		CATCH: func(typ, exception, expressions) {
			listStr("catch", typ, exception, expressions)
		},
		FINALLY: func{listStr("finally", ..)},
		NEW:     func{str(.., ".")},
		SHORTVARDECL:   func(identifier, expression) {
			if isPublic(identifier) {
				listStr("def", identifier, expression)
			} else {
				listStr("def", "^:private", identifier, expression)
			}
		},
		VARDECL:   func(identifier, expression) {
			if isPublic(identifier) {
				listStr("def", identifier, expression)
			} else {
				listStr("def", "^:private", identifier, expression)
			}
		} (identifier, typ, expression) {
		if isPublic(identifier) {
			listStr("def", "^" str typ, identifier, expression)
		} else {
			listStr("def",
				"^:private",
				"^" str typ,
				identifier,
				expression
			)
		}
		},
		FUNCTIONCALL:    func(function) {
			listStr(function)
		} (function, call) {
			listStr(function, call)
		},
		EXPRESSIONLIST: blankJoin,
		EXPRESSIONS:    blankJoin,
		CONSTS:         blankJoin,
		BLOCK: func (expr){
			expr
		} (expr0, exprRest...) {
			str("(do ",  " " s.join (expr0 cons exprRest),  ")")
		},
		INDEXED: func(xs, i){ listStr("nth", xs, i) },
		WITHCONST: declBlockFunc("let"),
		LOOP:      declBlockFunc("loop"),
		CONST: func(identifier, expression) {
			str(identifier, " ", expression)
		},
		VECDESTRUCT: vecStr,
		DICTDESTRUCT: func{str('{', (" " s.join ...), "}")}, 
		DICTDESTRUCTELEM: func(destruct, label) {
			str(destruct, " ", label)
		},
		VARIADICDESTRUCT:  func{str("& ", ..)},
		SYMBOL: func(identifier){
			identifier
		} (pkg, identifier) {
			if symbolTable symbols.HasPackage pkg {
				throw(new Exception(format(
					`package "%s" in %s.%s does not appear in imports %s`,
					pkg, pkg, identifier, symbols.Packages(symbolTable))))
			}
			str(pkg, "/", identifier)
		},
		BINARYOP: identity,
		MULOP: identity,
		ADDOP: identity,
		RELOP: identity,
		FUNCTIONDECL:   func(identifier, function) {
			const defn = if isPublic(identifier) { "defn" } else { "defn-" }
			listStr(defn, identifier, function)
		},
		FUNCLIKEDECL:   func(funclike, identifier, function) {
			listStr(funclike, identifier, function)
		},
		FUNCTIONLIT:    func{listStr("fn", ..)},
		SHORTFUNCTIONLIT:  func(expr) {
			if first(expr) == '(' && last(expr) == ')' {
				"#" str expr
			}else{
				listStr("fn", "[]", expr)
		}
		},
		STRUCTSPEC: func(javaIdentifier, fields...) {
			listStr("defrecord", javaIdentifier, str("[", blankJoin apply fields, "]"))
		},
		FIELDS: blankJoin,
		INTERFACESPEC: func(args...){
			listStr apply ("defprotocol" cons args)
		},
		VOIDMETHODSPEC: func(javaIdentifier) {
			listStr(javaIdentifier, "[this]")
		}(javaIdentifier, methodparams) {
			listStr(javaIdentifier, str("[this ", methodparams, "]"))
		},
		TYPEDMETHODSPEC: func(javaIdentifier, typ) {
			listStr("^" str typ, javaIdentifier, "[this]")
		} (javaIdentifier, methodparams, typ) {
			listStr("^" str typ, javaIdentifier, str("[this ", methodparams, "]"))
		},
		IMPLEMENTS: func(protocol, concrete, methodimpls...) {
			listStr apply concat(list("extend-type", concrete, protocol), methodimpls)
		},
		METHODIMPL: func(javaIdentifier, function) {
			listStr(javaIdentifier, function)
		},
		METHODPARAMETERS: blankJoin,
		METHODPARAM: func(symbol) {
			symbol
		} (symbol, typ) {
			str("^", typ, " ", symbol)
		},
		PERCENT: constantFunc("%"),
		PERCENTNUM: func{"%" str ..},
		PERCENTVARADIC: constantFunc("%&"),
		FUNCTIONPARTS:  func{str("(",  ") (" s.join ...,  ")")},
		FUNCTIONPART0:  func(expression) {
			"[] " str expression
		} (typ, expression) {
			str("^", typ, " [] ", expression)
		},
		VFUNCTIONPART0:  func(variadic, expression) {
			str("[", variadic, "] ", expression)
		} (variadic, typ, expression) {
			str("^", typ, " [", variadic, "] ", expression)
		},
		FUNCTIONPARTN:  func(parameters, expression) {
			str("[", parameters, "] ", expression)
		} (parameters, typ, expression) {
			str("^", typ, " [", parameters, "] ", expression)
		},
		VFUNCTIONPARTN: func(parameters, variadic, expression) {
			str("[", parameters, " ", variadic, "] ", expression)
		} (parameters, variadic, typ, expression) {
			str("^", typ, " [", parameters, " ", variadic, "] ", expression)
		},
		UNTYPEDMETHODIMPL: func(name, block) {
			listStr(name, str("[this]"), block)
		} (name, params, block) {
			listStr(name, str("[this ", params, "]"), block)
		},
		TYPEDMETHODIMPL: func(name, typ, block) {
			listStr("^" str typ, name, str("[this]"), block)
		} (name, params, typ, block) {
			listStr("^" str typ, name, str("[this ", params, "]"), block)
		},
		PARAMETERS:     blankJoin,
		VARIADIC:       func{"& " str ..},
		VECLIT:         func{str("[",  " " s.join ...,	"]")},
		DICTLIT:        func{str apply ...},
		DICTELEMENT:    func(key, value) {str(key, " ", value, " ")},
		SETLIT:         func{str("#{",  " " s.join ...,  "}")},
		LABEL:          func{str(":", s.replace(s.lowerCase(..), /_/, "-"))},
		IDENTIFIER:     func(string) {
			s.replace(
				string,
				/\p{Ll}\p{Lu}/,
				func{str(first(..), "-", s.lowerCase(last(..)))}
			)
		},
		TYPEDIDENTIFIER: func(identifier, typ) {
			str(`^`, typ, " ", identifier)
		},
		IMPORTED:         func{"." s.join ...},
		DECIMALLIT:    identity,
		BIGINTLIT:     str,
		BIGFLOATLIT:   str,
		FLOATLIT:      str,
		DECIMALS:      identity,
		EXPONENT:      str,
		REGEX:         func{str(`#"`,  str apply ...,  `"`)},
		ESCAPEDSLASH: constantFunc(`/`),
		INTERPRETEDSTRINGLIT: func{str(`"`,  str apply ...,  `"`)},
		CLOJUREESCAPE: identity,
		LITTLEUVALUE:  func(d1,d2,d3,d4){str(`\u`,d1,d2,d3,d4)},
		OCTALBYTEVALUE:  func(d1,d2,d3){str(`\o`,d1,d2,d3)},
		UNICODECHAR:   func{`\` str ..},
		NEWLINECHAR:   constantFunc(`\newline`),
		SPACECHAR:     constantFunc(`\space`),
		BACKSPACECHAR: constantFunc(`\backspace`),
		RETURNCHAR:    constantFunc(`\return`),
		TABCHAR:       constantFunc(`\tab`),
		BACKSLASHCHAR: constantFunc(`\\`),
		SQUOTECHAR:    constantFunc(`\'`),
		DQUOTECHAR:    constantFunc(`\"`),
		HEXDIGIT:      identity,
		OCTALDIGIT:    identity,
		RAWSTRINGLIT:  func{str(`"`, s.escape(.., charEscapeString), `"`)},
		DASHIDENTIFIER: func{ "-" str ..},
		ISIDENTIFIER:   func(initial, identifier) {
			str( s.lowerCase(initial), identifier, "?")
		},
		EQUALS: constantFunc("="),
		AND:    constantFunc("and"),
		OR:     constantFunc("or"),
		MUTIDENTIFIER:  func(initial, identifier) {
			str( s.lowerCase(initial), identifier, "!")
		},
		ESCAPEDIDENTIFIER:  identity,
		UNARYEXPR: func(operator, expression) { listStr(operator, expression) },
		NOTEQ:       constantFunc("not="),
		BITXOR:      constantFunc("bit-xor"),
		BITNOT:      constantFunc("bit-not"),
		SHIFTLEFT:   constantFunc("bit-shift-left"),
		SHIFTRIGHT:  constantFunc("bit-shift-right"),
		NOT:         constantFunc("not"),
		MOD:         constantFunc("mod"),
		DEREF:           func{"@"   str ..},
		SYNTAXQUOTE:     func{"`"   str ..},
		UNQUOTE:         func{"~"   str ..},
		UNQUOTESPLICING: func{ "~@" str ..},
		JAVAFIELD:      func(expression, identifier) {
			listStr(".", expression, identifier)
		},
		JAVASTATIC:      func{"/" s.join ...},
		TYPENAME:        func{"." s.join ...},
		UNDERSCOREJAVAIDENTIFIER: func{ "-" str ..},
		JAVAMETHODCALL: func(expression, identifier) {
			str("(. ", expression, " (", identifier, "))")
		} (expression, identifier, call) {
			str("(. ", expression, " (", identifier, " ", call, "))")
		},
		LONG: constantFunc("long"),
		DOUBLE: constantFunc("double")
	}
}

func packageclauseFunc(symbolTable, path String) {
	const (
		isGoscript    = path->endsWith(".gos")
		[parent, name] = splitPath(path)
	)
	if isGoscript {
		symbolTable symbols.AddPackage "js"
	}
	func(imported, importDecls) {
		const fullImported = parent str imported
		if imported != name {
			throw(new java.lang.Exception(str(
				`Got package "`, imported, `" instead of expected "`,
				name, `" in "`, path, `"`
			)))
		}
		if isGoscript {
			listStr("ns", fullImported, importDecls)
		} else {
			str(
				listStr("ns",
					fullImported,
					"(:gen-class)",
					importDecls
				),
				" (set! *warn-on-reflection* true)"
			)
		}
	}
}

// Return the Clojure code generated from the given parse tree.
func Generate(path String, parsed) {
	const (
		symbolTable = symbols.New()
		codeGen = assoc(
			codeGenerator(symbolTable),
			PACKAGECLAUSE,
			packageclauseFunc(symbolTable, path)
		)
	)
	insta.transform(codeGen, parsed)
}
	
