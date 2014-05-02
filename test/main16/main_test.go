package main_test

import(
        test "midje/sweet"
	"main16/main"
)

test.fact("main16",
	withOutStr(main.main()),
	=>, `Hello 世界
Happy 3.14 Day
Go rules? true
`)
