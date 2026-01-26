package router

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/windlant/go-frame/internal/controller"
)

func Register(s *ghttp.Server) {
	userCtrl := controller.NewUserController()

	s.Group("/users", func(group *ghttp.RouterGroup) {
		group.POST("/list", userCtrl.ListUsers)
		group.POST("/create", userCtrl.CreateUsers)
		group.POST("/get", userCtrl.GetUsers)
		group.POST("/update", userCtrl.UpdateUsers)
		group.POST("/delete", userCtrl.DeleteUsers)
	})
}
