// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/playbymail/fargo"
	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of this application",
	Long: `Displays the version number of this application.
Needed when opening bug reports.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", fargo.Version.String())
	},
}
