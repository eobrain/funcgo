package row_test

import (
        test "midje/sweet"
//	"funcgo/reference/contract"
	r "funcgo/reference/row"
)
ϵ := 1e-10
a := [2.0, 3.0, 4.0]
b := [3.0, 4.0, 5.0]

// contract.CheckPreconditions = true


test.fact("a vector supports addition",
	a  r.+  b, =>, [5.0, 7.0, 9.0]
)

test.fact("a vector supports dot product",
	a  r.*  b,  =>, test.roughly(6.0 + 12.0 + 20.0, ϵ)
)
