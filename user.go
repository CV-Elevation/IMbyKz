package main

import (
	"net"

	"golang.org/x/text/encoding/simplifiedchinese"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn //why lower : user send message itself
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}
	//start listener
	go user.ListenMessage()
	return user
}

// listen the conn,send to client
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		gbkBytes, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(msg + "\n"))
		u.conn.Write([]byte(gbkBytes))
		// u.conn.Write([]byte(msg + "\n"))
	}
}
