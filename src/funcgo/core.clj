(ns funcgo.core
  (:require [instaparse.core :as insta]))

(def funcgo-parser
     (insta/parser "

SourceFile     = PackageClause _  { _ Expression _ NL }
PackageClause  = <'package'> <__> identifier NL  _ ImportDecl _ NL
<ImportDecl>   = <'import'> _ <'('>  { _ ImportSpec _ NL } <')'>
ImportSpec     = [ '.' | identifier ] ImportPath
ImportPath     = identifier { <'.'> identifier }
Expression     = UnaryExpr         (* | Expression binary_op UnaryExpr *)
UnaryExpr      = PrimaryExpr           (* | unary_op UnaryExpr *)
PrimaryExpr = Operand        (*|
	Conversion |
	BuiltinCall |
	PrimaryExpr Selector |
	PrimaryExpr Index |
	PrimaryExpr Slice |
	PrimaryExpr TypeAssertion |
	PrimaryExpr Call *)
Operand     = Literal        (*| OperandName | MethodExpr | '(' Expression ')' *)
Literal     = BasicLit       (*| CompositeLit | FunctionLit *)
BasicLit    = int_lit        (*| float_lit | imaginary_lit | rune_lit | string_lit *)
int_lit     = decimal_lit    (*| octal_lit | hex_lit .*)
decimal_lit = #'[1-9][0-9]*'

identifier     = #'[\\p{L}_][\\p{L}_\\p{Digit}]*'  (* letter { letter | unicode_digit } *)
letter         = unicode_letter | '_'
unicode_letter = #'\\p{L}'
unicode_digit  = #'\\p{Digit}'
<_>            = <#'[ \\t\\x0B\\f\\r]*'>   (* optional non-newline whitespace *)
__             = #'[ \\t\\x0B\\f\\r]+'     (* non-newline whitespace *)
<NL>           = [ nl | comment ]
<nl>           = <#'\\s*[\\n;]\\s*'>       (* whitespace with at least one newline or semicolon *)
comment        = #'//[^\\n]*\\n'
"))


(defn funcgo-parse [fgo]
  (insta/transform
   {
    :PackageClause (fn [identifier & import-decl] (str "(ns " (second identifier) ")")) }
   (funcgo-parser fgo)))
