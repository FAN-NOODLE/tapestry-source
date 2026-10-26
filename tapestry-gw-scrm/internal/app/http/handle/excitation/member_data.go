package excitation

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/app/http/handle"
	excitationReq "gw-scrm/internal/app/http/request/excitation"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/helper"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"math"
	"strconv"
	"time"
)

type MemberDataHandle struct {
	handle.BaseHandle
	api microscrm.MemberDataApi
}

func NewMemberDataHandle() MemberDataHandle {
	return MemberDataHandle{
		api: microscrm.NewMemberDataApi(),
	}
}

func (h MemberDataHandle) List(c *gin.Context) {
	var request excitationReq.MemberDataListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.Page == 0 {
		request.Page = 1
	}

	if request.Size == 0 {
		request.Size = 10
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	res, err := h.api.List(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

// ListExport 后台列表 数据导出
func (h MemberDataHandle) ListExport(c *gin.Context) {
	var request excitationReq.MemberDataListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID

	request.Size = 1000
	request.Page = 1
	request.Total = 0

	list, err := h.api.List(ctx.Wrap(c), request)

	if err != nil {
		h.ResError(c, err)
		return
	}

	jsonStr, _ := json.Marshal(list)
	var data excitationReq.MemberDataListResp
	errJson := json.Unmarshal(jsonStr, &data)
	if errJson != nil {
		h.ResError(c, err)
		return
	}

	var datalist []excitationReq.MemberDataInfoResp
	for _, v := range data.List {
		datalist = append(datalist, v)
	}

	limitEnd := request.Page * request.Size
	if data.Total > limitEnd {

		maxPage := int(math.Ceil(float64(data.Total) / float64(request.Size)))
		for i := 2; i <= maxPage; i++ {
			request.Page = i
			newList, err1 := h.api.List(ctx.Wrap(c), request)
			if err1 != nil {
				h.ResError(c, err1)
				return
			}
			jsonStr1, _ := json.Marshal(newList)
			var data1 excitationReq.MemberDataListResp
			errJson1 := json.Unmarshal(jsonStr1, &data1)
			if errJson1 != nil {
				h.ResError(c, err)
				return
			}
			for _, v := range data1.List {
				datalist = append(datalist, v)
			}
		}
	}
	var isCoach bool
	if helper.GetEnterpriseConfig(c, enterpriseID) == "coach" {
		isCoach = true
	}
	excelFile := excelize.NewFile()
	sheetName := "导购等级数据"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	// 设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "H", 15)
	if nil != err {
		h.ResError(c, err)
		return
	}
	// 设置样式
	style, styleErr := excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"},"font":{"bold":true}}`)
	if nil != styleErr {
		h.ResError(c, styleErr)
		return
	}
	err = excelFile.SetCellStyle(sheetName, "A1", "N1", style)
	if nil != err {
		h.ResError(c, err)
		return
	}
	_ = excelFile.SetCellStr(sheetName, "A1", "导购ID")
	_ = excelFile.SetCellStr(sheetName, "B1", "导购姓名")
	_ = excelFile.SetCellStr(sheetName, "C1", "所在部门")
	_ = excelFile.SetCellStr(sheetName, "D1", "积分")
	_ = excelFile.SetCellStr(sheetName, "E1", "最近积分获取时间")
	if isCoach {
		// 写入表头
		_ = excelFile.SetCellStr(sheetName, "F1", "当前金币")
		_ = excelFile.SetCellStr(sheetName, "G1", "累计总金币")
		_ = excelFile.SetCellStr(sheetName, "H1", "累计消耗金币")
		_ = excelFile.SetCellStr(sheetName, "I1", "等级")
		_ = excelFile.SetCellStr(sheetName, "J1", "最近等级升级时间")
		_ = excelFile.SetCellStr(sheetName, "K1", "奖章数")
		_ = excelFile.SetCellStr(sheetName, "L1", "最近奖章获取时间")
		_ = excelFile.SetCellStr(sheetName, "M1", "历史最高等级")
		_ = excelFile.SetCellStr(sheetName, "N1", "最高等级获取时间")
	} else {
		_ = excelFile.SetCellStr(sheetName, "F1", "等级")
		_ = excelFile.SetCellStr(sheetName, "G1", "最近等级升级时间")
		_ = excelFile.SetCellStr(sheetName, "H1", "奖章数")
		_ = excelFile.SetCellStr(sheetName, "I1", "最近奖章获取时间")
		_ = excelFile.SetCellStr(sheetName, "J1", "历史最高等级")
		_ = excelFile.SetCellStr(sheetName, "K1", "最高等级获取时间")
	}

	startRow := 2
	// 填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		h.ResError(c, err)
		return
	}
	for _, item := range datalist {
		rowStartStr := strconv.Itoa(startRow)
		if item.PointUpdatedAt > 0 {
			item.PointUpdatedAtStr = time.Unix(int64(item.PointUpdatedAt), 0).Format("2006-01-02 15:04:05")
		}
		if item.LevelUpdatedAt > 0 {
			item.LevelUpdatedAtStr = time.Unix(int64(item.LevelUpdatedAt), 0).Format("2006-01-02 15:04:05")
		}
		if item.MedalUpdatedAt > 0 {
			item.MedalUpdatedAtStr = time.Unix(int64(item.MedalUpdatedAt), 0).Format("2006-01-02 15:04:05")
		}
		if item.MaxLevelTime > 0 {
			item.MaxLevelTimeStr = time.Unix(item.MaxLevelTime, 0).Format("2006-01-02 15:04:05")
		}
		points, _ := strconv.ParseFloat(item.PointValue, 64)
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.MemberUserid)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.MemberName)
		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, item.ThirdName)
		_ = excelFile.SetCellFloat(sheetName, "D"+rowStartStr, points, -1, 64)
		_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, item.PointUpdatedAtStr)
		if isCoach {
			_ = excelFile.SetCellFloat(sheetName, "F"+rowStartStr, item.Coins, -1, 64)
			_ = excelFile.SetCellFloat(sheetName, "G"+rowStartStr, item.TotalCoins, -1, 64)
			_ = excelFile.SetCellFloat(sheetName, "H"+rowStartStr, item.SpendCoins, -1, 64)
			_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, item.Level)
			_ = excelFile.SetCellStr(sheetName, "J"+rowStartStr, item.LevelUpdatedAtStr)
			_ = excelFile.SetCellInt(sheetName, "K"+rowStartStr, item.MedalNum)
			_ = excelFile.SetCellStr(sheetName, "L"+rowStartStr, item.MedalUpdatedAtStr)
			_ = excelFile.SetCellStr(sheetName, "M"+rowStartStr, item.Level)
			_ = excelFile.SetCellStr(sheetName, "N"+rowStartStr, item.MaxLevelTimeStr)
		} else {
			_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, item.Level)
			_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, item.LevelUpdatedAtStr)
			_ = excelFile.SetCellInt(sheetName, "H"+rowStartStr, item.MedalNum)
			_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, item.MedalUpdatedAtStr)
			_ = excelFile.SetCellStr(sheetName, "J"+rowStartStr, item.Level)
			_ = excelFile.SetCellStr(sheetName, "K"+rowStartStr, item.MaxLevelTimeStr)
		}

		startRow = startRow + 1
		err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "N"+rowStartStr, style)
		if nil != err {
			h.ResError(c, err)
			return
		}
	}
	excelFile.SetActiveSheet(index)
	fileName := time.Now().Format("200601021504") + ".xlsx"

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return

}
