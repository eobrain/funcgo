package funcgo.clojure_cookbook_test
import(
        test midje.sweet
        fgo funcgo.core
)

func add(x,y) {
        x + y
}
test.fact("Simple example",
        add(1,2),
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

