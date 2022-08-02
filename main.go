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
	webPort string
)

func init() {
	flag.StringVar(&webPort, "web", "9996", "TcpServer Port")
	flag.StringVar(&wsPort, "ws", "9998", "WebsocketServer Port")
	flag.StringVar(&tcpPort, "tcp", "9999", "TcpServer Port")
}

func main() {
	flag.Parse()
	utils.Async(func() {
		webServer := server.NewWebServer(":" + webPort)
		log.Panic(webServer.Serve())
	})

	utils.Async(func() {
		wsServer := server.NewWebsocketServer(":" + wsPort)
		log.Panic(wsServer.Serve())
	})

	tcpServer := server.NewTcpServer(":" + tcpPort)
	log.Panic(tcpServer.Serve())
}
