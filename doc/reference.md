# Funcgo Reference

_[incomplete]_

## Source File Structure

Source files names must end with either `.go` (if to be run on the
JVM) or `.gos` (if to be run in JavaScript).

A file called `name.go` must start with `package name` followed by
optional `import` statements, an optional `const` declaration, and a
list of expressions.

```go
package hello
println("Hello World")
```
As a minimal example, the code above is in a file `hello.go`.

```go
package larger

import (
	test "midje/sweet"
)

name := "Eamonn"

test.fact("can concatenate strings",
	str("Hello ",  name),     =>, "Hello Eamonn"
)

test.fact("can use infix when calling two-parameter-function",
	"Hello "  str  name,      =>, "Hello Eamonn"
)
```
The above slightly longer example is in a file called `larger.go`.

## Most things are Expression

In Funcgo most things are expression, including constructs like `if`
statements that are statements in Go.

```go
		smaller := if a < b {
			a
		} else {
			b
		}
		smaller
	=> 55
```
The above treats an `if`-`else` as an expression, setting `smaller`
to either `a` or `b`.

```go
		digits  := [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
		squares := for d := lazy digits {
			d * d
		}
		squares
	=> [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
```
And here the value returned from a `for` loop is actually the vector
of the values generated on each iteration.  (Called a _list
comprehension_ in some language.)

## Syntax

Unlike some languages, newlines can be significant in Funcgo.  This
happens when you have multiple expressions inside curly braces or at
the top level of the source file.

```go
		if a < b {
			println("Conclusion:")
			println(a, "is smaller than", b)
		}
	=>
Conclusion:
55 is smaller than 66
```

For example in the `if` statement above the two `println` expressions
must be separated by a newline.  (In this case we are ignoring the
values returned by the two expressions, the latter of which is
returned by the `if`.  Instead we are using these expressions for
their side-effects:


```go
		if a < b { println("Conclusion:"); println(a, "is smaller than", b) }
	=>
Conclusion:
55 is smaller than 66
```

If you really want to, you can use semicolons instead of newlines as
shown above, but for readability I recommend you avoid semicolons.

## Imports

You can directly use anything provided by the [clojure.core][1] API
without further specification.  However if you want to use anything
from any other library you have to explicitly import it at the top of
your file.  Depending on what you are you importing you use one of
these forms.

1. `import` (the most common case) for Clojure or Funcgo libraries

   ```go
   import (
           test "midje/sweet"
           fgo "funcgo/core"
           fgoc "funcgo/main"
           "clojure/string"
   )
   ...
   test.fact(...
   ...
   fgo.Parse(...
   ...
   string.trim(...

   ```

    As shown above an `import` statement can import multiple Clojure
    or Funcgo libraries.  It specifies the library as a string of
    slash-separated identifiers. Each library can be preceded by a
    short name by which the library is referred to in the body of the
    code. If no short name is specified, then the last identifier in
    the library name is used (for example `string` in the last example
    above).

    In the body of the code any function or variable referenced from
    an imported library must be qualified by short name.

    ```go
    import(
        _ "hiccups/runtime"
        "fgosite/code"
    )
    ```

    Sometimes you want to import a library only for the side-effect of
    importing it. To avoid getting a compile error complaining about
    an unused import, you can use `_` as shown above.

    ```go
    import "clojure/string"
    ```

    If you are only importing a single library you can use a short
    form without parentheses as shown above.

1. `import type` for JVM classes and interfaces

    ```go
    import type (
        java.io.{BufferedWriter, File, StringWriter}
        jline.console.ConsoleReader
    )
    ...
	... = new StringWriter()
    ...
    func compileTree(root File, opts) {
    ...
    ```

    JVM types, such as defined in Java (and sometimes in Clojure),
    have a different syntax for importing as show above.  Each type
    must be explicitly listed, though types in the same package can
    expressed using the compressed syntax shown above for `java.io`.

    Once imported, such types can be simply referenced by name,
    without qualification.

    Types from the base `java.lang` API do not need to be imported.

1. `import macros` (when targeting JavaScript only) for importing
ClojureScript macros

    ```go
    import macros (
        hiccups "hiccups/core"
    )
    ...
    func<hiccups.defhtml> pageTemplate(index) {
    ...
    ```

    When targeting the JavaScript runtime you sometimes need to import
    macro definitions in a special way as shown above.

1. `import extern` (advanced use only) needed when creating macros

    ```go
    import extern(
        produce
        bakery
    )
    ...
    ... quote(produce.onions) ...
    ...
    ```

    Occasionally you will need to refer to symbols in libraries that
    you cannot import.  As shown above you can declare them as
    `extern` libraries.

## Const

In Funcgo you should use constants for any value that is
set once and never changed.

```go
...
{
	cljText   := core.Parse(inPath, fgoText, EXPR)
	strWriter := new StringWriter()
	writer    := new BufferedWriter(strWriter)
	cljText writePrettyTo writer
	strWriter->toString()
}
```
As shown above, constants are defined using the `:=` operator.  

There can only be a single contiguous group of contant declarations in
each _block_ of expressions, and they must appear at the top of the
block.  A block is either to top-level code if a file after the
`import` statements, or some newline-separated expressions surrounded
in curly braces.  The constants you define in a block can only be used
inside that block.

```go
...
{
	const (
		cljText = core.Parse(inPath, fgoText, EXPR)
		strWriter = new StringWriter()
		writer = new BufferedWriter(strWriter)
	)
	cljText writePrettyTo writer
	strWriter->toString()
}
```

As shown above, there is also an alternative syntax using the `const`
keyword.  It too must be at the beginning of a block.

```go
...
{
	const consoleReader = new ConsoleReader()
	consoleReader->setPrompt("fgo=>     ")
	consoleReader
}
```

If there is a just a single constant, you can drop the parentheses.

## Looping with tail recursion

First, lets look at an ordinary (non-tail) recursion

```go
		func sumSquares(vec) {
			if isEmpty(vec) {
				0
			} else {
				x := first(vec)
				x * x + sumSquares(rest(vec))
			}
		}
		sumSquares([3, 4, 5, 10])
	=> 150
```

The above example shows the `sumSquares` function that returns the sum
of squares of a vector of numbers.  It is implemented as the square of
the first element plus the recursive sum of squares of the rest of the
vector.  This works fine for small vectors but for large vectors it
could cause an infamous _stack overflow_ exception.

```go
		func sumSquares(vec) {
			func sumSq(accum, v) {
				if isEmpty(v) {
					accum
				} else {
					x := first(v)
					recur(accum + x * x, rest(v))
				}
			}
			sumSq(0, vec)
		}
		sumSquares([3, 4, 5, 10])
	=> 150
```

The above example avoids this stack overflow by using the special
`recur` syntax to recursively call the containing function.  However
`recur` must be in _tail position_, which means that the function
needs to be re-arranged to add an inner recursive function that passes
down as accumulator variable.  This version can be called on
arbitrarily long vectors without blowing your stack.

There is also an equivalent way of getting the same result using the
`loop` construct.

```go
		func sumSquares(vec) {
			loop(accum=0, v=vec) {
				if isEmpty(v) {
					accum
				} else {
					const x int = first(v)
					recur(accum + x * x, rest(v))
				}
			}
		}
		sumSquares([3, 4, 5, 10])
	=> 150
```

The `loop` construct declares a set of iteration variables and sets
their initial values.  The `recur` calls the nearest enclosing `loop`
passing in updated iteration variables (which are actually constants
in each iteration).  The number of parameters in the `recur` must match the
number of parameters in the `loop`.

```go
	loop(vec=[], count = 0) {
		if count < 10 {
			v := vec  conj  count
			recur(v, count + 1)
		} else {
			vec
		}
	=> [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
```

And above is another simpler example of using `loop`, starting with an
empty vector and using the `conj` operator to add numbers to it.

## Curly Brace Block

Everywhere you can put an expression you can put a newline-separated
sequence of expressions in a curly braces block.  The result of the
last expression is returned as the result of the block.

```go
		product := {
			log->info("doing the multiplication")
			100 * 100
		}
		product
	=> 10000
```

Above is an example of the `product` constant being assigned the value
of the block, with the multiplication expression being preceded by a
logging statement that is executed only for its side-effects.

# Switch

There are three forms of switch statement.

```go
				switch count(remaining) {
				case 1: {
					[expr] := remaining
					str(acc, " :else ", expr, ")")
				}
				case 2:
					typeCase()  str  ")"
				default:
					recur(typeCase(), 2  drop  remaining)
				}
```

In the first form, shown above, the switch takes an expression and
matches execute whichever of its `case` sections match the result of
the expression.  This is the more efficient form of switch because the
dispatch to a case happens in constant time, but it has the restriction
that the `case` sections must have compile-time constants values.

```go
			switch {
			case isNil(t):
				new TreeNode(v, nil, nil)
			case v < VAL(t):
				new TreeNode(VAL(t), L(t)  xconj  v, R(t))
			default:
				new TreeNode(VAL(t), L(t), R(t)  xconj  v)
			}
```

The second form, shown above, is more general.  There is no expression
beside the `switch` but instead each `case` has an arbitrary Boolean
expression.  In general this form is slower because the dispatch
happens in linear time, each case expression being evaluated in turn
until one returns true.

The third form is the _type switch_ using the `.(type)`
suffix to indicate that we are switching on the type, and using
type names in the case statements.

```go
		func plus(a, b) {
			switch a.(type) {
			case Number:   a + b
			case String:   a  str  b
			case Iterable: vec(a  concat  b)
			default:       str("Unknown types for ", a, " and ", b)
			}
		}

		[
			2       plus  3,
			0.5     plus  0.75,
			[P, Q]  plus  [R, S, T],
			"foo"   plus  "bar",
                        FOO     plus  BAR
		]

	=> [
		5,
		1.25,
		[P, Q, R, S, T],
		"foobar",
		"Unknown types for :foo and :bar"
	]
```

In the above example we define a _plus_ function that does different
operations depending on the types of the first argument.  (A more
robust version would check both arguments.)

## Java Statics

To use static methods or fields from a Java class you use the `::`
operator after the Java class name (which you should import using the
`import type` syntax unless it is in the java.lang package).

```go
	2 * Double::MAX_VALUE          // => Double::POSITIVE_INFINITY
	Integer::parseInt("-42")       // => -42
	Math::round(2.999)             // => 3
	13 Integer::toString 2         // => "1101"
```

The first example shows how a static field is uses.  The remaining
three are static method invocations, with the last one showing how you
can use infix notation to invoke a static method that takes two
parameters.

## Identifiers

Clojure allows characters in identifiers that are not allowed in
Funcgo identifiers, therefore to allow inter-operation Funcgo
identifiers are mangled like so:

* camel-case is converted to dash-separated:
  * `fooBarBaz` &rarr; `foo-bar-baz`
  * `FooBarBaz` &rarr; `Foo-bar-baz`
* `is` prefix is converted to `?` suffix
  * `isEqual` &rarr; `equal?`
* `mutate` prefix is converted to `!` suffix
  * `mutateSort` &rarr; `sort!`
* underscore prefix is converted to dash prefix
  * `_main` &rarr; `-main`

However, identifiers referring to Java entities are *not* mangled.
These are any identifiers in `import type` statements, anything before
or after a `::` and anything after a `->`.

And finally you can avoid mangling by using a backslash and back-quotes:

```go
        const origDispatch = \`pprint/*print-pprint-dispatch*`
```

The above example uses this escaped identifier syntax to refer to the
`pprint/*print-pprint-dispatch*` Clojure identifier which has the
"earmuff" characters not allowed by Funcgo.

## Vars

If possible you should use consts because they are immutable, but
sometimes you need vars.  These are thread-local mutable storage.

The case of the name is significant.  If it begins with an upper-case
letter then it is exported and visible globally, otherwise it is
private to the file it is declared in.


```go
var initialBoard = [
	[EE, KW, EE],
	[EE, EE, EE],
	[EE, KB, EE]
]
```

They use a syntax similar to the `const` construct.  They can be
without parentheses as shown above.

```go
		var (
			pp = 111
			qq = 222
		)
		pp + qq
	=> 333
```

Or it can use the grouped version of the syntax as shown
above.

```go
		var pp = 111
		var qq = 222
		pp + qq
	=> 333
```

Alternatively each var can be separately declared as shown above.

Unlike consts, var declarations do not have to be at the beginning of
a curly-bracket block.

```go
		var tt int    = 111
		var uu string = "foo"
		uu  str  tt
	=> "foo111"
```

If you want you can add type hints as shown above.

## If-Else

```go
		filename := if isJvm { "main.go" } else { "main.gos" }
```
The above example shows the if-else expression.

```go
	if cmdLine(ERRORS) {
		println(cmdLine(ERRORS))
	}
```

If the `else` part is omitted the `if` expression returns nil when
the condition is false, though in the example above this is ignored.

```go
			if met := meta(o); met {
				print("^")
				if count(met) == 1 {
					if met(TAG) {
						origDispatch(met(TAG))
					} else {
						if met(PRIVATE) == true {
							origDispatch(PRIVATE)
						} else {
							origDispatch(met)
						}
					}
				} else {
					origDispatch(met)
				}
				print(" ")
				pprint.pprintNewline(FILL)
			}
```

Finally, the first line of the example above shows another format,
where a constant `met` is set and tested as part of the `if` line.
This constant can be used in the body of the if expression.  The above
example also emphasizes the fact that there is no "else-if" construct
-- you must use nested if-else expressions (or alternatively use the
switch expression).

## For loops

There are three types of `for` expression.

```go
		fib        := [1, 1, 2, 3, 5, 8]
		fibSquared := for x := lazy fib {
			x * x
		}
		fibSquared
	=> [1, 1, 4, 9, 25, 64]
```

The "lazy" version returns a sequence that is the same length as the
input sequence (given after `lazy`), with the body of the loop being
executed for each member of the input sequence.

```go
		fib        := [1, 1, 2, 3, 5, 8]
		fibSquared := func(x){ x * x }  map  fib
		fibSquared
	=> [1, 1, 4, 9, 25, 64]
```

As an aside, you can get the same result as the lazy `for` using the
`map` function as shown above.

```go
		fib := [1, 1, 2, 3, 5, 8]
		for x := lazy fib {
			print(" ", x)
		}
	=> ""
```

The reason that this construct is called "lazy", is shown in the example
above where the body of the `for` does not return a value, but instead
has a side-effect (writing on the console).  In this case the `print`
is *not* executed.

```go
		fib := [1, 1, 2, 3, 5, 8]
		for x := range fib {
			print(" ", x)
		}
	=> "  1  1  2  3  5  8"
```

To cause such a side-effect body to be executed you can use the
"range" form of the for loop as shown above.  It is not lazy, but will
execute the body of the loop for each member of the input sequence.

```go
		for x := times 10 {
			print(" ", x)
		}
	=> "  0  1  2  3  4  5  6  7  8  9"
```

The final form of the for loop is the "times" version, which executes
its body the number of times specified after `times` as shown above.

## Exceptions

Funcgo supports exceptions in a way similar to Java.

```go
		eval := try{
			main := loadString(clj(id))
			withOutStr(main())
		} catch Throwable e {
			str(e)
		}
```

The above example shows an example of catching an exception.  A
difference from Java is that try-catch is an expression, thus in the
above case if the exception is caught the `eval` constant will be set
to the value of `str(e)`.

```go
		try {

			throw(new AssertionError("foo"))

		} catch OutOfMemoryError e {
			"out of memory"
		} catch AssertionError e {
			"assertion failed: "  str  e->getMessage()
		} finally {
			"useless"
		}
	=> "assertion failed: foo"
```

The above example shows a more complete example, including how you can
throw your own exceptions.  Note, that although the `finally`
expression is evaluated, its result is ignored so it is not useful in
this case.

```go
				<-mutex   // grab mutex
				try {
					i := dangerous->get(0)
					dangerous->set(0, i + 1)
				} finally {
					mutex <- true   // release mutex
				}
```

Above is an example of a more useful applicaton of `finally` where we
are depending on the side-effect of evaluating its expression.

## Asynchronous Channels

```go
		c1 := make(chan, 1)
		c2 := make(chan, 1)
		thread {
			Thread::sleep(10)
			c1 <- 111
		}
		c2 <- 222
		select {
		case x = <-c1:
			x * 100
		case x = <-c2:
			x * 100
		}
	=> 22200
```

The example above uses the same syntax as Go, where for a channel `c`
the operation `<-c` is taking from a channel (blocking if necessary
until input arrives) and `c <- x` is sending the value `x` to the
channel.

The `select` construct allows you to block on multiple asynchronous
channel operations, such that the first one that unblocks will
activate.

When targeting JavaScript however we are restricted because
JavaScript is single-threaded, instead you have to use a different
syntax (using `<:` instead of `<-') for channel operations,

```go
		c1 := make(chan, 1)
		c2 := make(chan, 1)
		go {
			for i := times(10000) { var x = i }
			c1 <: 111
		}
		go {
			c2 <: 222
		}
		<-go {
			select {
			case x = <:c1:
				x * 100
			case x = <:c2:
				x * 100
			}
		}
	=> 22200
```

These operations are restricted to being directly inside a `go {
... }` block as shown above.

```go
		c1 := make(chan, 1)
		c2 := make(chan)
		thread {
			Thread::sleep(20)
			c1 <- 111
		}
		thread {
			Thread::sleep(10)
			<-c2
		}
		select {
		case x = <-c1:
			x * 100
		case c2 <- 222:
			"wrote to c2"
		default:
			"nothing ready"
		}
	=> "nothing ready"
```

Finally the example above shows some more features.  The `thread`
block is like a `go` block except that it fires up a real thread
rather than a goroutine, thus it can only be used when targeting the
JVM.

If there is a `default` clause in the `select` and all the `case`
clauses are blocked, then it will execute instead.

Finally note that this example has both types of `case` clauses, those
writing to channels and those reading from channels, both of which can
block.

## Infix functions

```go
	str("foo", "bar")
	=> "foobar"
```

You can call a function of two arguments in the normal prefix format
of `f(a,b)`.

```go
	"foo"  str  "bar"
	=> "foobar"
```

Alternatively you can call such a function in an infix format of `a  f
b`.

## Binary Operators and Precedence

This table shows all the built-in operators and how they group.  The
ones at the top bind most tightly.

6. unary expression
5. `*` `/` `%` `<<` `>>` `&` `&^`
4. `+` `-` `|` `^`
3. `==` `!=` `<` `>` `<=` `>=`
2. `&&`
1. `||`
0. inline function call

```go
	^a * b          // => (^a) * b,
	a * b - c       // => (a * b) - c
	a + b < c       // => (a + b) < c
	a < b && b < c  // => (a < b) && (b < c)
	p && q || r     // => (p && q) || r
	p || q  str  r  // => (p || q)  str  r
```

## Destructuring

You can declare multiple constants on the left-hand-side of the `=`
and put a vector on the right-hand-side.  Thus "unpacks" the vector
assigning each element to the corresponding constant.

```go
		vec          := [111, 222, 333, 444]
		[a, b, c, d] := vec
		b
    => 222
```

For example, above we unpack the vector `vec`, so that constant `b`
ends up with the value `222`.

```go
		func second([a, b, c, d]) {
			b
		}
		second(vec)
	=> 222
```

This also works for function arguments, where above we have used a
function to extract the second element from the vector.

```go
		vec              := [111, 222, 333, 444]
		[first, rest...] := vec
		rest
	=> [222, 333, 444]
```

For variable-length vectors you can use ellipses `...` after the
constant to match it to the remaining part of the vector.  So for
example, above `first` gets the the first element in the vector and
`rest` gets the remaining elements.

```go
		dict             := {AAA: 11,  BBB: 22,  CCC: 33,  DDD: 44}
		{c: CCC, a: AAA} := dict
		c
	=> 33
```

You can also destructure dicts using the syntax shown above, where on
the left-hand-side each match is specified as _constant_`:` _key_.

```go
		func extractCCC({c: CCC}) {
			c
		}
		extractCCC(dict)
	=> 33
```

Dict destructuring also works in function parameters as shown above.

```go
		planets := [
			{NAME: "Mercury", RADIUS_KM: 2440},
			{NAME: "Venus",   RADIUS_KM: 6052},
			{NAME: "Earth",   RADIUS_KM: 6371},
			{NAME: "Mars",    RADIUS_KM: 3390}
		]
		[_, _, {earthRadiusKm: RADIUS_KM}, _] := planets
		earthRadiusKm
	=> 6371
```

You can nest these destructurings to any depth.  For example the above
example plucks the `earthRadiusKm` constant from two-levels down inside
a vector of dicts.  We are using the convention of using the `_`
identifier for unused values.

## Mutable State

## Quoting and Unquoting

## Invoking Functions

## Java Object Fields and Methods

## Defining Functions

## Closures

## Function-Like Macros

## Interfaces

## Structs

## new

## Labels

## Literals

## Regular Expressions

## Vectors

## Dictionaries

## Type Hints

Funcgo is a _gradually typed_ language.  Unlike Go, you do not need to
specify any types, and in most cases the Clojure runtime can figure
them out. However sometimes you may get a runtime warning from the
Clojure runtime that it is using _reflection_ because it cannot figure
out your types.  To allow your code to run more efficiently you can
add types using the same syntax as the Go language.

In practice, you usually only need to add types in a very few places
in your code.

```go
	consoleReader ConsoleReader := newConsoleReader()
```

Above is an example of a constant being declared of `ConsoleReader`
type so future uses of the constant are more efficient.

```go
func compileFile(inFile File, root File, opts) {
```

And above is an example of the first two of the three function
parameters being declared to be of type `File`.

[1]: http://clojure.github.io/clojure/
