package router

import (
	"little-vote/pkg/api"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	path := r.Group("/api")
	path.GET("/graphql", api.GraphqlHandler())
	path.POST("/graphql", api.GraphqlHandler())
}
