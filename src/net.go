package net
import type (
	java.net.Socket
	java.io.{BufferedReader, InputStreamReader, PrintWriter}
)

func Dial(network, address) {
	const(
		socket = new Socket("127.0.0.1", 1234)
		out = new PrintWriter(socket->getOutputStream(), true)
		in = new BufferedReader(new InputStreamReader(socket->getInputStream()))
	)
	[in, out]
}
