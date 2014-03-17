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
	apply(str, header, interpose("\n", rows)),
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

