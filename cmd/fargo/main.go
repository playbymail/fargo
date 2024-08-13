// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package main implements the fargo game engine.
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
	Use:   "fargo",
	Short: "fargo game engine",
	Long:  `Run the fargo game engine.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Hello from root command\n")
	},
}
