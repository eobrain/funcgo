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

package parser

import (
	insta "instaparse/core"
)

var Parse = insta.parser(`
sourcefile = NL? packageclause (expressions|topwithconst) _
nonpkgfile = NL? (expressions|topwithconst) _
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
       importspec = ( Identifier _ )?  QQ imported QQ
         imported = Identifier {<'/'> Identifier}
 expressions = expr | expressions NL expr
   <expr>  = precedence0 | Vars | shortvardecl | ifelseexpr | letifelseexpr | tryexpr | forrange |
                   forlazy | fortimes | forcstyle | Blocky | ExprSwitchStmt | sendstmt
                                                                            | sendstmtingo

     <ExprSwitchStmt> = boolswitch | constswitch | letconstswitch | typeswitch
                        | selectstmtingo | selectstmt
       selectstmt = <'select'> _ <'{'> _ { CommClause _ } <'}'>
         <CommClause> = sendclause | recvclause | recvvalclause | defaultclause
           sendclause       = <'case'> _ expr _ <'<-'> _ expr                _ <':'> (_ expressions)?
           recvclause       = <'case'> _                       <'<-'> _ expr _ <':'> (_ expressions)?
           recvvalclause    = <'case'> _  identifier _ <'='> _ <'<-'> _ expr _ <':'>  _ expressions
           defaultclause    = <'default'>                                    _ <':'> (_ expressions)?
       selectstmtingo = <'select'> _ <'{'> _ { CommClauseInGo _ } <'}'>
         <CommClauseInGo> = sendclauseingo | recvclauseingo | recvvalclauseingo | defaultclause
           sendclauseingo    = <'case'> _ expr _ <'<:'> _ expr                _ <':'> (_ expressions)?
           recvclauseingo    = <'case'> _                       <'<:'> _ expr _ <':'> (_ expressions)?
           recvvalclauseingo = <'case'> _  identifier _ <'='> _ <'<:'> _ expr _ <':'>  _ expressions
       typeswitch = <'switch'> _ PrimaryExpr _ <'.'> _ <'('> _ <'type'> _ <')'> _  <'{'>
                         _   <'case'> _ typename _ <':'> _ expressions
                         {NL <'case'> _ typename _ <':'> _ expressions}
                         (NL <'default'>         _ <':'> _ expressions )?
                    _ <'}'>
       boolswitch = <'switch'> _ <'{'>  _ boolcaseclause { NL boolcaseclause } _ <'}'>
       constswitch = <'switch'> _ expr _ <'{'> _ constcaseclause { NL constcaseclause } _ <'}'>
       letconstswitch = <'switch'> _ Destruct _ <':='> _ expr _ <';'>
                                   _ expr _ <'{'> _ constcaseclause { NL constcaseclause } _ <'}'>
	 boolcaseclause = boolswitchcase _ <':'> _ expressions
	 constcaseclause = constswitchcase _ <':'> _ expressions
	   boolswitchcase = <'case'> _ expressionlist | <'default'>
	   constswitchcase = <'case'> _ constantlist | <'default'>
	     constantlist = expr {_ <','> _ expr}
	       <Constant> = label | BasicLit | veclit | dictlit | setlit | structlit
     <Blocky> = block | withconst | loop
       loop = <'loop'> _  <'('> _ {consts} _ <')'> _ ImpliedDo
       <ImpliedDo> =  <'{'> _ expressions _ <'}'> | withconst
       block = <'{'> _ expr {NL expr} _ <'}'>
       topwithconst =  <'const'> _ ( const NL | <'('> _ consts _ <')'> )  _ expressions
       withconst = <'{'> _ <'const'> _ ( const NL | <'('> _ consts _ <')'> )  _ expressions _ <'}'>
         consts = ( const {NL const} )?
           const = Destruct _ <'='> _ expr
	     <Destruct> = Identifier | typedidentifier | vecdestruct | dictdestruct
	       typedidentifiers = Identifier ({_ <','> _ Identifier })? _ typename
	       typedidentifier = Identifier _ typename
		 typename = JavaIdentifier {<'.'>  JavaIdentifier} | primitivetype | string
                   <primitivetype> = long | double | 'byte' | 'short' | 'char' | 'boolean'
                     long = <'int'> | <'long'>
                     double = <'float'> | <'float64'> 
                   string = <'string'>
	       vecdestruct = <'['> _ VecDestructElem _ {<','> _ VecDestructElem _ } <']'>
		 <VecDestructElem> = Destruct | variadicdestruct | label
		   variadicdestruct = Destruct Ellipsis
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
             chanops = '<-' | '<:'
             relop = equals | noteq | (!chanops '<') | '<=' | '>='           (* TODO(eob)  | '>' *)
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
     shortvardecl =  identifier _ <':='> _ expr
                   | identifier _ <','> _ identifier _ <':='> _ expr  _ <','> _ expr
                   | identifier _ <','> _ identifier _<','> _ identifier _
                                   <':='> _ expr  _ <','> _ expr  _ <','> _ expr
               (*   | Identifier _ ',' _ shortvardecl _ ',' _ expr *)
     sendstmtingo = expr _ <'<:'> _ expr
     sendstmt     = expr _ <'<-'> _ expr
     <Vars> = <'var'> _ ( <'('> _ VarDecl {NL VarDecl} _ <')'> | VarDecl )
     <VarDecl> = primarrayvardecl | arrayvardecl | vardecl1 | vardecl2
       primarrayvardecl = Identifier _ <'['> _ int_lit  _ <']'> _ primitivetype
       arrayvardecl = Identifier _ <'['> _ int_lit  _ <']'> _ typename
       vardecl1 = Identifier ( _ typename )? _ <'='> _ expr
       vardecl2 = Identifier  _ <','> _ Identifier ( _ typename )? _ <'='> _ expr _ <','> _ expr
     ifelseexpr = <'if'> _ expr _ Blocky ( _ <'else'> _ Blocky )?
     letifelseexpr = <'if'> _ Destruct _ <':='> _ expr _ <';'>
                            _ expr _ Blocky ( _ <'else'> _ Blocky )?
     forrange = <'for'> <__> Destruct _ <':='> _ <'range'> <_> expr _  Blocky
     forlazy = <'for'> <__> Destruct _ <':='> _ <'lazy'> <_> expr
               ( <__> <'if'> <__> expr )? _ Blocky
     fortimes = <'for'> <__> Identifier _ <':='> _ <'times'> <_> expr _ Blocky
     forcstyle = <'for'> <__> Identifier _ <':='> _ <'0'> _ <';'>
                         _ Identifier _ <'<'> _ expr _ <';'>
                         _ Identifier _ <'++'>
                         _ Blocky
     tryexpr = <'try'> _ ImpliedDo _ catches ( _ finally )?
       catches = ( catch {_ catch} )?
         catch = <'catch'> _ typename _ Identifier _ ImpliedDo
       finally = <'finally'> _ Blocky
     <UnaryExpr> = PrimaryExpr | javafield | ReaderMacro | assoc | dissoc | associn | unaryexpr
       assoc = expr _ <'+='> _ <'{'> _ associtem _ { <','> _ associtem _ } <'}'>
       dissoc = expr _ <'-='> _ <'{'> _ associtem _ { <','> _ associtem _ } <'}'>
	 associtem = expr _ <':'> _ expr
       associn = expr _ <'+='> _ <'{'> _ associnpath _ <':'> _ expr _ <'}'>
	 associnpath = expr _ expr {_ expr}
       unaryexpr = unary_op _ UnaryExpr
	 <unary_op> = '+' | '-' | '!' | not | (!and '&') | bitnot | take | takeingo
	   bitnot = <'^'>
	   not    = <'!'>
           takeingo = <'<:'>
           take     = <'<-'>
       <ReaderMacro> = deref | syntaxquote | unquote | unquotesplicing
       deref           = <'*'>               _ UnaryExpr
       syntaxquote     = <'syntax'>     _ UnaryExpr
       unquote         = <'unquote'>         _ UnaryExpr
       unquotesplicing = <'unquotes'> _ UnaryExpr
       javafield  = UnaryExpr _ <'->'> _ JavaIdentifier
       <PrimaryExpr> = Routine
                     | goroutine
                     | goblock
                     | chan
                     | Operand
                     | functiondecl
                     | TypeDecl
                     | implements
                     | funclikedecl
                     | indexed
                     | dropslice
                     | takeslice
                                                                (* Conversion |
                                                                BuiltinCall |
                                                                PrimaryExpr Selector |
                                                                PrimaryExpr Slice |
                                                                PrimaryExpr TypeAssertion | *)
         goroutine = <'go'> _ Routine
         goblock = <'go'> _ ImpliedDo
         chan      = <'make'> _ <'('> _ <'chan'> (_ <typename>)? (_ <','> _ expr)? _ <')'>
         <Routine> = functioncall
                     | MappedFunctionCall
                     | variadiccall
                     | typeconversion
                     | javamethodcall
         typeconversion = primitivetype _ <'('> _ expr _ <')'>
         indexed = PrimaryExpr _ <'['> _ expr _ <']'>
         takeslice = PrimaryExpr _ <'['> _ <':'> _ expr _ <']'>
         dropslice = PrimaryExpr _ <'['>  _ expr _ <':'> _ <']'>
         variadiccall = PrimaryExpr
                           <'('> _ ( ArgumentList _ <','> _ )? _ Ellipsis _ PrimaryExpr _ <')'>
         functioncall = PrimaryExpr Call
         <MappedFunctionCall> = len
           len = <'len'> Call
         javamethodcall = UnaryExpr _ <'->'> _ JavaIdentifier _ Call
           <Call> =  <'('> _ ( ArgumentList _ )? <')'>
             <ArgumentList> = expressionlist                                         (* [ _ Ellipsis ] *)
               expressionlist = expr {_ <','> _ expr} (_ <','>)?
         <TypeDecl> = <'type'> _ ( TypeSpec | <'('> _ ( TypeSpec NL )* <')'> )
	   <TypeSpec> = interfacespec | structspec
             structspec = JavaIdentifier _ <'struct'> _ <'{'>  _ (fields _)? <'}'>
               fields = Field
                        | fields NL Field
                 <Field> = Identifier | typedidentifiers
	     interfacespec = JavaIdentifier _ <'interface'> _ <'{'> _ ( MethodSpec NL )* <'}'>
	       <MethodSpec> = voidmethodspec | typedmethodspec
	       voidmethodspec = Identifier _ <'('> _ methodparameters? _ <')'>
	       typedmethodspec = Identifier _ <'('> _ methodparameters? _ <')'> _ typename
		 methodparameters = methodparam
				  | methodparameters _ <','> _ methodparam
		   methodparam = symbol (_ Identifier)?
         implements = <'implements'> _ typename _
                        <'func'> _ <'('> _ JavaIdentifier <')'> _ (
                          MethodImpl | <'('> _ MethodImpl ( NL MethodImpl )* _ <')'>
                        )
           <MethodImpl> = typedmethodimpl | untypedmethodimpl
             untypedmethodimpl = Identifier _ <'('>  _ parameters? _ <')'> _
                                   (ReturnBlock|Blocky)
             typedmethodimpl = Identifier _ <'('>  _ parameters? _ <')'> _ typename _
                                   (ReturnBlock|Blocky)
         functiondecl = <'func'> _ Identifier _ Function
         funclikedecl = <'func'> _ <'<'> _ symbol _ <'>'> _ Identifier _ Function
           <Function> = FunctionPart | functionparts
             functionparts = FunctionPart _ FunctionPart {_ FunctionPart}
               <FunctionPart> = functionpart0 | functionpartn | vfunctionpart0 | vfunctionpartn
                 functionpart0 = <'('> _ <')'>  ( _ typename )? _ (ReturnBlock|Blocky)
		 vfunctionpart0 = <'('> _ variadic _ <')'> ( _ typename )? _ (ReturnBlock|Blocky)
		 functionpartn  = <'('> _ parameters _ <')'> ( _ typename )? _ (ReturnBlock|Blocky)
		 vfunctionpartn = <'('> _ parameters _  <','> _ variadic _ <')'> ( _ typename )? _
                                 (ReturnBlock|Blocky)
                   parameters = Destruct {<','> _ Destruct}
                   variadic = Identifier Ellipsis
                   <ReturnBlock> = <'{'> _ <'return'> _ expr _ <'}'>
         <Operand> = Literal | OperandName | label | islabel | new  | <'('> expr <')'>       (*|MethodExpr*)
           label = #'\b\p{Lu}[\p{Lu}_0-9#\.]*\b'
	   islabel = <'IS_'> #'\b\p{Lu}[\p{Lu}_0-9#\.]*\b'
           <Literal> = BasicLit | veclit | dictlit | setlit | structlit | functionlit | shortfunctionlit
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
                 interpretedstringlit = QQ {#'[^\"]' | '\\"'} QQ
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
	     veclit =                               <'['> _ ( expr {_ <','> _ expr _} )? <']'>
                     |  <'['> _  <']'> _ <typename> <'{'> _ ( expr {_ <','> _ expr _} )? <'}'>
	     dictlit = '{' _ ( dictelement _ {<','> _ dictelement} )? (_ <','>)? _ '}'
               dictelement = expr _ <':'> _ expr
             NotType = 'func' | 'set'
             structlit = !NotType typename _ <'{'> ( _ expr _ {<','> _ expr} )? _ (<','> _)? <'}'> 
             setlit = <'set'> _ <'{'> ( _ expr _ {<','> _ expr} )? _ <'}'>
           new = <'new'> <__> typename
           <OperandName> = symbol | NonAlphaSymbol                           (*| QualifiedIdent*)
             <NonAlphaSymbol> = '=>' | '->>' | relop | addop | mulop | unary_op
                              | percentnum | percentvaradic
               percentnum     = <'$'> #'[1-9]'
               percentvaradic = <'$*'>
  <Ellipsis> = <'...'> | <'…'>
  <QQ> = <'"'> | <'“'>  | <'”'>
`)
