package middleware

import (
	"APM-server/internal/pkg/core"
	"APM-server/internal/pkg/errno"
	"APM-server/internal/pkg/known"
	"APM-server/pkg/token"
	"github.com/gin-gonic/gin"
)

// Authn 是认证中间件，用来从 gin.Context 中提取 token 并验证 token 是否合法，
// 如果合法则将 token 中的 sub 作为<邮箱>存放在 gin.Context 的 XUsernameKey 键中.
func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析 JWT Token
		email, err := token.ParseRequest(c)
		if err != nil {
			core.WriteResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()

			return
		}

		c.Set(known.XEmailKey, email)
		c.Next()
	}
}
