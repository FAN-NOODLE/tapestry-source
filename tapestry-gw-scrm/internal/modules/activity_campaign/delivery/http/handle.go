package http

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/activity_campaign/delivery"
	planTask "gw-scrm/internal/modules/plan_task"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"strconv"
	"time"
)

func activitySignInList(c *gin.Context) {
	var request delivery.SignReceiveInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	list, err := microscrm.GetActivitySignInList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)

}

func activitySignInDetail(c *gin.Context) {
	var request delivery.SignDetailInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	detail, err := microscrm.GetActivitySignInDetail(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, detail)
}

func activitySignInExport(c *gin.Context) {
	var request delivery.SignReceiveInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
	}
	request.Type = 1
	request.StartTime = 0
	request.EndTime = 0
	list, err := microscrm.GetActivitySignInList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var newList delivery.NewResultList
	byte, err := json.Marshal(list)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if err := json.Unmarshal(byte, &newList); err != nil {
		ginplus.ResError(c, err)
		return
	}
	excelFile := excelize.NewFile()
	sheetName := "签到数据明细"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	//设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "P", 15)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	//设置样式
	style, styleErr := excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"},"font":{"bold":true}}`)
	if nil != styleErr {
		ginplus.ResError(c, styleErr)
		return
	}
	err = excelFile.SetCellStyle(sheetName, "A1", "P1", style)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "参与用户")
	_ = excelFile.SetCellStr(sheetName, "B1", "手机号")
	_ = excelFile.SetCellStr(sheetName, "C1", "签到天数")
	_ = excelFile.SetCellStr(sheetName, "D1", "最后签到时间")

	startRow := 2
	//填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	for _, item := range newList.SignInList {
		rowStartStr := strconv.Itoa(startRow)
		var timeStyle = "2006-01-02 15:04:05"
		lastTime := time.Unix(item.LastSignTime, 0).Format(timeStyle)
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.Name)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.Phone)
		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, fmt.Sprintf("%d", item.Count))
		_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, lastTime)

		startRow = startRow + 1
		err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "P"+rowStartStr, style)
		if nil != err {
			ginplus.ResError(c, err)
			return
		}
	}
	excelFile.SetActiveSheet(index)
	fileName := sheetName + time.Now().Format("200601021504") + ".xlsx"

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())

	return
}

func ExportUserData(c *gin.Context) {
	params := planTask.ExportUserDataReq{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	userList, err := microscrm.ExportUserData(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	jsonStr, _ := json.Marshal(userList)
	var data []planTask.UserModel
	errJson := json.Unmarshal(jsonStr, &data)
	if errJson != nil {
		ginplus.ResError(c, err)
		return
	}
	excelFile := excelize.NewFile()
	sheetName := "报名活动"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "用户昵称")
	_ = excelFile.SetCellStr(sheetName, "B1", "unionid")
	_ = excelFile.SetCellStr(sheetName, "C1", "会员卡号")
	_ = excelFile.SetCellStr(sheetName, "D1", "收货人姓名")
	_ = excelFile.SetCellStr(sheetName, "E1", "收货人手机号")
	_ = excelFile.SetCellStr(sheetName, "F1", "收货地址")
	_ = excelFile.SetCellStr(sheetName, "G1", "报名时间")
	_ = excelFile.SetCellStr(sheetName, "H1", "审核状态")
	startRow := 2
	cstSh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		ginplus.ResError(c, err)
		cstSh = time.Local
	}
	for _, item := range data {
		rowStartStr := strconv.Itoa(startRow)
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.WxName)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.UnionId)
		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, item.MemberCode)
		_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, item.ReceiverName)
		_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, item.ReceiverPhone)
		_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, item.ReceiverAddress)
		_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, item.CreatedTimeText)
		_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, item.StatusText)
		startRow = startRow + 1
	}
	excelFile.SetActiveSheet(index)
	fileName := "报名活动" + time.Now().In(cstSh).Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return
}

func GetActivityList(c *gin.Context) {
	activityList, err := microscrm.GetActivityList(ctx.Wrap(c), nil)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, activityList)
}

func GetUserList(c *gin.Context) {
	params := planTask.UserListReq{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	activityList, err := microscrm.GetUserList(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, activityList)
}

func GetUserInfo(c *gin.Context) {
	params := planTask.UserInfoReq{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	activityList, err := microscrm.GetActivityUserInfo(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, activityList)
}

func VerifyUserInfo(c *gin.Context) {
	params := planTask.UserInfoReq{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	activityList, err := microscrm.VerifyActivityUserInfo(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, activityList)
}
