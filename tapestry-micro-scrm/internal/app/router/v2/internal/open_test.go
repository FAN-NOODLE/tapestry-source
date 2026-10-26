package internal

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOpenRegistersCommonDecryptRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	Open(router.Group("/open"))

	for _, route := range router.Routes() {
		if route.Method == "POST" && route.Path == "/open/common/decrypt" {
			return
		}
	}
	t.Fatal("未注册 POST /open/common/decrypt")
}

func TestOpenRegistersStoreMasterRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	Open(router.Group("/open"))

	for _, route := range router.Routes() {
		if route.Method == "POST" && route.Path == "/open/store_worker/store_master_list" {
			return
		}
	}
	t.Fatal("未注册 POST /open/store_worker/store_master_list")
}

func TestOpenRegistersEmployeeMasterRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	Open(router.Group("/open"))

	for _, route := range router.Routes() {
		if route.Method == "POST" && route.Path == "/open/store_worker/employee_master_list" {
			return
		}
	}
	t.Fatal("未注册 POST /open/store_worker/employee_master_list")
}
