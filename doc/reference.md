# Funcgo Reference

## Source File Structure

Source files names must end with either `.go` (if to be run on the
JVM) or `.gos` (if to be run in JavaScript).

A file called `_name_.go` must start with `package _name_` followed by
optional `import` statements, an optional `const` declaration, and a
list of expressions.

A minimal example:

```go
package hello
println("Hello World")
```

A slighttly longer example:

```go
package larger

import (
	test "midje/sweet"
)

const name = "Eamonn"

test.fact("can concatenate strings",
	str("Hello ",  name),     =>, "Hello Eamonn"
)

test.fact("can use infix when calling two-parameter-function",
	"Hello "  str  name,      =>, "Hello Eamonn"
)
```

