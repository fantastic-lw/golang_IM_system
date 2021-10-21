package main

import (
	"fmt"
	"net"
	"strings"
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
	//查询处理
	// fmt.Println(len(msg))
	if msg == "who" {
		OnlineUsers := user.getOnlineUsers()
		OnlineUsersCnt := len(OnlineUsers)
		var sendMsg string
		if OnlineUsersCnt == 1 {
			sendMsg = fmt.Sprintf("There is only 1 user online:%s", OnlineUsers[0])
		} else {

			sendMsg = fmt.Sprintf("There are  %d users online:%s", OnlineUsersCnt, strings.Join(OnlineUsers, ","))
		}
		user.WriteToClient(sendMsg)
	} else {
		user.server.Broadcast(user, msg)
	}

}
func (user *User) getOnlineUsers() (res []string) {
	user.server.mapLock.Lock()
	for name := range user.server.OnlineMap {
		res = append(res, name)
	}
	user.server.mapLock.Unlock()
	return
}
func (user *User) WriteToClient(msg string) {
	user.conn.Write([]byte(msg + "\n\r"))
}
func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n\r"))
	}

}
