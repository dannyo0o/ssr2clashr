package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VERSION string
var VERSION string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{

	Use:   "version",
	Short: "print version",
	Long:  `显示 ssr2clashr 的版本号（其实就是打包时间）`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ssr2clashr version:", VERSION)
	},
}
