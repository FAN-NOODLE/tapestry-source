package source

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/handle"
	"gw-work-scrm/internal/app/http/request/source"
	"gw-work-scrm/internal/pkg/errs"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type GroupRelationHandle struct {
	handle.BaseHandle
	groupRelationSdk microScrm.SourceGroupRelationApi
}

func NewGroupRelationHandle() GroupRelationHandle {
	return GroupRelationHandle{
		groupRelationSdk: microScrm.NewSourceGroupRelationApi(),
	}
}

// List 素材分组归属门店列表
func (s GroupRelationHandle) List(c *gin.Context) {
	var request source.GroupRelationListRequest

	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	request.EnterpriseID = s.EnterpriseID(c)

	res, err := s.groupRelationSdk.List(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// Save 素材分组归属门店保存
//func (s GroupRelationHandle) Save(c *gin.Context) {
//	var request source.GroupRelationSaveRequest
//	if err := c.ShouldBind(&request); err != nil {
//		s.ResError(c, errs.ParamsErr)
//		return
//	}
//	request.EnterpriseID = s.EnterpriseID(c)
//
//	res, err := s.groupRelationSdk.Save(ctx.Wrap(c), request)
//	if err != nil {
//		s.ResError(c, err)
//		return
//	}
//	s.ResSuccess(c, res)
//	return
//}

// GetListForStore List 根据门店获取素材分组归属门店列表
func (s GroupRelationHandle) GetListForStore(c *gin.Context) {
	var request source.GroupRelationFindRequest

	err := c.BindJSON(&request)
	if err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	request.EnterpriseID = s.EnterpriseID(c)
	memberInfo := c.Value("memberInfo").(*microScrm.Member)
	request.MemberUserId = memberInfo.Userid
	request.MemberID = memberInfo.ID

	res, err := s.groupRelationSdk.GetListForStore(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// SelectList 小程序端门店员工根据查看分组权限查询各素材列表
func (s GroupRelationHandle) SelectList(c *gin.Context) {
	var request source.RequestSelectList

	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	request.EnterpriseID = s.EnterpriseID(c)
	memberInfo := c.Value("memberInfo").(*microScrm.Member)
	request.MemberUserId = memberInfo.Userid
	request.MemberID = memberInfo.ID
	res, err := s.groupRelationSdk.SelectList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// GetNewSourceCount 小程序端素材上新数量
func (s GroupRelationHandle) GetNewSourceCount(c *gin.Context) {
	var request source.RequestNewSource

	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	request.EnterpriseID = s.EnterpriseID(c)
	memberInfo := c.Value("memberInfo").(*microScrm.Member)
	request.MemberUserId = memberInfo.Userid
	request.MemberID = memberInfo.ID
	request.Mobile = memberInfo.Mobile
	res, err := s.groupRelationSdk.GetNewSourceCount(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s GroupRelationHandle) CollectSelectList(c *gin.Context) {
	var request source.RequestSelectList

	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, errs.ParamsErr)
		return
	}

	request.EnterpriseID = s.EnterpriseID(c)
	memberInfo := c.Value("memberInfo").(*microScrm.Member)
	request.MemberUserId = memberInfo.Userid
	request.MemberID = memberInfo.ID
	request.Mobile = memberInfo.Mobile
	request.CollectSearch = true
	res, err := s.groupRelationSdk.CollectSelectList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}
