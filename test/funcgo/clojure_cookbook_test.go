package funcgo.clojureCookbookTest
import(
	test midje.sweet
	fgo funcgo.core
)

func add(x,y) {
	+(x, y)
}
test.fact( add(1,2), =>, 3)

	
