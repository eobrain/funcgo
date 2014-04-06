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

package  funcgo/parser
import (
        insta "instaparse/core"
)


var Parse = insta.parser(`
sourcefile = [ NL ] packageclause expressions _
 <_> =      <#'[ \t\x0B\f\r\n]*'> | comment+                                 (* optional whitespace *)
 <_nonNL> = <#'[ \t\x0B\f\r]*'>                                  (* optional non-newline whitespace *)
 <NL> = nl | comment+
   <nl> = <#'\s*[\n;]\s*'>                     (* whitespace with at least one newline or semicolon *)
   <comment> = <#'[;\s]*//[^\n]*\n\s*'>
 packageclause = <'package'> <__> imported NL importdecls
   __ =  #'[ \t\x0B\f\r\n]+' | comment+     (* whitespace *)
   importdecls = [importdecl NL] [macroimportdecl NL]
     importdecl = ( <'import'> _ <'('>  _ { ImportSpec _ } <')'> )
                 | ( <'import'>  _ ImportSpec )
     macroimportdecl = ( <'import'> _ <'macros'> _ <'('>  _ { ImportSpec _ } <')'> )
                     | ( <'import'> _ <'macros'> _ ImportSpec )
     <ImportSpec> = importspec
       importspec = ( Identifier _ )?  <'"'> imported <'"'>
         imported = Identifier { <'/'> Identifier }
 expressions = Expression | expressions NL Expression
   <Expression>  = precedence0 | Vars | shortvardecl | ifelseexpr | tryexpr | forrange |
                   forlazy | fortimes | withconst | block
     withconst = <'const'> _ ( const NL | <'('> _ { consts } _ <')'> _ ) ( expressions | <'{'> _ expressions _ <'}'> )
       consts = [ const { NL const } ]
         const = Destruct _ <'='> _ Expression
           <Destruct> = Identifier | typedidentifier | vecdestruct | dictdestruct
             typedidentifier = Identifier _ type
               type = JavaIdentifier ( <'.'>  JavaIdentifier )*
             vecdestruct = <'['> _ VecDestructElem _ { <','> _ VecDestructElem _  } <']'>
               <VecDestructElem> = Destruct | variadicdestruct | label
                 variadicdestruct = Destruct <'...'>
             dictdestruct = <'{'> dictdestructelem {  <','> _ dictdestructelem } <'}'>
               dictdestructelem = (Destruct|label) _ <':'> _ Expression
     precedence0 = precedence1 | ( precedence0 _nonNL symbol _nonNL precedence1 )
       symbol = (( Identifier <'.'> )? Identifier ) | javastatic
       precedence1 = precedence2 | ( precedence1 _ or _ precedence2 )
	 or = <'||'>
	 precedence2 = precedence3 | precedence2 _ and _ precedence3
	   and = <'&&'>
	   precedence3 = precedence4 | precedence3 _ relop  _ precedence4
             relop = equals | noteq | '<' | '<=' | '>='  (* TODO(eob)  | '>' *)
	       equals = <'=='>
               noteq  = <'!='>
	     precedence4 = precedence5 | precedence4 _ addop _ precedence5
	       addop = '+' | '-' | ( !or '|' ) | bitxor
                 bitxor = <'^'>
	       precedence5 = UnaryExpr | ( precedence5 _ mulop _ UnaryExpr )
	         mulop = '*' | (!comment '/') | '%' | shiftleft | shiftright | (!and '&') | '&^'
                   shiftleft = <'<<'>
                   shiftright = <'>>'>
	   javastatic = type _ <'::'> _ JavaIdentifier
	     <JavaIdentifier> = #'\b[\p{L}_][\p{L}_\p{Digit}]*\b' | underscorejavaidentifier
               underscorejavaidentifier = <'_'> JavaIdentifier
	   <Identifier> = !Keyword  (identifier | dashidentifier | isidentifier | mutidentifier |
			  escapedidentifier)
             Keyword = '\bconst\b' | '\bfor\b' | '\bnew\b' | '\bpackage\b' | '\brange\b' | '\bif\b'
	     identifier = #'\b[\p{L}_][\p{L}_\p{Digit}]*\b'
	     dashidentifier = <'_'> identifier
	     isidentifier = <'is'> #'\p{L}' identifier
	     mutidentifier = <'mutate'> #'\p{L}' identifier
	     escapedidentifier = <'\\'> #'\b[\p{L}_][\p{L}_\p{Digit}]*\b'
     shortvardecl = Identifier _ <':='> _ Expression
                  | Identifier _ ',' _ shortvardecl _ ',' _ Expression
     <Vars> = <'var'> _ ( <'('> _ vardecl ( NL vardecl )* _ <')'> | vardecl )
     vardecl = Identifier ( _ type )? _ <'='> _ Expression
     ifelseexpr = <'if'> _ Expression _ ( ( block _ <'else'> _ block ) |
                  ( _ <'{'> _ expressions _ <'}'> )   )
       block = <'{'> _ Expression { NL Expression } _ <'}'>
     forrange = <'for'> <__> Identifier _ <':='> _ <'range'> <_> Expression _
                <'{'> _ expressions _ <'}'>
     forlazy = <'for'> <__> Identifier _ <':='> _ <'lazy'> <_> Expression
               [ <__> <'if'> <__> Expression ] _ <'{'> _ expressions _ <'}'>
     fortimes = <'for'> <__> Identifier _ <':='> _ <'times'> <_> Expression _
                <'{'> _ expressions _ <'}'>
     tryexpr = <'try'> _ <'{'> _ expressions _ <'}'> _ catches ( _ finally )?
       catches = ( catch { _ catch } )?
         catch = <'catch'> _ Identifier _ Identifier _ <'{'> _ expressions _ <'}'>
       finally = <'finally'> _ <'{'> _ expressions _ <'}'>
     <UnaryExpr> = PrimaryExpr | javafield | ReaderMacro | unaryexpr
       unaryexpr = unary_op _ UnaryExpr
	 <unary_op> = '+' | '-' | '!' | not | (!and '&') | bitnot
	   bitnot = <'^'>
	   not    = <'!'>
       <ReaderMacro> = deref | syntaxquote | unquote | unquotesplicing
       deref           = <'*'>               _ UnaryExpr
       syntaxquote     = <'syntax'>     _ UnaryExpr
       unquote         = <'unquote'>         _ UnaryExpr
       unquotesplicing = <'unquotes'> _ UnaryExpr
       javafield  = UnaryExpr _ <'->'> _ JavaIdentifier
       <PrimaryExpr> = functioncall | javamethodcall | Operand | functiondecl | funclikedecl |  indexed
                                                                (*Conversion |
                                                                BuiltinCall |
                                                                PrimaryExpr Selector |
                                                                PrimaryExpr Slice |
                                                                PrimaryExpr TypeAssertion |*)
         indexed = PrimaryExpr _ <'['> _ Expression _ <']'>
         functioncall = PrimaryExpr Call
         javamethodcall = UnaryExpr _ <'->'> _ JavaIdentifier _ Call
           <Call> = <'('> _ ( ArgumentList _ )? <')'>
             <ArgumentList> = expressionlist                                         (* [ _ '...' ] *)
               expressionlist = Expression { _ <','> _ Expression }
         functiondecl = <'func'> _ Identifier _ Function
         funclikedecl = <'func'> _ <'<'> _ symbol _ <'>'> _ Identifier _ Function
           <Function> = FunctionPart | functionparts
             functionparts = FunctionPart _ FunctionPart { _ FunctionPart }
               <FunctionPart> = functionpart0 | functionpartn | vfunctionpart0 | vfunctionpartn
                 functionpart0 = <'('> _ <')'>  ( _ type )? _ <'{'> _ expressions _ <'}'>
		 vfunctionpart0 = <'('> _ variadic _ <')'> ( _ type )? _ <'{'> _ expressions _ <'}'>
		 functionpartn  = <'('> _ parameters _ <')'> ( _ type )? _ <'{'> _ expressions _ <'}'>
		 vfunctionpartn = <'('> _ parameters _  <','> _ variadic _ <')'> ( _ type )? _
                                  <'{'> _ expressions _ <'}'>
                   parameters = Destruct { <','> _ Destruct }
                   variadic = Identifier <'...'>
         <Operand> = Literal | OperandName | label | new  | ( <'('> Expression <')'> ) (*|MethodExpr*)
           label = #'\b\p{Lu}[\p{Lu}_0-9]*\b'
           <Literal> = BasicLit | veclit | dictlit | functionlit
             functionlit = <'func'> _ Function
             <BasicLit> = int_lit | bigintlit | string_lit | regex  | rune_lit | floatlit | bigfloatlit (*| imaginary_lit *)
               floatlit = decimals '.' [ decimals ] [ exponent ]
                        | decimals exponent
                        | '.' decimals [ exponent ]
                 decimals  = #'[0-9]+'
                 exponent  = ( 'e' | 'E' ) [ '+' | '-' ] decimals
               bigfloatlit = (floatlit | int_lit) 'M'
               <int_lit> = decimallit    (*| octal_lit | hex_lit .*)
		 decimallit = #'[1-9][0-9]*' | #'[0-9]'
               bigintlit = int_lit 'N'
	       regex = <'/'> #'[^/\n]+'<'/'>                               (* TODO: handle / escape *)
	       <string_lit> = rawstringlit | interpretedstringlit | clojureescape
                 rawstringlit = <#'\x60'> #'[^\x60]*' <#'\x60'>     (* \x60 is back quote character *)
                 interpretedstringlit = <#'\"'> #'[^\"]*' <#'\"'>     (* TODO: handle string escape *)
                 clojureescape = <'\\'> <#'\x60'> #'[^\x60]*' <#'\x60'>       (* \x60 is back quote *)
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
           new = <'new'> <__> type
           <OperandName> = symbol | NonAlphaSymbol                           (*| QualifiedIdent*)
             <NonAlphaSymbol> = '=>' | '->>' | relop | addop | mulop | unary_op
`)
