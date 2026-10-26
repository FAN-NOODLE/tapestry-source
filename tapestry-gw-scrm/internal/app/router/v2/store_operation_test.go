package v2

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterStoreOperationRoutesIncludesOrganizationReads(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	registerStoreOperationRoutes(router.Group("/store-operation"))
	routes := map[string]bool{}
	for _, route := range router.Routes() {
		routes[route.Method+" "+route.Path] = true
	}
	for _, expected := range []string{
		"POST /store-operation/organization/tree",
		"POST /store-operation/organization/store/list",
	} {
		if !routes[expected] {
			t.Fatalf("route %s is not registered", expected)
		}
	}
}
