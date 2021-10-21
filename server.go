package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	mapLock   sync.RWMutex
	Message   chan string
	OnlineMap map[string]*User
}

//新建服务器
func NewServer(ip string, port int) (server *Server) {
	server = &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return
}

//启动服务器
func (server *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("Listener listen err:", err)
		return
	}
	fmt.Println("服务器建立成功")
	//close listen socket
	defer listener.Close()
	go server.ListenMessage()
	//accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accept err:", err)
			continue
		}
		//转到处理go程进行逻辑处理

		go server.Handle(conn)
	}
	//do handle

}

//监听Message并发送

func (server *Server) ListenMessage() {
	for {
		msg := <-server.Message
		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.mapLock.Unlock()

	}
}

//广播消息
func (server *Server) Broadcast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s:%s", user.Addr, user.Name, msg)
	server.Message <- sendMsg

}
func (server *Server) Handle(conn net.Conn) {
	fmt.Println("连接成功")
	user := NewUser(server, conn)
	//广播用户已上线
	user.Online()
	//接收用户的信息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn read err:", err)
				return
			}
			//提取用户的消息
			msg := string(buf)
			user.DoMessage(msg)
		}
	}()

}
