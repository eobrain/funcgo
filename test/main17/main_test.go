package main_test

import(
        test "midje/sweet"
	"main17/main"
)

test.fact("main17",
	withOutStr(main.main()),
	=>, `21
0.2
`)

//	=>, `21
//0.2
//1.2676506002282295e+29
//`)
