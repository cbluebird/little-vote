package router

import (
	"little-vote/pkg/api"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	path := r.Group("/api")
	path.GET("/query", api.Query)
	path.GET("/cas", api.Cas)
	path.POST("/vote", api.Vote)
}
