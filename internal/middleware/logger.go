package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

func Logger(r *ghttp.Request) {
	glog.Infof(
		r.Context(),
		"%s %s %d %s",
		r.Method,
		r.URL.Path,
		r.Response.Status,
		r.GetClientIp(),
	)
	r.Middleware.Next()
}
