package controller

import (
	"github.com/gin-gonic/gin"
)

func installHttpRouter(g *gin.Engine) *gin.Engine {
	api := g.Group("/cache")
	{
		api.GET("/api", Query)
	}
	return g
}
