package async
import test "midje/sweet"

func printAfterDelay(s) {
	Thread::sleep(100)
	print(s)
}

test.fact("can use goroutines to execute code in parallel",

	// no parallelism
	withOutStr(printAfterDelay("foo")), =>, "foo",

	// don't wait for results
	withOutStr({
		go printAfterDelay("bar")
	}), =>, "",

	// wait for result
	withOutStr({
		go printAfterDelay("baz")
		Thread::sleep(200)
	}), =>, "baz"
)

func sum(x1, x2, c) {
	c <- x1 + x2
}

func primes(c) {
	c <- 2
	c <- 3
	c <- 5
	c <- 7
	c <- 11
}

test.fact("can read and write channels in parallel",

	{
		c := make(chan)
		go sum(3, 4, c)
		<-c
	}, =>, 7,

	{
		c := make(chan)
		go primes(c)
		[<-c, <-c, <-c, <-c]
	}, =>, [2, 3, 5, 7],

	{
		c := make(chan)
		go primes(c)
		{
			c2 := make(chan)
			go func{
				c2 <- [<-c, <-c, <-c, <-c]
			}()
			<-c2
		}
	}, =>, [2, 3, 5, 7]

)

test.fact("can read and write channels in parallel using buffered channels",

	{
		c := make(chan, 10)
		go sum(3, 4, c)
		<-c
	}, =>, 7,

	{
		c := make(chan, 10)
		go primes(c)
		[<-c, <-c, <-c, <-c]
	}, =>, [2, 3, 5, 7],

	{
		c := make(chan, 10)
		go primes(c)
		{
			c2 := make(chan, 10)
			go func{
				c2 <- [<-c, <-c, <-c, <-c]
			}()
			<-c2
		}
	}, =>, [2, 3, 5, 7]

)

test.fact("can read and write channels in parallel using lightweight processes",

	{
		c := make(chan, 10)
		go { c <: 3 + 4 }
		<-c
	}, =>, 7,

	{
		c := make(chan, 10)
		go {
			c <: 2
			c <: 3
			c <: 5
			c <: 7
			c <: 11
		}
		[<-c, <-c, <-c, <-c]
	}, =>, [2, 3, 5, 7],

	{
		c := make(chan, 10)
		go primes(c)
		{
			c2 := make(chan, 10)
			go {
				c2 <: [<:c, <:c, <:c, <:c]
			}
			<-c2
		}
	}, =>, [2, 3, 5, 7]

)

func fibonacci(c, quit) {
	loop(x=0, y=1){
		select {
		case c <- x:
			recur(y, x + y)
		case <-quit:
			println("quit")
		}
	}
}



test.fact("can use select to block on multiple things",

	withOutStr({
		c := make(chan int)
		quit := make(chan int)
		go func() {
			for i := 0; i < 10; i++ {
				println(<-c)
			}
			quit <- 0
		}()
		fibonacci(c, quit)
	}), =>, `0
1
1
2
3
5
8
13
21
34
quit
`,

	withOutStr({
		c := make(chan int)
		quit := make(chan int)
		go {
			for i := 0; i < 10; i++ {
				println(<:c)
			}
			quit <: 0
		}
		fibonacci(c, quit)
	}), =>, `0
1
1
2
3
5
8
13
21
34
quit
`
)

func goFibonacci(c, quit) {
	loop(x=0, y=1){
		select {
		case c <: x:
			recur(y, x + y)
		case <:quit:
			println("quit")
		}
	}
}
