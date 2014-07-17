package fundamental_test
import (
        test "midje/sweet"
        fgoc "funcgo/main"
)

func parse(expr) {
	fgoc.CompileString("foo.go", "package foo;" str expr)
}

test.fact("func",
	parse("func foo(b,c){d;e}"),       =>, parse("\\defn-\\(foo,[b,c],do(d,e))"),
	parse("func Foo(b,c){d;e}"),       =>, parse("defn(Foo,[b,c],do(d,e))"),
	parse("func<defn> Foo(b,c){d;e}"), =>, parse("func Foo(b,c){d;e}")
)

test.fact("funcform",
	parse("func<something> a(b,c){d;e}"), =>, parse("something(a,[b,c],do(d,e))")
)

test.fact("if",
	parse("if a{b;c}"),          =>, parse("when(a,do(b,c))"),
	parse("if a{b;c}else{d;e}"), =>, parse("if(a,do(b,c),do(d,e))"),
	parse("if a{b}else{d;e}"),   =>, parse("if(a,b,do(d,e))"),
	parse("if a(b,c){d;e}"),     =>, parse("when(a(b,c),do(d,e))"),
	parse("if a(b){d;e}"),       =>, parse("when(a(b),do(d,e))")
)

test.fact("const",
	parse("{const(a=b;c=d)e;f}"), =>, parse("let([a,b,c,d],e,f)")
)

test.fact(":=",
	parse("{a:=b;c:=d;e;f}"),      =>, parse("let([a,b, c,d],  e,f)"),
	parse("{x,y:=a,b;e}"),         =>, parse("let([x,a, y,b],  e)"),
	parse("{x,y,z:=a,b,c;e}"),     =>, parse("let([x,a, y,b, z,c],  e)"),
	parse("{x,y,z,w:=a,b,c,d;e}"), =>, parse("let([x,a, y,b, z,c, w,d],  e)"),
	parse("{a,b,c,d,e,f,g,h,i,j:=1,2,3,4,5,6,7,8,9,10;v}"),
	=>, parse("let([a,1, b,2, c,3, d,4, e,5, f,6, g,7, h,8, i,9, j,10],  v)"),
)
