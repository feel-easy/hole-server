/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"flag"
	"log"

	"github.com/feel-easy/hole-server/utils"
)

var (
	webPort string
)

func init() {
	flag.StringVar(&webPort, "web", "9996", "TcpServer Port")
}

func main() {
	flag.Parse()
	utils.Async(func() {
	})
	log.Println("hello")
}
