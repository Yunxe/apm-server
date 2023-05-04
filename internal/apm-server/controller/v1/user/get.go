package user

import (
	"APM-server/internal/pkg/core"
	"APM-server/internal/pkg/log"
	"github.com/gin-gonic/gin"
)

// Get 获取一个用户的详细信息.
func (ctrl *UserController) Get(c *gin.Context) {

	user, err := ctrl.b.Users().Get(c, c.Param("name"))
	if err != nil {
		core.WriteResponse(c, err, nil)

		log.C(c).Infow(err.Error())

		return
	}

	log.C(c).Infow("Get user function called")
	core.WriteResponse(c, nil, user)
}
