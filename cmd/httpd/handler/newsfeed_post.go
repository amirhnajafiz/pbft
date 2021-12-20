package handler

import (
	"cmd/platform/newsfeed"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewsFeedPost(feed *newsfeed.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"hello": "Found me",
		})
	}
}
