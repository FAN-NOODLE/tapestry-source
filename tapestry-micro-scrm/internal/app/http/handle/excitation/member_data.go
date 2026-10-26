package excitation

import (
	"github.com/gin-gonic/gin"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	"micro-scrm/internal/app/http/handle"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationSvc "micro-scrm/internal/app/service/excitation"
)

type MemberDataHandle struct {
	handle.BaseHandle
	memberDataSvc *excitationSvc.MemberDataService
}

func NewMemberDataHandle() MemberDataHandle {
	return MemberDataHandle{
		memberDataSvc: excitationSvc.NewMemberDataService(),
	}
}

//DataIndex 小程序导购积分&等级
func (h MemberDataHandle) DataIndex(c *gin.Context) {
	var req excitationReq.MemberDataIndex
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.memberDataSvc.DataIndex(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)
	return
}

//GradePoints 导购等级积分
func (h MemberDataHandle) GradePoints(c *gin.Context) {
	var req excitationReq.GradePoints
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.memberDataSvc.GradePoints(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)
	return
}

// GetPointsRules 小程序积分获取规则
func (h MemberDataHandle) GetPointsRules(c *gin.Context) {
	var req excitationReq.GetPointsRulesReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.memberDataSvc.GetPointsRules(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)
	return
}

// List 后台列表
func (h MemberDataHandle) List(c *gin.Context) {

	var req excitationReq.MemberDataListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	result, total, err := h.memberDataSvc.GetList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, req.Page, req.Size, result)

	return
}

func (h MemberDataHandle) Save(c *gin.Context) {

	var req excitationReq.UpdateDimensionFieldReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.memberDataSvc.SaveDimensionField(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)

	return
}

func (h MemberDataHandle) UpdateLevel(c *gin.Context) {

	var req excitationReq.UpdateLevelReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.memberDataSvc.UpdateLevel(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)

	return
}

// LevelRecord DataIndex 导购等级变更记录
func (h MemberDataHandle) LevelRecord(c *gin.Context) {
	var req excitationReq.LevelRecordReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.memberDataSvc.MemberLevelRecord(h.WarpContext(c), req.EnterpriseId, req.MemberId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	list, err := excitationDao.NewLevelRule().GetLevelList(h.WarpContext(c), req.EnterpriseId)
	if err != nil {
		return
	}
	var maxSign bool
	for _, levelInfo := range list {
		for i, record := range result {
			if record.Level == levelInfo.LevelName {
				if !maxSign {
					result[i].IsMax = 1
					maxSign = true
				}
				result[i].LevelIcon = levelInfo.LevelIcon
			}
		}
	}
	h.ResSuccess(c, result)
	return
}
