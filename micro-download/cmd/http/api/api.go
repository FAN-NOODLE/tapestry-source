package api

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	middleware "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus/mw"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-download/internal/config"
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
	router.Use(middleware.Logger())
	if config.CfgData.Restful.Cors.Enable {
		// TODO 设置跨域
		router.Use(middleware.Cors())
	}

	c := ctx.New()
	api.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c *gin.Context))
		path = docs.ColonPath(path)

		router.Handle(endpoint.Method, path, h)
		log.Info(c).Msgf("add route [%s]%s", endpoint.Method, path)
	})

	router.GET(api.DocEndpointPath(), gin.WrapH(api.Handler(true)))
	log.Info(c).Msgf("add route [%s]%s", "GET", api.DocEndpointPath())

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "please access spec url.")
	})

	// 心跳检测地址
	router.GET(api.BasePath+"/heart-beat", func(context *gin.Context) {
		context.String(http.StatusOK, "true")
	})

	log.Info(c).Msgf("add route [%s]%s", "GET", "/")

	return router
}
