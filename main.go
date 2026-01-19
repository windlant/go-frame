package main

import (
	"log"
	"os"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/windlant/go-frame/internal/controller"
	"github.com/windlant/go-frame/internal/middleware"
)

func init() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	logFile, err := os.OpenFile("logs/goframe-access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	middleware.SetLoggerOutput(logFile)
}

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
