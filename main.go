package main

import "github.com/l1huanyu/nomcar/port"

func main() {
	h := port.NewHTTPHandler()
	h.Run()
}
