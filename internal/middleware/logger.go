package middleware

import (
	"fmt"
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

	fmt.Printf("[GOFRAME LOG] %s %s -> %d (%v)\n", method, path, status, duration)
}
