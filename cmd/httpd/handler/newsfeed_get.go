package handler

import (
	"cmd/platform/newsfeed"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewsFeedGet(feed newsfeed.Getter) gin.HandlerFunc {
	return func(c *gin.Context) {
		results := feed.GetAll()
		c.JSON(http.StatusOK, results)
	}
}
