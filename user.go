package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func (user *User) Online() {
	// 用户上线，将用户加入OnlineMap
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	// 广播当前上线消息
	user.server.Broadcast(user, "已上线")
}

func (user *User) Offline() {
	//用户下线，将用户从OnlineMap删除
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	// 广播当前下线消息
	user.server.Broadcast(user, "已下线")
}

// 给当前User对应的客户端发送消息
func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (user *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户
		user.server.mapLock.Lock()
		for _, onlineUser := range user.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + onlineUser.Name + ":" + "在线...\n"
			user.SendMsg(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式: rename|name
		newName := strings.Split(msg, "|")[1]
		if _, ok := user.server.OnlineMap[newName]; ok {
			user.SendMsg("当前用户名已被使用\n")
		} else {
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newName] = user
			user.server.mapLock.Unlock()

			user.Name = newName
			user.SendMsg("您已经更新用户名为:" + user.Name + "\n")
		}
	} else if len(msg) > 3 && msg[:3] == "to|" {
		//消息格式：to|name|消息

		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			user.SendMsg("消息格式不正确，请使用\"to|name|message格式\"\n")
			return
		}

		remoteUser, ok := user.server.OnlineMap[remoteName]
		if !ok {
			user.SendMsg("该用户名不存在，请重发\n")
			return
		}

		remoteMsg := strings.Split(msg, "|")[2]
		if remoteMsg == "" {
			user.SendMsg("发送内容为空，请重发\n")
		}

		remoteUser.SendMsg(user.Name + "对你说:" + remoteMsg)
	} else {
		user.server.Broadcast(user, msg)
	}
}

// 建立一个新用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

// 监听当前user channel的方法
func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}
