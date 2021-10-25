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
	conn       net.Conn
	chartModel bool
	charObj    *User
}

//获取一个User
func NewUser(server *Server, conn net.Conn) (user *User) {
	userAddr := conn.RemoteAddr().String()
	user = &User{
		Name:       userAddr,
		Addr:       userAddr,
		C:          make(chan string),
		conn:       conn,
		server:     server,
		chartModel: false,
		charObj:    nil,
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

	//私聊模式
	if user.chartModel {
		if msg == "break" {
			user.chartModel = false
			user.charObj = nil
			user.SendMsg("quit chart model")
		} else {
			user.charObj.SendMsg(msg)
		}
		return
	}
	if msg == "who" {
		OnlineUsers := user.getOnlineUsers()
		OnlineUsersCnt := len(OnlineUsers)
		var sendMsg string
		if OnlineUsersCnt == 1 {
			sendMsg = fmt.Sprintf("There is only 1 user online:%s", OnlineUsers[0])
		} else {

			sendMsg = fmt.Sprintf("There are  %d users online:%s", OnlineUsersCnt, strings.Join(OnlineUsers, ","))
		}
		user.SendMsg(sendMsg)
	} else if len(msg) > 7 && msg[0:7] == "rename " {
		newName := strings.Split(msg, " ")[1]
		user.Rename(newName)
	} else if len(msg) > 7 && msg[0:5] == "chart" {
		args := strings.Split(msg, " ")
		if len(args) == 3 {
			targetName := args[1]
			chartmsg := args[2]
			user.CharTo(targetName, chartmsg)
		} else if len(args) == 2 {
			targetName := args[1]
			_, ok := user.server.OnlineMap[targetName]
			if !ok {
				user.SendMsg("Cannot find " + targetName + "\n")
				return
			}
			user.chartModel = true
			user.charObj = user.server.OnlineMap[targetName]
			user.SendMsg("enter private chart successfully")
		}

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
func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg + "\n\r"))
}
func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n\r"))
	}

}
func (user *User) Rename(newName string) bool {
	_, ok := user.server.OnlineMap[newName]
	if ok {
		user.SendMsg(newName + "has been used\n")
		return false
	}
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.Name = newName
	user.server.OnlineMap[newName] = user
	user.server.mapLock.Unlock()
	user.SendMsg("rename successfully:" + user.Name + "\n")
	return true
}
func (user *User) CharTo(targetName, msg string) bool {
	_, ok := user.server.OnlineMap[targetName]
	if !ok {
		user.SendMsg("Cannot find " + targetName + "\n")
		return false
	}
	targetUser := user.server.OnlineMap[targetName]
	targetUser.SendMsg(user.Name + ":" + msg)
	return true
}
