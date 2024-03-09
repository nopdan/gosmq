package server

import "testing"

func TestServer(t *testing.T) {
	Serve(7008, false, `D:\Code\go\gosmq\build`)
}
