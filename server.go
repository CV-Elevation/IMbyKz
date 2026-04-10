package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Server struct {
	Ip   string
	Port int

	//online map
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (s *Server) ListenManager() {
	for {
		msg := <-s.Message
		//remember the lock!
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	//service
	//fmt.Println("success")
	user := NewUser(conn, s)
	user.Online()

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			var msg string
			//windows
			//remove the char '\n'and '\r'
			if osPlatform == "windows" && charset == "GBK" {
				msg, _ = DecodeToGBK(buf[:n-2])
			} else {
				//linux mac
				msg = string(buf[:n-1])
			}

			user.DoMessage(msg)
		}

	}()

	//block
	select {}
}

func DecodeToGBK(input []byte) (string, error) {
	// 1. 创建一个解码器转换流 (GBK -> UTF-8)
	reader := transform.NewReader(bytes.NewReader(input), simplifiedchinese.GBK.NewDecoder())

	// 2. 读取所有转换后的数据
	output, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	//close
	defer listener.Close()
	//accept
	go s.ListenManager()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		//do handler
		go s.Handler(conn)
	}

}
