package main

import (
	"fmt"
	"net"

	"golang.org/x/text/encoding/simplifiedchinese"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn //why lower : user send message itself

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
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

func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	u.server.BroadCast(u, "当前用户已上线")
}

func (u *User) Offline() {
	u.server.BroadCast(u, "下线")
}

func (u *User) SendMsg(msg string) {
	gbkBytes, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(msg + "\n"))
	u.conn.Write([]byte(gbkBytes))
	//u.conn.Write([]byte(msg))
	//to meet the gbk...
}

// the service of handling user message
func (u *User) DoMessage(msg string) {
	fmt.Printf("%s send msg:%s\n", u.Name, msg)
	if msg == "who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else {
		u.server.BroadCast(u, msg)
	}
}
