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

	userCtrl := new(controller.UserController)

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/users", userCtrl.GetUsers)
		group.GET("/users/:id", userCtrl.GetUser)
		group.POST("/users", userCtrl.CreateUsers)
		group.PUT("/users", userCtrl.UpdateUsers)
		group.DELETE("/users", userCtrl.DeleteUsers)
	})

	s.Run()
}
