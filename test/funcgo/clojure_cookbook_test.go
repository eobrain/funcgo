package funcgo.clojure_cookbook_test
import(
        test midje.sweet
        fgo funcgo.core
        string clojure.string
)

func add(x,y) {
        x + y
}
test.fact("Simple example",
        add(1,2),
        =>, 3
)

test.fact("More complex example",
        into([],  \range(1, 20)),
        =>,  [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19]
)

test.fact("Any function of two arguments can be written infix",
        [] into \range(1, 20),
        =>,  [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19]
)

test.fact("Infix is most convenient for math operators.",
        1 + 2,
        =>, 3
)

test.fact("Dotted identifers are from other packages.",
        // import section includes
        //    string clojure.string
        string.isBlank(""),
        =>, true
)

test.fact("Capitalize first character in a string.",
        string.capitalize("this is a proper sentence."),
        =>,  "This is a proper sentence."
)

test.fact("Capitalize or lower-case all characters.",

        string.upperCase("loud noises!"),
        =>, "LOUD NOISES!",

        string.lowerCase("COLUMN_HEADER_ONE"),
        =>, "column_header_one",

        string.lowerCase("!&$#@#%^[]"),
        =>, "!&$#@#%^[]",

        string.upperCase("Dépêchez-vous, l'ordinateur!"),
        =>, "DÉPÊCHEZ-VOUS, L'ORDINATEUR!"
)

test.fact("Remove whitespace at beginning and end.",
        string.trim(" \tBacon ipsum dolor sit.\n"),
        =>, "Bacon ipsum dolor sit."
)

test.fact("Collapse whitespace into single whitespace",
        string.replace("Who\t\nput  all this\fwhitespace here?", /\s+/, " "),
        =>, "Who put all this whitespace here?"
)

test.fact("Windows to Unix line-endings",
        string.replace("Line 1\r\nLine 2", "\r\n", "\n"),
        =>, "Line 1\nLine 2"
)

test.fact("Trim only from one end",
        string.triml(" Column Header\t"),
        =>, "Column Header\t",

        string.trimr("\t\t* Second-level bullet.\n"),
        =>, "\t\t* Second-level bullet."
)

test.fact("Concatenate strings",
        str("John", " ", "Doe"),
        =>, "John Doe"
)

test.fact("Can concatenate consts.",
        const(
                firstName = "John"
                lastName = "Doe"
                age = 42
        )
        str(lastName, ", ", firstName, " - age: ", age),
        =>, "Doe, John - age: 42"
)

firstName := "John"
lastName := "Doe"
age := 42
test.fact("Can concatenate vars.",
        str(lastName, ", ", firstName, " - age: ", age),
        =>, "Doe, John - age: 42"
)

test.fact("turn characters into a string",
        apply(str, "ROT13: ", ['W', 'h', 'y', 'v', 'h', 'f', ' ', 'P', 'n', 'r', 'f', 'n', 'e']),
        =>, "ROT13: Whyvhf Pnrfne"
)

lines := [
        "#! /bin/bash\n",
        "du -a ./ | sort -n -r\n"
]
test.fact("make file from lines (with newlines)",
        str apply lines,
        =>,  "#! /bin/bash\ndu -a ./ | sort -n -r\n"
)

header := "first_name,last_name,employee_number\n"
rows := [
        "luke,vanderhart,1",
        "ryan,neufeld,2"
]
test.fact("Making CSV from header vector of rows",
        apply(str, header, ("\n" interpose rows)),
        =>, `first_name,last_name,employee_number
luke,vanderhart,1
ryan,neufeld,2`
)


foodItems := ["milk", "butter", "flour", "eggs"]
test.fact("Join can be easier",

        string.join(", ", foodItems),
        =>, "milk, butter, flour, eggs",

        ", " string.join foodItems,
        =>, "milk, butter, flour, eggs",

        string.join( [1, 2, 3, 4] ),
        =>, "1234"
)

test.fact("seq() exposes the characters in a string",
        seq("Hello, world!"),
        =>, ['H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', '!']
)

func isYelling(s) {
  isEvery(
          func(ch) { !Character::isLetter(ch) || Character::isUpperCase(ch) },
          s
  )
}

test.fact("Function taking a sequence will cooerce a string into a set of chars",

        // count of chars in string
        frequencies(string.lowerCase("An adult all about A's")),
        =>, {' ':4, 'a':5, 'b':1, 'd':1, '\'':1, 'l':3, 'n':1, 'o':1, 's':1, 't':2, 'u':2},

        //  every letter capitalized?
        isYelling("LOUD NOISES!"),
        =>, true,

        isYelling("Take a DEEP breath."),
        =>, false
)


test.fact("Can transform characters back into a string",
        str apply ['H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', '!'],
        =>, "Hello, world!"
)

test.fact("int function converts characters to integers",

        int('a'),
        =>, 97,

        int('ø'),
        =>, 248,

        int('α'),  // Greek letter alpha
        =>, 945,

        int('\u03B1'), // Greek letter alpha (by code point)
        =>, 945,

        int map "Hello, world!",
        =>, [72, 101, 108, 108, 111, 44, 32, 119, 111, 114, 108, 100, 33]
)

test.fact("char function does the opposite",
        char(97),
        =>, 'a',

        char(125),
        =>,  '}',

        char(945),
        =>,  'α',

        reduce(
                func(acc, i){str(acc, char(i))},
                "",
                [115, 101, 99, 114, 101, 116, 32, 109, 101, 115, 115, 97, 103, 101, 115]
        ),
        =>,  "secret messages"
)

me := {FIRST_NAME: "Ryan", FAVORITE_LANGUAGE: "Clojure"}

test.fact("str is the easiest way of formatting values into a string",
	
        str("My name is ", me[FIRST_NAME],
                ", and I really like to program in ", me[FAVORITE_LANGUAGE]),
        =>, "My name is Ryan, and I really like to program in Clojure",

        str apply (" " interpose [1, 2.000, 3/1, 4/9]),
        =>, "1 2.0 3 4/9"
)


// Produce a filename with a zero-padded sortable index
func filename(name, i) {
        format("%03d-%s", i, name)
}

test.fact("format is another way of constructing strings",
        filename("my-awesome-file.txt", 42),
        =>, "042-my-awesome-file.txt"
)

// Create a table using justification
func tableify(row) {
        apply(format, "%-20s | %-20s | %-20s", row)
}

header := ["First Name", "Last Name", "Employee ID"]
employees := [
        ["Ryan", "Neufeld", 2],
        ["Luke", "Vanderhart", 1]
]

mapv(
	println,
        map(
		tableify,
                concat([header], employees)
	)
)
// *out*
// First Name           | Last Name            | Employee ID
// Ryan                 | Neufeld              | 2
// Luke                 | Vanderhart           | 1

->>(
        concat([header], employees),
        map(tableify),
        mapv(println)
)
// *out*
// First Name           | Last Name            | Employee ID
// Ryan                 | Neufeld              | 2
// Luke                 | Vanderhart           | 1
