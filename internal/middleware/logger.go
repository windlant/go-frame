package middleware

import (
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Logger(r *ghttp.Request) {
	start := time.Now()
	method := r.Method
	path := r.URL.Path

	r.Middleware.Next()

	duration := time.Since(start)
	status := r.Response.Status

	println("[GOFRAE LOG]", method, path, "->", status, "(", duration, ")")
}
