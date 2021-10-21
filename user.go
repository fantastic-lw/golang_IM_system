package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	//当前所连接的服务器
	server *Server
	//与client的连接
	conn net.Conn
}

//获取一个User
func NewUser(server *Server, conn net.Conn) (user *User) {
	userAddr := conn.RemoteAddr().String()
	user = &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	//启动监听当前user channel的go程
	go user.ListenMessage()
	return
}

//用户上下线功能
func (user *User) Online() {
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()
	user.server.Broadcast(user, "Online")
}

func (user *User) Offline() {
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()
	user.server.Broadcast(user, "Offline")
}

func (user *User) DoMessage(msg string) {
	user.server.Broadcast(user, msg)
}
func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n\r"))
	}

}
