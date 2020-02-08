package config

import (
	"fmt"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func defaultInit() {

	defaultconfig := make(map[string]interface{})

	if indexAsset, e1 := Asset("config.yaml"); e1 == nil {
		if e2 := yaml.Unmarshal(indexAsset, &defaultconfig); e2 != nil {
			panic(e2)
		}
	} else {
		panic(e1)
	}

	for k, v := range defaultconfig {
		viper.SetDefault(k, v)
	}

	viper.BindEnv("envport", "PORT")
	if p := viper.GetInt("envport"); p != 0 {
		viper.SetDefault("port", p)
	}

}

// Execute func()
func Execute() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config.yaml")

	defaultInit()

	if e := viper.ReadInConfig(); e != nil {
		fmt.Println(e, "\n未找到配置文件，将使用默认配置")

	} else {
		fmt.Println("使用配置文件:", viper.ConfigFileUsed())
	}
}
