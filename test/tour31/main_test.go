package main_test

import(
        test "midje/sweet"
	"tour31/main"
)

test.fact("tour31",
	withOutStr(main.main()),
	=>, `Hello World
[Hello World]
`)
