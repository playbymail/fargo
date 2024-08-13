// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package main implements the web server for fargo.
package main

import (
	"github.com/spf13/cobra"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	if err := Execute(); err != nil {
		log.Fatal(err)
	}
}

func Execute() error {
	cmdRoot.AddCommand(cmdVersion)

	return cmdRoot.Execute()
}

var cmdRoot = &cobra.Command{
	Use:   "fargod",
	Short: "fargo web server",
	Long:  `Run the fargo web server.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Hello from root command\n")
	},
}
