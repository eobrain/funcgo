package funcgo/fundamental_test
import (
        test "midje/sweet"
        fgo "funcgo/core"
)

func parse(expr) {
	fgo.funcgoParse("foo.go", "package foo;" str expr)
}

test.fact("func",
	parse("func a(b,c){d;e}"), =>, parse("defn(a,[b,c],d,e)"),
	parse("func<defn> a(b,c){d;e}"), =>, parse("func a(b,c){d;e}")
)

test.fact("funcform",
	parse("func<something> a(b,c){d;e}"), =>, parse("something(a,[b,c],d,e)")
)

test.fact("if",
	parse("if a{b;c}"),          =>, parse("when(a,b,c)"),
	parse("if a{b;c}else{d;e}"), =>, parse("if(a,{b;c},{d;e})"),
	parse("if a{b}else{d;e}"),   =>, parse("if(a,{b},{d;e})"),
	parse("if a(b,c){d;e}"),     =>, parse("when(a(b,c),d,e)"),
	parse("if a(b){d;e}"),       =>, parse("when(a(b),d,e)")
)

test.fact("const",
	parse("const(a=b;c=d){e;f}"), =>, parse("let([a,b,c,d],e,f)")
)
