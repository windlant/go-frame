package main

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/windlant/go-frame/internal/controller"
	"github.com/windlant/go-frame/internal/middleware"
)

func main() {
	s := ghttp.GetServer()
	s.SetPort(8080)
	s.Use(middleware.Logger)

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(new(controller.UserController))
	})

	s.Run()
}
