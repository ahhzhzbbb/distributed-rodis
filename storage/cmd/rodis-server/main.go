package main

import (
	_ "net/http/pprof"
	"rodis/storage/internal/server"
)

func main() {
	s := server.NewServer(server.DefaultConfig())
	s.Start()
}
