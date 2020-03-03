package cmd

import (
	"time"

	"github.com/mritd/sgdict/pkg/download"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download dict",
	Run: func(cmd *cobra.Command, args []string) {
		download.Run()
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.PersistentFlags().StringVar(&download.BaseDir, "dir", "./dict", "sougou dict scel file download dir")
	downloadCmd.PersistentFlags().DurationVar(&download.Timeout, "timeout", 10*time.Second, "http client timeout")
	downloadCmd.PersistentFlags().IntVar(&download.RetryCount, "retry", 5, "auto retry count")
	downloadCmd.PersistentFlags().DurationVar(&download.RetryMaxWaitTime, "retrywaittime", 5*time.Second, "retry max wait time")
}
