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

package  funcgo.parser
import (
        insta instaparse.core
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
	 <unary_op> = '+' | '-' | '!' | not | '&' | bitnot
	   bitnot = <'^'>
	   not    = <'!'>
       deref = <'*'> _ UnaryExpr
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
