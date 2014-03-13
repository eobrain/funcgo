(ns funcgo.core
  (:gen-class)
  (:require [instaparse.core :as insta]
            [clojure.string :as string])
  (:require [instaparse.failure :as failure])
  (:require clojure.pprint)
  )

(def funcgo-parser
     (insta/parser "
sourcefile     = [ NL ] packageclause _ expressions _
packageclause  = <'package'> <__> dotted NL importdecl
expressions    = Expression { NL Expression }
importdecl     = <'import'> _ <'('>  _ { importspec _ } <')'>
importspec     = identifier _ dotted
<Expression>   = UnaryExpr | shortvardecl | ifelseexpr | tryexpr | forrange | forlazy | fortimes (* | Expression binary_op UnaryExpr *)
ifelseexpr     = <'if'> _ Expression _ ( ( block _ <'else'> _ block ) | ( _ <'{'> _ expressions _ <'}'> )   )
forrange       = <'for'> <__> identifier _ <':='> _ <'range'> <_> Expression _ <'{'> _ expressions _ <'}'>
forlazy        = <'for'> <__> identifier _ <':='> _ <'lazy'> <_> Expression [ <__> <'if'> <__> Expression ] _ <'{'> _ expressions _ <'}'>
fortimes       = <'for'> <__> identifier _ <':='> _ <'times'> <_> Expression _ <'{'> _ expressions _ <'}'>
tryexpr        = <'try'> _ <'{'> _ expressions _ <'}'> _ catches _ finally?
catches        = ( catch { _ catch } )?
catch          = <'catch'> _ identifier _ identifier _ <'{'> _ expressions _ <'}'>
finally        = <'finally'> _ <'{'> _ expressions _ <'}'>
block          = <'{'> _ Expression { NL Expression } _ <'}'>
<UnaryExpr>    = PrimaryExpr                                                (* | unary_op UnaryExpr *)
<PrimaryExpr>  = functioncall | Operand | functiondecl | withconst
                                                                         (*Conversion |
                                                                         BuiltinCall |
                                                                         PrimaryExpr Selector |
                                                                         PrimaryExpr Index |
                                                                         PrimaryExpr Slice |
                                                                         PrimaryExpr TypeAssertion |*)
withconst      = <'const'> _ <'('> _ { consts } _ <')'> _ expressions
consts         = [ const { NL const } ]
const          = identifier _ <'='> _ Expression 
functioncall   = PrimaryExpr Call
<Call>         = <'('> _ ( ArgumentList _ )? <')'>
<ArgumentList> = expressionlist                                                      (* [ _ '...' ] *)
expressionlist = Expression { _ <','> _ Expression }
<Operand>      = Literal | OperandName | label | new            (*| MethodExpr | '(' Expression ')' *)
new            = <'new'> <__> symbol
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
symbol         = ( identifier <'.'> )? !Keyword identifier
Keyword        = ( 'for' | 'range' )
identifier     = #'[\\p{L}_][\\p{L}_\\p{Digit}]*'              (* letter { letter | unicode_digit } *)
label          = #'\\p{Lu}[\\p{Lu}_]*'
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
    (if (insta/failure? parsed)
      (do
        (failure/pprint-failure parsed)
        (throw (Exception. "\"SYNTAX ERROR\"")))
      ;;(do
      ;;(clojure.pprint/pprint parsed)
      (insta/transform
       {
        :sourcefile     (fn [header body] (str header body "\n"))
        :packageclause  (fn [dotted  import-decl]
                          (str "(ns " dotted import-decl ")\n\n"))
        :importdecl     (fn [ & import-specs] (apply str import-specs))
        :importspec     (fn [identifier dotted]
                          (str "\n  (:require [" dotted " :as " identifier "])"))
        :ifelseexpr (fn
                      ([condition exprs] (str "(when " condition " " exprs ")"))
		      ([condition block1 block2] (str "(if " condition " " block1 " " block2 ")")))
        :forrange   (fn [identifier seq expressions] 
                                (str "(doseq ["  identifier " " seq "] " expressions ")"))
        :forlazy    (fn
                      ([identifier seq expressions] 
                                (str "(for ["  identifier " " seq "] " expressions ")"))
                      ([identifier seq condition expressions] 
                                (str "(for ["  identifier " " seq "] :when " condition " " expressions ")")))
        :fortimes   (fn [identifier count expressions] 
                                (str "(dotimes ["  identifier " " count "] " expressions ")"))
        :tryexpr (fn
                   ([expressions catches] (str "(try " expressions " " catches ")"))
		   ([expressions catches finally] (str "(try " expressions " " catches " " finally ")")))
        :catches (fn [& catches]
                   (reduce
                    (fn [acc catch] (str acc " " catch))
                    catches)
                   )
        :catch (fn [typ exception expressions] 
                 (str "(catch " typ " " exception " " expressions ")")
                 )
        :finally (fn [expressions] (str "(finally " expressions ")"))
        :new        (fn [symbol] (str symbol "."))
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
                           (fn [acc expr] (str acc " " expr))
                           expr0
                           expr-rest))
        :consts         (fn [& consts]
                          (reduce
                           (fn [acc const] (str acc "\n" const))
                           consts))
        :block          (fn
                          ([expr] expr)
                          ([expr0 & expr-rest]
                             (str "(do "
                                  (reduce
                                   (fn [acc expr] (str acc " " expr))
                                   expr0
                                   expr-rest)
                                  ")")))
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
        :identifier     (fn [s]
                          (clojure.string/replace
                           s
                           #"\p{Ll}\p{Lu}"
                            (fn [s] (str (first s) "-" (clojure.string/lower-case (last s))))))
        :dotted         (fn [idf0 & idf-rest]
                          (reduce
                           (fn [acc idf] (str acc "." idf))
                           idf0
                           idf-rest))
        :decimal_lit    (fn [s] s)
        :regex          (fn [s] (str "#\"" s "\""))
        :interpreted_string_lit (fn [s] (str "\"" s "\""))
        :raw_string_lit (fn [s] (str "\"" (string/escape s char-escape-string) "\""))}
       parsed))))
  
(defn -main
  "Convert funcgo to clojure."
  [& args]
  (try
    (let
        [clj (funcgo-parse (slurp (first args)))]
      ;;(println clj)
      (doseq
          [expr (read-string (str "[" clj "]"  ))]
        (clojure.pprint/pprint expr)
        (println)))
    (catch Exception e
      (println "\n" (.getMessage e)))))
