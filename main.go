package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/windlant/go-frame/internal/middleware"
	"github.com/windlant/go-frame/internal/router"
)

func main() {
	ctx := gctx.New()

	// 初始化数据库
	if _, err := g.DB().GetAll(ctx, "SELECT 1"); err != nil {
		g.Log().Fatal(ctx, "MySQL connect failed:", err)
	}

	// 初始化 Redis
	if _, err := g.Redis().Do(ctx, "PING"); err != nil {
		g.Log().Fatal(ctx, "Redis connect failed:", err)
	}

	s := g.Server()
	s.SetPort(8080)

	// 注册中间件
	s.Use(middleware.Logger)

	// 注册路由
	router.Register(s)

	// 启动服务
	s.Run()
}
