package server

import (
	"net"

	"github.com/feel-easy/hole-server/utils"
	"github.com/feel-easy/hole-server/utils/logs"
	"github.com/feel-easy/hole-server/utils/protocol"
)

type Tcp struct {
	addr string
}

func NewTcpServer(addr string) Tcp {
	return Tcp{addr: addr}
}

func (t Tcp) Serve() error {
	listener, err := net.Listen("tcp", t.addr)
	if err != nil {
		logs.Error(err)
		return err
	}
	logs.Infof("Tcp server listening on %s\n", t.addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			logs.Infof("listener.Accept err %v\n", err)
			continue
		}
		utils.Async(func() {
			err := handle(protocol.NewTcpReadWriteCloser(conn))
			if err != nil {
				logs.Error(err)
			}
		})
	}
}
