/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"flag"
	"log"

	"github.com/feel-easy/hole-server/server"
	"github.com/feel-easy/hole-server/utils"
)

var (
	wsPort  string
	tcpPort string
)

func init() {
	flag.StringVar(&wsPort, "w", "9998", "WebsocketServer Port")
	flag.StringVar(&tcpPort, "t", "9999", "TcpServer Port")
}

func main() {
	flag.Parse()
	utils.Async(func() {
		wsServer := server.NewWebsocketServer(":" + wsPort)
		log.Panic(wsServer.Serve())
	})

	tcpServer := server.NewTcpServer(":" + tcpPort)
	log.Panic(tcpServer.Serve())
}
