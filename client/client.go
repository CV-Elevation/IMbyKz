package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //user mode
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

func (c *Client) DealResponse() {
	//once read msg print to stdout
	io.Copy(os.Stdout, c.conn)

	// for {
	// 	buf := make([]byte, 4096)
	// 	c.conn.Read(buf)
	// 	fmt.Println(string(buf))
	// }
}

func (c *Client) menu() bool {
	var flag int = 999
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println(">>>>请输入合法范围的数字")
		return false
	}
}

func (c *Client) updateName() bool {
	fmt.Println(">>>>请输入用户名:")
	fmt.Scanln(&c.Name)

	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (c *Client) publicChat() {
	fmt.Println("请输入消息,exit表示退出:")
	var chatMsg string
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	chatMsg = strings.TrimRight(line, "\r\n")
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
		fmt.Println("请输入消息,exit表示退出")
		line, _ = reader.ReadString('\n')
		chatMsg = strings.TrimRight(line, "\r\n")
	}
}

func (c *Client) Run() {
	for c.flag != 0 {
		for {
			if c.menu() {
				break
			}
		}
		//handle different services
		switch c.flag {
		case 1:
			c.publicChat()
		case 2:
			fmt.Println("私聊模式")
		case 3:
			c.updateName()
		}
	}
}
