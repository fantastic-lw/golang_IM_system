package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	serverIp   string
	serverPort int
	conn       net.Conn
	model      int
}

func NewClient(ip string, port int) (client *Client) {
	client = &Client{
		serverIp:   ip,
		serverPort: port,
		conn:       nil,
		model:      999,
	}

	//connect to  the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {

		return nil
	}
	client.conn = conn

	return
}

var serverIp string
var severPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "the server ip you want to connect!")
	flag.IntVar(&severPort, "port", 8787, "the server port you want to connect!")
}
func main() {
	flag.Parse()
	client := NewClient(serverIp, severPort)
	if client == nil {
		fmt.Println("connect err!")
		return
	}
	fmt.Println("connect successfully!")
	go client.dealResponse()
	client.Run()
}
func (client *Client) dealResponse() {
	io.Copy(os.Stdout, client.conn)
}
func (client *Client) menu() bool {
	var model int
	for client.model != 0 {
		fmt.Println("1.公聊模式\n")
		fmt.Println("2.私聊模式\n")
		fmt.Println("3.修改用户名\n")
		fmt.Println("0.退出\n")
		fmt.Scanln(&model)
		if model >= 0 && model <= 3 {
			client.model = model
			return true
		} else {
			fmt.Println("invalid model")
			return false
		}

	}
	return false
}
func (client *Client) Run() {
	for client.model != 0 {
		for client.menu() != true {
		}
		switch client.model {

		case 1:
			fmt.Println("1.公聊模式\n")
			client.GroupChart()
			break
		case 2:
			fmt.Println("2.私聊模式\n")
			client.PrivateChart()
			break

		case 3:
			fmt.Println("3.修改用户名\n")
			client.updateName()
			break
		}

	}
}

func (client *Client) GroupChart() {

}
func (client *Client) PrivateChart() {

}
func (client *Client) updateName() bool {
	fmt.Println("input your new name!")
	newName := ""
	instru := "rename "
	fmt.Scanln(&newName)
	sendMsg := instru + newName
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client.conn err")
		return false
	}
	return true
}
