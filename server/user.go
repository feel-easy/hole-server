package server

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
	// 用户对应的服务器
	server *Server
}

// 监听channel功能
func (u *User) ListenMessage() {
	for {
		mes := <-u.C
		// 传入mes，并将其转换为字节数组
		fmt.Printf("%v", mes)
		u.conn.Write([]byte(mes + "\n"))
	}
}

// 上线功能
func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	// fmt.Printf("%v", u.server.OnlineMap)
	u.server.BroadCast(u, "上线")
}

// 下线功能
func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	u.server.BroadCast(u, "下线")
}

// 仅对当前客户端发送消息
func (u *User) SendMsgSimple(msg string) {
	u.conn.Write([]byte(msg + "\n"))
}

// 消息处理
func (u *User) DoMessage(msg string) {
	// 规定:当接收到的msg为whos时，显示当前有哪些用户在线
	if msg == "whos" {
		for _, client := range u.server.OnlineMap {
			onlineMsg := "[" + client.Addr + "]" + client.Name + "在线"
			u.SendMsgSimple(onlineMsg)
		}
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 规定当接收到rename|新名称指令时，更改用户名
		newName := strings.Split(msg, "|")[1]
		fmt.Println(newName)
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsgSimple("当前用户名已被占用")
			return
		} else {
			u.server.mapLock.Lock()
			// fmt.Printf("%v", u.server.OnlineMap)
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			// fmt.Printf("%v", u.server.OnlineMap)
			u.server.mapLock.Unlock()
			u.SendMsgSimple("更新用户名成功")
			u.Name = newName
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 规定私聊的信息为 to|用户名|消息
		// 获取用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMsgSimple("您的消息格式有误，请按照\"to|张三|你好\"的格式进行私聊消息发送")
			return
		}
		// 根据用户名获取User对象
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMsgSimple("对方用户不存在或已下线")
			return
		}
		// 调用对方用户的私发消息方法
		remotemsg := strings.Split(msg, "|")[2]
		if remotemsg == "" {
			u.SendMsgSimple("发送消息不能为空")
			return
		} else {
			remoteUser.SendMsgSimple(u.Name + "对您说:" + remotemsg)
		}
	} else {
		u.server.BroadCast(u, msg)
	}
}

// 当有新的TCP连接时，从该连接中获得User
func GetUser(conn net.Conn, server *Server) *User {
	// 从连接获得连接地址并转换为字符串
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// 启动监听
	go user.ListenMessage()
	return user
}
