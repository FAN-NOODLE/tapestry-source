package v1

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
)

func InitRouter(api *swagger.API, router *gin.Engine) {

	router.GET(api.DocEndpointPath(), gin.WrapH(api.Handler(true)))

	api.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c *gin.Context))
		path = docs.ColonPath(path)
		router.Handle(endpoint.Method, path, h)
	})
}
