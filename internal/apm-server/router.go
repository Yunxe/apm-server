package apm_server

import (
	"APM-server/internal/apm-server/controller/v1/user"
	"APM-server/internal/apm-server/controller/v2/apm"
	"APM-server/internal/apm-server/store"
	"APM-server/internal/pkg/core"
	"APM-server/internal/pkg/errno"
	"APM-server/internal/pkg/log"
	"APM-server/internal/pkg/middleware"
	"APM-server/pkg/auth"

	"github.com/gin-gonic/gin"
)

// installRouters 安装 miniblog 接口路由.
func installRouters(g *gin.Engine) error {
	// 注册 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	// 注册 /healthz handler.
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")

		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	authz, err := auth.NewAuthz(store.S.DB())
	if err != nil {
		return err
	}

	uc := user.New(store.S, authz)

	g.POST("/login", uc.Login)

	//创建 v1 路由分组
	v1 := g.Group("/v1")
	v2 := g.Group("/v2")
	{
		// 创建 users 路由分组
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create) // 创建用户
			userv1.Use(middleware.Authn(), middleware.Authz(authz))
			userv1.GET(":name", uc.Get) // 获取用户详情
			//userv1.PUT(":name", uc.Update)    // 更新用户
			//userv1.GET("", uc.List)           // 列出用户列表，只有 root 用户才能访问
			//userv1.DELETE(":name", uc.Delete) // 删除用户
		}
		apmv2 := v2.Group("/apm")
		{
			apmv2.GET(":addr", apm.WebsocketHandler)
		}
	}

	return nil
}
