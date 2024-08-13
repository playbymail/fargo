// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package main implements the fargo game engine.
package main

import (
	"github.com/playbymail/fargo"
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

	cmdRoot.PersistentFlags().StringVar(&argsRoot.seed, "seed", "", "optional seed for the PRNG")

	if argsRoot.seed != "" {
		fargo.WithSeed(argsRoot.seed, true)
	}

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

var argsRoot = struct {
	e    *fargo.Engine
	seed string
}{
	seed: "0xdeadbeef^0xcafebabe",
}
