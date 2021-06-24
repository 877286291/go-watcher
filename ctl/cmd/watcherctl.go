package cmd

import (
	"github.com/spf13/cobra"
)

var auth Auth
var rootCmd = &cobra.Command{
	Use: "watcherctl",
	RunE: func(cmd *cobra.Command, args []string) error {
		//tcpRequest := NewTcpRequest("connection", "", auth)
		TerminalUI()
		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
	auth = DefaultAuth()
	rootCmd.Flags().StringVar(&auth.Host, "host", "127.0.0.1", "")
	rootCmd.Flags().IntVar(&auth.Port, "port", 2021, "")
	rootCmd.Flags().StringVar(&auth.UserName, "username", "", "")
	rootCmd.Flags().StringVar(&auth.Password, "password", "", "")
}
