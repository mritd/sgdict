package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert scel to RIME format",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("convert called")
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
