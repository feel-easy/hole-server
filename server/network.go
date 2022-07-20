package server

import (
	"github.com/feel-easy/hole-server/state"
	"github.com/feel-easy/hole-server/utils/logs"
	"github.com/feel-easy/hole-server/utils/protocol"
)

// Network is interface of all kinds of network.
type Network interface {
	Serve() error
}

func handle(rwc protocol.ReadWriteCloser) error {
	// 给新进入的用户分配资源
	c := protocol.Wrapper(rwc)
	defer func() {
		err := c.Close()
		if err != nil {
			logs.Error(err)
		}
	}()
	logs.Info("new user connected! ")
	authInfo, err := loginAuth(c)
	if err != nil || authInfo.ID == 0 {
		_ = c.Write(protocol.ErrorPacket(err))
		return err
	}
	user := Connected(c, authInfo)
	logs.Infof("user auth accessed, ip %s, %d:%s\n", user.IP, authInfo.ID, authInfo.Name)
	go state.Run(user)
	defer user.Offline()
	return user.Listening()
}
