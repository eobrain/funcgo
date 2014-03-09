(ns funcgo.core
  (:gen-class)
  (:require [instaparse.core :as insta]
            [clojure.string :as string])
  (:require clojure.pprint)
  )

(def funcgo-parser
     (insta/parser "
sourcefile     = [ NL ] packageclause <NL> expressions NL
packageclause  = <'package'> <__> dotted NL importdecl
expressions    = Expression { NL Expression }
importdecl     = <'import'> _ <'('>  NL { importspec NL } <')'>
importspec     = identifier _ dotted
<Expression>   = UnaryExpr | shortvardecl                       (* | Expression binary_op UnaryExpr *)
<UnaryExpr>    = PrimaryExpr                                                (* | unary_op UnaryExpr *)
<PrimaryExpr>  = functioncall | Operand | functiondecl | withconst
                                                                         (*Conversion |
                                                                         BuiltinCall |
                                                                         PrimaryExpr Selector |
                                                                         PrimaryExpr Index |
                                                                         PrimaryExpr Slice |
                                                                         PrimaryExpr TypeAssertion |*)
withconst      = <'const'> _ <'('> NL { const NL } <')'> NL expressions
const          = identifier _ <'='> _ Expression 
functioncall   = PrimaryExpr Call
<Call>         = <'('> _ ( ArgumentList _ )? <')'>
<ArgumentList> = expressionlist                                                      (* [ _ '...' ] *)
expressionlist = Expression { _ <','> _ Expression }
<Operand>      = Literal | OperandName | label                  (*| MethodExpr | '(' Expression ')' *)
<OperandName>  = symbol                                                           (*| QualifiedIdent*)
<Literal>      = BasicLit | dictlit | functionlit
<BasicLit>     = int_lit | string_lit | regex              (*| float_lit | imaginary_lit | rune_lit *)
shortvardecl   = identifier _ <':='> _ Expression
functiondecl   = <'func'> _ identifier _ Function
functionlit    = <'func'> _ Function
<Function>     = FunctionPart | functionparts
functionparts  = FunctionPart _ FunctionPart { _ FunctionPart }
<FunctionPart> = functionpart0 | functionpartn | vfunctionpart0 | vfunctionpartn
functionpart0  = <'('> _ <')'> _ <'{'> _ Expression _ <'}'>
vfunctionpart0 = <'('> _ varadic _ <')'> _ <'{'> _ Expression _ <'}'>
functionpartn  = <'('> _ parameters _ <')'> _ <'{'> _ Expression _ <'}'>
vfunctionpartn = <'('> _ parameters _  <','> _ varadic _ <')'> _ <'{'> _ Expression _ <'}'>
parameters     = identifier { <','> _ identifier }
varadic        = <'&'> identifier
dictlit        = '{' _ ( dictelement _ { <','> _ dictelement } )? _ '}'
dictelement    = Expression _ <':'> _ Expression
<int_lit>      = decimal_lit    (*| octal_lit | hex_lit .*)
decimal_lit    = #'[1-9][0-9]*'
regex          = <'/'> #'[^/]*'<'/'>   (* TODO: handle / escape *)
<string_lit>   = raw_string_lit   | interpreted_string_lit
raw_string_lit = <#'\\x60'> #'[^\\x60]*' <#'\\x60'>      (* \\x60 is back quote character *)
interpreted_string_lit = <#'\"'> #'[^\\\"]*' <#'\"'>      (* TODO: handle string escape *)
dotted         = identifier { <'.'> identifier }
symbol         = ( identifier <'.'> )? identifier
identifier     = #'[\\p{L}_][\\p{L}_\\p{Digit}]*'              (* letter { letter | unicode_digit } *)
label          = #'[\\p{Lu}]+'
letter         = unicode_letter | '_'
unicode_letter = #'\\p{L}'
unicode_digit  = #'\\p{Digit}'
<_>            = <#'[ \\t\\x0B\\f\\r\\n]*'> | comment  (* optional whitespace *)
__             =  #'[ \\t\\x0B\\f\\r\\n]+' | comment     (* whitespace *)
<NL>           = nl | comment
<nl>           = <#'\\s*[\\n;]\\s*'>       (* whitespace with at least one newline or semicolon *)
<comment>      = <#'[;\\s]*//[^\\n]*\\n\\s*'>
"))


(defn funcgo-parse [fgo]
  (let
      [parsed (funcgo-parser fgo)]
    ;;(clojure.pprint/pprint parsed)
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
      :functioncall   (fn
                        ([function]           (str "(" function ")"))
                        ([function call] (str "(" function " " call ")")))
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
      :withconst      (fn [& xs]
                         (let [
                             consts (butlast xs)
                             expressions (last xs)]
                           (str "(let ["
                                (reduce (fn [acc konst] (str acc " " konst)) consts)
                                "] "
                                expressions
                                ")")))
      :const          (fn [identifier expression] (str identifier " " expression))
      :symbol         (fn
                        ([identifier]        identifier)
                        ([package identifier] (str package "/" identifier)))
      :functiondecl   (fn [identifier function] (str "(defn " identifier " " function ")"))
      :functionlit    (fn [function] (str "(fn " function ")"))
      :functionparts   (fn [& functionpart]
                         (str "("
                              (reduce
                               (fn [acc fp] (str acc ") (" fp))
                               functionpart)
                              ")"))
      :functionpart0  (fn [expression]
                        (str "[] " expression))
      :vfunctionpart0 (fn [varadic expression]
                        (str "[" varadic "] " expression))
      :functionpartn  (fn [parameters expression]
                        (str "[" parameters "] " expression))
      :vfunctionpartn (fn [parameters varadic expression]
                        (str "[" parameters " " varadic "] " expression))
      :parameters     (fn [arg0 & args-rest]
                          (reduce
                           (fn [acc arg] (str acc " " arg))
                           arg0
                           args-rest))
      :varadic        (fn [parameter] (str "& " parameter))
      :dictlit        (fn [& dict-elems] (apply str dict-elems))
      :dictelement    (fn [key value] (str key " " value " "))
      :label          (fn [s] (str ":" (string/lower-case s)))
      :identifier     (fn [s] (clojure.string/replace s #"L"
                               (fn [s] (str "-" (clojure.string/lower-case s)))))
      :dotted         (fn [idf0 & idf-rest]
                        (reduce
                         (fn [acc idf] (str acc "." idf))
                         idf0
                         idf-rest))
      :decimal_lit    (fn [s] s)
      :regex          (fn [s] (str "#\"" s "\""))
      :interpreted_string_lit (fn [s] (str "\"" s "\""))
      :raw_string_lit (fn [s] (str "\"" (string/escape s char-escape-string) "\""))}
     parsed)))
  
(defn -main
  "Convert funcgo to clojure."
  [& args]
  (doseq
      [expr (read-string (str "["
                              (funcgo-parse (slurp (first args)))
                              "]"  ))]
    (clojure.pprint/pprint expr)
    (println)))
