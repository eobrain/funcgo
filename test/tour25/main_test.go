package main_test

import(
        test "midje/sweet"
	"tour25/main"
)

test.fact("tour25",
	withOutStr(main.main()),
	=>, `10.000000000000007
`)
