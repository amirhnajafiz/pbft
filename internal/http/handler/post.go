package handler

import (
	"fmt"
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
		var requestBody newsfeedPostRequest

		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
		}

		fmt.Println(requestBody)

		item := provider.Item{
			Title: requestBody.Title,
			Post:  requestBody.Post,
		}

		feed.Add(&item)

		c.Status(http.StatusNoContent)
	}
}
