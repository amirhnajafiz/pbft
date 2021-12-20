package main

import (
	"cmd/cmd/httpd/handler"
	"cmd/platform/newsfeed"
	"github.com/gin-gonic/gin"
)

func main() {
	feed := newsfeed.New()
	r := gin.Default()

	r.GET("/ping", handler.PingGet())
	r.GET("/newsfeed", handler.NewsFeedGet(feed))
	r.POST("/newsfeed", handler.NewsFeedPost(feed))

	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
