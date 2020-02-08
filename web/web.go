package web

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/heiha/ssr2clashr/api"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/spf13/viper"
)

var excludeExtensions = [...]string{
	".js",
	".css",
	".jpg",
	".png",
	".ico",
	".svg",
}

var accessLog *os.File = nil

func onError(ctx iris.Context) {
	switch ctx.GetStatusCode() {
	case iris.StatusNotFound, iris.StatusOK:
		errorAsset, err := Asset("public/errors/404.html")
		if err != nil {
			fmt.Println(err)
		}
		_, err = ctx.Write(errorAsset)
		if err != nil {
			fmt.Println(err)
		}
	default:
		errorAsset, err := Asset("public/errors/errors.html")
		if err != nil {
			fmt.Println(err)
		}
		_, err = ctx.Write(errorAsset)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func debugRequestLogger() iris.Handler {
	c := logger.Config{
		Status:  true,
		IP:      true,
		Method:  true,
		Path:    true,
		Columns: true,
	}

	if accessLog == nil {
		var err error
		accessLog, err = os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	}

	c.LogFunc = func(now time.Time, latency time.Duration, status, ip, method, path string, message interface{}, headerMessage interface{}) {
		output := logger.Columnize(now.Format("2006/01/02 - 15:04:05"), latency, status, ip, method, path, message, headerMessage)
		fmt.Println(output)
		_, err := accessLog.Write([]byte(output))
		if err != nil {
			fmt.Println(err)
		}

	}

	c.AddSkipper(func(ctx iris.Context) bool {
		path := ctx.Path()
		for _, ext := range excludeExtensions {
			if strings.HasSuffix(path, ext) {
				return true
			}
		}
		return false
	})
	return logger.New(c)
}

func newApp() *iris.Application {

	app := iris.New()

	app.Get("/sub", func(ctx iris.Context) {

		switch key := ctx.URLParam("key"); key {
		case viper.GetString("key"):
			_, err := ctx.Write(api.Execute(ctx.URLParam("url")))
			if err != nil {
				fmt.Println(err)
			}

		default:
			_, err := ctx.WriteString("未通过验证！")
			if err != nil {
				fmt.Println(err)
			}
		}

	})

	app.OnAnyErrorCode(onError)

	if viper.GetBool("debug") {
		app.UseGlobal(debugRequestLogger())
		fmt.Println("Debug")
	} else {
		app.Get("/test", onError)
		fmt.Println("欢迎使用！")
	}

	app.HandleDir("/", "./public", iris.DirOptions{
		Asset:      Asset,
		AssetInfo:  AssetInfo,
		AssetNames: AssetNames,
		Gzip:       true,
	})

	return app
}

// Execute func()
func Execute() {
	err := newApp().Run(iris.Addr(":"+viper.GetString("port")), iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		panic(err)
	}
	defer accessLog.Close()
}
