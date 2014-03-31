# funcgo

Funcgo is a compiler that converts Functional Go into Clojure.

The compiler itself is written in Functional Go. (Turtles all the way down!)

## Introduction to the Funcgo Language

### Why a new language?

The goal of Funcgo is to combine the readability of the Go language
with the semantics of Clojure.

1. Go is a language that has been well designed to be very
readable. However it is best as a low-level system programming
language (replacing C) and it is missing many of the higher-level
features that programmers expect for working further up the stack, for
example in web applications.

2. Clojure is a variety of Lisp that inter-operates with Java or
JavaScript.  It encourages a functional programming style with
efficient immutable containers, combined with a thread-safe model for
mutable state called software transactional memory. However, for
programmers unfamiliar with Lisp syntax Clojure is very difficult to
read.

### Examples

In this section are Funcgo versions of some of the Clojure examples
from the [Clojure Cookbook][cookbook].

```go
		func add(x, y) {
			x + y
		}
		add(1, 2)

        => 3
```
Here we define a function `foo` and then call it.  If you are a Go
programmer this should look familiar. However you might notice that
the types are missing and that there is no `return` statement.

Funcgo does not require types, though as we will see later, in certain
cases when performance is important you can specify types at a few
strategic locations.

Funcgo does not have a `return` statement, rather a function simply
returns the value of its last expression (often its only expression).

```go
package example
import(
        "clojure/string"
)
```

Here we see what the top of a Funcgo source file called `example.go`
might look like.  Here we import in a Clojure
[string utility package][string] to be used in this file.

```go
        string.isBlank("")

        => true
```

Because of the `import` statement at the top we can now access
functions in the `string` package provide by Clojure.  One little
wrinkle is that the Clojure function is actually [`blank?`][isblank],
with a `?` character that is illegal in Funcgo. Similarly many Clojure
functions have `-` characters in their name that Funcgo does not
allow. So we automatically _mangle_ identifiers so that `isSomething`
becomes `something?` and `thisIsAnIdentifier` becomes
`this-is-an-identifier`. This is important, because you will often
have to refer to the [Clojure documentation of its library][apidoc].

```go
        string.capitalize("this is a proper sentence.")

        =>  "This is a proper sentence."
```

```go
        string.upperCase("Dépêchez-vous, l'ordinateur!")

        => "DÉPÊCHEZ-VOUS, L'ORDINATEUR!"
```

```go
        string.replace("Who\t\nput  all this\fwhitespace here?", /\s+/, " ")

        => "Who put all this whitespace here?"
```

The last example above shows that string escapes are familiar-looking
to most programmers. It also introduces the syntax for _regular
expression literals_, which are written between a pair of `/`
characters.

```go
        str("John", " ", "Doe")

        => "John Doe"
```

Funcgo does *not* concatenate strings using a `+` operator like other
languages you may be familiar with.  Instead you use the [`str`][str]
function. This is one of the many functions defined in [`clojure.core`][ccore]
that can be used without needing an `import` statement.

```go
        const(
                firstName = "John"
                lastName  = "Doe"
                age       = 42
        )
        str(lastName, ", ", firstName, " - age: ", age)
        
        =>, "Doe, John - age: 42"
```

In keeping with its orientation as a functional programming language,
Funcgo does *not* have mutable local variables. Instead, inside
functions and other scopes you should create constants (whose values
can not be changes.

```
		var firstName = "John"
		var lastName = "Doe"
		var age = 42
		str(lastName, ", ", firstName, " - age: ", age)

        => "Doe, John - age: 42"
```

You can create mutable variables using `var`, but these are global and
changes are not propagated between threads, so you should avoid using
them if possible.

Note that in the previous two examples there is a single grouped
`const` declaration but multiple separate individual `var`
declarations. Actually either type of declaration can be grouped or
individual, but for `const` it is better to use the grouped version
for multiple declarations because it generates more efficient Clojure.

```go
        into([], range(1, 20))
        
        =>  [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19]
```

Here we see an example of using the [`range`][range] function to create
a lazy sequence of integers and then using the [`into`][into] function
to create a vector with the same values.

This example also introduces vector literals, with the empty vector
being passed as the first parameter of `into`.

```go
        [] into range(1, 20)
        
        =>,  [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19]
```

This example has the exact same effect as the previous example, but we
are taking advantage of another feature of Funcgo any function that
takes two parameters `foo(param1, param2)` can alternatively be
written in _infix_ notation as `param1 foo param2`.  This can
sometimes lead to cleaner and more readable code.

```go
		const me = {FIRST_NAME: "Eamonn", FAVORITE_LANGUAGE: "Funcgo"}
		str("My name is ", me(FIRST_NAME),
			", and I really like to program in ", me(FAVORITE_LANGUAGE))
    
        => "My name is Eamonn, and I really like to program in Funcgo",
```

The above example introduces a number of new language features.

First note the _dictionary literal_ which in this case creates a dictionary
with two entries.

In this case the keys are _keywords_ which in Funcgo are distinguished
by being all-uppercase.  Unlike symbols that evaluate to something
else, keywords just evaluate to themselves and are most commonly used
like this as dictionary keys.

Note that to extract values from the dictionary you treat it as if it
were a function, using the key as the parameter to the function.

```go
        str apply (" " interpose [1, 2.000, 3/1, 4/9])
        
        => "1 2.0 3 4/9"
```

This example shows two nested infix expressions.

The inner ones uses the [`interpose`] function to take the vector
`[1, 2.000, 3/1, 4/9]` and create a new vector with blanks inserted
between `[1, " ", 2.000, " ", 3/1, " ", 4/9]`.

The outer infix expression shows an example of Funcgo being used as a
functional programming language. The [`apply`][apply] function is an
example of a function that takes a function as a parameter.  In this
case [`str`][str]) is passed as the first argument.

```go
func isYelling(utterance String) {
  isEvery(
          func(ch Character) { !Character::isLetter(ch) || Character::isUpperCase(ch) },
          utterance
  )
}
```

This example shows an example of Java interoperability. In this case
symbols are not mangled but are passed on unchanged to Java.  The `::`
specifies access to a static function.

This is also the first time we have specified a type for a value, in
this case the `String` type on the outer function parameter.  This is
optional, but doing so in this case avoids using Java reflection,
making for a more efficient implementation.

We also see here an example of an anonymous function, in this case a
predicate (function returning Boolean) that tests if a character is a
non-letter or an upper-case letter.

The [`isEvery`][isevery] function tests whether this predicate is true
for every character in the string.


## Building and Development

You need Leiningen (the Clojure build tool) to build the compiler.
(Note that if you are on Ubuntu, as of March 2014 the version in the
standard Ubuntu package manager is too old to work with this project.
Instead download the `lein` script from the
[Leiningen web site](http://leiningen.org/#install) and put in your
PATH.

To create a new compiler JAR execute ...

```sh
lein fgoc
lein uberjar
```

... which will compile the compiler and generate a JAR file
<code>target/funcgo-<i>x</i>.<i>y</i>.<i>z</i>-standalone.jar</code>

You can run the unit tests by doing

```sh
lein midje
```

## License

The Funcgo code is distributed under the Eclipse Public License either
version 1.0 or (at your option) any later version.

<a rel="license" href="http://creativecommons.org/licenses/by/4.0/"><img alt="Creative Commons License" style="border-width:0" src="http://i.creativecommons.org/l/by/4.0/80x15.png" /></a><br /><span xmlns:dct="http://purl.org/dc/terms/" href="http://purl.org/dc/dcmitype/Text" property="dct:title" rel="dct:type">Funcgo Documentation</span> by <span xmlns:cc="http://creativecommons.org/ns#" property="cc:attributionName">Eamonn O'Brien-Strain</span> is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by/4.0/">Creative Commons Attribution 4.0 International License</a>.


[cookbook]: http://clojure-cookbook.com/
[string]: http://clojure.github.io/clojure/clojure.string-api.html
[isblank]: http://clojure.github.io/clojure/clojure.string-api.html#clojure.string/blank?
[apidoc]: http://clojure.github.io/clojure/index.html
[ccore]: http://clojure.github.io/clojure/clojure.core-api.html
[str]: http://clojure.github.io/clojure/clojure.core-api.html#clojure.core/str
[into]: http://clojuredocs.org/clojure_core/clojure.core/into
[range]: http://clojuredocs.org/clojure_core/clojure.core/range
[apply]: http://clojuredocs.org/clojure_core/clojure.core/apply
[interpose]: http://clojuredocs.org/clojure_core/clojure.core/interpose
[isevery]: http://clojuredocs.org/clojure_core/clojure.core/every_q
