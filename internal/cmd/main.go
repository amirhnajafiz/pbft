package main

import (
	"github.com/amirhnajafiz/news-feeder/internal/http/handler"
	"github.com/amirhnajafiz/news-feeder/internal/provider"
	"github.com/gin-gonic/gin"
)

func main() {
	// creating a new provider
	feed := provider.New()
	// get gin package default listener
	r := gin.Default()

	// creating our endpoints
	r.GET("/newsfeed", handler.NewsFeedGet(feed))
	r.POST("/newsfeed", handler.NewsFeedPost(feed))

	// starting gin service
	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
