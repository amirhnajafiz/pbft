package handler

import (
	"net/http"

	"github.com/amirhnajafiz/news-feeder/internal/provider"
	"github.com/gin-gonic/gin"
)

type newsfeedPostRequest struct {
	Title string `json:"title"`
	Post  string `json:"post"`
}

// NewsFeedPost
// input is an interface.
func NewsFeedPost(feed provider.Adder) gin.HandlerFunc {
	return func(c *gin.Context) {
		// creating an empty instance
		var requestBody newsfeedPostRequest

		// binding the json request
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
		}

		item := provider.Item{
			Title: requestBody.Title,
			Post:  requestBody.Post,
		}

		feed.Add(&item)

		c.Status(http.StatusNoContent)
	}
}
