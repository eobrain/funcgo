package main_test

import(
        test "midje/sweet"
	"tour26/main"
)

test.fact("tour26",
	withOutStr(main.main()),
	=>, `{1 2}
`)
