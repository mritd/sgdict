package cmd

import (
	"github.com/mritd/sgdict/pkg/download"

	"github.com/spf13/cobra"
)

var downloadDir string

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download dict",
	Run: func(cmd *cobra.Command, args []string) {
		download.Run(downloadDir)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.PersistentFlags().StringVar(&downloadDir, "dir", "./dict", "sougou dict scel file download dir")
}
