// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var coswidCmd = &cobra.Command{
	Use:   "coswid",
	Short: "coswid manipulation",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help() // nolint: errcheck
			os.Exit(0)
		}
	},
}

func init() {
	fmt.Println("Initializing coswid command")
	rootCmd.AddCommand(coswidCmd)
}
