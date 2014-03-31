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
programmers unfamiliar with Lisp syntax, Clojure is very difficult to
read.

### Examples

The following are Funcgo versions of some of the Clojure examples
given in the Clojure Cookbook.

```go
		func add(x, y) {
			x + y
		}
        
		add(1,2)
```
```
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

