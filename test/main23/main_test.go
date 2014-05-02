package main_test

import(
        test "midje/sweet"
	"main23/main"
)

test.fact("main23",
	withOutStr(main.main()),
	=>, `9 20
`)
