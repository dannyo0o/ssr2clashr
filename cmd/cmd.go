package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/heiha/ssr2clashr/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(cmdCmd)

	cmdCmd.PersistentFlags().StringP(
		"url",
		"u",
		viper.GetString("url"),
		"订阅链接",
	)

	cmdCmd.PersistentFlags().StringP(
		"path",
		"p",
		viper.GetString("path"),
		"文件存储路径",
	)

	cmdCmd.PersistentFlags().StringP(
		"template",
		"t",
		viper.GetString("template"),
		"套用的模板地址",
	)

	viper.BindPFlag("url", cmdCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("path", cmdCmd.PersistentFlags().Lookup("path"))
	viper.BindPFlag("template", cmdCmd.PersistentFlags().Lookup("template"))

}

var cmdCmd = &cobra.Command{

	Use:   "cmd",
	Short: "ssr2clashr cli",
	Long:  `直接生成配置文件`,
	Run: func(cmd *cobra.Command, args []string) {
		api.InitRules()
		body := api.Execute(viper.GetString("url"))
		if body == nil || len(body) == 0 {
			fmt.Println("订阅内容为空。")
		}
		ioutil.WriteFile(viper.GetString("path"), body, 0777)
	},
}
