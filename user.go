package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

//获取一个User
func NewUser(conn net.Conn) (user *User) {
	userAddr := conn.RemoteAddr().String()
	user = &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	//启动监听当前user channel的go程
	go user.ListenMessage()
	return
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
		fmt.Println("Send successfully")
	}

}
