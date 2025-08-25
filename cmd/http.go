package cmd

import (
	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:     "http",
	Aliases: []string{},
	Short:   "HTTP tools",
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
