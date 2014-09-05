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

import insta "instaparse/core"

var Parse = insta.parser(`
sourcefile = packageclause (expressions|topwithconst|topwithassign)
nonpkgfile = (expressions|topwithconst|topwithassign)
 packageclause = <'package'> imported <NL>
                 importdecls
   <NL> = ';' | '\n'
   importdecls = {AnyImportDecl}
     <AnyImportDecl> = importdecl | macroimportdecl | externimportdecl | typeimportdecl | exclude
     exclude = <'exclude' '('>
                  (Identifier|operator) { <','> (Identifier|operator) } <')'>
     importdecl = <'import' '('>
                    ImportSpec {ImportSpec}
                  <')'>
                | <'import'>  ImportSpec
     macroimportdecl = <'import' 'macros' '('>
                         ImportSpec {ImportSpec} <')'>
                     | <'import' 'macros'> ImportSpec
     <externimportdecl> = <'import' 'extern' '('>
                              externimportspec {externimportspec} <')'>
                     | <'import' 'extern'> externimportspec
       externimportspec = identifier
     typeimportdecl = <'import' 'type' '('>
                        typeimportspec {typeimportspec} <')'>
                     | <'import' 'type'> typeimportspec
       typeimportspec = typepackageimportspec <'.'> (
                                          JavaIdentifier
                                        | <'{'> JavaIdentifier {<','> JavaIdentifier}  <'}'>
                         )
         typepackageimportspec = JavaIdentifier {<'.'>  JavaIdentifier}
     <ImportSpec> = importspec
       importspec = ( Identifier )?  QQ imported QQ
         imported = Identifier {<'/'> Identifier}
 expressions = expr | expressions <NL> expr
   <expr>  = precedence00 | Vars | (*shortvardecl |*) ifelseexpr | letifelseexpr | tryexpr | forrange |
                   forlazy | fortimes | forcstyle | Blocky | ExprSwitchStmt


     <Blocky> = block | withconst | withassign | loop
       withconst  = <'{' 'const'> ( const <NL> | <'('> consts <')'> ) expressions <'}'>
       withassign = <'{'> assigns <NL> expressions <'}'>
         consts  = ( const {<NL> const} )?
         assigns = assign {<NL> assign}
           const  = Destruct <'='> expr
           assign = Destruct {<','> Destruct} ':=' expr {<','> expr}
	     <Destruct> = Identifier | typedidentifier | vecdestruct | dictdestruct
	       typedidentifiers = Identifier ({ <','> Identifier })? typename
	       typedidentifier = Identifier typename
		 typename = JavaIdentifier {<'.'>  JavaIdentifier} | primitivetype | string
                   <primitivetype> = long | double | 'byte' | 'short' | 'char' | 'boolean'
                     long = <'int'> | <'long'>
                     double = <'float'> | <'float64'>
                   string = <'string'>
	       vecdestruct = <'['> VecDestructElem {<','> VecDestructElem } <']'>
		 <VecDestructElem> = Destruct | variadicdestruct | label
		   variadicdestruct = Destruct Ellipsis
	       dictdestruct = <'{'> dictdestructelem { <','> dictdestructelem} <'}'>
		 dictdestructelem = Destruct <':'> expr
       loop = <'loop' '('> commaconsts <')'> ImpliedDo
         commaconsts = ( const { <','> const} )?
       <ImpliedDo> =  <'{'> expressions <'}'> | withconst | withassign
       block = <'{'> expr {<NL> expr} <'}'>
       topwithconst  =  <'const'> ( const | <'('> consts <')'> )  expressions
       topwithassign =  assigns expressions
     <ExprSwitchStmt> = boolswitch | constswitch | letconstswitch | typeswitch
                        | selectstmtingo | selectstmt
       selectstmt = <'select' '{'> (CommClause {<NL> CommClause})? <'}'>
         <CommClause> = sendclause | recvclause | recvvalclause | defaultclause
           sendclause       = <'case'> UnaryExpr        <    '<-'> UnaryExpr <':'> expressions?
           recvclause       = <'case'                   '<-'> expr <':'> expressions?
           recvvalclause    = <'case'>  identifier <'=' '<-'> expr <':'> expressions
           defaultclause    = <'default'                            ':'> expressions?
       selectstmtingo = <'select' '{'> CommClauseInGo {<NL> CommClauseInGo} <'}'>
         <CommClauseInGo> = sendclauseingo | recvclauseingo | recvvalclauseingo | defaultclause
           sendclauseingo    = <'case'> UnaryExpr        <    '<:'> UnaryExpr <':'> expressions?
           recvclauseingo    = <'case'                   '<:'> expr <':'> expressions?
           recvvalclauseingo = <'case'>  identifier <'=' '<:'> expr <':'> expressions
       typeswitch = <'switch'> PrimaryExpr <'.' '(' 'type' ')' '{'
                           'case'> typename <':'> expressions
                         {<NL 'case'> typename <':'> expressions}
                         (<NL 'default'          ':'> expressions )? <'}'>
       boolswitch = <'switch' '{'> boolcaseclause {<NL> boolcaseclause} <'}'>
       constswitch = <'switch'> expr <'{'> constcaseclause {<NL> constcaseclause} <'}'>
       letconstswitch = <'switch'> Destruct <':='> expr <NL>
                                   expr <'{'> constcaseclause {<NL> constcaseclause} <'}'>
	 boolcaseclause = boolswitchcase <':'> expressions
	 constcaseclause = constswitchcase <':'> expressions
	   boolswitchcase = <'case'> expressionlist | <'default'>
	   constswitchcase = <'case'> constantlist | <'default'>
	     constantlist = expr { <','> expr}
	       <Constant> = label | BasicLit | veclit | dictlit | setlit | structlit
     operator =
                 or
                |and
                |'<-'|'<:'|equals|noteq|'<'|'<='|'>='|'>'
                |'+'|!'->' '-'|bitor|bitxor
                |'*'|'/'|mod|shiftleft|shiftright|bitand|bitandnot
                |'+='|'-='
     precedence00 = precedence0
                 | precedence00 SendOp precedence0
       <SendOp> = sendop | sendopingo
         sendop     = <'<-'>
         sendopingo = <'<:'>
     precedence0 = precedence1
                 | precedence0 symbol precedence1
       symbol = Identifier
              | Identifier <'.'>  Identifier
              | Identifier <'.'>  operator
              | javastatic
       precedence1 = precedence2
                   | precedence1 or precedence2
	 or = <'||'>
	 precedence2 = precedence3
                     | precedence2 and precedence3
	   and = <'&&'>
	   precedence3 = precedence4
                       | precedence3 relop  precedence4
             chanops = '<-' | '<:'
             relop = equals | noteq | (!chanops '<') | '<=' | '>=' | '>'
	       equals = <'=='>
               noteq  = <'!='>
	     precedence4 = precedence5
                         | precedence4 addop precedence5
	       addop = '+' | !'->' '-' | ( !or bitor ) | bitxor
                 bitor = <'|'>
                 bitxor = <'^'>
	       precedence5 = UnaryExpr
                           | precedence5 mulop UnaryExpr
	         mulop = '*' | (!'//' '/') | mod | shiftleft | shiftright | bitand | bitandnot
                   shiftleft = <'<<'>
                   shiftright = <'>>'>
                   mod = <'%'>
                   bitand = !and <'&'>
                   bitandnot = !and <'&^'>
	   javastatic = typename <'::'> JavaIdentifier
	     <JavaIdentifier> = #'\b[\p{L}_][\p{L}_\p{Nd}]*\b'
                              | underscorejavaidentifier
               underscorejavaidentifier = <'_'> JavaIdentifier
	   <Identifier> = !Keyword !hexlit (identifier | dashidentifier | isidentifier | mutidentifier |
			  escapedidentifier)
             Keyword = '\bconst\b' | '\bfor\b' | '\bnew\b' | '\bpackage\b' | '\brange\b' | '\bif\b'
	     identifier = #'[\p{L}_[\p{S}&&[^\p{Punct}]]][\p{L}_[\p{S}&&[^\p{Punct}]]\p{Nd}]*'
	     dashidentifier = <'_'> identifier
	     isidentifier = <'is'> #'\p{L}' identifier
	     mutidentifier = <'mutate'> #'\p{L}' identifier
	     (* escapedidentifier = <'\\'> #'\b[\p{L}_\p{Sm}][\p{L}_\p{Sm}\p{Nd}]*\b' *)
	     escapedidentifier = #'\\[^\\]+\\'
     <Vars> = <'var'> ( <'('> VarDecl {VarDecl} <')'> | VarDecl )
     <VarDecl> = primarrayvardecl | arrayvardecl | vardecl1 | vardecl2
       primarrayvardecl = Identifier <'['> int_lit  <']'> primitivetype
       arrayvardecl = Identifier <'['> int_lit  <']'> typename
       vardecl1 = Identifier ( typename )? <'='> expr
       vardecl2 = Identifier  <','> Identifier ( typename )? <'='> expr <','> expr
     ifelseexpr = <'if'> expr Blocky ( <'else'> Blocky )?
     letifelseexpr = <'if'> Destruct <':='> expr <NL>
                            expr Blocky ( <'else'> Blocky )?
     forrange = <'for'> Destruct <':=' 'range'> expr  Blocky
     forlazy = <'for'> Destruct <':=' 'lazy'> expr
               (<'if'> expr )? Blocky
     fortimes = <'for'> Identifier <':=' 'times'> expr Blocky
     forcstyle = <'for'> Identifier <':=' '0' ';'>
                         Identifier <'<'> expr <';'>
                         Identifier <'++'>
                         Blocky
     tryexpr = <'try'> ImpliedDo catches ( finally )?
       catches = ( catch {catch} )?
         catch = <'catch'> typename Identifier ImpliedDo
       finally = <'finally'> Blocky
     <UnaryExpr> = unaryexpr  (* TODO(eob) remove this indirection *)
       assoc = expr <'+=' '{'> associtem { <','> associtem } <'}'>
       dissoc = expr <'-=' '{'> associtem { <','> associtem } <'}'>
	 associtem = expr <':'> expr
       associn = expr <'+=' '{'> associnpath <':'> expr <'}'>
	 associnpath = expr expr {expr}
       unaryexpr = unary_op unaryexpr
                 | PrimaryExpr | javafield | ReaderMacro | assoc | dissoc | associn | prefixedblock
	 <unary_op> = '+' | !'->' '-' | '!' | not | (!and '&') | bitnot | take | takeingo
	   bitnot = <'^'>
	   not    = <'!'>
           takeingo = <'<:'>
           take     = <'<-'>
       <ReaderMacro> = deref | syntaxquote | unquote | unquotesplicing
       deref           = <'*'>               UnaryExpr
       syntaxquote     = <'syntax'>     UnaryExpr
       unquote         = <'unquote'>         UnaryExpr
       unquotesplicing = <'unquotes'> UnaryExpr
       javafield  = UnaryExpr <'->'> JavaIdentifier
       <PrimaryExpr> = Routine
                     | prefixedroutine
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
         prefixedroutine = prefix Routine
         prefixedblock   = prefix ImpliedDo
           prefix = asyncprefix | 'dosync'
             asyncprefix = 'go' | 'thread'
         threadblock = <'thread'> ImpliedDo
         chan      = <'make' '(' 'chan'> ( <typename>)? ( <','> expr)? <')'>
         <Routine> = functioncall
                     | MappedFunctionCall
                     | variadiccall
                     | typeconversion
                     | javamethodcall
         typeconversion = primitivetype <'('> expr <')'>
         indexed = PrimaryExpr <'['> expr <']'>
         takeslice = PrimaryExpr <'[' ':'> expr <']'>
         dropslice = PrimaryExpr <'['>  expr <':' ']'>
         variadiccall = PrimaryExpr
                           <'('> ( ArgumentList <','> )? Ellipsis PrimaryExpr <')'>
         functioncall = PrimaryExpr Call
         <MappedFunctionCall> = len
           len = <'len'> Call
         javamethodcall = UnaryExpr <'->'> JavaIdentifier Call
           <Call> =  <'('> ArgumentList? <')'>
             <ArgumentList> = expressionlist                                         (* [ Ellipsis ] *)
               expressionlist = expr { <','> expr} ( <','>)?
         <TypeDecl> = <'type'> ( TypeSpec | <'('> {TypeSpec} <')'> )
	   <TypeSpec> = interfacespec | structspec
             structspec = JavaIdentifier <'struct' '{'> (fields )? <'}'>
               fields = Field
                        | fields <NL> Field
                 <Field> = Identifier | typedidentifiers
	     interfacespec = JavaIdentifier <'interface' '{'> {MethodSpec} <'}'>
	       <MethodSpec> = voidmethodspec | typedmethodspec
	       voidmethodspec = Identifier <'('> methodparameters? <')'>
	       typedmethodspec = Identifier <'('> methodparameters? <')'> typename
		 methodparameters = methodparam
				  | methodparameters <','> methodparam
		   methodparam = symbol ( Identifier)?
         implements = <'implements'> typename <'func' '('> JavaIdentifier <')'> (
                          MethodImpl | <'('>  MethodImpl {<NL> MethodImpl}   <')'>
                        )
           <MethodImpl> = typedmethodimpl | untypedmethodimpl
             untypedmethodimpl = Identifier <'('>  parameters? <')'>
                                   (ReturnBlock|Blocky)
             typedmethodimpl = Identifier <'('>  parameters? <')'> typename
                                   (ReturnBlock|Blocky)
         functiondecl = <'func'> (Identifier|operator) Function
         funclikedecl = <'func' '<'> symbol <'>'> Identifier Function
           <Function> = FunctionPart | functionparts
             functionparts = FunctionPart FunctionPart { FunctionPart}
               <FunctionPart> = functionpart0 | functionpartn | vfunctionpart0 | vfunctionpartn
                 functionpart0 = <'(' ')'>  ( typename )? (ReturnBlock|Blocky)
		 vfunctionpart0 = <'('> variadic <')'> ( typename )? (ReturnBlock|Blocky)
		 functionpartn  = <'('> parameters <')'> ( typename )? (ReturnBlock|Blocky)
		 vfunctionpartn = <'('> parameters  <','> variadic <')'> ( typename )?
                                 (ReturnBlock|Blocky)
                   parameters = Destruct {<','> Destruct}
                   variadic = Identifier Ellipsis
                   <ReturnBlock> = <'{' 'return'> expr <'}'>
         <Operand> = Literal | OperandName | label | islabel | new  | <'('> expr <')'>       (*|MethodExpr*)
           label = #'\b\p{Lu}[\p{Lu}_\p{Nd}#\.]*\b'
	   islabel = <'IS_'> #'\b\p{Lu}[\p{Lu}_\p{Nd}#\.]*\b'
           <Literal> = BasicLit | veclit | dictlit | setlit | structlit | functionlit | shortfunctionlit
             functionlit = <'func'> Function
             shortfunctionlit = <'func' '{'> expr <'}'>
             <BasicLit> = int_lit | bigintlit | regex | string_lit | rune_lit | floatlit | bigfloatlit (*| imaginary_lit *)
               floatlit = #'([0-9]+\.[0-9]*([eE]?[\+\-]?[0-9]+)?|[0-9]+[eE]?[\+\-]?[0-9]+|\.[0-9]+[eE]?[\+\-]?[0-9]*)'
               (*floatlit = decimals '.' decimals? exponent?
                        | decimals exponent
                        | '.' decimals exponent?
                 decimals  = #'[0-9]+'
                 exponent  = ( 'e' | 'E' ) ( '+' | !'->' '-' )? decimals*)
               bigfloatlit = (floatlit | int_lit) 'M'
               <int_lit> = decimallit | octal_lit | hexlit
		 decimallit = #'[1-9][0-9]*' | #'[0-9]'
		 <octal_lit>  = #'0[0-7]+'
		 hexlit    = <'0x'> #'[0-9a-fA-F]+'
               bigintlit = int_lit 'N'
	       (* regex = <'/'> ( #'[^/\n]' | escapedslash )+ <'/'> *)
               regex =  #'/[^/\\]+(\\.[^/\\]*)*/' | #'/[^/\\]*(\\.[^/\\]*)+/'
                 (* escapedslash = <'\\/'> *)
	       <string_lit> = interpretedstringlit | clojureescape
                 interpretedstringlit = #'["“”](?:[^"\\]|\\.)*["“”]'
                 clojureescape = <'\\' #'\x60'> #'[^\x60]*' <#'\x60'>       (* \x60 is back quote *)
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
	     veclit =                               <'['> ( expr {<','> expr} )? <']'>
                     |  <'[' ']' typename '{'> ( expr { <','> expr } )? <'}'>
	     dictlit = '{' ( dictelement {<','> dictelement} )? ( <','>)? '}'
               dictelement = expr <':'> expr
             NotType = 'func' | 'set' | prefix
             structlit = !NotType typename <'{'> ( expr {<','> expr} )? (<','> )? <'}'>
             setlit = <'set' '{'> ( expr {<','> expr} )? <'}'>
           new = <'new'> typename
           <OperandName> = symbol | NonAlphaSymbol                           (*| QualifiedIdent*)
             <NonAlphaSymbol> = '=>' | '->>' | relop | addop | mulop | unary_op
                              | percentnum | percentvaradic
               percentnum     = <'$'> #'[1-9]'
               percentvaradic = <'$*'>
  <Ellipsis> = <'...'> | <'…'>
  <QQ> = <'"'> | <'“'>  | <'”'>
`,
	AUTO_WHITESPACE, STANDARD, //whitespace,
	NO_SLURP, true,  // for App Engine compatibility
)
