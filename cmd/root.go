package cmd

import (
	"encoding/base64"
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var bannerBase64 = "ICwgX18gICAgIF8gICxfXyBfXyAgICBfX18gICAgXyAgICAgICBfX18KL3wvICBcICAgfCB8L3wgIHwgIHwgIC8gKF8pICAoX3wgICB8IC8gKF8pICAoKSB8CiB8X19fLyAgIHwgfCB8ICB8ICB8ICBcX18gICAgICB8ICAgfCBcX18gICAgL1wgfAogfCBcICAgXyB8LyAgfCAgfCAgfCAgLyAgICAgICAgfCAgIHwgLyAgICAgLyAgXHwKIHwgIFxfL1xfL1wvIHwgIHwgIHxfL1xfX18vICAgICBcXy98L1xfX18vLyhfXy9vCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAvfAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgXHwK"

var versionTpl = `%s
Name: sgdict
Version: %s
Arch: %s
BuildDate: %s
CommitID: %s
`

var (
	Version   string
	BuildDate string
	CommitID  string
)

var rootCmd = &cobra.Command{
	Use:     "sgdict",
	Version: Version,
	Short:   "sougou dict spider",
	Run:     func(cmd *cobra.Command, args []string) { _ = cmd.Help() },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initLog)
	banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
	rootCmd.SetVersionTemplate(fmt.Sprintf(versionTpl, banner, Version, runtime.GOOS+"/"+runtime.GOARCH, BuildDate, CommitID))
}

func initLog() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
