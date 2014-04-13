package joy
import(
        test "midje/sweet"
)

matrix := [
	[1,2,3],
	[4,5,6],
	[7,8,9]
]

test.fact("getIn",
	matrix getIn [1,2],
	=>,  6
)

test.fact("assocIn",
	assocIn(matrix, [1,2], X),
	=>, [
		[1,2,3],
		[4,5,X],
		[7,8,9]
	]
)

test.fact("updateIn",
	updateIn(matrix, [1,2], *, 100),
	=>, [
		[1,2,3],
		[4,5,600],
		[7,8,9]
	]
)


func neighbors(size, yx) {
	neighbors([[-1,0], [1,0], [0,-1], [0,1]], size, yx)
} (deltas, size, yx) {
	func(newYx) {
		func{-1 < % && % < size} isEvery newYx
	} filter map(func{map(+, yx, %)}, deltas)
}

test.fact("neighbors works",
	func{matrix getIn %} map neighbors(3, [0,0]),
	=>, [4,2]
)
	

__pool__ := java.util.concurrent.Executors::newFixedThreadPool(
	2 + Runtime::getRuntime()->availableProcessors()
)

func mutateDothreads(f, {threadCount:THREADS, execCount:TIMES} ) {
	for _ := times threadCount {
		__pool__->submit( func(){
			for _ := times execCount { f() }
		})
	}
}


initialBoard := [
	[EE, KW, EE],
	[EE, EE, EE],
	[EE, KB, EE]
]

func boardMap(f, bd) {
	vec(
		func(_){
			vec(for s := lazy _ {f(s)})
		} map bd
	)
}

func boardMap(f, bd) {
	vec(
		func{
			vec(for s := lazy % {f(s)})
		} map bd
	)
}

func doReset() {
	board := boardMap(ref, initialBoard)
	toMove := ref([[BK, [2, 2]], [WK, [0,1]]])
	numMoves := ref(0)
}

kingMoves := partial(
	neighbors,
	[[-1,-1], [-1,0], [-1,1], [0,-1], [0,1], [1,-1], [1,0], [1,1]],
	3
)

func isGoodMove(to, enemySq){
	if (to != enemySq) {
		to
	}
}

func chooseMove([[mover, mpos], [_, enemyPos]]) {
	[
		mover,
		func(_){_ isGoodMove enemyPos} some shuffle(kingMoves(mpos))
	]
}

func chooseMoveNoShuffle([[mover, mpos], [_, enemyPos]]) {
	[
		mover,
		func(_){_ isGoodMove enemyPos} some kingMoves(mpos)
	]
}

// test.fact("chooseMove works",
// 	{
// 		doReset()
// 		5 take repeatedly(func(){chooseMoveNoShuffle(*toMove)})
// 	},
// 	=>,
// 	[
// 		[BK, [1,1]], [BK, [1,1]], [BK, [1,0]], [BK, [1,0]], [BK, [2,0]]
// 	]
// )
