# Funcgo Reference

## Source File Structure

Source files names must end with either `.go` (if to be run on the
JVM) or `.gos` (if to be run in JavaScript).

A file called `name.go` must start with `package name` followed by
optional `import` statements, an optional `const` declaration, and a
list of expressions.

As a minimal example, the file `hello.go` contains
```go
package hello
println("Hello World")
```

As a slightly longer example, the file `larger.go` contains
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

## Everything is an Expression

Unlike Go, in Funcgo everything is an expression, including constructs
like if statements.

The following treats an if-else as an expression, setting `smaller`
to either `a` or `b`.
```go
		const smaller = if a < b {
			a
		} else {
			b
		}
```

And here the value returned from a for loop is actually the vector of
the values generated on each iteration.  (Called a _list
comprehension_ in some language.)

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

## Syntax

Unlike some languages, newlines can be significant in Funcgo.  This
happens when you have multiple expressions inside each construct.  For
example in the `if` statement below the two `println` expressions must
be separated by a newline.  (In this case we are ignoring the values
returned by the two expressions, the latter of which is returned by
the `if`.  Instead we are using these expressions for their
side-effects/

```go
		if a < b {
			println("Conclusion:")
			println(a, "is smaller than", b)
		}
	=>
Conclusion:
55 is smaller than 66
```

If you really want to, you can use semicolons instead of newlines, but
for readability I recommend you avoid semicolons.

```go
		if a < b { println("Conclusion:"); println(a, "is smaller than", b) }
	=>
Conclusion:
55 is smaller than 66
```
