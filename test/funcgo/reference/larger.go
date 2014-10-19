package larger

import (
	test "midje/sweet"
)

test.fact("can concatenate strings", {
	greeting := "Hello "
	name := "Eamonn"
	str(greeting,  name)
},     =>, "Hello Eamonn"
)

test.fact("can use infix when calling two-parameter-function", {
	greeting := "Hello "
	name := "Eamonn"
	greeting  str  name
},      =>, "Hello Eamonn"
)
