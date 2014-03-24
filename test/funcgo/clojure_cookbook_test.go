package funcgo.clojure_cookbook_test
import(
        test midje.sweet
        fgo funcgo.core
        string clojure.string
	inf inflections.core
)

test.fact("Simple example",
	{
		func add(x,y) {
			x + y
		}
		add(1,2)
	},
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

test.fact("Can concatenate vars.",
	{
		firstName := "John"
		lastName := "Doe"
		age := 42
		str(lastName, ", ", firstName, " - age: ", age)
	},
        =>, "Doe, John - age: 42"
)

test.fact("turn characters into a string",
        apply(str, "ROT13: ", ['W', 'h', 'y', 'v', 'h', 'f', ' ', 'P', 'n', 'r', 'f', 'n', 'e']),
        =>, "ROT13: Whyvhf Pnrfne"
)

test.fact("make file from lines (with newlines)",
	{
		lines := [
			"#! /bin/bash\n",
			"du -a ./ | sort -n -r\n"
		]
		str apply lines
	},
        =>,  "#! /bin/bash\ndu -a ./ | sort -n -r\n"
)

test.fact("Making CSV from header vector of rows",
	{
		header := "first_name,last_name,employee_number\n"
		rows := [
			"luke,vanderhart,1",
			"ryan,neufeld,2"
		]
		apply(str, header, ("\n" interpose rows))
	},
        =>, `first_name,last_name,employee_number
luke,vanderhart,1
ryan,neufeld,2`
)


test.fact("Join can be easier",
	{
		foodItems := ["milk", "butter", "flour", "eggs"]
		string.join(", ", foodItems)
	},
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


test.fact("str is the easiest way of formatting values into a string",
	
	{
		me := {FIRST_NAME: "Eamonn", FAVORITE_LANGUAGE: "Funcgo"}
		str("My name is ", me[FIRST_NAME],
			", and I really like to program in ", me[FAVORITE_LANGUAGE])
	},
        =>, "My name is Eamonn, and I really like to program in Funcgo",

        str apply (" " interpose [1, 2.000, 3/1, 4/9]),
        =>, "1 2.0 3 4/9"
)



test.fact("format is another way of constructing strings",
	{
		// Produce a filename with a zero-padded sortable index
		func filename(name, i) {
			format("%03d-%s", i, name)
		}
		"my-awesome-file.txt" filename 42
	},
        =>, "042-my-awesome-file.txt",

	"%07.3f" format 0.005,
	=>, "000.005"
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

println mapv (tableify map ([header] concat employees))
// *out*
// First Name           | Last Name            | Employee ID
// Ryan                 | Neufeld              | 2
// Luke                 | Vanderhart           | 1

->>(
        [header] concat employees,
        map(tableify),
        mapv(println)
)
// *out*
// First Name           | Last Name            | Employee ID
// Ryan                 | Neufeld              | 2
// Luke                 | Vanderhart           | 1

test.fact("Regular expressions, using reFind",

	/\d+/ reFind "I've just finished reading Fahrenheit 451",
	=>, "451",

	/Bees/ reFind "Beads aren't cheap.",
	=>, nil
)

test.fact("To match only the whole string use reMatches",

	/\w+/ reFind "my-param",
	=>, "my",

	/\w+/ reMatches "my-param",
	=>, nil,

	/\w+/ reMatches "justLetters",
	=>, "justLetters"
)

test.fact("Extract strings from a larger string using reSeq",
	
        /\w+/ reSeq "My Favorite Things",
	=>, ["My", "Favorite", "Things"],
	
	/\d{3}-\d{4}/ reSeq "My phone number is 555-1234.",
	=>, ["555-1234"],
	
        {
		// Extract Twitter identifiers in a tweet
		func mentions(tweet) {
			/(@|#)(\w+)/ reSeq tweet
		}
		
		mentions("So long, @earth, and thanks for all the #fish. #goodbyes")
	},
        =>, [["@earth", "@", "earth"], ["#fish", "#", "fish"], ["#goodbyes", "#", "goodbyes"]],

        {
		// Capture and decompose a phone number and its title
		rePhoneNumber := /(\w+): \((\d{3})\) (\d{3}-\d{4})/  
		rePhoneNumber reSeq "Home: (919) 555-1234, Work: (191) 555-1234"
	},
        =>, [["Home: (919) 555-1234", "Home", "919", "555-1234"],
		["Work: (191) 555-1234", "Work", "191", "555-1234"]]

)


test.fact("simple string replacement via string.replace",
	{
		aboutMe := "My favorite color is green!"
		string.replace(aboutMe, "green", "red")
	},
	=>, "My favorite color is red!",
	{
		func deCanadianize(s) {
			string.replace(s, "ou", "o")
		}
		deCanadianize(str(
			"Those Canadian neighbours have coloured behaviour",
			" when it comes to word endings"))
	},
	=>, "Those Canadian neighbors have colored behavior when it comes to word endings"
)

test.fact("More complex string replacement requires regular expressions",
	{
		// Add Markdown-style links for any GitHub issue numbers present in comment
		func linkifyComment(repo, comment) {
			string.replace(
				comment,
				/#(\d+)/,
				str("[#$1](https://github.com/", repo, "/issues/$1)")
			)
		}
		linkifyComment(
			"next/big-thing",
			"As soon as we fix #42 and #1337 we should be set to release!"
		)
	},
	=>, "As soon as we fix [#42](https://github.com/next/big-thing/issues/42) and [#1337](https://github.com/next/big-thing/issues/1337) we should be set to release!"
)


test.fact("Use string.split to split strings",
	
	"HEADER1,HEADER2,HEADER3" string.split /,/,
	=>, ["HEADER1", "HEADER2", "HEADER3"],
	
	"Spaces   Newlines\n\n" string.split /\s+/,
	=>, ["Spaces", "Newlines"],
	
	// whitespace splitting with implicit trim
	"field1    field2 field3   "  string.split /\s+/,
	=>, ["field1", "field2", "field3"],
	
	// avoid implicit trimming by adding limit of -1

	string.split("ryan,neufeld,", /,/, -1),
	=>, ["ryan", "neufeld", ""],
	
	{
		dataDelimiters := /[ :-]/
			
		//No-limit split on any delimiter
		"2013-04-05 14:39" string.split dataDelimiters
	},
	=>, ["2013", "04", "05", "14", "39"],

	// Limit of 1 - functionally: return this string in a collection
	string.split("2013-04-05 14:39", dataDelimiters, 1),
	=>, ["2013-04-05 14:39"],

	// Limit of 2
	string.split("2013-04-05 14:39", dataDelimiters, 2),
	=>, ["2013", "04-05 14:39"],

	// Limit of 100
	string.split("2013-04-05 14:39", dataDelimiters, 100),
	=>, ["2013", "04", "05", "14", "39"]
)


// The following requires the following in leinigen dependencies
//    [inflections "0.9.5"]
test.fact("can use inf.pluralize to with word labelling counts",
	// In import have
	//      inf inflections.core

	1 inf.pluralize "monkey",
	=>, "1 monkey",

	12 inf.pluralize "monkey",
	=>, "12 monkeys",

	// Can provide non-standard pluralization as an arg

	inf.pluralize(1, "box", "boxen"),
	=>, "1 box",

	inf.pluralize(3, "box", "boxen"),
	=>, "3 boxen",

	// Or you can add your own rules
	inf.plural("box"),
	=>, "boxes",

	{
		// Words ending in 'ox' pluralize with 'en' (and not 'es')
		/(ox)(?i)$/ inf.mutatePlural "$1en"
		
		inf.plural("box")
	},
	=>, "boxen",

	// plural is also the basis for pluralize...
	2 inf.pluralize "box",
	=>, "2 boxen",

	// Convert "snake_case" to "CamelCase"
	inf.camelize("my_object"),
	=>, "MyObject",

	// Clean strings for usage as URL parameters
	inf.parameterize("My most favorite URL!"),
	=>, "my-most-favorite-url",

	// Turn numbers into ordinal numbers
	inf.ordinalize(42),
	=>, "42nd"
)

test.fact("Can convert between different types of language things (note Funcgo mangling).",

	symbol("valid?"),
	=>, quote(isValid),

	str(quote(isValid)),
	=>, "valid?",

	name(TRIUMPH),
	=>, "triumph",

	str(TRIUMPH),
	=>, ":triumph",

	keyword("fantastic"),
	=>, FANTASTIC,

	keyword(quote(fantastic)),
	=>, FANTASTIC,

	symbol(name(WONDERFUL)),
	=>, quote(wonderful),

	// If you only want the name part of a keyword.
	// (We have to escape into Clojure for this.)
	name(\`:user/valid?`),
	=>, "valid?",

	// If you only want the namespace
	namespace(\`:user/valid?`),
	=>, "user",

	str(\`:user/valid?`),
	=>, ":user/valid?",

	str(\`:user/valid?`)->substring(1),
	=>, "user/valid?",

	keyword(quote(produce.onions)),
	=>, \`:produce/onions`,

	symbol(str(\`:produce/onions`)->substring(1)),
	=>, quote(produce.onions),

	// keyword and symbol also have 2-argument (infix) versions
	{
		shoppingArea := "bakery"
		shoppingArea keyword "bagels"
	},
	=>, \`:bakery/bagels`,

	shoppingArea symbol "cakes",
	=>, quote(bakery.cakes)
)

test.fact("Funcgo has numbers",

	// Avogadro's number
	6.0221413e23,
	=>, 6.0221413E23,

	// 1 Angstrom in meters
	1e-10,
	=>, 1.0E-10,

	// Size-bounded integers can overflow
	try {
		9999 * 9999 * 9999 * 9999 * 9999

	} catch \ArithmeticException e {
		e->getMessage
	},
	=>, "integer overflow",

	// which you can avoid using Big integers
	9999N * 9999 * 9999 * 9999 * 9999,
	=>, 99950009999000049999N,

	2 * Double::MAX_VALUE,
	=>, Double::POSITIVE_INFINITY,

	2 * bigdec(Double::MAX_VALUE),
	=>, 3.5953862697246314E+308M,

	// Result of integer division is a ratio type
	type(1 / 3),
	=>, \`clojure.lang.Ratio`,

	3 * (1 / 3),
	=>, 1N,

	(1 / 3) + 0.3,
	=>, 0.6333333333333333,

	// Avoid losing precision
	rationalize(0.3),
	=>, 3/10,

	(1 / 3) + rationalize(0.3),
	=>, 19/30
)

test.fact("Can parse numbers from strings.",

	Integer::parseInt("-42"),
	=>, -42,

	Double::parseDouble("3.14"),
	=>, 3.14,

	bigdec("3.141592653589793238462643383279502884197"),
	=>, 3.141592653589793238462643383279502884197M,

	bigint("122333444455555666666777777788888888999999999"),
	=>, 122333444455555666666777777788888888999999999N

)

test.fact("Can coerce numbers.",

	int(2.0001),
	=>, 2,

	int(2.999999999),
	=>, 2,

	Math::round(2.0001),
	=>, 2,

	Math::round(2.999),
	=>, 3,

	int(2.99 + 0.5),
	=>, 3,

	Math::ceil(2.0001),
	=>, 3.0,

	Math::floor(2.999),
	=>, 2.0,

	3 withPrecision (7M / 9),
	=>, 0.778M,

	1 withPrecision (7M / 9),
	=>, 0.8M,

	withPrecision(1, ROUNDING, \FLOOR, (7M / 9)),
	=>, 0.7M,

	// note non-big arithmetic not effected by withPrecision
	3 withPrecision (1 / 3),
	=>, 1/3,

	3 withPrecision (bigdec(1) / 3),
	=>, 0.333M
)

test.fact("Easy to implement fuzzy equality",

	{
		func fuzzyEq(tolerance, x, y) {
			const(
				diff = Math::abs(x - y)
			)
			diff < tolerance
		}
		fuzzyEq(0.01, 10, 10.000000000001)
	},
	=>, true,

	fuzzyEq(0.01, 10, 10.1),
	=>, false,

	0.22 - 0.23,
	=>, -0.010000000000000009,

	0.23 - 0.24,
	=>, -0.009999999999999981,

	{
		isEqualWithinTen := partial(fuzzyEq, 10)
		100 isEqualWithinTen 109
	},
	=>, true,

	100 isEqualWithinTen 110,
	=>, false
)

test.fact("Can sort with fuzzy equality",
	{
		func fuzzyComparator(tolerance) {
			func(x, y) {
				if fuzzyEq(tolerance, x, y) {
					0
				} else {
					x compare y
				}
			}
		}
		fuzzyComparator(10) sort [100, 11, 150, 10, 9]
	},
	=>, [11, 10, 9, 100, 150]  // 100 and 150 have moved, but not 11, 10, and 9
)


test.fact("Can do trig",

	{
		// Calculating sin(a + b). The formula for this is
		// sin(a + b) = sin a * cos b + sin b cos a
		func sinPlus(a, b) {
			Math::sin(a) * Math::cos(b) + Math::sin(b) * Math::cos(a)
		}
		sinPlus(0.1, 0.3)
	},
	=>, 0.38941834230865047,

	{
		// Calculating the distance in kilometers between two points on Earth
		earthRadius := 6371.009

		func degreesToRadians(point) {
			func(x){Math::toRadians(x)} mapv point
		}

		// Calculate the distance in km between two points on Earth. Each
		// point is a pair of degrees latitude and longitude, in that order.
		func distanceBetween(p1, p2) {
			distanceBetween(p1, p2, earthRadius)
		} (p1, p2, radius) {
			const(
				[lat1, long1] = degreesToRadians(p1)
				[lat2, long2] = degreesToRadians(p2)
			)
			radius * Math::acos(
					Math::sin(lat1) * Math::sin(lat2)
					+
					Math::cos(lat1) * Math::cos(lat2) * Math::cos(long1 - long2)
				)
		}

		[49.2000, -98.1000] distanceBetween [35.9939, -78.8989]
	},
	=>, 2139.42827188432
)


