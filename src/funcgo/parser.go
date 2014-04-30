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

package parser

import (
	insta "instaparse/core"
)

var Parse = insta.parser(`
sourcefile = NL? packageclause expressions _
 <_> =      <#'[ \t\x0B\f\r\n]*'> | comment+                                 (* optional whitespace *)
 <_nonNL> = <#'[ \t\x0B\f\r]*'>                                  (* optional non-newline whitespace *)
 <NL> = nl | comment+
   <nl> = <#'\s*[\n;]\s*'>                     (* whitespace with at least one newline or semicolon *)
   <comment> = <#'[;\s]*//[^\n]*\n\s*'>
 packageclause = <'package'> <__> imported NL importdecls
   __ =  #'[ \t\x0B\f\r\n]+' | comment+     (* whitespace *)
   importdecls = (AnyImportDecl NL)*
     <AnyImportDecl> = importdecl | macroimportdecl | externimportdecl | typeimportdecl
     importdecl = <'import'> _ <'('>  _ {ImportSpec _} <')'>
                | <'import'>  _ ImportSpec
     macroimportdecl = <'import'> _ <'macros'> _ <'('>  _ {ImportSpec _} <')'>
                     | <'import'> _ <'macros'> _ ImportSpec
     <externimportdecl> = <'import'> _ <'extern'> _ <'('>  _ {externimportspec _} <')'>
                     | <'import'> _ <'extern'> _ externimportspec
       externimportspec = identifier
     typeimportdecl = <'import'> _ <'type'> _ <'('>  _ {typeimportspec _} <')'>
                     | <'import'> _ <'type'> _ typeimportspec
       typeimportspec = typepackageimportspec <'.'> _ (
                                          JavaIdentifier
                                        | <'{'> _ JavaIdentifier _ (<','> _ JavaIdentifier)* _ <'}'>
                         )
         typepackageimportspec = JavaIdentifier {<'.'>  JavaIdentifier} 
     <ImportSpec> = importspec
       importspec = ( Identifier _ )?  <'"'> imported <'"'>
         imported = Identifier {<'/'> Identifier}
 expressions = expr | expressions NL expr
   <expr>  = precedence0 | Vars | shortvardecl | ifelseexpr | letifelseexpr | tryexpr | forrange |
                   forlazy | fortimes | Blocky | ExprSwitchStmt | assoc | associn
     assoc = expr _ <'+='> _ <'{'> _ associtem _ { <','> _ associtem _ } <'}'>
     dissoc = expr _ <'-='> _ <'{'> _ associtem _ { <','> _ associtem _ } <'}'>
       associtem = expr _ <':'> _ expr 
     associn = expr _ <'+='> _ <'{'> _ associnpath _ <':'> _ expr _ <'}'>
       associnpath = expr _ expr {_ expr}
     <ExprSwitchStmt> = boolswitch | constswitch
       boolswitch = <'switch'> _ <'{'>  _ boolcaseclause { NL boolcaseclause } _ <'}'>
       constswitch = <'switch'> _ expr _ <'{'> _ constcaseclause { NL constcaseclause } _ <'}'>
	 boolcaseclause = boolswitchcase _ <':'> _ expressions
	 constcaseclause = constswitchcase _ <':'> _ expressions
	   boolswitchcase = <'case'> _ expressionlist | <'default'>
	   constswitchcase = <'case'> _ constantlist | <'default'>
	     constantlist = expr {_ <','> _ expr}
	       <Constant> = label | BasicLit | veclit | dictlit | setlit
     <Blocky> = block | withconst | loop
       loop = <'loop'> _  <'('> _ {consts} _ <')'> _ ImpliedDo
       <ImpliedDo> =  <'{'> _ expressions _ <'}'> | withconst
       block = <'{'> _ expr {NL expr} _ <'}'>
       withconst = <'{'> _ <'const'> _ ( const NL | <'('> _ consts _ <')'> )  _ expressions _ <'}'>
         consts = ( const {NL const} )?
           const = Destruct _ <'='> _ expr
	     <Destruct> = Identifier | typedidentifier | vecdestruct | dictdestruct
	       typedidentifier = Identifier _ typename
		 typename = JavaIdentifier {<'.'>  JavaIdentifier} | primitivetype
                   <primitivetype> = long | double
                     long = <'int'>
                     double = <'float'>
	       vecdestruct = <'['> _ VecDestructElem _ {<','> _ VecDestructElem _ } <']'>
		 <VecDestructElem> = Destruct | variadicdestruct | label
		   variadicdestruct = Destruct <'...'>
	       dictdestruct = <'{'> dictdestructelem { _ <','> _ dictdestructelem} <'}'>
		 dictdestructelem = Destruct _ <':'> _ expr
     precedence0 = precedence1
                 | precedence0 _nonNL symbol _nonNL precedence1
       symbol = ( Identifier <'.'> )? Identifier
              | javastatic
       precedence1 = precedence2
                   | precedence1 _ or _ precedence2
	 or = <'||'>
	 precedence2 = precedence3
                     | precedence2 _ and _ precedence3
	   and = <'&&'>
	   precedence3 = precedence4
                       | precedence3 _ relop  _ precedence4
             relop = equals | noteq | '<' | '<=' | '>='               (* TODO(eob)  | '>' *)
	       equals = <'=='>
               noteq  = <'!='>
	     precedence4 = precedence5
                         | precedence4 _ addop _ precedence5
	       addop = '+' | '-' | ( !or '|' ) | bitxor
                 bitxor = <'^'>
	       precedence5 = UnaryExpr
                           | precedence5 _ mulop _ UnaryExpr
	         mulop = '*' | (!comment '/') | mod | shiftleft | shiftright | mod | (!and '&') | '&^'
                   shiftleft = <'<<'>
                   shiftright = <'>>'>
                   mod = <'%'>
	   javastatic = typename _ <'::'> _ JavaIdentifier
	     <JavaIdentifier> = #'\b[\p{L}_][\p{L}_\p{Digit}]*\b'
                              | underscorejavaidentifier
               underscorejavaidentifier = <'_'> JavaIdentifier
	   <Identifier> = !Keyword  (identifier | dashidentifier | isidentifier | mutidentifier |
			  escapedidentifier)
             Keyword = '\bconst\b' | '\bfor\b' | '\bnew\b' | '\bpackage\b' | '\brange\b' | '\bif\b'
	     identifier = #'\b[\p{L}_][\p{L}_\p{Digit}]*\b'
	     dashidentifier = <'_'> identifier
	     isidentifier = <'is'> #'\p{L}' identifier
	     mutidentifier = <'mutate'> #'\p{L}' identifier
	     escapedidentifier = <'\\'> #'\b[\p{L}_][\p{L}_\p{Digit}]*\b'
     shortvardecl = Identifier _ <':='> _ expr
               (*   | Identifier _ ',' _ shortvardecl _ ',' _ expr *)
     <Vars> = <'var'> _ ( <'('> _ vardecl {NL vardecl} _ <')'> | vardecl )
     vardecl = Identifier ( _ typename )? _ <'='> _ expr
     ifelseexpr = <'if'> _ expr _ Blocky ( _ <'else'> _ Blocky )?
     letifelseexpr = <'if'> _ Destruct _ <':='> _ expr _ <';'>_ expr _ Blocky ( _ <'else'> _ Blocky )?
     forrange = <'for'> <__> Destruct _ <':='> _ <'range'> <_> expr _  Blocky
     forlazy = <'for'> <__> Destruct _ <':='> _ <'lazy'> <_> expr
               ( <__> <'if'> <__> expr )? _ Blocky
     fortimes = <'for'> <__> Identifier _ <':='> _ <'times'> <_> expr _ Blocky
     tryexpr = <'try'> _ ImpliedDo _ catches ( _ finally )?
       catches = ( catch {_ catch} )?
         catch = <'catch'> _ typename _ Identifier _ ImpliedDo
       finally = <'finally'> _ Blocky
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
       <PrimaryExpr> = functioncall
                     | javamethodcall
                     | Operand
                     | functiondecl
                     | TypeDecl
                     | implements
                     | funclikedecl
                     | indexed
                                                                (* Conversion |
                                                                BuiltinCall |
                                                                PrimaryExpr Selector |
                                                                PrimaryExpr Slice |
                                                                PrimaryExpr TypeAssertion | *)
         indexed = PrimaryExpr _ <'['> _ expr _ <']'>
         functioncall = PrimaryExpr Call
         javamethodcall = UnaryExpr _ <'->'> _ JavaIdentifier _ Call
           <Call> = <'('> _ ( ArgumentList _ )? <')'>
             <ArgumentList> = expressionlist                                         (* [ _ '...' ] *)
               expressionlist = expr {_ <','> _ expr}
         <TypeDecl> = <'type'> _ ( TypeSpec | <'('> _ ( TypeSpec NL )* <')'> )
	   <TypeSpec> = interfacespec | structspec
             structspec = JavaIdentifier _ <'struct'> _ <'{'>  _ (fields _)? <'}'>
               fields = Field
                        | fields NL Field
                 <Field> = Identifier | typedidentifier
	     interfacespec = JavaIdentifier _ <'interface'> _ <'{'> _ ( MethodSpec NL )* <'}'>
	       <MethodSpec> = voidmethodspec | typedmethodspec
	       voidmethodspec = JavaIdentifier _ <'('> _ methodparameters? _ <')'>
	       typedmethodspec = JavaIdentifier _ <'('> _ methodparameters? _ <')'> _ typename
		 methodparameters = methodparam
				  | methodparameters _ <','> _ methodparam
		   methodparam = symbol (_ JavaIdentifier)?
         implements = <'implements'> _ typename _ 
                        <'func'> _ <'('> _ JavaIdentifier <')'> _ (
                          MethodImpl | <'('> _ MethodImpl ( NL MethodImpl )* _ <')'>
                        )
           <MethodImpl> = typedmethodimpl | untypedmethodimpl
             untypedmethodimpl = JavaIdentifier _ <'('>  _ parameters? _ <')'> _ Blocky
             typedmethodimpl = JavaIdentifier _ <'('>  _ parameters? _ <')'> _ typename _ Blocky
         functiondecl = <'func'> _ Identifier _ Function
         funclikedecl = <'func'> _ <'<'> _ symbol _ <'>'> _ Identifier _ Function
           <Function> = FunctionPart | functionparts
             functionparts = FunctionPart _ FunctionPart {_ FunctionPart}
               <FunctionPart> = functionpart0 | functionpartn | vfunctionpart0 | vfunctionpartn
                 functionpart0 = <'('> _ <')'>  ( _ typename )? _ Blocky
		 vfunctionpart0 = <'('> _ variadic _ <')'> ( _ typename )? _ Blocky
		 functionpartn  = <'('> _ parameters _ <')'> ( _ typename )? _ Blocky
		 vfunctionpartn = <'('> _ parameters _  <','> _ variadic _ <')'> ( _ typename )? _
                                 Blocky
                   parameters = Destruct {<','> _ Destruct}
                   variadic = Identifier <'...'>
         <Operand> = Literal | OperandName | label | new  | <'('> expr <')'>       (*|MethodExpr*)
           label = #'\b\p{Lu}[\p{Lu}_0-9]*\b'
           <Literal> = BasicLit | veclit | dictlit | setlit | functionlit | shortfunctionlit
             functionlit = <'func'> _ Function
             shortfunctionlit = <'func'> _ <'{'> _ expr _ <'}'>
             <BasicLit> = int_lit | bigintlit | string_lit | regex  | rune_lit | floatlit | bigfloatlit (*| imaginary_lit *)
               floatlit = decimals '.' decimals? exponent?
                        | decimals exponent
                        | '.' decimals exponent?
                 decimals  = #'[0-9]+'
                 exponent  = ( 'e' | 'E' ) ( '+' | '-' )? decimals
               bigfloatlit = (floatlit | int_lit) 'M'
               <int_lit> = decimallit | octal_lit | hex_lit
		 decimallit = #'[1-9][0-9]*' | #'[0-9]'
		 <octal_lit>  = #'0[0-7]+'
		 <hex_lit>    = #'0x[0-9a-fA-F]+'
               bigintlit = int_lit 'N'
	       regex = <'/'> ( #'[^/\n]' | escapedslash )+ <'/'>
                 escapedslash = <'\\/'>
	       <string_lit> = rawstringlit | interpretedstringlit | clojureescape
                 rawstringlit = <#'\x60'> #'[^\x60]*' <#'\x60'>     (* \x60 is back quote character *)
                 interpretedstringlit = <#'\"'> {#'[^\"]' | '\\"'} <#'\"'>
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
	     veclit = <'['> _ ( expr {_ <','> _ expr _} )? <']'>
	     dictlit = '{' _ ( dictelement _ {<','> _ dictelement} )? _ '}'
               dictelement = expr _ <':'> _ expr
             setlit = <'set'> _ <'{'> ( _ expr _ {<','> _ expr} )? _ <'}'>
           new = <'new'> <__> typename
           <OperandName> = symbol | NonAlphaSymbol                           (*| QualifiedIdent*)
             <NonAlphaSymbol> = '=>' | '->>' | relop | addop | mulop | unary_op
                              | percent| percentnum | percentvaradic
               percent        = <'..'>
               percentnum     = <'..'> #'[1-9]'
               percentvaradic = <'...'>
`)
