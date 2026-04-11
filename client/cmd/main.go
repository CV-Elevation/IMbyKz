package main

import (
	"flag"
	"fmt"

	"github.com/CV-Elevation/IMbyKz/client"
)

var serverIp string
var port int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set server ip")
	flag.IntVar(&port, "port", 8888, "set server port")
}

func main() {
	flag.Parse()
	c := client.NewClient(serverIp, port)
	if c == nil {
		fmt.Println(">>>>>>>connect to server error......")
		return
	}
	fmt.Println(">>>>>>>connect to server succeed......")

}
