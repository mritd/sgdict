package cmd

import (
	"github.com/mritd/sgdict/pkg/converter"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert scel to RIME format",
	Run: func(cmd *cobra.Command, args []string) {
		converter.Convert()
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.PersistentFlags().StringVar(&converter.BaseDir, "dir", "./dict", "sougou dict scel file download dir")
}
