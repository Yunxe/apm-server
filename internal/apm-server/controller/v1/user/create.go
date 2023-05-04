package user

import (
	"APM-server/pkg/api/apm-server/v1"
	"github.com/asaskevich/govalidator"

	"APM-server/internal/pkg/core"
	"APM-server/internal/pkg/errno"
	"APM-server/internal/pkg/log"
	"github.com/gin-gonic/gin"
)

const defaultMethods = "(GET)|(POST)|(PUT)|(DELETE)"

func (ctrl *UserController) Create(c *gin.Context) {
	var r v1.CreateUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		log.C(c).Infow(errno.ErrBind.Message)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		errt := errno.ErrInvalidParameter.SetMessage(err.Error())
		core.WriteResponse(c, errt, nil)

		log.C(c).Infow(errt.Message)

		return
	}

	if err := ctrl.b.Users().Create(c, &r); err != nil {
		core.WriteResponse(c, err, nil)

		log.C(c).Infow(err.Error())

		return
	}

	if _, err := ctrl.a.AddNamedPolicy("p", r.Email, "/v1/users/"+r.Email, defaultMethods); err != nil {
		core.WriteResponse(c, err, nil)

		log.C(c).Infow(err.Error())

		return
	}

	log.C(c).Infow("Create user function called", "New-User:", r.Username)
	core.WriteResponse(c, nil, nil)
}
