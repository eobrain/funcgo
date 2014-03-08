(ns funcgo.core
  (:gen-class)
  (:require [instaparse.core :as insta]
            [clojure.string :as string]))

(def funcgo-parser
     (insta/parser "
sourcefile     = [ NL ] packageclause NL expressions NL
packageclause  = <'package'> <__> dotted NL  _ importdecl _ NL
expressions    = Expression { NL Expression }
importdecl     = <'import'> _ <'('>  NL { _ importspec _ NL } <')'>
importspec     = identifier _ dotted
<Expression>   = UnaryExpr | shortvardecl                       (* | Expression binary_op UnaryExpr *)
<UnaryExpr>    = primaryexpr                                                (* | unary_op UnaryExpr *)
primaryexpr    = Operand | functiondecl |
                                                                         (*Conversion |
                                                                         BuiltinCall |
                                                                         PrimaryExpr Selector |
                                                                         PrimaryExpr Index |
                                                                         PrimaryExpr Slice |
                                                                         PrimaryExpr TypeAssertion |*)
	         primaryexpr Call
<Call>         = <'('> _ [ ArgumentList ] _ <')'>
<ArgumentList> = expressionlist                                                      (* [ _ '...' ] *)
expressionlist = Expression { _ <','> _ Expression }
<Operand>      = Literal | OperandName | label                  (*| MethodExpr | '(' Expression ')' *)
<OperandName>  = symbol                                                           (*| QualifiedIdent*)
<Literal>      = BasicLit | dictlit | functionlit
<BasicLit>     = int_lit | string_lit                      (*| float_lit | imaginary_lit | rune_lit *)
shortvardecl   = identifier _ <':='> _ Expression
functiondecl   = <'func'> _ identifier _ function
functionlit    = <'func'> _ function
function       = <'('> _ parameters _ <')'> _ <'{'> _ Expression _ <'}'>
parameters     = ( identifier [ <','> _ identifier ]  )? ( _ varadic)?
varadic        = identifier _ <'...'>
dictlit        = '{' _ ( dictelement _ { <','> _ dictelement } )? _ '}'
dictelement    = Expression _ <':'> _ Expression
<int_lit>      = decimal_lit    (*| octal_lit | hex_lit .*)
decimal_lit    = #'[1-9][0-9]*'
<string_lit>   = raw_string_lit   | interpreted_string_lit
raw_string_lit = <#'\\x60'> #'[^\\x60]*' <#'\\x60'>      (* \\x60 is back quote character *)
interpreted_string_lit = <#'\"'> #'[^\\\"]*' <#'\"'>      (* TODO: handle string escape *)
dotted         = identifier { <'.'> identifier }
symbol         = ( identifier <'.'> )? identifier
<identifier>   = #'[\\p{L}_][\\p{L}_\\p{Digit}]*'              (* letter { letter | unicode_digit } *)
label          = #'[\\p{Lu}]+'
letter         = unicode_letter | '_'
unicode_letter = #'\\p{L}'
unicode_digit  = #'\\p{Digit}'
<_>            = <#'[ \\t\\x0B\\f\\r\\n]*'>   (* optional whitespace *)
__             =  #'[ \\t\\x0B\\f\\r\\n]+'    (* whitespace *)
<NL>           = [ nl | comment ]
<nl>           = <#'\\s*[\\n;]\\s*'>       (* whitespace with at least one newline or semicolon *)
comment        = #'//[^\\n]*\\n'
"))


(defn funcgo-parse [fgo]
  (insta/transform
   {
    :sourcefile     (fn [header body] (str header body "\n"))
    :packageclause  (fn [dotted  import-decl]
                      (str "(ns " dotted import-decl ")\n\n"))
    :importdecl     (fn [ & import-specs] (apply str import-specs))
    :importspec     (fn [identifier dotted]
                      (str "\n  (:require [" dotted " :as " identifier "])"))
    :shortvardecl   (fn [identifier expression]
                    (str "(def " identifier " " expression ")"))
    :primaryexpr    (fn
                      ([operand]           operand)
                      ([primary-expr call] (str "(" primary-expr " " call ")")))
    :expressionlist (fn [expr0 & expr-rest]
                      (reduce
                       (fn [acc expr] (str acc " " expr))
                       expr0
                       expr-rest))
    :expressions    (fn [expr0 & expr-rest]
                      (reduce
                       (fn [acc expr] (str acc "\n\n" expr))
                       expr0
                       expr-rest))
    :symbol         (fn
                      ([identifier]        identifier)
                      ([package identifier] (str package "/" identifier)))
    :functiondecl   (fn [identifier function] (str "(defn " identifier function ")"))
    :functionlit    (fn [function] (str "(fn" function ")"))
    :function       (fn [parameters expression] (str " [" parameters "]\n  " expression))
    :parameters     (fn [& args]
                      (when (seq args)
                        (reduce
                         (fn [acc arg] (str acc " " arg))
                         args)))
    :dictlit        (fn [& dict-elems] (apply str dict-elems))
    :dictelement    (fn [key value] (str key " " value " "))
    :label          (fn [s] (str ":" (string/lower-case s)))
    :dotted         (fn [idf0 & idf-rest]
                      (reduce
                       (fn [acc idf] (str acc "." idf))
                       idf0
                       idf-rest))
    :decimal_lit    (fn [s] s)
    :interpreted_string_lit (fn [s] (str "\"" s "\""))
    :raw_string_lit (fn [s] (str "\"" (string/escape s char-escape-string) "\""))}
   (funcgo-parser fgo)))

(defn -main
  "Convert funcgo to clojure."
  [& args]
  (println (funcgo-parse (slurp (first args)))))