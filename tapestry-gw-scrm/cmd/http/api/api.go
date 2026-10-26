package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	middleware "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus/mw"
	v1 "gw-scrm/cmd/http/api/v1"
	v2 "gw-scrm/internal/app/router/v2"
	"gw-scrm/internal/config"
	interMiddleware "gw-scrm/internal/pkg/middleware"
	"gw-scrm/internal/pkg/tracing"
	"net/http"
)

// Api api 描述
func Api(mode string) http.Handler {
	api := docs.New(
		docs.Version(config.Version),
		docs.Host(config.CfgData.Restful.Host),
		docs.BasePath(config.CfgData.Restful.BasePath),
		docs.Endpoints(endpoints()...),
		docs.Title(config.CfgData.Restful.Doc.Title),
		docs.Description(config.CfgData.Restful.Doc.Description),
		docs.ContactEmail(config.CfgData.Restful.Doc.Contact),
		docs.TermsOfService(config.CfgData.Restful.Doc.TermOfService),
		docs.Tag("trace", "跟踪服务"),
	)

	return handler(api, mode)
}

// handler swagger.API 生成 http.Handler 服务
func handler(api *swagger.API, mode string) http.Handler {

	gin.SetMode(mode)
	router := gin.New()
	router.Use(middleware.Recovery())
	router.Use(middleware.SkyWalkingTracer(router, tracing.SkyTracer()))
	router.Use(interMiddleware.Access())
	router.Use(middleware.Logger())
	router.Use(interMiddleware.CostTime())
	if config.CfgData.Restful.Cors.Enable {
		// 设置跨域
		router.Use(middleware.Cors(func(c *cors.Config) {
			c.AllowAllOrigins = config.CfgData.Restful.Cors.AllowAllOrigins
			c.AllowOrigins = config.CfgData.Restful.Cors.AllowOrigins
			c.AllowMethods = config.CfgData.Restful.Cors.AllowMethods
			c.AllowHeaders = config.CfgData.Restful.Cors.AllowHeaders
			c.AllowCredentials = config.CfgData.Restful.Cors.AllowCredentials
			c.ExposeHeaders = config.CfgData.Restful.Cors.ExposeHeaders
			c.MaxAge = config.CfgData.Restful.Cors.MaxAge
			c.AllowWebSockets = config.CfgData.Restful.Cors.AllowWebSockets
			c.AllowFiles = config.CfgData.Restful.Cors.AllowFiles
			c.AllowWildcard = true
		}))
	}
	router.Use(interMiddleware.AuthCheck(interMiddleware.PathNotCheckAuth()))

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, config.CfgData.App)
	})
	// 心跳检测地址
	router.GET(api.BasePath+"/heart-beat", func(context *gin.Context) {
		context.String(http.StatusOK, "true")
	})

	// 原框架swagger路由
	v1.InitV1Router(api, router)

	// 使用原生路由 放弃swagger
	microCrm := router.Group(api.BasePath)

	v2.InitV2Router(microCrm)

	return router

}
