package server

import (
	"time"

	"github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/hole-server/models"
	"github.com/feel-easy/hole-server/utils"
	"github.com/feel-easy/hole-server/utils/logs"
	"github.com/feel-easy/hole-server/utils/protocol"
)

func Connected(conn *protocol.Conn, info *protocol.AuthInfo) *models.User {
	user := &models.User{
		ID:   info.ID,
		IP:   conn.IP(),
		Name: info.Name,
	}
	user.Conn(conn) // 初始化user对象
	return user
}

// 登陆验签
func loginAuth(c *protocol.Conn) (*protocol.AuthInfo, error) {
	authChan := make(chan *protocol.AuthInfo)
	defer close(authChan)
	utils.Async(func() {
		packet, err := c.Read()
		if err != nil {
			logs.Error(err)
			return
		}
		authInfo := &protocol.AuthInfo{}
		err = packet.Unmarshal(authInfo)
		if err != nil {
			logs.Error(err)
			return
		}
		authChan <- authInfo
	})
	select {
	case authInfo := <-authChan:
		return authInfo, nil
	case <-time.After(3 * time.Second):
		return nil, consts.ErrorsAuthFail
	}
}
