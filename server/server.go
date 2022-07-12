package server

// Server是服务器端，负责监听端口，当server启动后，就会开始轮询监听，有连接进来时，就会分出一个协程单独处理业务
import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// 用户接受信息的缓冲区长度
var Length int = 4096

// server有两个成员变量IP和端口
type Server struct {
	IP   string
	port string

	// 在线用户列表与其 读写锁
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// 创建server
func NewServer(ip, port string) *Server {
	return &Server{
		IP:        ip,
		port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// 监听广播消息channel的协程，有消息则发送出去
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, value := range s.OnlineMap {
			value.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// 广播当前消息的方法
func (s *Server) BroadCast(user *User, msg string) {
	// 对OnlineMap进行消息广播
	sendMsg := fmt.Sprintf("%s：->%s", user.Name, msg)
	s.Message <- sendMsg
}

// server的socket编程，这时表示用户上线了
func (s *Server) Handler(conn net.Conn) {
	// 对于一个新加入的连接，创建一个新的User
	user := GetUser(conn, s)
	// 定义一个作为flag的channel，每当有消息写入时，向channel传入值
	isAlive := make(chan int)
	// 接收用户的信息，新开一个go程
	go func() {
		buf := make([]byte, Length)
		for {
			// socket读入Client端信息
			n, err := conn.Read(buf)
			if n == 0 {
				// 当n=0时说明连接关闭了，则令其下线，有时好像无法读出，也许需要设置心跳计数器
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn 读入出错")
				return
			}
			// 以上条件均不满足，则正常读入，同时去除在末尾的可能的回车符
			msg := string(buf[:n-1])
			user.DoMessage(msg)
			isAlive <- 1
		}
	}()
	// 某线程死亡时，所有子线程也会死亡，所以该线程需要阻塞
	for {
		select {
		case <-isAlive:
			// 这里什么都不用写，根据select的语法，isAlive读入时，所有channel都会被执行/刷新
		case <-time.After(time.Minute * 30):
			{
				// 将用户下线，并从服务器端释放资源
				user.SendMsgSimple("您已下线")
				user.conn.Close()
				close(user.C)
				s.mapLock.Lock()
				delete(s.OnlineMap, user.Name)
				s.mapLock.Unlock()
				// 退出当前Handler
				return
			}
		}
	}

}

// 开启Server，轮询监听
func (s *Server) Start() {
	// Listen方法有两个参数，第一个参数为协议，第二个参数为ip+port
	Listener, err := net.Listen("tcp", s.IP+":"+s.port)
	fmt.Printf("服务启动 tcp %s:%s ", s.IP, s.port)
	// 如果没有正常获取Listener
	if err != nil {
		fmt.Println("Listener has err!")
		return
	}
	// 在结束时关闭Listener
	defer Listener.Close()

	// 启动Message监听协程
	go s.ListenMessage()

	// 使用for死循环 循环监听Listener
	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("conn has err")
			continue
		}
		go s.Handler(conn)
	}

}
