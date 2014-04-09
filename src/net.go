package net

func Dial(network, address) {
	const(
		socket = new java.net.Socket("127.0.0.1", 1234)
		out = new java.io.PrintWriter(socket->getOutputStream(), true)
		in = new java.io.BufferedReader(new java.io.InputStreamReader(socket->getInputStream()))
	)
	[in, out]
}
