package main

import (
	"esp8266/src/server"
	"time"
)

func main() {
	server.Serve()
	time.Sleep(3*time.Second)
}
