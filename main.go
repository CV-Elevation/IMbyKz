package main

import (
	"flag"
	"fmt"
)

var osPlatform string
var charset string

func main() {
	flag.StringVar(&osPlatform, "os", "windows", "os")
	flag.StringVar(&charset, "charset", "GBK", "charset")
	flag.Parse()
	fmt.Println(osPlatform, charset)

	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
