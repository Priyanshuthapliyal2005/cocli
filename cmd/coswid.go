package cmd

import (
    "github.com/spf13/cobra"
    "fmt"
)

var coswidCmd = &cobra.Command{
    Use:   "coswid",
    Short: "A brief description of your command",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("coswid command executed")
        return nil
    },
}

func init() {
    rootCmd.AddCommand(coswidCmd)
}