package main

import (
	"fmt"
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
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("Listener listen err:", err)
		return
	}
	fmt.Println("服务器建立成功")
	//close listen socket
	defer listener.Close()
	go this.ListenMessage()
	//accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accept err:", err)
			continue
		}
		//转到处理go程进行逻辑处理

		go this.Handle(conn)
	}
	//do handle

}

//监听Message并发送

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()

	}
}

//广播消息
func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s:%s", user.Addr, user.Name, msg)
	this.Message <- sendMsg

}
func (this *Server) Handle(conn net.Conn) {
	fmt.Println("连接成功")
	user := NewUser(conn)
	//将用户加入OnlineMap中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	//广播用户已上线
	this.Broadcast(user, "has online\n")

}
