package server

import (
	"net/http"

	"github.com/feel-easy/hole-server/utils/logs"
	"github.com/feel-easy/hole-server/utils/protocol"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

type Message struct {
	Message string `json:"message"`
}

type Websocket struct {
	addr string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebsocketServer(addr string) Websocket {
	return Websocket{addr: addr}
}

func (w Websocket) Serve() error {
	http.HandleFunc("/ws", serveWs)
	logs.Infof("Websocket server listener on %s\n", w.addr)
	return http.ListenAndServe(w.addr, nil)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Error(err)
		return
	}
	err = handle(protocol.NewWebsocketReadWriteCloser(conn))
	if err != nil {
		logs.Error(err)
	}
}
