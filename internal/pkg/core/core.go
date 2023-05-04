package core

import (
	"APM-server/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

type ErrResponse struct {
	// Code 指定了业务错误码.
	Code int `json:"code"`

	// Message 包含了可以直接对外展示的错误信息.
	Message string `json:"message"`

	Data interface{} `json:"data,omitempty"`
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	// if err != nil {
	// 	hcode, code, message := errno.Decode(err)
	// 	c.JSON(hcode, ErrResponse{
	// 		Code:    code,
	// 		Message: message,
	// 	})

	// 	return
	// }
	hcode, code, message := errno.Decode(err)
	c.JSON(hcode, ErrResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
