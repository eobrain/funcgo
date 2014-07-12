# Frequently Asked Questions

## Are other targets being considered?

The Funcgo compiler is likely to only ever emit Clojure.  However,
that means it already has good support for the JVM and JavaScript
targets, and in addition there are projects to have Clojure target
Android and the CLR. and perhaps there will be more targets in future.
If there is enough interest, Funcgo can tweak its Clojure output
appropriately for new targets, as it already does for Clojure versus
ClojureScript.

## Does it support numerical programming?

Funcgo can readily call Java code, so it supports numerical
programming or linking in C++ code to the same extent that Java does.
It can also of course directly use any Clojure numerical programming
library such as [core.matrix][1].

## Is this a full implementation of Go as the "Functional Go" name implies?

No.  Perhaps it would be better named "Glojure".  The language is
fundamentally just syntactical sugar on top of Clojure, though I did
try to keep the sugar as close to Go as I could.

## Can a programmer use macros (as in Clojure) to modify the language itself.

I do not plan to make it easy to _define_ macros in Funcgo.  Inspired
by Go's philosophy, I want to limit the scope of the language to keep
it somewhat simple. Funcgo however can _use_ macros, so a sufficiently
motivated programmer can write macros in Clojure and use them in
Funcgo.

## So Funcgo is really just Lisp without parentheses... again.

Basically yes, or put it another way, it is Lisp with more syntax. The
fact that Lisp has so little syntax is beautifully elegant from a
mathematical point of view, but that is I believe at the expense of
readability, at least for programmers coming from other programming
traditions. Of course readability is somewhat subjective and I
understand that for long-time Lisp programmers, Lisp is perfectly
readable.

## How does the user keeps track of the source location (source map)?

I have [plans to implement source maps][3] when targeting JavaScript.
I have not yet figured out the best way to do something equivalent
when targeting the JVM.  This is one of several tool-chain issues. In
addition to stabilizing the language itself, getting a smooth
tool-chain is the biggest thing left to do before a 1.0 release of
Funcgo.  (Though it is worth noting that Coffeescript, an analogous
language, got widespread adoption before tackling the source location
issue.)


## Credit

Some of these questions are paraphrased from some threads on
[Hacker News][2]. Thanks to the posters there.


[1]: https://github.com/mikera/core.matrix
[2]: https://news.ycombinator.com/item?id=8017588
[3]: https://github.com/eobrain/funcgo/issues/19
