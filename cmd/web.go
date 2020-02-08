package cmd

import (
	"github.com/heiha/ssr2clashr/api"
	"github.com/heiha/ssr2clashr/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(webCmd)

	webCmd.PersistentFlags().StringP(
		"port",
		"p",
		viper.GetString("port"),
		"请填写服务器端口！",
	)
	webCmd.PersistentFlags().StringP(
		"key",
		"k",
		viper.GetString("password"),
		"请填写服务器密码！",
	)
	webCmd.PersistentFlags().StringP(
		"url",
		"u",
		viper.GetString("url"),
		"请填写默认订阅链接！",
	)
	webCmd.PersistentFlags().StringP(
		"template",
		"t",
		viper.GetString("template"),
		"套用的模板地址",
	)

	viper.BindPFlag("port", webCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("key", webCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("url", webCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("template", webCmd.PersistentFlags().Lookup("template"))

}

var webCmd = &cobra.Command{

	Use:   "web",
	Short: "ssr2clashr web",
	Long:  `启动网页服务`,
	Run: func(cmd *cobra.Command, args []string) {
		api.InitRules()
		web.Execute()
	},
}
