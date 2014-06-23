# Funcgo Reference

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

const (
	name = "Eamonn"
)

test.fact("can concatenate strings",
	str("Hello ",  name),     =>, "Hello Eamonn"
)

test.fact("can use infix when calling two-parameter-function",
	"Hello "  str  name,      =>, "Hello Eamonn"
)
```
The above slightly longer example is in a file called `larger.go`.

## Everything is an Expression

Unlike Go, in Funcgo everything is an expression, including constructs
like `if` statements.

```go
		const smaller = if a < b {
			a
		} else {
			b
		}
```
The above treats an `if`-`else` as an expression, setting `smaller`
to either `a` or `b`.

```go
		const (
			digits = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
			squares = for d := lazy digits {
				d * d
			}
		)
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

In Funcgo you should use `const` declarations for any value that is
set once and never changed.

```go
... {
	const (
		cljText = core.Parse(inPath, fgoText, EXPR)
		strWriter = new StringWriter()
		writer = new BufferedWriter(strWriter)
	)
	cljText writePrettyTo writer
	strWriter->toString()
}
```

There can only be a single `const` section in each _block_ of
expressions, where a block is either to top-level code if a file after
the `import` statements, or some newline-separated expressions
surrounded in curly braces.  The constants you define can only be used
inside that block.

```go
... {
	const consoleReader = new ConsoleReader()
	consoleReader->setPrompt("fgo=>     ")
	consoleReader
}
```

If there is a just a single constant, you can drop the parentheses.

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
	const consoleReader ConsoleReader = newConsoleReader()
```

Above is an example of a constant being declared of `ConsoleReader`
type so future uses of the constant are more efficient.

```go
func compileFile(inFile File, root File, opts) {
```

And above is an example of the first two of the three function
parameters being declared to be of type `File`.

[1]: http://clojure.github.io/clojure/
