package main

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the gateway",
	Run:   StartServer,
}

var Port string

func init() {
	// add the configuration paramters for the start command
	startCmd.Flags().StringVarP(&Port, "port", "p", "8080", "the port to listen on.")

	// add the start command to the root executable
	rootCmd.AddCommand(startCmd)
}

// StartServer begins an http server running the gateway
func StartServer(cmd *cobra.Command, args []string) {
	ListenAndServe()
}
