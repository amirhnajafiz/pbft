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

func NewsFeedPost(feed provider.Adder) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := newsfeedPostRequest{}

		_ = c.Bind(&requestBody)

		item := provider.Item{
			Title: requestBody.Title,
			Post:  requestBody.Post,
		}

		feed.Add(item)

		c.Status(http.StatusNoContent)
	}
}
