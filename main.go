package main

import (
	"github.com/Unknwon/macaron"
	"github.com/weisd/goapi/modules/middleware"
	"github.com/weisd/goapi/routers"
	"github.com/weisd/goapi/routers/account"
	"github.com/weisd/goapi/routers/sso"
	"github.com/weisd/goapi/routers/token"
	// "net/http"
	"log"
)

const APP_VER = "0.0.1.1106"

func main() {
	routers.GlobalInit()

	m := macaron.Classic()
	// 要加renderer才能用 - -
	m.Use(macaron.Renderer())
	m.Use(middleware.Contexter())

	// header验证
	// m.Use(func(ctx *macaron.Context) {
	// 	if ctx.Req.Header.Get("X-API-KEY") != "ktkt" {
	// 		ctx.Resp.WriteHeader(http.StatusUnauthorized)
	// 	}
	// })

	//首页，作测试连接用，返回当前版本号等信息
	m.Get("/", func(ctx *macaron.Context, logger *log.Logger) string {
		logger.Println("the path is :", ctx.Req.RequestURI)
		return "wellcome to goapi v1 !"
	})

	m.Group("/token", func() {
		m.Get("/create", token.Create)
		m.Get("/update", token.Update)
		m.Get("/auth", token.Auth)
	})

	m.Group("/account", func() {
		m.Post("/create", account.Create)
		m.Get("/create/test", account.Create)
		m.Get("/update", account.Update)
		m.Get("/auth/test", account.Auth)
		m.Post("/auth", account.Auth)
		m.Get("/delete", account.Delete)
		m.Get("/info/test", account.Info)
		m.Post("/info", account.Info)
		m.Get("/list", account.Accounts)
	})

	m.Group("/sso", func() {
		m.Post("/auth", sso.Auth)
		m.Get("/auth/test", sso.Auth)
		m.Get("/list", sso.List)
		m.Get("/client", sso.Client)
	})

	// m.NotFound(func(ctx *macaron.Context, logger *log.Logger) {
	// 	logger.Println("404 page")
	// })

	m.Run()
}
