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
        insta instaparse.core
	failure instaparse.failure
        string clojure.string
        pprint clojure.pprint
)

funcgoParser := insta.parser(`
sourcefile     = [ NL ] packageclause _ expressions _
packageclause  = <'package'> <__> dotted NL importdecl
expressions    = Expression { NL Expression }
importdecl     = <'import'> _ <'('>  _ { importspec _ } <')'>
importspec     = Identifier _ dotted
<Expression>   = UnaryExpr | withconst | shortvardecl | ifelseexpr | tryexpr | forrange | forlazy | fortimes (* | Expression binary_op UnaryExpr *)
ifelseexpr     = <'if'> _ Expression _ ( ( block _ <'else'> _ block ) | ( _ <'{'> _ expressions _ <'}'> )   )
forrange       = <'for'> <__> Identifier _ <':='> _ <'range'> <_> Expression _ <'{'> _ expressions _ <'}'>
forlazy        = <'for'> <__> Identifier _ <':='> _ <'lazy'> <_> Expression [ <__> <'if'> <__> Expression ] _ <'{'> _ expressions _ <'}'>
fortimes       = <'for'> <__> Identifier _ <':='> _ <'times'> <_> Expression _ <'{'> _ expressions _ <'}'>
tryexpr        = <'try'> _ <'{'> _ expressions _ <'}'> _ catches _ finally?
catches        = ( catch { _ catch } )?
catch          = <'catch'> _ Identifier _ Identifier _ <'{'> _ expressions _ <'}'>
finally        = <'finally'> _ <'{'> _ expressions _ <'}'>
block          = <'{'> _ Expression { NL Expression } _ <'}'>
<UnaryExpr>    = PrimaryExpr | javafield  (* | unary_op UnaryExpr *)
<PrimaryExpr>  = functioncall | javamethodcall | Operand | functiondecl
                                                                         (*Conversion |
                                                                         BuiltinCall |
                                                                         PrimaryExpr Selector |
                                                                         PrimaryExpr Index |
                                                                         PrimaryExpr Slice |
                                                                         PrimaryExpr TypeAssertion |*)
withconst      = <'const'> _ <'('> _ { consts } _ <')'> _ expressions
consts         = [ const { NL const } ]
const          = Identifier _ <'='> _ Expression 
functioncall   = PrimaryExpr Call
javamethodcall = Expression _ <'->'> _ JavaIdentifier _ Call
<Call>         = <'('> _ ( ArgumentList _ )? <')'>
<ArgumentList> = expressionlist                                                      (* [ _ '...' ] *)
expressionlist = Expression { _ <','> _ Expression }
<Operand>      = Literal | OperandName | label | new            (*| MethodExpr | '(' Expression ')' *)
new            = <'new'> <__> symbol
<OperandName>  = symbol                                             (*| QualifiedIdent*)
<Literal>      = BasicLit | dictlit | functionlit
<BasicLit>     = int_lit | string_lit | regex              (*| float_lit | imaginary_lit | rune_lit *)
shortvardecl   = Identifier _ <':='> _ Expression
functiondecl   = <'func'> _ Identifier _ Function
functionlit    = <'func'> _ Function
<Function>     = FunctionPart | functionparts
functionparts  = FunctionPart _ FunctionPart { _ FunctionPart }
<FunctionPart> = functionpart0 | functionpartn | vfunctionpart0 | vfunctionpartn
functionpart0  = <'('> _ <')'> _ <'{'> _ Expression _ <'}'>
vfunctionpart0 = <'('> _ varadic _ <')'> _ <'{'> _ Expression _ <'}'>
functionpartn  = <'('> _ parameters _ <')'> _ <'{'> _ Expression _ <'}'>
vfunctionpartn = <'('> _ parameters _  <','> _ varadic _ <')'> _ <'{'> _ Expression _ <'}'>
parameters     = Identifier { <','> _ Identifier }
varadic        = <'&'> Identifier
dictlit        = '{' _ ( dictelement _ { <','> _ dictelement } )? _ '}'
dictelement    = Expression _ <':'> _ Expression
<int_lit>      = decimallit    (*| octal_lit | hex_lit .*)
decimallit    = #'[1-9][0-9]*'
regex          = <'/'> #'[^/]+'<'/'>   (* TODO: handle / escape *)
<string_lit>   = rawstringlit   | interpretedstringlit
rawstringlit = <#'\x60'> #'[^\x60]*' <#'\x60'>      (* \x60 is back quote character *)
interpretedstringlit = <#'\"'> #'[^\"]*' <#'\"'>      (* TODO: handle string escape *)
dotted         = Identifier { <'.'> Identifier }
symbol         = ( Identifier <'.'> )? !Keyword Identifier
javafield      = Expression _ <'->'> _ JavaIdentifier
Keyword        = ( 'for' | 'range' )
<Identifier>     = identifier | dashidentifier | isidentifier | mutidentifier
identifier     = #'[\p{L}_][\p{L}_\p{Digit}]*'
<JavaIdentifier> = #'[\p{L}_][\p{L}_\p{Digit}]*'
dashidentifier = <'_'> identifier
isidentifier   = <'is'> #'\p{L}' identifier
mutidentifier  = <'mutate'> #'\p{L}' identifier
label          = #'\p{Lu}[\p{Lu}_0-9]*'
letter         = unicode_letter | '_'
unicode_letter = #'\p{L}'
unicode_digit  = #'\p{Digit}'
<_>            = <#'[ \t\x0B\f\r\n]*'> | comment+  (* optional whitespace *)
__             =  #'[ \t\x0B\f\r\n]+' | comment+     (* whitespace *)
<NL>           = nl | comment+
<nl>           = <#'\s*[\n;]\s*'>       (* whitespace with at least one newline or semicolon *)
<comment>      = <#'[;\s]*//[^\n]*\n\s*'>
`)

func funcgoParse(fgo) {
        const(
                parsed = funcgoParser(fgo)
        )
        if insta.isFailure(parsed) {
                failure.pprintFailure(parsed)
                throw(new Exception(`"SYNTAX ERROR"`))
        } else {
            insta.transform(
                {
                        SOURCEFILE:     func(header, body) {str(header, body, "\n")},
                        PACKAGECLAUSE:  func(dotted, importDecl) {
                                str("(ns ", dotted, " (:gen-class)", importDecl, ")\n\n")
                        },
                        IMPORTDECL:     func(&importSpecs) {apply(str, importSpecs)},
                        IMPORTSPEC:     func(identifier, dotted) {
                                str("\n  (:require [", dotted, " :as ", identifier, "])")
                        },
                        IFELSEEXPR: func(condition, exprs) {
                                str("(when ", condition, " ", exprs, ")")
                        } (condition, block1, block2) {
                                str("(if ", condition, " ", block1, " ", block2, ")")
                        },
                        FORRANGE: func(identifier, seq, expressions) {
                                str("(doseq [", identifier, " ", seq, "] ", expressions, ")")
                        },
                        FORLAZY: func(identifier, seq, expressions) {
                                str("(for [", identifier, " ", seq, "] ", expressions, ")")
                        } (identifier, seq, condition, expressions) {
                                str("(for [", identifier, " ", seq, "] :when ", condition, " ", expressions, ")")
                        },
                        FORTIMES: func(identifier, count, expressions) {
                                str("(dotimes [", identifier, " ", count, "] ", expressions, ")")
                        },
                        TRYEXPR: func(expressions, catches) {
                                str("(try ", expressions, " ", catches, ")")
                        } (expressions, catches, finally) {
                                str("(try ", expressions, " ", catches, " ", finally, ")" )
                        },
                        CATCHES: func(&catches){
                                reduce(
                                        func(acc,catch) {str(acc, " ", catch)},
                                        catches)
                        },
                        CATCH: func(typ, exception, expressions) {
                                str("(catch ", typ, " ", exception, " ", expressions,")")
                        },
                        FINALLY: func(expressions) {
                                str("(finally ", expressions, ")")
                        },
                        NEW: func(symbol){str(symbol, ".")},
                        SHORTVARDECL:   func(identifier, expression) {
                                str("(def ", identifier, " ", expression, ")")
                        },
                        FUNCTIONCALL:    func(function) {
                                str("(", function, ")")
                        } (function, call) {
                                str("(", function, " ", call, ")")
                        },
                        EXPRESSIONLIST: func(expr0, &exprRest){
                                reduce(
                                        func(acc, expr) {str(acc, " ", expr)},
                                        expr0,
                                        exprRest)
                        },
                        EXPRESSIONS: func(expr0, &exprRest){
                                reduce(
                                        func(acc, expr) {str(acc, " ", expr)},
                                        expr0,
                                        exprRest)
                        },
                        CONSTS:  func(&consts) {
                          reduce(
                                  func(acc,konst) {str(acc, "\n", konst)},
                                  consts)
                        },
                        BLOCK: func (expr){
                                expr
                        } (expr0, &exprRest) {
                                str("(do ",
                                  reduce(
                                          func(acc, expr) {str(acc, " ", expr)},
                                          expr0,
                                          exprRest),
                                        ")")
                        },
                        WITHCONST: func(&xs){
                                const(
                                        consts = butlast(xs)
                                        expressions = last(xs)
                                )
                                str("(let [",
                                        reduce(func(acc,konst) {str(acc, " ", konst)}, consts),
                                        "] ",
                                        expressions,
                                        ")")
                        },
                        CONST: func(identifier, expression) {str(identifier, " ", expression)},
                        SYMBOL: func(identifier){
                                identifier
                        } (pkg, identifier) {
                                str(pkg, "/", identifier)
                        },
                        FUNCTIONDECL:   func(identifier, function) {
                                str("(defn ", identifier, " ", function, ")")
                        },
                        FUNCTIONLIT:    func(function) {str("(fn ", function, ")")},
                        FUNCTIONPARTS:  func(&functionpart) {
                                str("(",
                                        reduce(
                                                func(acc, fp) {str(acc, ") (", fp)},
                                                functionpart),
                                        ")")
                        },
                        FUNCTIONPART0:  func(expression) {
                                str("[] ", expression)
                        },
                        VFUNCTIONPART0:  func(varadic, expression) {
                                str("[", varadic, "] ", expression)
                        },
                        FUNCTIONPARTN:  func(parameters, expression) {
                                str("[", parameters, "] ", expression)
                        },
                        VFUNCTIONPARTN: func(parameters, varadic, expression) {
                        str("[", parameters, " ", varadic, "] ", expression)
                        },
                        PARAMETERS:     func(arg0, &argsRest) {
                                reduce(
                                        func(acc, arg) {str(acc, " ", arg)},
                                        arg0,
                                        argsRest)
                        },
                        VARADIC:        func(parameter) {str("& ", parameter)},
                        DICTLIT:        func(&dictElems) {apply(str, dictElems)},
                        DICTELEMENT:    func(key, value) {str(key, " ", value, " ")},
                        LABEL:          func(s) {
				str(":", string.replace(string.lowerCase(s), /_/, "-"))
			},
                        IDENTIFIER:     func(s) {
                                string.replace(s, /\p{Ll}\p{Lu}/, func(s){
                                        str(first(s), "-", string.lowerCase(last(s)))
                                })
                        },
                        DOTTED:         func(idf0, &idfRest){
                                reduce(
                                        func(acc, idf) {str(acc, ".", idf)},
                                        idf0,
                                        idfRest)
                        },
                        DECIMALLIT:    func(s){s},
                        REGEX:          func(s){str(`#"`,s,`"`)},
                        INTERPRETEDSTRINGLIT: func(s){str(`"`, s, `"`)},
                        RAWSTRINGLIT: func(s){str(`"`, string.escape(s, charEscapeString), `"`)},
			DASHIDENTIFIER: func(s){ str("-", s) },
			ISIDENTIFIER:   func(initial, identifier) {
				str( string.lowerCase(initial), identifier, "?")
			},
			MUTIDENTIFIER:  func(initial, identifier) {
				str( string.lowerCase(initial), identifier, "!")
			},
			JAVAFIELD:      func(expression, identifier) {
				str("(. ", expression, " ", identifier, ")")
			},
			JAVAMETHODCALL: func(expression, identifier) {
				str("(. ", expression, " (", identifier, "))")
			} (expression, identifier, call) {
				str("(. ", expression, " (", identifier, " ", call, "))")
			}
                },
                parsed
            )
        }
}

// Convert funcgo to clojure
func _main(&args) {
  try {
	  if not(seq(args)) {
		  println("Compiling all go files")
	  }else{
		  const(
			  clj = funcgoParse(slurp(first(args)))
		  )
		  for expr := range readString( str("[", clj, "]")) {
			  pprint.pprint(expr)
			  println()
		  }
	  }
  } catch Exception e {
          println("\n", e->getMessage())
  }
}