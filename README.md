# funcgo

Funcgo is a compiler that converts Functional Go into Clojure.

The compiler itself is written in Functional Go. (Turtles all the way down!)

## Usage

To create a new compiler JAR execute ...

```sh
lein fgoc
lein uberjar
```

... which will compile the compiler and generate a JAR file
`target/funcgo-<i>x</i>.<i>y</i>.<i>z</i>-standalone.jar`

## License

The Funcgo code is dtributed under the Eclipse Public License either
version 1.0 or (at your option) any later version.

<a rel="license" href="http://creativecommons.org/licenses/by/4.0/"><img alt="Creative Commons License" style="border-width:0" src="http://i.creativecommons.org/l/by/4.0/80x15.png" /></a><br /><span xmlns:dct="http://purl.org/dc/terms/" href="http://purl.org/dc/dcmitype/Text" property="dct:title" rel="dct:type">Funcgo Documentation</span> by <span xmlns:cc="http://creativecommons.org/ns#" property="cc:attributionName">Eamonn O'Brien-Strain</span> is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by/4.0/">Creative Commons Attribution 4.0 International License</a>.

