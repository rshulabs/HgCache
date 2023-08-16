package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rshulabs/HgCache/pkg/errorx"
)

type Response struct {
	Code      int         `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Reference string      `json:"reference,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func WriteResponse(c *gin.Context, err error, data any) {
	if err != nil {
		coder := errorx.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), Response{
			Code:      coder.Code(),
			Message:   coder.String(),
			Reference: coder.Reference(),
			Data:      data,
		})
		// fmt.Println(coder.Code())
	}
	c.JSON(http.StatusOK, Response{Data: data})
}
