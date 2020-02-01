package main

import (
	"nc_chat/server"
)

func main() {
	s := server.Default()
	s.Run()
}
