package shop_member_http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/shop_member"
)

var validate *validator.Validate
var shopMemberService *shop_member.Service

//  等级列表
func levelList(c *gin.Context) {
	var request shop_member.MemberLevelListRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	list, err := shopMemberService.LevelList(ctx.Wrap(c), request.EnterpriseID, []int64{})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var response = shop_member.ConfigGuideListResponse{
		List: list,
	}
	ginplus.ResSuccess(c, response)
	return
}
