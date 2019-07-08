package main

import (
	"go-hotel/configs"
	"go-hotel/eths"
	"go-hotel/routes"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var e *echo.Echo //echo框架对象全局定义

func main() {

	e = echo.New()             // 创建echo对象
	e.Use(middleware.Logger()) // 安装日志中间件
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	go eths.EventSubscribe("ws://localhost:8546", configs.Config.Eth.HotelAddr)

	e.GET("/ping", routes.PingHandler)
	e.GET("/mint", routes.UpLoad)                       // UpLoad
	e.GET("/check", routes.Check)                       // GetPic, /content/titleValues
	e.Logger.Fatal(e.Start(configs.Config.Common.Port)) // 启动服务
}
