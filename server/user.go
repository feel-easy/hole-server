package server

import (
	"net"
	"strings"

	"github.com/feel-easy/hole-server/global"
	"go.uber.org/zap"
)

type User struct {
	Name     string `yaml:"name"`
	Addr     string `yaml:"addr"`
	PassWord string `yaml:"password"`
	Email    string `yaml:"email"`
	Logined  bool
	C        chan string
	conn     net.Conn
	// 用户对应的服务器
	server *Server
}

// 监听channel功能
func (u *User) ListenMessage() {
	for {
		mes := <-u.C
		// 不接收 自己发的消息
		if strings.HasPrefix(mes, u.Name) {
			continue
		}
		global.LOG.Info("消息", zap.String("msg", mes))
		u.conn.Write([]byte(mes + "\n"))
	}
}

// 上线功能
func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
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

func (u *User) whos() {
	for _, client := range u.server.OnlineMap {
		if !client.Logined {
			continue
		}
		onlineMsg := "[" + client.Addr + "]" + client.Name + "在线"
		u.SendMsgSimple(onlineMsg)
	}
}

func (u *User) rename(newName string) {
	// 规定当接收到rename|新名称指令时，更改用户名
	// fmt.Println(newName)
	_, ok := u.server.OnlineMap[newName]
	if ok {
		u.SendMsgSimple("当前用户名已被占用")
		return
	} else {
		u.server.mapLock.Lock()
		delete(u.server.OnlineMap, u.Name)
		u.server.OnlineMap[newName] = u
		u.server.mapLock.Unlock()
		u.SendMsgSimple("更新用户名成功")
		u.Name = newName
	}
}

func (u *User) to(remoteName, remotemsg string) {
	// 规定私聊的信息为 to|用户名|消息
	// 获取用户名
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
	if remotemsg == "" {
		u.SendMsgSimple("发送消息不能为空")
		return
	} else {
		remoteUser.SendMsgSimple(u.Name + "对您说:" + remotemsg)
	}
}

func (u *User) login(username, password string) {
	user, ok := u.server.OnlineMap[username]
	if ok && user.PassWord == password {
		u.Logined = true
		u.Online()
		return
	}
	u.SendMsgSimple("账号密码错误")
}

func (u *User) register(msg []string) {
	// _, email, username, password := msg...
}

// 消息处理
func (u *User) DoMessage(msg string) {
	msgList := strings.Split(msg, "|")
	if len(msgList) == 0 {
		u.SendMsgSimple("您的消息格式有误")
		return
	}
	if !u.Logined && msgList[0] != "login" && msgList[0] != "register" {
		u.SendMsgSimple(`客户端未登录，请登录或着注册
			登录：login|用户名｜密码
			注册：register|邮箱|用户名｜密码
			`)
		return
	}
	switch msgList[0] {
	default:
		u.server.BroadCast(u, msg)
	case "whos":
		u.whos()
	case "login":
		u.login(msgList[1], msgList[2])
	case "register":
		u.register(msgList)
	case "rename":
		u.rename(msgList[1])
	case "to":
		u.to(msgList[1], msgList[2])
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
