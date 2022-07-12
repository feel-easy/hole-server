/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/feel-easy/hole-server/global"
	"github.com/feel-easy/hole-server/server"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "启动服务",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := server.NewServer(global.CONFIG.System.Host, global.CONFIG.System.Port)
		s.Start()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
