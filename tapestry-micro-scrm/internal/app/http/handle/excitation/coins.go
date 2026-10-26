package excitation

import (
	"github.com/gin-gonic/gin"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/http/handle"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	"micro-scrm/internal/pkg/helper"
)

// CoinsHandle 导购积分模块
type CoinsHandle struct {
	handle.BaseHandle
	coinsRecordSvc *excitationSvc.CoinsService
}

func NewCoinsHandle() CoinsHandle {
	return CoinsHandle{
		coinsRecordSvc: excitationSvc.NewCoinsService(),
	}
}

// AddManualCoins  手动添加金币
func (h CoinsHandle) AddManualCoins(c *gin.Context) {
	var req excitationReq.ManualCoinsReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	var userIds []string
	for _, data := range req.ManualCoinsData {
		userIds = append(userIds, data.UserId)
	}
	members, err := membersDao.NewMembers().GetMany(h.WarpContext(c), map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"userid":        userIds,
	})
	if err != nil {
		return
	}

	var dataList []excitationReq.ManualCoinsRow //要写入的积分数据
	var userToMember = make(map[string]int64)
	for _, member := range members {
		userToMember[member.Userid] = member.ID
	}
	for _, data := range req.ManualCoinsData {
		if _, ok := userToMember[data.UserId]; ok {
			dataList = append(dataList, excitationReq.ManualCoinsRow{
				MemberId:    userToMember[data.UserId],
				UserId:      data.UserId,
				CoinsRemark: data.CoinsRemark,
				Coins:       helper.MathRounded(data.Coins, 1),
			})
		}
	}
	err = excitationSvc.NewCoinsService().ManualCoins(h.WarpContext(c), req.EnterpriseId, dataList)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

// GetCoinsList  小程序端导购金币记录列表
func (h CoinsHandle) GetCoinsList(c *gin.Context) {
	var req excitationReq.GetCoinsListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := excitationSvc.NewCoinsService().GetCoinsList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)

	return
}

// GetCoinsSettingsInfo     小程序端商城开启信息
func (h CoinsHandle) GetCoinsSettingsInfo(c *gin.Context) {
	var request excitationReq.GetCoinSettingInfoReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.coinsRecordSvc.GetCoinsSettingInfo(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)

	return
}

// SaveCoinsSettings   金币相关配置设置
func (h CoinsHandle) SaveCoinsSettings(c *gin.Context) {
	var req excitationReq.SaveCoinsSettingsReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.coinsRecordSvc.SaveCoinsSetting(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)

	return
}

// RecordList   获取导购积分获取记录列表
func (h CoinsHandle) RecordList(c *gin.Context) {

	var req excitationReq.CoinsRecordListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := excitationSvc.NewCoinsService().GetRecordList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)

	return
}
