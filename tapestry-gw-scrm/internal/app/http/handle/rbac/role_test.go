package rbac

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gw-scrm/internal/pkg/common"
)

func TestConfigurablePermissionsRequestUsesAuthenticatedContext(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)
	c.Set(string(common.ContextKeyEnum.TokenInfo), &common.TokenInfo{
		EnterpriseId: 3363603667487424,
		AdminId:      12345,
	})

	request := configurablePermissionsRequest(c, "store_operation_admin")

	if request.EnterpriseID != 3363603667487424 {
		t.Fatalf("unexpected enterprise id: %d", request.EnterpriseID)
	}
	if request.AdminID != 12345 {
		t.Fatalf("unexpected admin id: %d", request.AdminID)
	}
	if request.AppCode != "store_operation_admin" {
		t.Fatalf("unexpected app code: %s", request.AppCode)
	}
}

func TestConfigurablePermissionsRequestDefaultsToWebAdmin(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)
	request := configurablePermissionsRequest(c, "")
	if request.AppCode != "web_admin" {
		t.Fatalf("unexpected default app code: %s", request.AppCode)
	}
}
