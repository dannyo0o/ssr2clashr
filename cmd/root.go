package cmd

import (
	"os"

	"github.com/heiha/ssr2clashr/config"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ssr2clashr",
		Short: "ssr to clashr",
		Long:  `将 ssr配置 转换为 clashr 配置`,
	}
)

func init() {
	cobra.MousetrapHelpText = ""
	cobra.OnInitialize(config.Execute)
}

// Execute func()
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
