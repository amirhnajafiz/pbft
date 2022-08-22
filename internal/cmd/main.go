package main

import (
	"github.com/amirhnajafiz/news-feeder/internal/http/handler"
	"github.com/amirhnajafiz/news-feeder/internal/provider"
	"github.com/gin-gonic/gin"
)

func main() {
	feed := provider.New()
	r := gin.Default()

	r.GET("/newsfeed", handler.NewsFeedGet(feed))
	r.POST("/newsfeed", handler.NewsFeedPost(feed))

	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
