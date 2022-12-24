package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //client模式
}

func NewClient(ip string, port int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
		flag:       999,
	}
	// 连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("net Dial err:", err)
		return nil
	}

	client.conn = conn

	// 返回对象
	return client
}

func (client *Client) DealMessage() {
	// 一旦client.conn有数据，就直接copy至Stdout，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int
	var temp string

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&temp)

	flag, err := strconv.Atoi(temp)
	if err != nil {
		fmt.Println(">>>>>>请输入合法数字<<<<<<")
		return false
	}

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>请输入合法数字<<<<<<")
		return false
	}
}

func (client *Client) UpdateUsername() bool {
	fmt.Println(">>>>>请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"

	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err", err)
		return false
	}

	return true
}

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>>>>>>>>>请输入聊天内容, exit退出")
	fmt.Scanln(&chatMsg)
	for (chatMsg) != "exit" {

		// 消息不为空则发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>>>>>>请输入聊天内容, exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>>>>>>请输入要私聊的用户名称:")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>>>>>>>>请输入消息内容，exit退出:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// 消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write err", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>>>>>>>>>请输入消息内容，exit退出:")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println(">>>>>>>>>请输入要私聊的用户名称:")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		// 根据不同模式处理业务
		switch client.flag {
		case 1:
			// 公聊模式
			client.PublicChat()
			break
		case 2:
			// 私聊模式
			client.PrivateChat()
			break
		case 3:
			// 更新用户名
			client.UpdateUsername()
			break
		case 0:
			// 退出
		}

	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "47.98.154.163", "设置服务器IP地址")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认8888)")
	// 设置命令行解析
	flag.Parse()
}

func main() {

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("服务器连接失败>>>>>>>>>>>>>>")
		return
	}

	// 单独创建一个goroutine去处理server的回执消息
	fmt.Println("服务器连接成功>>>>>>>>>>>>>>")

	go client.DealMessage()

	// 启动客户端的业务
	client.Run()
}
