package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

//新建服务器
func NewServer(ip string, port int) (server *Server) {
	server = &Server{
		Ip:   ip,
		Port: port,
	}
	return
}

//启动服务器
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("Listener listen err:", err)
		fmt.Println("服务器创立成功")
		return
	}
	//close listen socket
	defer listener.Close()
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

func (this *Server) Handle(conn net.Conn) {
	fmt.Println("连接成功")
}
