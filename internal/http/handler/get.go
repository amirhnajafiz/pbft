package handler

import (
	"net/http"

	"github.com/amirhnajafiz/news-feeder/internal/provider"
	"github.com/gin-gonic/gin"
)

// NewsFeedGet
// input is an interface.
func NewsFeedGet(feed provider.Getter) gin.HandlerFunc {
	return func(c *gin.Context) {
		results := feed.GetAll()

		c.JSON(http.StatusOK, results)
	}
}
