package server

import (
	"net/http"

	"github.com/feel-easy/hole-server/utils/logs"
)

type Web struct {
	addr string
}

func NewWebServer(addr string) Web {
	return Web{addr: addr}
}

func (w Web) Serve() error {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	logs.Infof("Websocket server listener on %s\n", w.addr)
	return http.ListenAndServe(w.addr, nil)
}
