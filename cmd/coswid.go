package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var coswidCmd = &cobra.Command{
    Use:   "coswid",
    Short: "CoSWID manipulation",
    // Removed the Run function to allow subcommand delegation
}

func init() {
    fmt.Println("Initializing coswid command")
    rootCmd.AddCommand(coswidCmd)
}