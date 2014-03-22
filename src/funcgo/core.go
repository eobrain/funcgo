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
        pprint clojure.pprint
        failure instaparse.failure
        string clojure.string
)

funcgoParser := insta.parser(`
sourcefile = [ NL ] packageclause _ expressions _
 <_> = <#'[ \t\x0B\f\r\n]*'> | comment+  (* optional whitespace *)
 <NL> = nl | comment+
   <nl> = <#'\s*[\n;]\s*'>                     (* whitespace with at least one newline or semicolon *)
   <comment> = <#'[;\s]*//[^\n]*\n\s*'>
 packageclause = <'package'> <__> dotted NL importdecl
   __ =  #'[ \t\x0B\f\r\n]+' | comment+     (* whitespace *)
   importdecl = <'import'> _ <'('>  _ { importspec _ } <')'>
     importspec = Identifier _ dotted
       dotted = Identifier { <'.'> Identifier }
 expressions = Expression { NL Expression }
   <Expression>  = precedence0 | withconst | shortvardecl | ifelseexpr | tryexpr | forrange |
                   forlazy | fortimes
     precedence0 = precedence1 | ( precedence0 _ symbol _ precedence1 )
       symbol = (( Identifier <'.'> )? !Keyword Identifier ) | javastatic | '=>' | '->>' | '->'
         Keyword = ( '\bfor\b' | '\brange\b' )
       precedence1 = precedence2 | ( precedence1 _ or _ precedence2 )
	 or = <'||'>
	 precedence2 = precedence3 | precedence2 _ and _ precedence3
	   and = <'&&'>
	   precedence3 = precedence4 | precedence3 _ relop  _ precedence4
             relop = equals | noteq | '<' | '<=' | '>='  (* TODO(eob)  | '>' *)
	       equals = <'=='>
               noteq  = <'!='>
	     precedence4 = precedence5 | precedence4 _ addop _ precedence5
	       addop = '+' | '-' | '|' | bitxor
                 bitxor = <'^'>
	       precedence5 = UnaryExpr | ( precedence5 _ mulop _ UnaryExpr )
	         mulop = '*' | (!comment '/') | '%' | '<<' | '>>' | '&' | '&^'
	   javastatic = JavaIdentifier _ <'::'> _ JavaIdentifier
	     <JavaIdentifier> = #'\b[\p{L}_][\p{L}_\p{Digit}]*\b'
	   <Identifier> = identifier | dashidentifier | isidentifier | mutidentifier |
			  escapedidentifier
	     identifier = #'\b[\p{L}_][\p{L}_\p{Digit}]*\b'
	     dashidentifier = <'_'> identifier
	     isidentifier = <'is'> #'\p{L}' identifier
	     mutidentifier = <'mutate'> #'\p{L}' identifier
	     escapedidentifier = <'\\'> #'\b[\p{L}_][\p{L}_\p{Digit}]*\b'
     shortvardecl   = Identifier _ <':='> _ Expression
     ifelseexpr = <'if'> _ Expression _ ( ( block _ <'else'> _ block ) |
                  ( _ <'{'> _ expressions _ <'}'> )   )
       block = <'{'> _ Expression { NL Expression } _ <'}'>
     forrange = <'for'> <__> Identifier _ <':='> _ <'range'> <_> Expression _
                <'{'> _ expressions _ <'}'>
     forlazy = <'for'> <__> Identifier _ <':='> _ <'lazy'> <_> Expression
               [ <__> <'if'> <__> Expression ] _ <'{'> _ expressions _ <'}'>
     fortimes = <'for'> <__> Identifier _ <':='> _ <'times'> <_> Expression _
                <'{'> _ expressions _ <'}'>
     tryexpr = <'try'> _ <'{'> _ expressions _ <'}'> _ catches _ finally?
       catches = ( catch { _ catch } )?
         catch = <'catch'> _ Identifier _ Identifier _ <'{'> _ expressions _ <'}'>
       finally = <'finally'> _ <'{'> _ expressions _ <'}'>
     <UnaryExpr> = PrimaryExpr | javafield | unaryexpr | deref
       unaryexpr = unary_op _ UnaryExpr
	 <unary_op> = '+' | '-' | '!' | not | '*' | '&' | bitnot
	   bitnot = <'^'>
	   not    = <'!'>
       deref = <'<-'> _ UnaryExpr
       javafield  = Expression _ <'->'> _ JavaIdentifier
       <PrimaryExpr> = functioncall | javamethodcall | Operand | functiondecl |  indexed
                                                                (*Conversion |
                                                                BuiltinCall |
                                                                PrimaryExpr Selector |
                                                                PrimaryExpr Slice |
                                                                PrimaryExpr TypeAssertion |*)
         indexed = PrimaryExpr _ <'['> _ Expression _ <']'>
         functioncall = PrimaryExpr Call
         javamethodcall = Expression _ <'->'> _ JavaIdentifier _ Call
           <Call> = <'('> _ ( ArgumentList _ )? <')'>
             <ArgumentList> = expressionlist                                         (* [ _ '...' ] *)
               expressionlist = Expression { _ <','> _ Expression }
         functiondecl = <'func'> _ Identifier _ Function
           <Function> = FunctionPart | functionparts
             functionparts = FunctionPart _ FunctionPart { _ FunctionPart }
               <FunctionPart> = functionpart0 | functionpartn | vfunctionpart0 | vfunctionpartn
                 functionpart0 = <'('> _ <')'> _ <'{'> _ Expression _ <'}'>
		 vfunctionpart0 = <'('> _ varadic _ <')'> _ <'{'> _ Expression _ <'}'>
		 functionpartn  = <'('> _ parameters _ <')'> _ <'{'> _ Expression _ <'}'>
		 vfunctionpartn = <'('> _ parameters _  <','> _ varadic _ <')'> _
                                  <'{'> _ Expression _ <'}'>
                   parameters = Identifier { <','> _ Identifier }
                   varadic = <'&'> Identifier
         <Operand> = Literal | OperandName | label | new  | ( <'('> Expression <')'> ) (*|MethodExpr*)
           label = #'\b\p{Lu}[\p{Lu}_0-9]*\b'
           <Literal> = BasicLit | veclit | dictlit | functionlit
             functionlit = <'func'> _ Function
             <BasicLit> = int_lit | string_lit | regex  | rune_lit | floatlit (*| imaginary_lit *)
               floatlit = decimals '.' [ decimals ] [ exponent ]
                        | decimals exponent
                        | '.' decimals [ exponent ]
                 decimals  = #'[0-9]+'
                 exponent  = ( 'e' | 'E' ) [ '+' | '-' ] decimals
               <int_lit> = decimallit    (*| octal_lit | hex_lit .*)
		 decimallit = #'[1-9][0-9]*' | #'[0-9]'
	       regex = <'/'> #'[^/]+'<'/'>                                 (* TODO: handle / escape *)
	       <string_lit> = rawstringlit   | interpretedstringlit
                 rawstringlit = <#'\x60'> #'[^\x60]*' <#'\x60'>     (* \x60 is back quote character *)
                 interpretedstringlit = <#'\"'> #'[^\"]*' <#'\"'>     (* TODO: handle string escape *)
	       <rune_lit> = <'\''> ( unicode_value | byte_value ) <'\''> 
		 <unicode_value> = unicodechar | littleuvalue | escaped_char
                   unicodechar = #'[^\n ]'
                   <escaped_char> = newlinechar | spacechar | backspacechar | returnchar | tabchar |
                                    backslashchar | squotechar| dquotechar
		     newlinechar   = <'\\n'>
		     spacechar     = <' '>
		     backspacechar = <'\\b'>
		     returnchar    = <'\\r'>
		     tabchar       = <'\\t'>
		     backslashchar = <'\\\\'>
		     squotechar    = <'\\\''>
		     dquotechar    = <'\\"'>
		 <byte_value> = octalbytevalue                                  (* | hex_byte_value *)
                   octalbytevalue = <'\\'> octaldigit octaldigit octaldigit
                     octaldigit = #'[0-7]'
                   littleuvalue = <'\\u'> hexdigit hexdigit hexdigit hexdigit
                     hexdigit = #'[0-9a-fA-F]'
	     veclit = <'['> _ (( Expression { _ <','> _ Expression _ } )? )? <']'>
	     dictlit = '{' _ ( dictelement _ { <','> _ dictelement } )? _ '}'
               dictelement = Expression _ <':'> _ Expression
           new = <'new'> <__> symbol
           <OperandName> = symbol                                                 (*| QualifiedIdent*)
     withconst = <'const'> _ <'('> _ { consts } _ <')'> _ expressions
       consts = [ const { NL const } ]
         const = Identifier _ <'='> _ Expression 
`)

func infix(expression) {
	expression
} (left, operator, right) {
	str("(", operator, " ", left, " ", right, ")")
}

func listStr(&item) {
	str("(", string.join(" ", item), ")")
}

func funcgoParse(fgo) {
        const(
                parsed = funcgoParser(string.replace(fgo, /\t/, "        "))
        )
        if insta.isFailure(parsed) {
                failure.pprintFailure(parsed)
                throw(new Exception(`"SYNTAX ERROR"`))
        } else {
	    // pprint.pprint(parsed)
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
                        PRECEDENCE0: infix,
                        PRECEDENCE1: infix,
                        PRECEDENCE2: infix,
                        PRECEDENCE3: infix,
                        PRECEDENCE4: infix,
                        PRECEDENCE5: infix,
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
                                string.join(" ", catches)
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
                                string.join(" ", expr0 cons exprRest)
                        },
                        EXPRESSIONS: func(expr0, &exprRest){
                                string.join(" ", expr0 cons exprRest)
                        },
                        CONSTS:  func(&consts) {
                                "\n" string.join consts
                        },
                        BLOCK: func (expr){
                                expr
                        } (expr0, &exprRest) {
                                str("(do ",
                                        string.join(" ", expr0 cons exprRest),
                                        ")")
                        },
                        INDEXED: func(xs, i){ str("(", xs, " ", i, ")") },
                        WITHCONST: func(&xs){
                                const(
                                        consts = butlast(xs)
                                        expressions = last(xs)
                                )
                                str("(let [",
                                        " " string.join consts,
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
			BINARYOP: identity,
			MULOP: identity,
			ADDOP: identity,
			RELOP: identity,
                        FUNCTIONDECL:   func(identifier, function) {
                                str("(defn ", identifier, " ", function, ")")
                        },
                        FUNCTIONLIT:    func(function) {str("(fn ", function, ")")},
                        FUNCTIONPARTS:  func(&functionpart) {
                                str("(",
                                        ") (" string.join functionpart,
                                        ")")
                        },
                        FUNCTIONPART0:  func(expression) {
                                "[] " str expression
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
                                string.join(" ", arg0 cons argsRest)
                        },
                        VARADIC:        func(parameter) {str("& ", parameter)},
                        VECLIT:         func() {
                                "[]"
                        } (&expressions) {
                                str(
                                        "[",
                                        " " string.join expressions,
                                        "]"
                                )
                        },
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
                                string.join(".", idf0 cons idfRest)
                        },
                        DECIMALLIT:    identity,
			FLOATLIT:      str,
			DECIMALS:      identity,
			EXPONENT:      str,
                        REGEX:          func(s){str(`#"`, s, `"`)},
                        INTERPRETEDSTRINGLIT: func(s){str(`"`, s, `"`)},
                        LITTLEUVALUE:  func(d1,d2,d3,d4){str(`\u`,d1,d2,d3,d4)},
                        OCTALBYTEVALUE:  func(d1,d2,d3){str(`\o`,d1,d2,d3)},
                        UNICODECHAR:   func(s){`\` str s},
                        NEWLINECHAR:   func(){`\newline`},
                        SPACECHAR:     func(){`\space`},
                        BACKSPACECHAR: func(){`\backspace`},
                        RETURNCHAR:    func(){`\return`},
                        TABCHAR:       func(){`\tab`},
                        BACKSLASHCHAR: func(){`\\`},
                        SQUOTECHAR:    func(){`\'`},
                        DQUOTECHAR:    func(){`\"`},
                        HEXDIGIT:      identity,
                        OCTALDIGIT:    identity,
                        RAWSTRINGLIT: func(s){str(`"`, string.escape(s, charEscapeString), `"`)},
                        DASHIDENTIFIER: func(s){ "-" str s },
                        ISIDENTIFIER:   func(initial, identifier) {
                                str( string.lowerCase(initial), identifier, "?")
                        },
                        EQUALS: func() { "=" },
                        AND: func() { "and" },
                        OR: func() { "or" },
                        MUTIDENTIFIER:  func(initial, identifier) {
                                str( string.lowerCase(initial), identifier, "!")
                        },
                        ESCAPEDIDENTIFIER:  func(identifier) { identifier },
                        NOTEQ:   func() { "not=" },
			UNARYEXPR: func(operator, expression) { listStr(operator, expression) },
			BITXOR: func(){ "bit-xor" },
			BITNOT: func(){ "bit-not" },
			NOT: func(){ "not" },
			DEREF: func(expression) { str("@", expression) },
                        JAVAFIELD:      func(expression, identifier) {
                                str("(. ", expression, " ", identifier, ")")
                        },
                        JAVASTATIC:      func(clazz, identifier) {
                                str(clazz, "/", identifier)
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

