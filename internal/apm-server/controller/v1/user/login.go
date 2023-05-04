package user

import (
	"APM-server/internal/pkg/core"
	"APM-server/internal/pkg/errno"
	"APM-server/internal/pkg/log"
	v1 "APM-server/pkg/api/apm-server/v1"
	"github.com/gin-gonic/gin"
)

// Login 登录 miniblog 并返回一个 JWT Token.
func (ctrl *UserController) Login(c *gin.Context) {
	var r v1.LoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		log.C(c).Infow(errno.ErrBind.Message)

		return
	}

	resp, err := ctrl.b.Users().Login(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

		log.C(c).Infow(err.Error())

		return
	}

	log.C(c).Infow("Login function called", "Login-User:", r.Email)
	core.WriteResponse(c, nil, resp)
}
