package golang
import test "midje/sweet"

func printAfterDelay(s) {
	Thread::sleep(100)
	print(s)
}

test.fact("no parallelism",
	withOutStr(printAfterDelay("foo")), =>, "foo"
)

test.fact("don't wait",
	withOutStr({
		go printAfterDelay("bar")
	}), =>, ""
)


test.fact("wait",
	withOutStr({
		go printAfterDelay("baz")
		Thread::sleep(200)
	}), =>, "baz"
)

func sum(x1, x2, c) {
	c <- x1 + x2
}

test.fact("channel",
	{
		const c = chan
		go sum(3, 4, c)
		<-c
	}, =>, 7
)
