package main

import (
	"github.com/Unknwon/macaron"
	"github.com/weisd/goapi/modules/middleware"
	"github.com/weisd/goapi/routers"
	"github.com/weisd/goapi/routers/account"
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

	m.Group("/token", func(r *macaron.Router) {
		r.Get("/create", token.Create)
		r.Get("/auth", token.Auth)
		r.Get("/delete", token.Delete)
	})

	m.Group("/account", func(r *macaron.Router) {
		r.Get("/create", account.Create)
		r.Get("/update", account.Update)
		r.Get("/auth", account.Auth)
		r.Get("/delete", account.Delete)
		r.Get("/info", account.Info)
		r.Get("/accounts", account.Accounts)
	})

	// m.NotFound(func(ctx *macaron.Context, logger *log.Logger) {
	// 	logger.Println("404 page")
	// })

	m.Run()
}
