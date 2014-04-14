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

func onBoard(size) {
	func(yx) {
		func{-1 < .. && .. < size} isEvery yx
	}
}

func neighbors(size, yx) {
	neighbors([[-1,0], [1,0], [0,-1], [0,1]], size, yx)
} (deltas, size, yx) {
	const (
		addYx = func{map(+, yx, ..)}
		unfiltered = addYx map deltas
	)
	onBoard(size) filter unfiltered
}

test.fact("neighbors works",
	func{matrix getIn ..} map neighbors(3, [0,0]),
	=>, [4,2]
)
	

var pool java.util.concurrent.ExecutorService = java.util.concurrent.Executors::newFixedThreadPool(
	2 + Runtime::getRuntime()->availableProcessors()
)

func mutateDothreads(f, {threadCount:THREADS, execCount:TIMES} ) {
	for _ := times threadCount {
		const multipleCalls Runnable = func{
			for _ := times execCount { f() }
		}
		pool->submit(multipleCalls)
	}
}

initialBoard := [
	[EE, KW, EE],
	[EE, EE, EE],
	[EE, KB, EE]
]

func boardMap(f, bd) {
	vec(
		func{vec(for s := lazy .. { f(s) })} map bd
	)
}

func doReset() {
	board := boardMap(ref, initialBoard)
	toMove := ref([[KB, [2, 1]], [KW, [0,1]]])
	//toMove := &[[KB, [2, 2]], [KW, [0,1]]]
	numMoves := ref(0)
	//numMoves := &0
}

func kingMoves(yx){
	neighbors(
		[[-1,-1], [-1,0], [-1,1], [0,-1], [0,1], [1,-1], [1,0], [1,1]],
		3,
		yx
	)
}

func isGoodMove(to, enemySq){
	if (to != enemySq) {
		to
	}
}

rotateCount := ref(0)
func fakeShuffle(xs) {
	dosync(rotateCount alter inc)
	const shift = (*rotateCount) % count(xs)
	(shift drop xs) concat (shift take xs)
}

// Fake shuffle to make test deterministic
shuffle := fakeShuffle

test.fact("fake shuffle is actually rotate",
	shuffle([111,222,333,444]), =>, [222,333,444,111],
	shuffle([111,222,333,444]), =>, [333,444,111,222],
	shuffle([111,222,333,444]), =>, [444,111,222,333],
	shuffle([111,222,333,444]), =>, [111,222,333,444],
	shuffle([111,222,333,444]), =>, [222,333,444,111],
	shuffle([111,222,333,444]), =>, [333,444,111,222],
	shuffle([111,222,333,444]), =>, [444,111,222,333],
	shuffle([111,222,333,444]), =>, [111,222,333,444]
)


func chooseMove([[mover, mpos], [_, enemyPos]]) {
	[
		mover,
		func{.. isGoodMove enemyPos} some shuffle(kingMoves(mpos))
	]
}

doReset()
test.fact("initial state",
	boardMap(deref, board),
	=>, [
		[EE, KW, EE],
		[EE, EE, EE],
		[EE, KB, EE]
	],
	*toMove,
	=>, [[KB, [2, 1]], [KW, [0,1]]],

        *numMoves,
	=>, 0
)


test.fact("Coordinated, synchronous change using alter",
	5 take repeatedly(func{chooseMove(*toMove)}),
	=>,  [  // starting at [KB, [2,1]]
		[KB, [2,2]],  
		[KB, [1,0]],
		[KB, [1,1]],
		[KB, [1,2]],
		[KB, [2,0]]
	]
)



func place(from, to){to}

func movePiece([piece, dest], [[_, src], _]) {
	getIn(board, dest) alter func{place(.., piece)}
	getIn(board, src) alter func{place(.., EE)}
	numMoves alter inc
}

func updateToMove(move) {
	toMove alter func{vector(second(..), move)}
}

func makeMove() {
	dosync(
		const move = chooseMove(*toMove)
		movePiece(move, *toMove)
		updateToMove(move)
	)
}

doReset()

test.fact("using alter to update a Ref",
	makeMove(),
	=>, [[KW, [0,1]], [KB, [2,2]]],

	boardMap(deref, board),
	=>, [
		[EE, KW, EE],
		[EE, EE, EE],
		[EE, EE, KB]
	],

	*numMoves,
	=>, 1
)
