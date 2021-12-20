package main

import (
	"cmd/cmd/httpd/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ping", handler.PingGet())

	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
