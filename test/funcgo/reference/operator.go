package operator
import test "midje/sweet"
exclude ( ^, + )

func ^(x, y) {
	Math::pow(x, y)
}

func +(x, y) {
	x  str  y
}

test.fact("Can redefine existing operators",

	2 ^ 3, =>, 8.0,

	10 ^ 2, =>, 100.0,

	"foo" + "bar", =>, "foobar"
)

func \**\(x, y) {
    Math::pow(x, y)
}

test.fact("Can use new operators",

    2 \**\ 3,  =>, 8.0,
    10 \**\ 2, =>, 100.0
)
