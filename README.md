> [Tour][fgotour] - [FAQ][faq] - [Reference][ref]

# funcgo

Funcgo is a compiler that converts Functional Go into Clojure, to run
on the JVM or as JavaScript.

## Try It Out

Without installing anything you can try the [online tour][fgotour]
where you can type Funcgo and see how it converts to Clojure
and evaluates.

(By the way the online tour is itself an example web application that
[uses Funcgo][fgosite] for both it server side (JVM) and its
client side (JS).

## Quick Start

#### 1. Set up Clojure development environment.

Follow the [install instructions for Leiningen][lein]

On the command line, go to an empty directory and type ...
```sh
lein new app hello
cd hello
lein run
```
This should print out `Hello, World!`, showing that your Clojure
environment is configured correctly.

#### 2. Convert your project into a Funcgo project.

[![Clojars Project](http://clojars.org/org.eamonn.funcgo/funcgo-lein-plugin/latest-version.svg)](http://clojars.org/org.eamonn.funcgo/funcgo-lein-plugin)

In your favorite text editor, edit the file `project.clj` and insert
`[org.eamonn.funcgo/funcgo-lein-plugin "0.4.1"]` (or the latest
version from clojars shown above) in *both* the dependencies section
*and* the plugins section. You should end up with something like ...

```clj
(defproject hello "0.1.0-SNAPSHOT"
  :description "FIXME: write description"
  :url "http://example.com/FIXME"
  :license {:name "Eclipse Public License"
            :url "http://www.eclipse.org/legal/epl-v10.html"}
  :dependencies [[org.clojure/clojure "1.5.1"]
                 [org.eamonn.funcgo/funcgo-lein-plugin "0.3.0"]]
  :plugins [ [org.eamonn.funcgo/funcgo-lein-plugin "0.3.0"]]
  :main ^:skip-aot hello.core
  :target-path "target/%s"
  :profiles {:uberjar {:aot :all}})
```

Delete the existing Clojure file `src/hello/core.clj`

#### 3. Write Funcgo

Create a new Funcgo file `src/hello/core.go` containing ...
```go
package core

func _main(args...) {
	println("Hello, World from Funcgo")
}
```

On the command line (in the same directory as the `project.clj`
file) do ...
```sh
lein fgoc
lein run
```
This should print out `Hello, World from Funcgo`.

Congratulations, you have just written and executed your first
Funcgo program!

## Next Steps

You can get a better feel for the language by reading the Introduction
to the Funcgo Language section below.

To dive deeper, see [Funcgo Reference][ref] doc.

To browse some actual working code, the biggest and most complex
program so far written in Funcgo is its own compiler. (Turtles
all the way down!)  You might start at the `main.go` file in
[the source directory][src].

A smaller set of working code is [fgolib][fgolib].  In addition to
looking at the Funcgo code there, you can also examine the
`project.clj` file which is a working example of using the Leiningen
plugin.

If you want to see a complete web app, that generates both
Clojurescript and Clojure, see
[the source for www.funcgo.org][fgosite].

There is also do `lein fgoc --repl` to bring up the beginnings of a
REPL that you can use to explore...
```
$ lein fgoc --repl
test
src

fgo=>     2+3
Clojure:  (+ 2 3)
Result:   5

fgo=>     func{10 * $1} map [1,2,3,4,5,6]
Clojure:  (map #(* 10 %) [1 2 3 4 5 6])
Result:   (10 20 30 40 50 60)

fgo=>
```

### Not Using Leiningen?

The preferred way to use this compiler is via the
[Leiningen Plugin][plugin] as described in the Quick Start section.

If you are not using Leiningen you can use `java -jar
bin/funcgo-compiler-*-standalone.jar directory ...` to compile.

## Introduction to the Funcgo Language

### Why a new language?

The goal of Funcgo is to combine the readability of the Go language
with the semantics of Clojure.

1. Go is a language that has been well designed to be very
readable. However it is best as a low-level system programming
language (replacing C) and it is missing many of the higher-level
features that programmers expect for working further up the stack, in
for example in web applications.

2. Clojure is a variety of Lisp that inter-operates with Java or
JavaScript.  It encourages a functional programming style with
efficient immutable containers, combined with a thread-safe model for
mutable state called software transactional memory. However, Clojure
is difficult to read for programmers unfamiliar with Lisp syntax.

### Examples for Clojure Programmers

In this section are Funcgo versions of some of the Clojure examples
from the [Clojure Cookbook][cookbook].

#### Defining and using a function
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

Funcgo does not need a `return` statement, rather a function simply
returns the value of its last expression (often its only expression).

#### Adding a file header
```go
package example
import(
        "clojure/string"
)
```

Here we see what the top of a Funcgo source file called `example.go`
might look like.  Here we import in a Clojure
[string utility package][string] to be used in this file.

#### Using symbols from other packages
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

#### Specifying string escapes and regular expressions
```go
        string.replace("Who\t\nput  all this\fwhitespace here?", /\s+/, " ")

        => "Who put all this whitespace here?"
```

The example above shows that string escapes are familiar-looking to
most programmers. It also introduces the syntax for _regular
expression literals_, which are written between a pair of `/`
characters.

#### Concatenating strings
```go
        str("John", " ", "Doe")

        => "John Doe"
```

Funcgo does *not* concatenate strings using a `+` operator like other
languages you may be familiar with.  Instead you use the [`str`][str]
function. This is one of the many functions defined in [`clojure.core`][ccore]
that can be used without needing an `import` statement.

#### Specifying local (immutable) variables
```go
		firstName, lastName, age := "John", "Doe", 42
		str(lastName, ", ", firstName, " - age: ", age)
	=> "Doe, John - age: 42"
```

In keeping with its orientation as a functional programming language,
Funcgo does *not* have mutable local variables. Instead, inside
functions and other scopes you should create constants (whose values
can not be changed.

#### Specifying global (mutable) variables
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

#### Using vectors
```go
        into([], range(1, 20))

        =>  [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19]
```

Here we see an example of using the [`range`][range] function to create
a lazy sequence of integers and then using the [`into`][into] function
to create a vector with the same values.

This example also introduces vector literals, with the empty vector
being passed as the first parameter of `into`.

#### Getting cleaner syntax using infix notation
```go
        []  into  range(1, 20)

        =>,  [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19]
```

This example has the exact same effect as the previous example, but we
are taking advantage of another feature of Funcgo any function that
takes two parameters `foo(param1, param2)` can alternatively be
written in _infix_ notation as `param1 foo param2`.  This can
sometimes lead to cleaner and more readable code.

#### Specifying keyword and dictionary literals
```go
		me := {FIRST_NAME: "Eamonn", FAVORITE_LANGUAGE: "Funcgo"}
		str("My name is ", me(FIRST_NAME),
			", and I really like to program in ", me(FAVORITE_LANGUAGE))
	=> "My name is Eamonn, and I really like to program in Funcgo"
```

The above example introduces a number of new language features.

First note the _dictionary literal_ which creates a dictionary with
two entries.

Here the keys are _keywords_ which in Funcgo are distinguished by
being all-uppercase.  Unlike symbols that evaluate to something else,
keywords just evaluate to themselves and are most commonly used like
this as dictionary keys.

Note that to extract values from the dictionary you treat it as if it
were a function, using the key as the parameter to the function.

#### Combining infix and functional programming
```go
        str  apply  (" "  interpose  [1, 2.000, 3/1, 4/9])

        => "1 2.0 3 4/9"
```

This example shows two nested infix expressions.

The inner ones uses the [`interpose`][interpose] function to take the
vector `[1, 2.000, 3/1, 4/9]` and create a new vector with blanks
inserted between `[1, " ", 2.000, " ", 3/1, " ", 4/9]`.

The outer infix expression shows an example of Funcgo being used as a
functional programming language. The [`apply`][apply] function is an
example of a function that takes a function as a parameter.  Here
[`str`][str] is passed as the first argument.

#### Calling function variadically
```go
        str(...(" "  interpose  [1, 2.000, 3/1, 4/9]))

        => "1 2.0 3 4/9"
```

This example is equivalent to the previous one, but it shows some
syntactic sugar for the `apply` function in a way that echoes how
variadic functions are declared.  Essentially if you have `const args
= [a, b, c]` then calling `foo(...args)` is the same as calling
`foo(a, b, c)`.


#### Inter-operating with Java or JavaScript
```go
func isYelling(utterance String) {
  isEvery(
          func(ch Character) { !Character::isLetter(ch) || Character::isUpperCase(ch) },
          utterance
  )
}
```

This example shows an example of Java interoperability.  The `::`
specifies access to a static function (with symbol names not being
mangled, but passed to Java as-is).

This is also the first time we have specified a type for a value,
specifying the `String` type on the outer function's parameter.  This
is optional, but doing so in this case avoids Java reflection, making
for a more efficient implementation.

We also see here an example of an anonymous function, here a predicate
(function returning Boolean) that tests if a character is a non-letter
or an upper-case letter.

The [`isEvery`][isevery] function tests whether this predicate is true
for every character in the string.

### Examples for Go Programmers

In this section are Funcgo versions of some of the Go examples
from the [A Tour of Go][tour].

#### Placement of constant definitions
```go
package main

import "fmt"

Pi := 3.14

func main() {
	World := "世界"
	fmt.Println("Hello", World)
	fmt.Println("Happy", Pi, "Day")
	{
		Truth := true
		fmt.Println("Go rules?", Truth)
	}
}


    => Hello 世界
Happy 3.14 Day
Go rules? true
```

One constraint on `:=` definitions is that, except for at the top
level, they have to be at the beginning of a curly-brace block. So
above we had to add an extra level of curlies to allow `Truth` to be
defined at the bottom of the function.


#### Go primitive types
```go
package main

import (
	"fmt"
	"math"
)

func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		v
	} else {
		lim
	}
}

func main() {
	fmt.Println(
		pow(3, 2, 10),
		pow(3, 3, 20),
	)
}


    => 9 20
```

For compatibility with Go, you can use Go-style primitive types, but they are mapped to JVM
primitive types that may have different bit sizes.

#### Optional `return`

```go
package main

import (
	"fmt"
)

func newton(n int, x, z float64) float64 {
	if n == 0 {
		z
	} else {
		newton(n-1, x, z-(z*z-x)/(2*x))
	}
}

func Sqrt(x float64) float64 {
	return newton(500, x, x/2)
}

func main() {
	fmt.Println(Sqrt(100))
}


    => 10.000000000000007
```

For compatibility with Go, you can add a cosmetic `return` to a
function, but only in the special case of returning the top level
expression of a function.

#### Data structures

```go
package main

import "fmt"

type Vertex struct {
	X int
	Y int
}

func main() {
	fmt.Println(Vertex{1, 2})
}

    => {1 2}
```

You can go a long way in Funcgo just using the built in
dictionary and vector types, but you can also create data structures
that are implemented as Java classes.


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

## Thanks

Funcgo is built on the folder of giants.

Thanks to Rich Hickey and the Clojure contributors, to Thompson, Pike,
and Griesemer and the Go contributors, and to Mark Engelberg for the
instaparse parsing library.

## License

The Funcgo code is distributed under the Eclipse Public License either
version 1.0 or (at your option) any later version.

<a rel="license" href="http://creativecommons.org/licenses/by/4.0/"><img alt="Creative Commons License" style="border-width:0" src="http://i.creativecommons.org/l/by/4.0/80x15.png" /></a><br /><span xmlns:dct="http://purl.org/dc/terms/" href="http://purl.org/dc/dcmitype/Text" property="dct:title" rel="dct:type">Funcgo Documentation</span> by <span xmlns:cc="http://creativecommons.org/ns#" property="cc:attributionName">Eamonn O'Brien-Strain</span> is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by/4.0/">Creative Commons Attribution 4.0 International License</a>.

[lein]: http://leiningen.org/
[cookbook]: http://clojure-cookbook.com/
[tour]: http://tour.golang.org
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
[plugin]: https://github.com/eobrain/funcgo-lein
[src]: https://github.com/eobrain/funcgo/tree/master/src/funcgo
[fgolib]: https://github.com/eobrain/fgolib
[fgosite]: https://github.com/eobrain/fgosite
[fgotour]: http://tour.funcgo.org
[ref]: doc/reference.md
[faq]: doc/FAQ.md
