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

whitespaceOrComments := insta.parser(`
    ws-or-comments = #'(\s|(//[^\n]*\n))+'
`, NO_SLURP, true,  // for App Engine compatibility
);

var Parse = insta.parser(`
sourcefile = packageclause (expressions|topwithconst|topwithassign)
nonpkgfile = (expressions|topwithconst|topwithassign)
 packageclause = <#'\bpackage\b'> pkg <NL> importdecls
   pkg =  Identifier {<'/'> Identifier}
   <NL> = #'\s*[;\n]\s*' | #'\s*//[^\n]*\n\s*'
   importdecls = {AnyImportDecl}
     <AnyImportDecl> = importdecl | macroimportdecl | externimportdecl | typeimportdecl | exclude
     exclude = <#'\bexclude\b' '('>
                  (Identifier|operator) { <','> (Identifier|operator) } <')'>
     importdecl = <#'\bimport\b' '('>
                    ImportSpec {ImportSpec}
                  <')'>
                | <#'\bimport\b'>  ImportSpec
     macroimportdecl = <#'\bimport\b' #'\bmacros\b' '('>
                         ImportSpec {ImportSpec} <')'>
                     | <#'\bimport\b' #'\bmacros\b'> ImportSpec
     <externimportdecl> = <#'\bimport\b' #'\bextern\b' '('>
                              externimportspec {externimportspec} <')'>
                     | <#'\bimport\b' #'\bextern\b'> externimportspec
       externimportspec = identifier
     typeimportdecl = <#'\bimport\b' #'\btype\b' '('>
                        typeimportspec {typeimportspec} <')'>
                     | <#'\bimport\b' #'\btype\b'> typeimportspec
       typeimportspec = typepackageimportspec <'.'> (
                                          JavaIdentifier
                                        | <'{'> JavaIdentifier {<','> JavaIdentifier}  <'}'>
                         )
         typepackageimportspec = JavaIdentifier {<'.'>  JavaIdentifier}
     <ImportSpec> = importspec
       importspec = ( Identifier )?  string_lit
 expressions = expr
              | expressions <NL> expr
   <expr>  = precedence00 | Vars | (*shortvardecl |*) ifelseexpr | letifelseexpr | tryexpr | forrange |
                   forlazy | fortimes | forcstyle | Blocky | ExprSwitchStmt
                     | functiondecl


     <Blocky> = block | withconst | withassign | loop
       withconst  = <'{' #'\bconst\b'> ( const <NL> | <'('> consts <')'> ) expressions <'}'>
       withassign = <'{'> assigns <NL> expressions <'}'>
         consts  = ( const {<NL> const} )?
         assigns = assign {<NL> assign}
           const  = Destruct <'='> expr
           assign = Destruct {<','> Destruct} ':=' expr {<','> expr}
	     <Destruct> = Identifier | typedidentifier | vecdestruct | dictdestruct
	       typedidentifiers = Identifier ({ <','> Identifier })? typename
	       typedidentifier = Identifier typename
		 typename = JavaIdentifier {<'.'>  JavaIdentifier} | primitivetype | string
                   <primitivetype> = long | double | #'\bbyte\b' | #'\bshort\b' | #'\bchar\b' | #'\bboolean\b'
                     long = <#'\bint\b'> | <#'\blong\b'>
                     double = <#'\bfloat\b'> | <#'\bfloat64\b'>
                   string = <#'\bstring\b'>
	       vecdestruct = <'['> VecDestructElem {<','> VecDestructElem } <']'>
		 <VecDestructElem> = Destruct | variadicdestruct | label
		   variadicdestruct = Destruct Ellipsis
	       dictdestruct = <'{'> dictdestructelem { <','> dictdestructelem} <'}'>
		 dictdestructelem = Destruct <':'> expr
       loop = <#'\bloop\b' '('> commaconsts <')'> ImpliedDo
         commaconsts = ( const { <','> const} )?
       <ImpliedDo> =  <'{'> expressions <'}'> | withconst | withassign
       block = <'{'> expr {<NL> expr} <'}'>
       topwithconst  =  <#'\bconst\b'> ( const | <'('> consts <')'> )  expressions
       topwithassign =  assigns <NL> expressions
     <ExprSwitchStmt> = boolswitch | constswitch | letconstswitch | typeswitch
                        | selectstmtingo | selectstmt
       selectstmt = <#'\bselect\b' '{'> (CommClause {<NL> CommClause})? <'}'>
         <CommClause> = sendclause | recvclause | recvvalclause | defaultclause
           sendclause       = <#'\bcase\b'> UnaryExpr        <    '<-'> UnaryExpr <':'> expressions?
           recvclause       = <#'\bcase\b'                   '<-'> expr <':'> expressions?
           recvvalclause    = <#'\bcase\b'>  identifier <'=' '<-'> expr <':'> expressions
           defaultclause    = <#'\bdefault\b'                            ':'> expressions?
       selectstmtingo = <#'\bselect\b' '{'> CommClauseInGo {<NL> CommClauseInGo} <'}'>
         <CommClauseInGo> = sendclauseingo | recvclauseingo | recvvalclauseingo | defaultclause
           sendclauseingo    = <#'\bcase\b'> UnaryExpr        <    '<:'> UnaryExpr <':'> expressions?
           recvclauseingo    = <#'\bcase\b'                   '<:'> expr <':'> expressions?
           recvvalclauseingo = <#'\bcase\b'>  identifier <'=' '<:'> expr <':'> expressions
       typeswitch = <#'\bswitch\b'> PrimaryExpr <'.' '(' #'\btype\b' ')' '{'
                           #'\bcase\b'> typename <':'> expressions
                         {<NL #'\bcase\b'> typename <':'> expressions}
                         (<NL #'\bdefault\b'          ':'> expressions )? <'}'>
       boolswitch = <#'\bswitch\b' '{'> boolcaseclause {<NL> boolcaseclause} <'}'>
       constswitch = <#'\bswitch\b'> expr <'{'> constcaseclause {<NL> constcaseclause} <'}'>
       letconstswitch = <#'\bswitch\b'> Destruct <':='> expr <NL>
                                   expr <'{'> constcaseclause {<NL> constcaseclause} <'}'>
	 boolcaseclause = boolswitchcase <':'> expressions
	 constcaseclause = constswitchcase <':'> expressions
	   boolswitchcase = <#'\bcase\b'> expressionlist | <#'\bdefault\b'>
	   constswitchcase = <#'\bcase\b'> constantlist | <#'\bdefault\b'>
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
                 | assoc | dissoc | associn
       assoc = precedence0 <'+=' '{'> associtem { <','> associtem } <'}'>
       dissoc = precedence0 <'-=' '{'> associtem { <','> associtem } <'}'>
	 associtem = precedence0 <':'> precedence0
       associn = precedence0 <'+=' '{'> associnpath <':'> precedence0 <'}'>
	 associnpath = precedence0 precedence0 {precedence0}
       <SendOp> = sendop | sendopingo
         sendop     = <'<-'>
         sendopingo = <'<:'>
     precedence0 = precedence1
                 | precedence0 <DoubleSpace> symbol <DoubleSpace>precedence1
       DoubleSpace = <#'[ \t][ \t]'>
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
             relop = equals | noteq | (!SendOp '<') | '<=' | '>=' | '>'
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
                   bitand = !and !bitandnot <'&'>
                   bitandnot = !and <'&^'>
	   javastatic = typename <'::'> JavaIdentifier
	     <JavaIdentifier> = #'\b[\p{L}_][\p{L}_\p{Nd}]*\b'
                              | underscorejavaidentifier
               underscorejavaidentifier = #'\b_[\p{L}_][\p{L}_\p{Nd}]*\b'
	   <Identifier> = !(Keyword | hexlit) (identifier | isidentifier | mutidentifier |
			  escapedidentifier)
             Keyword = #'\bcase\b'
                     | #'\bconst\b'
                     | #'\bfor\b'
                     | #'\bif\b'
                     | #'\bnew\b'
                     | #'\bpackage\b'
                   (*| #'\brange\b'*)
                     | #'\bselect\b'
	     identifier = #'[\p{L}_[\p{S}&&[^\p{Punct}]]][\p{L}_[\p{S}&&[^\p{Punct}]]\p{Nd}]*'
	     isidentifier = <#'\bis'> #'\p{L}' identifier
	     mutidentifier = <#'\bmutate'> #'\p{L}' identifier
	     escapedidentifier = #'\\[^\n\\]+\\'
     <Vars> = <#'\bvar\b'> ( <'('> VarDecl+ <')'> | VarDecl )
     <VarDecl> = primarrayvardecl | arrayvardecl | vardecl1 | vardecl2
       primarrayvardecl = Identifier <'['> int_lit  <']'> primitivetype
       arrayvardecl = Identifier <'['> int_lit  <']'> typename
       vardecl1 = Identifier ( typename )? <'='> expr
       vardecl2 = Identifier  <','> Identifier ( typename )? <'='> precedence00 <','> precedence00
     ifelseexpr = <#'\bif\b'> expr Blocky ( <#'\belse\b'> Blocky )?
     letifelseexpr = <#'\bif\b'> Destruct <':='> expr <NL>
                            expr Blocky ( <#'\belse\b'> Blocky )?
     forrange = <#'\bfor\b'> Destruct <':=' #'\brange\b'> expr  Blocky
     forlazy = <#'\bfor\b'> Destruct <':=' #'\blazy\b'> expr
               (<#'\bif\b'> expr )? Blocky
     fortimes = <#'\bfor\b'> Identifier <':=' #'\btimes\b'> expr Blocky
     forcstyle = <#'\bfor\b'> Identifier <':=' '0' ';'>
                         Identifier <'<'> expr <';'>
                         Identifier <'++'>
                         Blocky
     tryexpr = <#'\btry\b'> ImpliedDo catches finally?
       catches = {catch}
         catch = <#'\bcatch\b'> typename Identifier ImpliedDo
       finally = <#'\bfinally\b'> ImpliedDo
     <UnaryExpr> = unaryexpr  (* TODO(eob) remove this indirection *)
       unaryexpr = unary_op unaryexpr
                 | PrimaryExpr | javafield | ReaderMacro | prefixedblock
	 <unary_op> = '+' | !'->' '-' | '!' | not | bitnot | take | takeingo
	   bitnot = <'^'>
	   not    = <'!'>
           takeingo = <'<:'>
           take     = <'<-'>
       <ReaderMacro> = deref | syntaxquote | unquote | unquotesplicing
       deref           = <'*'>               UnaryExpr
       syntaxquote     = <#'\bsyntax\b'>     UnaryExpr
       unquote         = <#'\bunquote\b'>         UnaryExpr
       unquotesplicing = <#'\bunquotes\b'> UnaryExpr
       javafield  = UnaryExpr <'->'> JavaIdentifier
       <PrimaryExpr> = Routine
                     | prefixedroutine
                     | chan
                     | Operand
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
           prefix = asyncprefix | #'\bdosync\b'
             asyncprefix = #'\bgo\b' | #'\bthread\b'
         threadblock = <#'\bthread\b'> ImpliedDo
         chan      = <#'\bmake\b' '(' #'\bchan\b'> ( <typename>)? ( <','> expr)? <')'>
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
           len = <#'\blen\b'> Call
         javamethodcall = UnaryExpr <'->'> JavaIdentifier Call
           <Call> =  <'('> ArgumentList? <')'>
             <ArgumentList> = expressionlist                                      (* [ Ellipsis ] *)
               expressionlist = expr { <','> expr} ( <','>)?
         <TypeDecl> = <#'\btype\b'> ( TypeSpec | <'('> {TypeSpec} <')'> )
	   <TypeSpec> = interfacespec | structspec
             structspec = JavaIdentifier <#'\bstruct\b' '{'> (fields )? <'}'>
               fields = Field
                        | fields <NL> Field
                 <Field> = Identifier | typedidentifiers
	     interfacespec = JavaIdentifier <#'\binterface\b' '{'> {MethodSpec} <'}'>
	       <MethodSpec> = voidmethodspec | typedmethodspec
	       voidmethodspec = Identifier <'('> methodparameters? <')'>
	       typedmethodspec = Identifier <'('> methodparameters? <')'> typename
		 methodparameters = methodparam
				  | methodparameters <','> methodparam
		   methodparam = symbol ( Identifier)?
         implements = <#'\bimplements\b'> typename <#'\bfunc\b' '('> JavaIdentifier <')'> (
                          MethodImpl | <'('>  MethodImpl {<NL> MethodImpl}   <')'>
                        )
           <MethodImpl> = typedmethodimpl | untypedmethodimpl
             untypedmethodimpl = Identifier <'('>  parameters? <')'>
                                   (ReturnBlock|Blocky)
             typedmethodimpl = Identifier <'('>  parameters? <')'> typename
                                   (ReturnBlock|Blocky)
         functiondecl = <#'\bfunc\b'> (Identifier|operator) Function
         funclikedecl = <#'\bfunc\b' '<'> symbol <'>'> Identifier Function
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
                   <ReturnBlock> = <'{' #'\breturn\b'> expr <'}'>
         <Operand> = Literal | OperandName | label | islabel | new  | <'('> expr <')'> (*|MethodExpr*)
           label = #'\b\p{Lu}[\p{Lu}_\p{Nd}#\.]*\b'
	   islabel = <#'\bIS_'> #'\p{Lu}[\p{Lu}_\p{Nd}#\.]*\b'
           <Literal> = BasicLit | veclit | dictlit | setlit | structlit | functionlit | shortfunctionlit
             functionlit = <#'\bfunc\b'> Function
             shortfunctionlit = <#'\bfunc\b' '{'> expr <'}'>
             <BasicLit> = int_lit | bigintlit | regex | string_lit | rune_lit | floatlit | bigfloatlit (*| imaginary_lit *)

               (* http://stackoverflow.com/a/2509752/978525 *)
               floatlit = FloatLitA | FloatLitB
               <FloatLitA> = #'([0-9]+\.[0-9]*|\.[0-9]+)([eE][+-]?[0-9]+)?'
               <FloatLitB> = #'[0-9]+[eE][+-]?[0-9]+'

               bigfloatlit = (floatlit | int_lit) #'M\b'
               <int_lit> = decimallit | octal_lit | hexlit
		 decimallit = #'[1-9][0-9]*' | #'[0-9]'
		 <octal_lit>  = #'0[0-7]+'
		 hexlit    = <'0x'> #'[0-9a-fA-F]+'
               bigintlit = int_lit #'N\b'
               regex = #'/([^\/\n\\]|\\.)+/'
	       <string_lit> = interpretedstringlit | rawstringlit | clojureescape
                 interpretedstringlit = #'["“”](?:[^"\\]|\\.)*["“”]'
                 rawstringlit = <#'\x60'> #'[^\x60]*' <#'\x60'>     (* \x60 is back quote character *)
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
             NotType = #'\bfunc\b' | #'\bset\b' | prefix
             structlit = !NotType typename <'{'> ( expr {<','> expr} )? (<','> )? <'}'>
             setlit = <#'\bset\b' '{'> ( expr {<','> expr} )? <'}'>
           new = <#'\bnew\b'> typename
           <OperandName> = symbol | NonAlphaSymbol                           (*| QualifiedIdent*)
             <NonAlphaSymbol> = '=>' | '->>' | relop | addop | mulop | unary_op
                              | percentnum | percentvaradic
               percentnum     = <'$'> #'[1-9]'
               percentvaradic = <'$*'>
  <Ellipsis> = <'...'> | <'…'>
`,
	AUTO_WHITESPACE, whitespaceOrComments, // STANDARD, //whitespace,
	NO_SLURP, true,  // for App Engine compatibility
)
