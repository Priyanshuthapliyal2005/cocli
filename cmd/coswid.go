// filepath: /d:/opensource/cocli/cmd/coswid.go
package cmd

import (
    "github.com/spf13/cobra"
)

var coswidCmd = &cobra.Command{
    Use:   "coswid",
    Short: "CoSWID manipulation",

    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            cmd.Help() // nolint: errcheck
        }
    },
}

func init() {
    rootCmd.AddCommand(coswidCmd)
}