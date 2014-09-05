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

package	 codegen
import (
	s     "clojure/string"
	insta "instaparse/core"
	symbols "funcgo/symboltable"
)
import type (
	java.util.List
	java.io.IOException
)

kAsyncRules := set{
	ASYNCPREFIX,
	CHAN,
	TAKE,
	TAKEINGO,
	SENDSTMT,
	SENDSTMTINGO,
	SELECTSTMT,
	SELECTSTMTINGO
}

// Returns a map of parser targets to functions that generate the
// corresponding Clojure code.
func codeGenerator(symbolTable, isGoscript) {

	func noDot(s String) {
		!(/\./ reFind s)
	}

	func isJavaClass(clazz) {
		try{
			Class::forName(clazz)
			true
		}catch Exception e { //ClassNotFoundException e {
			false
		}
	}

	func hasType(typ String) {
		(symbolTable symbols.HasType typ)
		|| isGoscript && typ->startsWith("js.")
		|| !isGoscript && noDot(typ) && isJavaClass("java.lang." str typ)
	}

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
		// not lowercode
		!(/^\p{Ll}/ reFind identifier) || identifier == "main" ||(/^bit-/ reFind identifier)
	}

	// Return a function that always returns the given constant string.
	func constantFunc(s) {
		func{s}
	}

	func splitPath(path String) {
		slash       := path->lastIndexOf(int('/'))
		beforeSlash := subs(path, 0, slash + 1)
		afterSlash  := subs(path, slash + 1)
		[
			s.replace(beforeSlash, '/', '.'),
			s.replace(afterSlash, /\.gos?$/, "")
		]
	}

	func declBlockFunc(typ) {
		func(xs...){
			consts      := butlast(xs)
			expressions := last(xs)
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
		if identifier == "_" {
			// package imported for sideeffect only
			vecStr(imported)
		} else {
			// normal import
			symbolTable symbols.PackageImported identifier
			vecStr(imported, ":as", identifier)
		}
	}
	func externImportSpec(identifier) {
		symbolTable symbols.PackageImported identifier
		""
	}

	func vardecl(identifier, expression) {
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
	}

	func def_(identifier, expression) {
		if isPublic(identifier) {
			listStr("def", identifier, expression)
		} else {
			listStr("def", "^:private", identifier, expression)
		}
	}

	func sendClause(channel, val, expr) {
		vecStr(channel  vecStr  val)  blankJoin   expr
	}

	func doStr(expressions) {
		"do"  listStr  expressions
	}

	func stripQuotes(literal string) string{
		literal->substring(1, literal->length() - 1)
	}

	// Mapping from parse tree to generators of CLJ code.
	{
		SOURCEFILE:  blankJoin,
		NONPKGFILE:  identity,
		IMPORTDECLS: blankJoin,
		IMPORTSPEC: importSpec,
		EXTERNIMPORTSPEC: externImportSpec,
		EXCLUDE: func(symbols...) {
			listStr(":refer-clojure", ":exclude", vecStr(...symbols))
		},
		TYPEIMPORTDECL: func() {
			""
		} (importSpecs...) {
			listStr(":import", ...importSpecs)
		},
		TYPEIMPORTSPEC: func(typepackage, typeclasses...) {
			for typeclass := range typeclasses {
				symbolTable symbols.TypeImported typeclass
			}
			listStr(typepackage, ...typeclasses)
		},
		TYPEPACKAGEIMPORTSPEC: func{
			"." s.join $*
		},
		TYPECLASSESIMPORTSPEC: blankJoin,
		PRECEDENCE00: infix,
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
		ASSOC: func(symbol, items...) {
			listStr("assoc", symbol, ...items)
		},
		DISSOC: func(symbol, items...) {
			listStr("dissoc", symbol, ...items)
		},
		ASSOCIN: func(symbol, path, value) {
			listStr("assoc-in", symbol, path, value)
		},
		ASSOCITEM: blankJoin,
		ASSOCINPATH: vecStr,
		BOOLSWITCH: func(clauses...) {
			listStr("cond", ...clauses)
		},
		BOOLCASECLAUSE: blankJoin,
		BOOLSWITCHCASE: func(){
			":else"
		} (cond) {
			cond
		},
		SELECTSTMT: func(clauses...){
			listStr("alt!!", ...clauses)
		},
		SENDCLAUSE: func(channel, value) {
			sendClause(channel, value, "nil")
		} (channel, value, expressions) {
			sendClause(channel, value, doStr(expressions))
		},
		RECVCLAUSE: func(channel) {
			channel  blankJoin  "nil"
		} (channel, expressions) {
			channel  blankJoin  doStr(expressions)
		},
		RECVVALCLAUSE: func(identifier, channel, expressions) {
			channel  blankJoin  (vecStr(identifier)  listStr   expressions)
		},
		DEFAULTCLAUSE: func() {
			":default"
		} (expessions) {
			":default"  blankJoin  doStr(expessions)
		},
		SELECTSTMTINGO: func(clauses...){
			listStr("alt!", ...clauses)
		},
		SENDCLAUSEINGO: func(channel, value) {
			sendClause(channel, value, "nil")
		} (channel, value, expressions) {
			sendClause(channel, value, doStr(expressions))
		},
		RECVCLAUSEINGO: func(channel) {
			channel  blankJoin  "nil"
		} (channel, expressions) {
			channel  blankJoin  doStr(expressions)
		},
		RECVVALCLAUSEINGO: func(identifier, channel, expressions) {
			channel  blankJoin  (vecStr(identifier)  listStr   expressions)
		},
		TYPESWITCH: func(x, args...) {
			loop(acc="(cond", remaining=args) {
				func typeCase() {
					typ  := first(remaining)
					expr := second(remaining)
					str(acc, " ", listStr("instance?", typ, x), " ", expr)
				}
				switch count(remaining) {
				case 1: {
					[expr] := remaining
					str(acc, " :else ", expr, ")")
				}
				case 2:
					typeCase()  str  ")"
				default:
					recur(typeCase(), 2  drop  remaining)
				}
			}
		},
		CONSTSWITCH: func(expr, clauses...) {
			listStr("case", expr, ...clauses)
		},
		LETCONSTSWITCH: func(lhs, rhs, expr, clauses...) {
			str("(let [", lhs, " ", rhs, "] ",
				listStr("case", expr, ...clauses),
				")"
			)
		},
		CONSTCASECLAUSE: blankJoin,
		CONSTANTLIST: func(c) {
			c
		}(c0, c...){
			listStr(c0, ...c)
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
		FORCSTYLE: func(ident, identAgain, count, identYetAgain, expressions) {
			if ident != identAgain || ident != identYetAgain {
				throw(new IOException(
					`cannot mix different identifiers in c-style for loop`
				))
			}
			str("(dotimes [", ident, " ", count, "] ", expressions, ")")
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
		FINALLY: func{listStr("finally", $1)},
		NEW:	 func{str($1, ".")},
		SHORTVARDECL:	func(identifier, expression) {
			def_(identifier, expression)
		} (ident1, ident2, expr1, expr2) {
			def_(ident1, expr1)  blankJoin  def_(ident2, expr2)
		} (ident1, ident2, ident3, expr1, expr2, expr3) {
			blankJoin(
				def_(ident1, expr1),
				def_(ident2, expr2),
				def_(ident3, expr3)
			)
		},
		PRIMARRAYVARDECL: func(identifier, number, primtype) {
			elements := blankJoin(...(for _ := times readString(number) {"0"}))
			listStr("def", identifier, listStr("vector-of", ":" str primtype, elements))
		},
		ARRAYVARDECL: func(identifier, number, typ) {
			elements := blankJoin(...(for _ := times readString(number) {"nil"}))
			listStr("def", identifier, listStr("vector", elements))
		},
		VARDECL1: vardecl,
		VARDECL2: func(identifier1, identifier2, expression1, expression2) {
			blankJoin(
				vardecl(identifier1, expression1),
				vardecl(identifier2, expression2)
			)
		} (identifier1, identifier2, typ, expression1, expression2) {
			blankJoin(
				vardecl(identifier1, typ, expression1),
				vardecl(identifier2, typ, expression2)
			)
		},
		PREFIXEDROUTINE: listStr,
		PREFIXEDBLOCK: listStr,
		PREFIX: identity,
		ASYNCPREFIX: identity,
		VARIADICCALL: func(function, params...) {
			listStr("apply", function, ...params)
		},
		FUNCTIONCALL:	 listStr,
		LEN: func(call) {
			listStr("count", call)
		},
		CHAN:		func() {
			"(chan)"
		} (n) {
			listStr("chan", n)
		},
		EXPRESSIONLIST: blankJoin,
		EXPRESSIONS:	blankJoin,
		CONSTS:	blankJoin,
		ASSIGNS:	blankJoin,
		COMMACONSTS:	blankJoin,
		BLOCK: func (expr){
			expr
		} (expr0, exprRest...) {
			str("(do ",  " " s.join (expr0 cons exprRest),	")")
		},
		TYPECONVERSION: listStr,
		INDEXED: func(xs, i){ listStr("nth", xs, i) },
		TAKESLICE: func(xs, i){ listStr("take", i, xs) },
		DROPSLICE: func(xs, i){ listStr("drop", i, xs) },
		TOPWITHCONST: declBlockFunc("let"),
		TOPWITHASSIGN: declBlockFunc("let"),
		WITHCONST: declBlockFunc("let"),
		WITHASSIGN: declBlockFunc("let"),
		LOOP:	   declBlockFunc("loop"),
		CONST: func(identifier, expression) {
			str(identifier, " ", expression)
		},
		ASSIGN: func(args...) {
			vArgs List := vec(args)
			opPos      := vArgs->indexOf(":=")
			n          := vArgs->size()
			if  n % 2 != 1 || (n - 1) / 2 != opPos {
				throw(new IOException(
					"LHS and RHS of := do not  match" str blankJoin(vArgs)
				))
			} else {
				" " s.join (for i := lazy \`range`(opPos) {
					str(vArgs[i], " ", vArgs[opPos + 1 + i])
				})
			}
		},
		VECDESTRUCT: vecStr,
		DICTDESTRUCT: func{str('{', (" " s.join $*), "}")},
		DICTDESTRUCTELEM: func(destruct, label) {
			str(destruct, " ", label)
		},
		VARIADICDESTRUCT:  func{str("& ", $1)},
		SYMBOL: func(identifier){
			identifier
		} (pkg, identifier) {
			if !(symbolTable symbols.HasPackage pkg) {
				throw(new IOException(format(
					`package "%s" in %s.%s does not appear in imports %s`,
					pkg, pkg, identifier, symbols.Packages(symbolTable))))
			}
			str(pkg, "/", identifier)
		},
		BINARYOP: identity,
		MULOP: identity,
		ADDOP: identity,
		RELOP: identity,
		OPERATOR: identity,
		FUNCTIONDECL:	func(identifier, function) {
			defn := if isPublic(identifier) { "defn" } else { "defn-" }
			listStr(defn, identifier, function)
		},
		FUNCLIKEDECL:	func(funclike, identifier, function) {
			listStr(funclike, identifier, function)
		},
		FUNCTIONLIT:	func{listStr("fn", $1)},
		SHORTFUNCTIONLIT:  func(expr) {
			if first(expr) == '(' && last(expr) == ')' {
				"#" str expr
			}else{
				listStr("fn", "[]", expr)
		}
		},
		STRUCTSPEC: func(javaIdentifier, fields...) {
			symbolTable symbols.TypeCreated javaIdentifier
			listStr("defrecord",
				javaIdentifier,
				vecStr(...fields),
				if isEmpty(fields) {
					""
				} else {
					fs     := fields[0] s.split / +/
					fsOnly := func(s String){!s->startsWith("^")} filter fs
					str(
						"Object (toString [this] ",
						listStr("str", `"{"`, ` " " ` s.join fsOnly, `"}"`),
						")"
					)
				}
			)
		},
		FIELDS: blankJoin,
		INTERFACESPEC: func(args...){
			symbolTable symbols.TypeCreated first(args)
			listStr("defprotocol", ...args)
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
			symbolTable symbols.TypeCreated concrete
			listStr(...concat(list("extend-type", concrete, protocol), methodimpls))
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
		PERCENTNUM: func{"%" str $1},
		PERCENTVARADIC: constantFunc("%&"),
		FUNCTIONPARTS:	func{str("(",  ") (" s.join $*,	 ")")},
		FUNCTIONPART0:	func(expression) {
			"[] " str expression
		} (typ, expression) {
			str("^", typ, " [] ", expression)
		},
		VFUNCTIONPART0:	 func(variadic, expression) {
			str("[", variadic, "] ", expression)
		} (variadic, typ, expression) {
			str("^", typ, " [", variadic, "] ", expression)
		},
		FUNCTIONPARTN:	func(parameters, expression) {
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
		PARAMETERS:	blankJoin,
		VARIADIC:	func{"& " str $1},
		VECLIT:		vecStr,
		DICTLIT:	func{str(...$*)},
		DICTELEMENT:	func(key, value) {str(key, " ", value, " ")},
		SETLIT:		func{str("#{",	" " s.join $*,	"}")},
		STRUCTLIT:	func(typ, exprs...) {
			listStr(typ str ".", ...exprs)
		},
		LABEL:		func{":" str s.replace(s.lowerCase($1), /_/, "-")},
		ISLABEL:	func{str(":", s.replace(s.lowerCase($1), /_/, "-"), "?")},
		IDENTIFIER:	func(string) {
			s.replace(
				string,
				/\p{Ll}\p{Lu}/,
				func{str(first($1), "-", s.lowerCase(last($1)))}
			)
		},
		TYPEDIDENTIFIER: func(identifier, typ) {
			str(`^`, typ, " ", identifier)
		},
		TYPEDIDENTIFIERS: func(args...) {
			typ         := last(args)
			identifiers := butlast(args)
			decls       := for identifier := lazy identifiers {
					str(`^`, typ, " ", identifier)
			}
			blankJoin(...decls)
		},
		IMPORTED:	  func{"." s.join $*},
		DECIMALLIT:    identity,
		BIGINTLIT:     str,
		BIGFLOATLIT:   str,
		FLOATLIT:      identity, //str,
		HEXLIT:        func(s string){
				Integer::parseInt(s, 16)
		},
		REGEX:	func(regex string){
			str(
				`#"`,
				stripQuotes(regex)->replace(`\/`, `/`)->replace(`"`, `\"`),
				`"`
			)
		},
		INTERPRETEDSTRINGLIT: func(literal) {
			str(`"`, stripQuotes(literal), `"`)
		},
		CLOJUREESCAPE: identity,
		LITTLEUVALUE:  func(d1,d2,d3,d4){str(`\u`,d1,d2,d3,d4)},
		OCTALBYTEVALUE:	 func(d1,d2,d3){str(`\o`,d1,d2,d3)},
		UNICODECHAR:   func{`\` str $1},
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
		DASHIDENTIFIER: func{ "-" str $1},
		ISIDENTIFIER:	func(initial, identifier) {
			str( s.lowerCase(initial), identifier, "?")
		},
		EQUALS: constantFunc("="),
		AND:	constantFunc("and"),
		OR:	constantFunc("or"),
		MUTIDENTIFIER:	func(initial, identifier) {
			str( s.lowerCase(initial), identifier, "!")
		},
		ESCAPEDIDENTIFIER:  func{ stripQuotes($1) },
		UNARYEXPR: func(e) {
			e
		} (operator, expression){
			listStr(operator, expression)
		},
		NOTEQ:	     constantFunc("not="),
		BITAND:	     constantFunc("bit-and"),
		BITANDNOT:	     constantFunc("bit-and-not"),
		BITOR:	     constantFunc("bit-or"),
		BITXOR:	     constantFunc("bit-xor"),
		BITNOT:	     constantFunc("bit-not"),
		TAKE:	     constantFunc("<!!"),
		TAKEINGO:    constantFunc("<!"),
		SENDOP:      constantFunc(">!!"),
		SENDOPINGO:  constantFunc(">!"),
		SHIFTLEFT:   constantFunc("bit-shift-left"),
		SHIFTRIGHT:  constantFunc("bit-shift-right"),
		NOT:	     constantFunc("not"),
		MOD:	     constantFunc("mod"),
		DEREF:		 func{"@"   str $1},
		SYNTAXQUOTE:	 func{"`"   str $1},
		UNQUOTE:	 func{"~"   str $1},
		UNQUOTESPLICING: func{ "~@" str $1},
		JAVAFIELD:	func(expression, identifier) {
			listStr(".", expression, identifier)
		},
		JAVASTATIC:	 func{"/" s.join $*},
		TYPENAME:	 func(segments...){
			typ := "." s.join segments
			if !hasType(typ) {
				throw(new IOException(format(
					`type "%s" does not appear in type imports %s`,
					typ, symbols.Types(symbolTable))))
			}
			typ
		},
		UNDERSCOREJAVAIDENTIFIER: func{ "-" str $1},
		JAVAMETHODCALL: func(expression, identifier) {
			str("(. ", expression, " (", identifier, "))")
		} (expression, identifier, call) {
			str("(. ", expression, " (", identifier, " ", call, "))")
		},
		LONG: constantFunc("long"),
		DOUBLE: constantFunc("double"),
		STRING: constantFunc("String")
	}
}

func syncImports(isGoscript, isSync) {
	if isSync {
		[]
	} else {
		if isGoscript {
			[vecStr(
				"cljs.core.async", ":as", "async", ":refer",
				"[chan <! >! alt!]"
			)]
		} else {
			[vecStr(
				"clojure.core.async", ":as", "async", ":refer",
				"[chan go thread <! >! alt! <!! >!! alt!!]"
			)]
		}
	}
}

func macroSyncImports(isGoscript, isSync) {
	if isSync || !isGoscript {
		[]
	} else {
		[vecStr("cljs.core.async.macros", ":as", "async", ":refer", "[go]")]
	}
}

// Return an empty list if tail is empty, otherwise return listStr(head, ...tail)
func req(head, tail) {
	if isEmpty(tail) {
		[]
	} else {
		[listStr(head, ...tail)]
	}
}

func packageclauseFunc(symbolTable, path String, isGoscript, isSync) {
	[parent, name] := splitPath(path)
	if isGoscript {
		symbolTable symbols.PackageCreated "js"
	}
	func(imported, importDecls String) {
		fullImported     := parent str imported
		hasImports       := importDecls->contains(":require ")
		hasMacroImports  := importDecls->contains(":require-macros ")
		xtraImports      := if hasImports {
			[]
		} else {
			req(":require", syncImports(isGoscript, isSync))
		}
		xtraMacroImports := if hasMacroImports {
			[]
		} else {
			req(":require-macros", macroSyncImports(isGoscript, isSync))
		}
		imports          := concat([importDecls], xtraMacroImports, xtraImports)
		if imported != name {
			throw(new IOException(str(
				`Got package "`, imported, `" instead of expected "`,
				name, `" in "`, path, `"`
			)))
		}
		if isGoscript {
			listStr("ns", fullImported, ...imports)
		} else {
			str(
				listStr("ns", fullImported, "(:gen-class)", ...imports),
				" (set! *warn-on-reflection* true)"
			)
		}
	}
}

func importDeclFunc(isGoscript, isSync) {
	func() {
		""
	} (importSpecs...) {
		imports := importSpecs concat syncImports(isGoscript, isSync)
		listStr(":require", ...imports)
	}
}

func macroImportDeclFunc(isGoscript, isSync) {
	func() {
		""
	} (importSpecs...) {
		imports := importSpecs concat macroSyncImports(isGoscript, isSync)
		listStr(":require-macros", ...imports)
	}
}

func usesAsync(parsed) {
	func walk(vector) {
		if isEmpty(vector) {
			false
		} else {
			f := first(vector)
			if isVector(f) && usesAsync(f) {
				true
			} else {
				recur(rest(vector))
			}
		}
	}
	if kAsyncRules  isContains  first(parsed) {
		true
	} else {
		walk(rest(parsed))
	}
}

// Return the Clojure code generated from the given parse tree.
func Generate(path String, parsed, isSync) {
	symbolTable := symbols.New()
	isGoscript  := path->endsWith(".gos")
	isSync      := !usesAsync(parsed)
	codeGen     := codeGenerator(symbolTable, isGoscript) += {
		PACKAGECLAUSE:   packageclauseFunc(symbolTable, path, isGoscript, isSync),
		IMPORTDECL:      importDeclFunc(isGoscript, isSync) ,
		MACROIMPORTDECL: macroImportDeclFunc(isGoscript, isSync)
	}
	clj         := insta.transform(codeGen, parsed)
	symbols.CheckAllUsed(symbolTable)
	clj
}
