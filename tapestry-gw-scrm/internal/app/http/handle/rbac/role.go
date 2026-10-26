package rbac

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microauthority"
)

type RoleHandle struct {
	handle.BaseHandle
}

type configurablePermissionsReq struct {
	EnterpriseID int64  `json:"enterprise_id"`
	AdminID      int64  `json:"admin_id"`
	AppCode      string `json:"app_code"`
}

func NewRoleHandle() RoleHandle {
	return RoleHandle{}
}

func configurablePermissionsRequest(c *gin.Context, appCode string) configurablePermissionsReq {
	enterpriseID, adminID := common.GetEnterprise(c)
	if appCode == "" {
		appCode = "web_admin"
	}
	return configurablePermissionsReq{
		EnterpriseID: enterpriseID,
		AdminID:      adminID,
		AppCode:      appCode,
	}
}

func (h RoleHandle) ConfigurablePermissions(c *gin.Context) {
	var input struct {
		AppCode string `json:"app_code"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := microauthority.ConfigurablePermissions(
		h.WarpContext(c),
		configurablePermissionsRequest(c, input.AppCode),
	)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)
}
