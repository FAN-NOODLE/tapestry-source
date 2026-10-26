package handle

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type BaseHandle struct {
}

func (b BaseHandle) ResSuccess(c *gin.Context, v interface{}) {
	ginplus.ResSuccess(c, v)
	return
}

func (b BaseHandle) ResError(c *gin.Context, err error) {
	ginplus.ResError(c, err)
	return
}

func (b BaseHandle) ResOk(c *gin.Context) {
	ginplus.ResOK(c)
	return
}

func (b BaseHandle) ResJson(c *gin.Context, status int, v interface{}) {
	ginplus.ResJSON(c, status, v)
	return
}

func (b BaseHandle) EnterpriseID(c *gin.Context) int64 {
	return c.Value("enterpriseId").(int64)
}

func (b BaseHandle) MemberInfo(c *gin.Context) *microScrm.Member {
	return c.Value("memberInfo").(*microScrm.Member)
}
