package controller

import (
	"github.com/gin-gonic/gin"

	v2 "github.com/rshulabs/HgCache/internal/cache/v2"
	"github.com/rshulabs/HgCache/internal/pkg/code"
	_ "github.com/rshulabs/HgCache/internal/pkg/code"
	"github.com/rshulabs/HgCache/internal/pkg/response"
	"github.com/rshulabs/HgCache/pkg/errorx"
)

func Query(c *gin.Context) {
	key := c.Query("key")
	g := v2.GetGroup("test")
	view, err := g.Get(key)
	// fmt.Println("query error: ", err.Error())
	if err != nil {
		response.WriteResponse(c, errorx.WithCode(code.ErrNotFound, err.Error()), nil)
		return
	}
	response.WriteResponse(c, nil, view.String())
}
