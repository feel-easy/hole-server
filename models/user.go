package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/hole-server/utils"
	"github.com/feel-easy/hole-server/utils/logs"
	"github.com/feel-easy/hole-server/utils/protocol"
)

type User struct {
	ID     int    `json:"id"`
	IP     string `json:"ip"`
	Name   string `json:"name"`
	Score  int    `json:"score"`
	Mode   int    `json:"mode"`
	Type   int    `json:"type"`
	RoomID int    `json:"roomId"`

	conn   *protocol.Conn
	data   chan *protocol.Packet
	read   bool
	state  consts.StateID
	online bool
}

func (u *User) UserID() int {
	return u.ID
}

func (u *User) UserName() string {
	return u.Name
}

func (u *User) Write(bytes []byte) error {
	return u.conn.Write(protocol.Packet{
		Body: bytes,
	})
}

func (u *User) Offline() {
	u.online = false
	_ = u.conn.Close()
	close(u.data)
	room := getRoom(u.RoomID)
	if room != nil {
		room.Lock()
		defer room.Unlock()
		room.broadcast(fmt.Sprintf("%s lost connection! \n", u.Name))
		if room.State == consts.Waiting {
			room.removeUser(u)
		}
		room.Cancel()
	}
}

func (u *User) Listening() error {
	for {
		pack, err := u.conn.Read()
		if err != nil {
			logs.Error(err)
			return err
		}
		if u.read {
			u.data <- pack
		}
	}
}

// 向客户端发生消息
func (u *User) WriteString(data string) error {
	return u.conn.Write(protocol.Packet{
		Body: []byte(data),
	})
}

func (u *User) WriteObject(data interface{}) error {
	return u.conn.Write(protocol.Packet{
		Body: utils.MarshalJson(data),
	})
}

func (u *User) WriteError(err error) error {
	if err == consts.ErrorsExist {
		return err
	}
	return u.conn.Write(protocol.Packet{
		Body: []byte(err.Error() + "\n"),
	})
}

func (u *User) AskForPacket(timeout ...time.Duration) (*protocol.Packet, error) {
	u.StartTransaction()
	defer u.StopTransaction()
	return u.askForPacket(timeout...)
}

func (u *User) askForPacket(timeout ...time.Duration) (*protocol.Packet, error) {
	var packet *protocol.Packet
	if len(timeout) > 0 {
		select {
		case packet = <-u.data:
		case <-time.After(timeout[0]):
			return nil, consts.ErrorsTimeout
		}
	} else {
		packet = <-u.data
	}
	if packet == nil {
		return nil, consts.ErrorsChanClosed
	}
	single := strings.ToLower(packet.String())
	if single == "exit" {
		return nil, consts.ErrorsExist
	}
	return packet, nil
}

func (u *User) AskForInt(timeout ...time.Duration) (int, error) {
	packet, err := u.AskForPacket(timeout...)
	if err != nil {
		return 0, err
	}
	return packet.Int()
}

func (u *User) AskForint(timeout ...time.Duration) (int, error) {
	packet, err := u.AskForPacket(timeout...)
	if err != nil {
		return 0, err
	}
	num, err := packet.Int64()
	return int(num), err
}

func (u *User) AskForString(timeout ...time.Duration) (string, error) {
	packet, err := u.AskForPacket(timeout...)
	if err != nil {
		return "", err
	}
	return packet.String(), nil
}

func (u *User) AskForStringWithoutTransaction(timeout ...time.Duration) (string, error) {
	packet, err := u.askForPacket(timeout...)
	if err != nil {
		return "", err
	}
	return packet.String(), nil
}

func (u *User) StartTransaction() {
	u.read = true
	_ = u.WriteString(consts.IsStart)
}

func (u *User) StopTransaction() {
	u.read = false
	_ = u.WriteString(consts.IsStop)
}

func (u *User) State(s consts.StateID) {
	u.state = s
}

func (u *User) GetState() consts.StateID {
	return u.state
}

func (u *User) Conn(conn *protocol.Conn) {
	u.conn = conn
	u.data = make(chan *protocol.Packet, 8)
	u.online = true
	setUser(u)
}

func (u *User) String() string {
	return fmt.Sprintf("%s[%d]", u.Name, u.ID)
}

func (user *User) BroadcastChat(msg string, exclude ...int) {
	logs.Infof("chat msg, user %s[%d] %s say: %s\n", user.Name, user.ID, user.IP, strings.TrimSpace(msg))
	Broadcast(user.RoomID, msg, exclude...)
}
