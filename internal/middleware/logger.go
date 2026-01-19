package middleware

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

var (
	logWriter io.Writer = os.Stdout
	logMu     sync.Mutex
)

// SetLoggerOutput 设置日志输出目标（例如文件）
func SetLoggerOutput(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	logWriter = w
}

// Logger 是 GoFrame 中间件，记录请求日志
func Logger(r *ghttp.Request) {
	start := time.Now()
	method := r.Method
	path := r.URL.Path

	r.Middleware.Next()

	duration := time.Since(start)
	status := r.Response.Status

	// 格式化日志内容
	msg := fmt.Sprintf("[GOFRAME LOG] %s %s -> %d (%v)\n", method, path, status, duration)
	fmt.Printf("[GOFRAME LOG] %s %s -> %d (%v)\n", method, path, status, duration)
	logMu.Lock()
	defer logMu.Unlock()
	if _, err := logWriter.Write([]byte(msg)); err != nil {
		log.Printf("Failed to write access log: %v", err)
	}
}
