package controller

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/rshulabs/HgCache/internal/pkg/middlewares"
)

func installHttpRouter(g *gin.Engine) *gin.Engine {
	api := g.Group("/cache").Use(middleware.Cors())
	{
		api.GET("/api", Query)
	}
	return g
}
