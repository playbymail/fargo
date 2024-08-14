// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/spf13/cobra"
)

var argsCreate = struct {
	seed string
}{}

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a new cluster or other game object",
	Long:  `Create a new game with the supplied parameters.`,
}
