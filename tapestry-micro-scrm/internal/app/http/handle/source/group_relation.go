package source

import (
	"errors"
	"fmt"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/request/source"
	sourceResp "micro-scrm/internal/app/http/response/source"
	sourceModel "micro-scrm/internal/app/models/source"
	cacheSvc "micro-scrm/internal/app/service/cache"
	departmentService "micro-scrm/internal/app/service/department"
	membersSvc "micro-scrm/internal/app/service/member"
	sourceService "micro-scrm/internal/app/service/source"
	"micro-scrm/internal/config"
	sourceModules "micro-scrm/internal/modules/source"
	"micro-scrm/internal/orm"
	bossModel "micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/assist"
	"micro-scrm/internal/pkg/helper"

	"github.com/gin-gonic/gin"
)

type GroupRelationHandle struct {
	handle.BaseHandle
	GroupRelationService     sourceService.GroupRelationService
	MemberDepartmentsService departmentService.MemberDepartmentsService
	DepartmentsService       departmentService.DepartmentsService
	SourceGroupService       sourceService.Service
	MembersService           *membersSvc.MembersService
}

func NewGroupRelationHandle() GroupRelationHandle {
	return GroupRelationHandle{
		GroupRelationService:     sourceService.NewGroupRelationService(),
		MemberDepartmentsService: departmentService.NewMemberDepartmentsService(),
		DepartmentsService:       departmentService.NewDepartmentsService(),
		SourceGroupService:       sourceService.NewService(),
		MembersService:           membersSvc.NewMembersService(),
	}
}

// List 素材分组归属门店列表
func (h GroupRelationHandle) List(c *gin.Context) {
	var request source.GroupRelationListRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	if request.SourceGroupId == 0 {
		h.ResError(c, errors.New("素材分组ID不能为空"))
		return
	}
	err, result, total := h.GroupRelationService.GroupRelationGetList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	// h.ResList(c, total, request.Page, request.Size, result)
	h.ResList(c, total, 0, 0, result)
	return
}

// Save 素材分组归属门店
func (h GroupRelationHandle) Save(c *gin.Context) {
	var request source.GroupRelationSaveRequest
	err := c.ShouldBind(&request)
	fmt.Print(request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	if len(request.ThirdIds) == 0 {
		h.ResError(c, errors.New("门店或部门不能为空"))
		return
	}

	err = h.GroupRelationService.GroupRelationSave(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}

// SaveAll 素材分组归属所有门店
func (h GroupRelationHandle) SaveAll(c *gin.Context) {
	var request source.GroupRelationSaveRequest
	err := c.ShouldBind(&request)
	fmt.Print(request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	if len(request.ThirdIds) == 0 {
		h.ResError(c, errors.New("门店或部门不能为空"))
		return
	}

	err = h.GroupRelationService.GroupRelationSaveAll(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}

// GetListForStore 根据门店查看素材分组列表
func (h GroupRelationHandle) GetListForStore(c *gin.Context) {

	var params source.GroupRelationFindRequest
	err := c.ShouldBind(&params)
	if err != nil {
		h.ResError(c, err)
		return
	}

	thirdsInfoList, err := h.MembersService.GetMemberAndDepartmentsByMemberId(h.WarpContext(c), params.EnterpriseID, params.MemberID)

	if err != nil {
		h.ResError(c, err)
		return
	}
	for _, thirdsInfo := range thirdsInfoList {
		params.ThirdIds = append(params.ThirdIds, int(thirdsInfo.DepartmentThirdId))
	}

	// 查看分组
	if err1, result := h.GroupRelationService.SourceListForId(h.WarpContext(c), params); err1 != nil {
		h.ResError(c, err1)
		return
	} else {
		h.ResSuccess(c, result)
	}

}

// SelectList 小程序端门店员工根据查看分组权限查询各素材列表
func (h GroupRelationHandle) SelectList(c *gin.Context) {
	var request source.RequestSelectList
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	go func() {
		// 更新已过期的素材状态
		err = h.SourceGroupService.SelectEndListAndUpdate(h.WarpContext(c), request.EnterpriseID)
		if err != nil {
			h.ResError(c, err)
			return
		}
	}()

	if request.MemberID == 0 {
		// 导购信息
		memberInfo, err1 := h.MembersService.GetInfoByMobile(h.WarpContext(c), request.EnterpriseID, request.Mobile)
		if err1 != nil {
			h.ResError(c, err1)
			return
		}
		request.MemberID = memberInfo.ID
	}

	request.ThirdIds, err = h.MembersService.MembersDepartmentsDao.GetThirdsByMemberIds(h.WarpContext(c), request.EnterpriseID, request.MemberID)
	if err != nil {
		return
	}

	request.GroupIDList = []int64{}
	if request.SourceGroupID != 0 {
		var requestChild source.RequestChildList
		requestChild.ParentID = request.SourceGroupID
		requestChild.EnterpriseID = request.EnterpriseID
		_, request.GroupIDList = h.SourceGroupService.GetChild(h.WarpContext(c), requestChild)
	} else {
		request.IgnoreIndex = true
	}

	var (
		result interface{}
		total  int64
	)

	// 七天前的起止时间戳
	sevenDaysStartTime, _ := helper.YesterdayPeriodOfTimestamp(7)

	switch request.Type {
	// 类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)
	case sourceModel.TypeText:
		result, total, err = h.GroupRelationService.TextSelectList(h.WarpContext(c), request, false)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceModel.ResponseTextList)
		if len(list) == 0 {
			result = []sourceModel.ResponseTextList{}
		}
		for k, v := range list {
			// 七天内上新素材
			if v.StartTime > sevenDaysStartTime {
				list[k].IsNew = 1
			}
		}
	case sourceModel.TypeImg:
		result, total, err = h.GroupRelationService.ImageSelectList(h.WarpContext(c), request, false)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceResp.ImageSelectList)
		if len(list) == 0 {
			result = []sourceResp.ImageSelectList{}
		}
		for k, v := range list {
			// 七天内上新素材
			if v.StartTime > sevenDaysStartTime {
				list[k].IsNew = 1
			}
		}
	case sourceModel.TypeLink:
		result, total, err = h.GroupRelationService.LinkSelectList(h.WarpContext(c), request, false)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceModel.ResponseLinkList)
		if len(list) == 0 {
			result = []sourceModel.ResponseLinkList{}
		}
		if total > 0 {

			corp, err := orm.BossService.Enterprise.GetInfo(bossModel.EnterpriseModel{ID: request.EnterpriseID})
			if err != nil {
				h.ResError(c, err)
				return
			}
			corpHash := corp.Hash
			for k, v := range list {
				// 七天内上新素材
				if v.StartTime > sevenDaysStartTime {
					list[k].IsNew = 1
				}
				if v.LinkType == sourceModules.LinkUrlTypeSelf {
					domain := config.CfgData.Service.WebH5.Domain
					linkUrl := fmt.Sprintf(sourceModules.SelfLinkUrl, domain, corpHash) + assist.StrVal(v.ID)
					list[k].LinkUrl = linkUrl
				}
			}
		}
	case sourceModel.TypeMiniProgram:
		result, total, err = h.GroupRelationService.MiniProgramSelectList(h.WarpContext(c), request, false)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceModel.ResponseMiniProgramList)
		if len(list) == 0 {
			result = []sourceModel.ResponseMiniProgramList{}
		}
		var pathMapArr = make(map[int64]string)
		for k, v := range list {
			// 七天内上新素材
			if v.StartTime > sevenDaysStartTime {
				list[k].IsNew = 1
			}
			pathMapArr[v.ID] = v.Path
		}
		// sw小程序路径添加参数
		resConfig, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(h.WarpContext(c), request.EnterpriseID, "store_name")
		if resConfig.Name == "sw" {
			newPathMap, err := membersSvc.NewMembersService().SetMiniProgramLinkParams(h.WarpContext(c), request.EnterpriseID, request.MemberID, pathMapArr)
			if err != nil {
				return
			}
			for k, v := range list {
				list[k].Path = newPathMap[v.ID]
			}
		}
	case sourceModel.TypeVideo:
		result, total, err = h.GroupRelationService.VideoSelectList(h.WarpContext(c), request, false)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceModel.ResponseVideoList)
		if len(list) == 0 {
			result = []sourceModel.ResponseVideoList{}
		}
		for k, v := range list {
			// 七天内上新素材
			if v.StartTime > sevenDaysStartTime {
				list[k].IsNew = 1
			}
		}

	case sourceModel.TypeFile:
		result, total, err = h.GroupRelationService.FileSelectList(h.WarpContext(c), request, false)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceModel.ResponseFileList)
		if len(list) == 0 {
			result = []sourceModel.ResponseFileList{}
		}
		for k, v := range list {
			// 七天内上新素材
			if v.StartTime > sevenDaysStartTime {
				list[k].IsNew = 1
			}
		}
	case sourceModel.TypeFriend:
		result, total, err = h.GroupRelationService.FriendSelectList(h.WarpContext(c), request, false)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceResp.FriendSelectList)
		if len(list) == 0 {
			result = []sourceResp.FriendSelectList{}
		}
		for k, v := range list {
			// 七天内上新素材
			if v.StartTime > sevenDaysStartTime {
				list[k].IsNew = 1
			}
			if len(v.MergeImage) > 0 {
				list[k].OriginalURL = ""
			}
		}
	default:
		h.ResError(c, err)
		return
	}

	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, request.Page, request.Size, result)
	return

}

// ConsoleList 后台管理查询各素材列表
func (h GroupRelationHandle) ConsoleList(c *gin.Context) {
	var request source.RequestConsoleList
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)

	go func() {
		// 更新已过期的素材状态
		err = h.SourceGroupService.SelectEndListAndUpdate(cc, request.EnterpriseID)
		if err != nil {
			h.ResError(c, err)
			return
		}
	}()

	if request.ProductName != "" || request.SkuItem != "" || request.SkuColor != "" {
		if !helper.InArray(request.Type, []int{sourceModel.TypeLink, sourceModel.TypeImg, sourceModel.TypeVideo, sourceModel.TypeFriend}) {
			h.ResList(c, 0, request.Page, request.Size, nil)
			return
		}
	}
	request.GroupIDList = []int64{}
	if request.SourceGroupID != 0 {
		var requestChild source.RequestChildList
		requestChild.ParentID = request.SourceGroupID
		requestChild.EnterpriseID = request.EnterpriseID
		_, request.GroupIDList = h.SourceGroupService.GetChild(cc, requestChild)
	}

	var (
		result interface{}
		total  int64
	)

	switch request.Type {
	// 类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件 7 朋友圈)
	case sourceModel.TypeText:
		err, result, total = h.SourceGroupService.TextConsoleList(cc, request)
		if err != nil {
			h.ResError(c, err)
			return
		}
	case sourceModel.TypeImg:
		err, result, total = h.SourceGroupService.ImageConsoleList(cc, request)
		if err != nil {
			h.ResError(c, err)
			return
		}
	case sourceModel.TypeLink:
		err, result, total = h.SourceGroupService.LinkConsoleList(cc, request)
		if err != nil {
			h.ResError(c, err)
			return
		}
	case sourceModel.TypeMiniProgram:
		err, result, total = h.SourceGroupService.MiniProgramConsoleList(cc, request)
		if err != nil {
			h.ResError(c, err)
			return
		}
	case sourceModel.TypeVideo:
		err, result, total = h.SourceGroupService.VideoConsoleList(cc, request)
		if err != nil {
			h.ResError(c, err)
			return
		}
	case sourceModel.TypeFile:
		err, result, total = h.SourceGroupService.FileConsoleList(cc, request)
		if err != nil {
			h.ResError(c, err)
			return
		}
	case sourceModel.TypeFriend:
		err, result, total = h.SourceGroupService.FriendConsoleList(cc, request)
		if err != nil {
			h.ResError(c, err)
			return
		}
	default:
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, result)
	return

}

// GetNewSourceNum 小程序端素材上新数量
func (h GroupRelationHandle) GetNewSourceNum(c *gin.Context) {

	var request source.RequestNewSource
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	thirdsInfoList, err := h.MembersService.GetMemberInfoAndAllDepartmentsByMobile(h.WarpContext(c), request.EnterpriseId, request.Mobile)
	if err != nil {
		h.ResError(c, err)
		return
	}
	for _, thirdsInfo := range thirdsInfoList {
		request.ThirdIds = append(request.ThirdIds, int(thirdsInfo.DepartmentThirdId))
	}

	// 昨天的起止时间戳
	startTime, endTime := helper.YesterdayPeriodOfTimestamp(1)

	result, err := h.SourceGroupService.GetNewSourceCount(h.WarpContext(c), request.EnterpriseId, startTime, endTime, request.ThirdIds)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
	return

}

func (h GroupRelationHandle) CollectSelectList(c *gin.Context) {
	var request source.RequestSelectList
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	var (
		result interface{}
		total  int64
	)

	request.ThirdIds, err = h.MembersService.MembersDepartmentsDao.GetThirdsByMemberIds(h.WarpContext(c), request.EnterpriseID, request.MemberID)
	if err != nil {
		return
	}

	switch request.Type {
	// 类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)
	case sourceModel.TypeText:
		result, total, err = h.GroupRelationService.TextSelectList(h.WarpContext(c), request, true)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceModel.ResponseTextList)
		if len(list) == 0 {
			result = []sourceModel.ResponseTextList{}
		}
	case sourceModel.TypeImg:
		result, total, err = h.GroupRelationService.ImageSelectList(h.WarpContext(c), request, true)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceResp.ImageSelectList)
		if len(list) == 0 {
			result = []sourceResp.ImageSelectList{}
		}

	case sourceModel.TypeLink:
		result, total, err = h.GroupRelationService.LinkSelectList(h.WarpContext(c), request, true)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceModel.ResponseLinkList)
		if len(list) == 0 {
			result = []sourceModel.ResponseLinkList{}
		}
		if total > 0 {

			corp, err := orm.BossService.Enterprise.GetInfo(bossModel.EnterpriseModel{ID: request.EnterpriseID})
			if err != nil {
				h.ResError(c, err)
				return
			}
			corpHash := corp.Hash
			for k, v := range list {
				if v.LinkType == sourceModules.LinkUrlTypeSelf {
					domain := config.CfgData.Service.WebH5.Domain
					linkUrl := fmt.Sprintf(sourceModules.SelfLinkUrl, domain, corpHash) + assist.StrVal(v.ID)
					list[k].LinkUrl = linkUrl
				}
			}
		}
	case sourceModel.TypeMiniProgram:
		result, total, err = h.GroupRelationService.MiniProgramSelectList(h.WarpContext(c), request, true)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceModel.ResponseMiniProgramList)
		if len(list) == 0 {
			result = []sourceModel.ResponseMiniProgramList{}
		}
		var pathMapArr = make(map[int64]string)
		for _, v := range list {
			pathMapArr[v.ID] = v.Path
		}
	case sourceModel.TypeVideo:
		result, total, err = h.GroupRelationService.VideoSelectList(h.WarpContext(c), request, true)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceModel.ResponseVideoList)
		if len(list) == 0 {
			result = []sourceModel.ResponseVideoList{}
		}
	case sourceModel.TypeFile:
		result, total, err = h.GroupRelationService.FileSelectList(h.WarpContext(c), request, true)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceModel.ResponseFileList)
		if len(list) == 0 {
			result = []sourceModel.ResponseFileList{}
		}
	case sourceModel.TypeFriend:
		result, total, err = h.GroupRelationService.FriendSelectList(h.WarpContext(c), request, true)
		if err != nil {
			h.ResError(c, err)
			return
		}

		list := result.([]sourceResp.FriendSelectList)
		if len(list) == 0 {
			result = []sourceResp.FriendSelectList{}
		}
	case sourceModel.TypeProduct:
		result, err = sourceService.NewBrandService().CollectList(h.WarpContext(c), source.BrandListReq{
			BrandCacheReq: source.BrandCacheReq{
				EnterpriseId: request.EnterpriseID,
				Mobile:       request.Mobile,
				Sign:         false,
				Accurate:     false,
			},
			Size: 100,
		}, request.Keyword)
		if err != nil {
			h.ResError(c, err)
			return
		}
		list := result.([]sourceResp.BrandProductInfo)
		if len(list) == 0 {
			result = []sourceResp.BrandProductInfo{}
		}
	default:
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, request.Page, request.Size, result)
	return
}
