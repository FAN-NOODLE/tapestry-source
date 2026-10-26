package reserve

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"gw-scrm/internal/app/http/handle"
	reserveReq "gw-scrm/internal/app/http/request/reserve"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/helper"
	"gw-scrm/internal/pkg/sdk/downloadcenter"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

const enterpriseNameCoach = "coach"
const enterpriseNameKS = "ks"
const enterpriseNameSW = "sw"

type ReserveHandle struct {
	handle.BaseHandle
	reserveSdk microscrm.ReserveApi
}

func NewReserveHandle() ReserveHandle {
	return ReserveHandle{
		reserveSdk: microscrm.NewReserveApi(),
	}
}

// Lists 预约列表
func (s ReserveHandle) Lists(c *gin.Context) {

	var request reserveReq.OrderListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	res, err := s.reserveSdk.OrderList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// OrderListsExport 导出预约单
func (s ReserveHandle) OrderListsExport(c *gin.Context) {
	var request reserveReq.OrderListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	// 获取企业配置
	enterpriseConfigName := helper.GetEnterpriseConfig(c, enterpriseID)
	fileName := time.Now().Format("20060102150405") + ".xlsx"
	// uniqueMap := make(map[string]string)
	excelFile := excelize.NewFile()
	sheetName := "导出预约订单"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	excelFile.SetActiveSheet(index)
	err := excelFile.SaveAs(fileName)
	defer os.Remove(fileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	stream, err := excelFile.NewStreamWriter(sheetName)
	if err != nil {
		return
	}
	// 设置样式
	// style, styleErr := excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"},"font":{"bold":true}}`)
	// if nil != styleErr {
	// 	ginplus.ResError(c, styleErr)
	// 	return
	// }
	// err = excelFile.SetCellStyle(sheetName, "A1", "AA1", style)
	// if nil != err {
	// 	ginplus.ResError(c, err)
	// 	return
	// }
	// 写入表头
	row := make([]any, 0)
	if enterpriseConfigName == enterpriseNameCoach {
		row = []any{
			"编号", "顾客昵称", "预约服务", "到店时间", "邀请来源", "预约门店", "门店地址", "预约订单是否是分配单", "预约导购", "导购工号", "预约单的状态", "预约服务的备注信息(导购)",
			"预约服务的备注信息(用户)", "用户是否打开预约单", "用户是否确认预约单", "邀约次数", "是否成单", "成单时间", "成单金额", "是否到店购买", "顾客类型", "是否评价", "评价时间",
			"活动安排", "整体服务", "评价内容", "创建时间",
		}
		// _ = excelFile.SetCellStr(sheetName, "A1", "编号")
		// _ = excelFile.SetCellStr(sheetName, "B1", "顾客昵称")
		// _ = excelFile.SetCellStr(sheetName, "C1", "预约服务")
		// _ = excelFile.SetCellStr(sheetName, "D1", "到店时间")
		// _ = excelFile.SetCellStr(sheetName, "E1", "邀请来源")
		// _ = excelFile.SetCellStr(sheetName, "F1", "预约门店")
		// _ = excelFile.SetCellStr(sheetName, "G1", "门店地址")
		// _ = excelFile.SetCellStr(sheetName, "H1", "预约订单是否是分配单")
		// _ = excelFile.SetCellStr(sheetName, "I1", "预约导购")
		// _ = excelFile.SetCellStr(sheetName, "J1", "导购工号")
		// _ = excelFile.SetCellStr(sheetName, "K1", "预约单的状态")
		// _ = excelFile.SetCellStr(sheetName, "L1", "预约服务的备注信息(导购)")
		// _ = excelFile.SetCellStr(sheetName, "M1", "预约服务的备注信息(用户)")
		// _ = excelFile.SetCellStr(sheetName, "N1", "用户是否打开预约单")
		// _ = excelFile.SetCellStr(sheetName, "O1", "用户是否确认预约单")
		// _ = excelFile.SetCellStr(sheetName, "P1", "邀约次数")
		// _ = excelFile.SetCellStr(sheetName, "Q1", "是否成单")
		// _ = excelFile.SetCellStr(sheetName, "R1", "成单时间")
		// _ = excelFile.SetCellStr(sheetName, "S1", "成单金额")
		// _ = excelFile.SetCellStr(sheetName, "T1", "是否到店购买")
		// _ = excelFile.SetCellStr(sheetName, "U1", "顾客类型")
		// _ = excelFile.SetCellStr(sheetName, "V1", "是否评价")
		// _ = excelFile.SetCellStr(sheetName, "W1", "评价时间")
		// _ = excelFile.SetCellStr(sheetName, "X1", "活动安排")
		// _ = excelFile.SetCellStr(sheetName, "Y1", "整体服务")
		// _ = excelFile.SetCellStr(sheetName, "Z1", "评价内容")
		// _ = excelFile.SetCellStr(sheetName, "AA1", "创建时间")
		// _ = excelFile.SetCellStr(sheetName, "AA1", "门店类型")
	} else {
		row = []any{
			"编号", "顾客昵称", "预约服务", "到店时间", "邀请来源", "预约门店", "门店地址", "预约订单是否是分配单", "预约导购", "导购工号", "预约单的状态", "预约服务的备注信息(导购)",
			"预约服务的备注信息(用户)", "用户是否打开预约单", "用户是否确认预约单", "邀约次数", "签到状态", "是否成单", "成单时间", "成单金额", "是否到店购买", "顾客类型", "是否评价", "评价时间",
			"活动安排", "整体服务", "评价内容",
		}
		// _ = excelFile.SetCellStr(sheetName, "A1", "编号")
		// _ = excelFile.SetCellStr(sheetName, "B1", "顾客昵称")
		// _ = excelFile.SetCellStr(sheetName, "C1", "预约服务")
		// _ = excelFile.SetCellStr(sheetName, "D1", "到店时间")
		// _ = excelFile.SetCellStr(sheetName, "E1", "邀请来源")
		// _ = excelFile.SetCellStr(sheetName, "F1", "预约门店")
		// _ = excelFile.SetCellStr(sheetName, "G1", "门店地址")
		// _ = excelFile.SetCellStr(sheetName, "H1", "预约订单是否是分配单")
		// _ = excelFile.SetCellStr(sheetName, "I1", "预约导购")
		// _ = excelFile.SetCellStr(sheetName, "J1", "导购工号")
		// _ = excelFile.SetCellStr(sheetName, "K1", "预约单的状态")
		// _ = excelFile.SetCellStr(sheetName, "L1", "预约服务的备注信息(导购)")
		// _ = excelFile.SetCellStr(sheetName, "M1", "预约服务的备注信息(用户)")
		// _ = excelFile.SetCellStr(sheetName, "N1", "用户是否打开预约单")
		// _ = excelFile.SetCellStr(sheetName, "O1", "用户是否确认预约单")
		// _ = excelFile.SetCellStr(sheetName, "P1", "邀约次数")
		// _ = excelFile.SetCellStr(sheetName, "Q1", "签到状态")
		// _ = excelFile.SetCellStr(sheetName, "R1", "是否成单")
		// _ = excelFile.SetCellStr(sheetName, "S1", "成单时间")
		// _ = excelFile.SetCellStr(sheetName, "T1", "成单金额")
		// _ = excelFile.SetCellStr(sheetName, "U1", "是否到店购买")
		// _ = excelFile.SetCellStr(sheetName, "V1", "顾客类型")
		// _ = excelFile.SetCellStr(sheetName, "W1", "是否评价")
		// _ = excelFile.SetCellStr(sheetName, "X1", "评价时间")
		// _ = excelFile.SetCellStr(sheetName, "Y1", "活动安排")
		// _ = excelFile.SetCellStr(sheetName, "Z1", "整体服务")
		// _ = excelFile.SetCellStr(sheetName, "AA1", "评价内容")
	}
	if enterpriseConfigName == enterpriseNameSW {
		row = []any{
			"编号", "顾客昵称", "预约服务", "到店时间", "邀请来源", "预约门店", "门店地址", "预约订单是否是分配单", "预约导购", "导购工号", "预约单的状态", "预约服务的备注信息(导购)",
			"预约服务的备注信息(用户)", "用户是否打开预约单", "是否签到", "邀约次数", "签到状态", "是否成单", "成单时间", "成单金额", "是否到店购买", "顾客类型", "是否评价", "评价时间",
			"活动安排", "整体服务", "评价内容",
		}
		// _ = excelFile.SetCellStr(sheetName, "O1", "是否签到")
	}
	err = stream.SetRow("A1", row)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	startRow := int64(2)

	// 填充数据
	// style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	// if nil != err {
	// 	ginplus.ResError(c, err)
	// 	return
	// }
	request.Size = 10000
	request.Page = 0
	wg := threading.RoutineGroup{}

	for {
		request.Page++
		newList, err1 := s.reserveSdk.OrderExportList(ctx.Wrap(c), request)
		if err1 != nil {
			s.ResError(c, err1)
			return
		}
		jsonStr1, _ := json.Marshal(newList)
		var data1 reserveReq.ResultOrderList
		errJson1 := json.Unmarshal(jsonStr1, &data1)
		if errJson1 != nil {
			s.ResError(c, err)
			return
		}
		if len(data1.List) == 0 {
			break
		}
		datalist := data1.List
		wg.RunSafe(func() {
			for _, item := range datalist {
				rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
				if item.Status == 1 {
					item.StatusStr = "邀约成功"
				}
				if item.Status == 2 {
					item.StatusStr = "邀约完成"
				}
				if enterpriseConfigName == enterpriseNameSW {
					if item.Status == 0 {
						item.StatusStr = "邀约生效中"
					}
					if item.Status == 4 {
						item.StatusStr = "邀约匹配失败"
					}
				}
				switch item.OrderOrigin {
				case 1:
					item.OrderOriginStr = "到店邀约"
				case 2:
					item.OrderOriginStr = "邀请函"
				}
				customerStatus := "-"
				switch item.CustomerStatus {
				case 1:
					item.CustomerStatusStr = "新客"
					customerStatus = "1"
				case 2:
					item.CustomerStatusStr = "老客"
					customerStatus = "0"
				case 4:
					item.CustomerStatusStr = "沉睡客"
				default:
					item.CustomerStatusStr = "--"

				}
				if item.VipType == 1 {
					item.CustomerStatusStr += ",VIP客"
				}
				var ReserveOrderPriceFloat decimal.Decimal
				if item.ReserveOrderPrice > 0 {
					ReserveOrderPriceFloat = decimal.NewFromInt(int64(item.ReserveOrderPrice)).Div(decimal.NewFromInt(100))
				}
				priceFloat, _ := ReserveOrderPriceFloat.Float64()
				isPay := item.IsPay
				if item.IsPay != 1 {
					isPay = 0
				}
				isStoreStatus := "-"
				isStore := item.IsStore
				switch item.IsStore {
				case 1:
					isStoreStatus = "1"
				case 2:
					isStoreStatus = "0"
				}
				if item.IsStore != 1 {
					isStore = 0
				}

				if enterpriseConfigName != enterpriseNameCoach {
					row = []any{
						item.ReserveNumber, item.UserName, item.ReserveName, item.ReserveTimeStr, item.OrderOriginStr, item.StoreName,
						item.StoreAddress, item.ReserveTypeStr, item.MemberName, item.MemberUserID, item.StatusStr, item.ReserveRemark, item.UserRemark,
						item.IsPayStr, item.IsCommitStr, item.InviteNum, item.IsSignInStr, isPay, item.ReserveOrderTimeStr, priceFloat, isStore, item.CustomerStatusStr,
						item.IsCommentStr, item.CommentTimeStr, fmt.Sprintf("%d星", item.ServiceScore), fmt.Sprintf("%d星", item.StoreScore), item.Remark,
					}

				} else {
					row = []any{
						item.ReserveNumber, item.UserName, item.ReserveName, item.ReserveTimeStr, item.OrderOriginStr, item.StoreName,
						item.StoreAddress, item.ReserveTypeStr, item.MemberName, item.MemberUserID, item.StatusStr, item.ReserveRemark, item.UserRemark,
						item.IsPayStr, item.IsCommitStr, item.InviteNum, isPay, item.ReserveOrderTimeStr, priceFloat, isStoreStatus, customerStatus, item.IsCommentStr,
						item.CommentTimeStr, fmt.Sprintf("%d星", item.ServiceScore), fmt.Sprintf("%d星", item.StoreScore), item.Remark, item.CreateTimeStr,
					}
				}
				err = stream.SetRow("A"+rowStartStr, row)
				// _ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.ReserveNumber)
				// _ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.UserName)
				// _ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, item.ReserveName)
				// _ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, item.ReserveTimeStr)
				// _ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, item.OrderOriginStr)
				// _ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, item.StoreName)
				// _ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, item.StoreAddress)
				// _ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, item.ReserveTypeStr)
				// _ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, item.MemberName)
				// _ = excelFile.SetCellStr(sheetName, "J"+rowStartStr, item.MemberUserID)
				// _ = excelFile.SetCellStr(sheetName, "K"+rowStartStr, item.StatusStr)
				// _ = excelFile.SetCellStr(sheetName, "L"+rowStartStr, item.ReserveRemark)
				// _ = excelFile.SetCellStr(sheetName, "M"+rowStartStr, item.UserRemark)
				// _ = excelFile.SetCellStr(sheetName, "N"+rowStartStr, item.IsPayStr)
				// _ = excelFile.SetCellStr(sheetName, "O"+rowStartStr, item.IsCommitStr)
				// _ = excelFile.SetCellValue(sheetName, "P"+rowStartStr, item.InviteNum)
				// if enterpriseConfigName != enterpriseNameCoach {
				// 	_ = excelFile.SetCellStr(sheetName, "Q"+rowStartStr, item.IsSignInStr)
				// 	_ = excelFile.SetCellValue(sheetName, "R"+rowStartStr, isPay)
				// 	_ = excelFile.SetCellStr(sheetName, "S"+rowStartStr, item.ReserveOrderTimeStr)
				// 	_ = excelFile.SetCellValue(sheetName, "T"+rowStartStr, priceFloat)
				// 	_ = excelFile.SetCellValue(sheetName, "U"+rowStartStr, isStore)
				// 	_ = excelFile.SetCellStr(sheetName, "V"+rowStartStr, item.CustomerStatusStr)
				// 	_ = excelFile.SetCellStr(sheetName, "W"+rowStartStr, item.IsCommentStr)
				// 	_ = excelFile.SetCellStr(sheetName, "X"+rowStartStr, item.CommentTimeStr)
				// 	_ = excelFile.SetCellStr(sheetName, "Y"+rowStartStr, fmt.Sprintf("%d星", item.ServiceScore))
				// 	_ = excelFile.SetCellStr(sheetName, "Z"+rowStartStr, fmt.Sprintf("%d星", item.StoreScore))
				// 	_ = excelFile.SetCellStr(sheetName, "AA"+rowStartStr, item.Remark)
				// } else {
				// 	_ = excelFile.SetCellValue(sheetName, "Q"+rowStartStr, isPay)
				// 	_ = excelFile.SetCellStr(sheetName, "R"+rowStartStr, item.ReserveOrderTimeStr)
				// 	_ = excelFile.SetCellValue(sheetName, "S"+rowStartStr, priceFloat)
				// 	_ = excelFile.SetCellValue(sheetName, "T"+rowStartStr, isStore)
				// 	_ = excelFile.SetCellValue(sheetName, "U"+rowStartStr, customerStatus)
				// 	_ = excelFile.SetCellStr(sheetName, "V"+rowStartStr, item.IsCommentStr)
				// 	_ = excelFile.SetCellStr(sheetName, "W"+rowStartStr, item.CommentTimeStr)
				// 	_ = excelFile.SetCellStr(sheetName, "X"+rowStartStr, fmt.Sprintf("%d星", item.ServiceScore))
				// 	_ = excelFile.SetCellStr(sheetName, "Y"+rowStartStr, fmt.Sprintf("%d星", item.StoreScore))
				// 	_ = excelFile.SetCellStr(sheetName, "Z"+rowStartStr, item.Remark)
				// 	_ = excelFile.SetCellStr(sheetName, "AA"+rowStartStr, item.CreateTimeStr)
				// }

				atomic.AddInt64(&startRow, 1)
				// err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "AA"+rowStartStr, style)
				// if nil != err {
				// 	ginplus.ResError(c, err)
				// 	return
				// }
			}
		})
	}
	wg.Wait()
	err = stream.Flush()
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	err = excelFile.SaveAs(fileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	des, err := os.ReadFile(fileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	_, _ = c.Writer.Write(des)
	return

}

func (s ReserveHandle) OrderListsExportNew(c *gin.Context) {
	var request reserveReq.OrderListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	enterpriseId := tokenInfo.EnterpriseId
	adminId := tokenInfo.AdminId
	enterpriseConfigName := helper.GetEnterpriseConfig(c, enterpriseId)
	request.EnterpriseID = enterpriseId

	sourceTitle := make([]string, 0)
	sourceTitleAlias := make(map[string]string, 0)
	sourceTitleValue := make([]string, 0)
	if enterpriseConfigName != enterpriseNameCoach {
		sourceTitle = []string{
			"reserve_number", "user_name", "reserve_name", "reserve_time_str", "order_origin_str", "store_name", "store_address", "reserve_type_str", "member_name",
			"member_user_id", "status_str", "reserve_remark", "user_remark", "is_open_str", "is_commit_str", "invite_num", "is_pay_str", "reserve_order_time_str",
			"order_price", "is_store_str", "customer_status_str", "is_comment_str", "comment_time_str", "service_score_str", "store_score_str", "remark", "create_time_str",
		}
		sourceTitleAlias = map[string]string{
			"reserve_number": "编号", "user_name": "顾客昵称", "reserve_name": "预约服务", "reserve_time_str": "到店时间", "order_origin_str": "邀请来源", "store_name": "预约门店", "store_address": "门店地址",
			"reserve_type_str": "预约订单是否是分配单", "member_name": "预约导购", "member_user_id": "导购工号", "status_str": "预约单的状态", "reserve_remark": "预约服务的备注信息(导购)",
			"user_remark": "预约服务的备注信息(用户)", "is_open_str": "用户是否打开预约单", "is_commit_str": "用户是否确认预约单", "invite_num": "邀约次数", "is_pay_str": "是否成单", "reserve_order_time_str": "成单时间",
			"order_price": "成单金额", "is_store_str": "是否到店购买", "customer_status_str": "顾客类型", "is_comment_str": "是否评价", "comment_time_str": "评价时间", "service_score_str": "活动安排",
			"store_score_str": "整体服务", "remark": "评价内容", "create_time_str": "创建时间",
		}
		sourceTitleValue = []string{
			"order_price", "invite_num",
		}

	} else {
		sourceTitle = []string{
			"reserve_number", "user_name", "reserve_name", "reserve_time_str", "order_origin_str", "store_name", "store_address", "reserve_type_str", "member_name",
			"member_user_id", "status_str", "reserve_remark", "user_remark", "is_open_str", "is_commit_str", "invite_num", "is_pay", "reserve_order_time_str",
			"order_price", "relation_order_num", "is_store_str", "customer_status_str", "create_time_str",
		}
		sourceTitleAlias = map[string]string{
			"reserve_number": "编号", "user_name": "顾客昵称", "reserve_name": "预约服务", "reserve_time_str": "到店时间", "order_origin_str": "邀请来源", "store_name": "预约门店", "store_address": "门店地址",
			"reserve_type_str": "预约订单是否是分配单", "member_name": "预约导购", "member_user_id": "导购工号", "status_str": "预约单的状态", "reserve_remark": "预约服务的备注信息(导购)",
			"user_remark": "预约服务的备注信息(用户)", "is_open_str": "用户是否打开预约单", "is_commit_str": "用户是否确认预约单", "invite_num": "邀约次数", "is_pay": "是否成单(1是，0否)", "reserve_order_time_str": "成单时间",
			"order_price": "成单金额", "relation_order_num": "关联订单数", "is_store_str": "是否到店购买(1是,0否)", "customer_status_str": "顾客类型(1新客，0老客)", "create_time_str": "创建时间",
		}
		sourceTitleValue = []string{
			"order_price", "is_pay", "invite_num", "relation_order_num",
		}
	}
	// 邀约单导出
	_, err := downloadcenter.DC.CreateExport(ctx.Wrap(c), downloadcenter.CreateExportRequest{
		FileName:           "到店邀约--邀约单列表导出",
		Type:               downloadcenter.FileTypeExcel,
		StoreId:            1,
		EnterpriseId:       enterpriseId,
		OperationId:        adminId,
		OperationName:      "adminId",
		SourceType:         downloadcenter.SourceTypeRestful,
		SourceTitleAlias:   sourceTitleAlias,
		SourceTitle:        sourceTitle,
		SourceTitleValue:   sourceTitleValue,
		SourceHeaderPrefix: nil,
		Label:              downloadcenter.ReserveOrderExport,
		CreateData:         downloadcenter.CreateData{},
		CreateSql:          downloadcenter.CreateSql{},
		CreateRestful: downloadcenter.CreateRestful{
			SourceRestfulPath:             microscrm.GetUrl("%s/micro-scrm/v1/reserve/order/export-lists"),
			SourceRestfulMethod:           "POST",
			SourceRestfulHeaders:          map[string]interface{}{"enterprise-id": strconv.FormatInt(enterpriseId, 10)},
			SourceRestfulParams:           helper.StructToMap(request),
			SourceSuffix:                  nil,
			SourceRestfulCountField:       "data.total",
			SourceRestfulDataField:        "data.list",
			SourceRestfulRequestPageField: "page",
			SourceRestfulRequestSizeField: "size",
		},
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, "下载成功")
}

// Detail 预约单详情
func (s ReserveHandle) Detail(c *gin.Context) {

	var request reserveReq.OrderDetailReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	// 获取企业配置
	// params := common.GetRawMap(c)
	// params["enterprise_id"] = enterpriseID
	// params["key"] = enterpriseNameKey
	// enterpriseConfigs, err := microboss.GetEnterpriseName(ctx.Wrap(c), params)
	// var enterpriseConfigArr []microboss.Config
	// var enterpriseConfig microboss.Config
	// if err != nil {
	//	enterpriseConfig = microboss.Config{
	//	}
	// }
	//
	// result , _ :=json.Marshal(enterpriseConfigs)
	// err = json.Unmarshal(result, &enterpriseConfigArr)
	// if err != nil {
	//	enterpriseConfig.Name = ""
	// }
	// for _,v := range enterpriseConfigArr {
	//	enterpriseConfig.Name = v.Name
	//	enterpriseConfig.Key = v.Key
	// }
	request.EnterpriseName = helper.GetEnterpriseConfig(c, enterpriseID)

	res, err := s.reserveSdk.OrderDetail(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s ReserveHandle) ServiceCreate(c *gin.Context) {

	var request reserveReq.ReserveServiceCreateReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	if request.ServiceType == 0 {
		if request.StartTime == 0 || request.EndTime == 0 {
			s.ResError(c, errors.New("有效期的时间段格式配置错误"))
			return
		}
	}
	res, err := s.reserveSdk.ServiceCreate(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s ReserveHandle) ServiceUpdate(c *gin.Context) {

	var request reserveReq.ReserveServiceUpdateReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID
	if request.ServiceType == 0 {
		if request.StartTime == 0 || request.EndTime == 0 {
			s.ResError(c, errors.New("有效期的时间段格式配置错误"))
			return
		}
	}
	res, err := s.reserveSdk.ServiceUpdate(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s ReserveHandle) ServiceDelete(c *gin.Context) {

	var request reserveReq.ReserveServiceDeleteReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	res, err := s.reserveSdk.ServiceDelete(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

func (s ReserveHandle) ServiceList(c *gin.Context) {

	var request reserveReq.ReserveServiceListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	res, err := s.reserveSdk.ServiceList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// ScoreList 用户评分列表
func (s ReserveHandle) ScoreList(c *gin.Context) {

	var request reserveReq.ServiceScoreListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	res, err := s.reserveSdk.ScoreList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, res)
	return
}

// ScoreExport 导出评价列表
func (s ReserveHandle) ScoreExport(c *gin.Context) {
	var request reserveReq.ServiceScoreListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)

	// 获取企业配置
	enterpriseConfigName := helper.GetEnterpriseConfig(c, enterpriseID)
	excelFile := excelize.NewFile()
	sheetName := "导出评价列表"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	// 设置宽度
	err := excelFile.SetColWidth(sheetName, "A", "H", 15)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	// 设置样式
	style, styleErr := excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"},"font":{"bold":true}}`)
	if nil != styleErr {
		ginplus.ResError(c, styleErr)
		return
	}
	err = excelFile.SetCellStyle(sheetName, "A1", "J1", style)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "预约单编号")
	_ = excelFile.SetCellStr(sheetName, "B1", "用户名称")
	_ = excelFile.SetCellStr(sheetName, "C1", "服务导购名称")
	_ = excelFile.SetCellStr(sheetName, "D1", "预约门店名称")
	_ = excelFile.SetCellStr(sheetName, "E1", "到店时间")
	if enterpriseConfigName == enterpriseNameCoach {
		_ = excelFile.SetCellStr(sheetName, "F1", "评分时间")
		_ = excelFile.SetCellStr(sheetName, "G1", "评论内容")

	} else if enterpriseConfigName == enterpriseNameSW {
		_ = excelFile.SetCellStr(sheetName, "F1", "活动安排")
		_ = excelFile.SetCellStr(sheetName, "G1", "整体服务")
		_ = excelFile.SetCellStr(sheetName, "H1", "评分时间")
		_ = excelFile.SetCellStr(sheetName, "I1", "评论内容")
	} else if enterpriseConfigName == enterpriseNameKS {
		_ = excelFile.SetCellStr(sheetName, "F1", "店铺环境")
		_ = excelFile.SetCellStr(sheetName, "G1", "服务流程")
		_ = excelFile.SetCellStr(sheetName, "H1", "店员服务")
		_ = excelFile.SetCellStr(sheetName, "I1", "评分时间")
		_ = excelFile.SetCellStr(sheetName, "J1", "评论内容")
	}

	startRow := int64(2)
	// 填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	wg := threading.RoutineGroup{}
	request.EnterpriseID = enterpriseID
	request.Size = 1000
	request.Page = 0
	request.Total = int64(0)
	for {
		request.Page++
		list, err := s.reserveSdk.ScoreList(ctx.Wrap(c), request)
		if err != nil {
			s.ResError(c, err)
			return
		}

		jsonStr, _ := json.Marshal(list)
		var data reserveReq.ResultScoreList

		errJson := json.Unmarshal(jsonStr, &data)
		if errJson != nil {
			s.ResError(c, err)
			return
		}
		if len(data.List) == 0 {
			break
		}
		datalist := data.List
		wg.RunSafe(func() {
			for _, item := range datalist {
				rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
				item.CreateTimeStr = time.Unix(item.CreatedAt, 0).Format("2006-01-02 15:04:05")
				item.ReserveTimeStr = time.Unix(int64(item.ReserveTime), 0).Format("2006-01-02 15:04:05")
				_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.ReserveNumber)
				_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.UserName)
				_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, item.MemberName)
				_ = excelFile.SetCellStr(sheetName, "D"+rowStartStr, item.StoreName)
				_ = excelFile.SetCellStr(sheetName, "E"+rowStartStr, item.ReserveTimeStr)

				if enterpriseConfigName == enterpriseNameCoach {
					_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, item.CreateTimeStr)
					_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, item.Remark)

				} else if enterpriseConfigName == enterpriseNameSW {
					_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, strconv.Itoa(int(item.ServiceScore)))
					_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, strconv.Itoa(int(item.StoreScore)))
					_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, item.CreateTimeStr)
					_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, item.Remark)
				} else if enterpriseConfigName == enterpriseNameKS {
					_ = excelFile.SetCellStr(sheetName, "F"+rowStartStr, strconv.Itoa(int(item.FlowScore)))
					_ = excelFile.SetCellStr(sheetName, "G"+rowStartStr, strconv.Itoa(int(item.ServiceScore)))
					_ = excelFile.SetCellStr(sheetName, "H"+rowStartStr, strconv.Itoa(int(item.StoreScore)))
					_ = excelFile.SetCellStr(sheetName, "I"+rowStartStr, item.CreateTimeStr)
					_ = excelFile.SetCellStr(sheetName, "J"+rowStartStr, item.Remark)
				}

				atomic.AddInt64(&startRow, 1)
				err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "J"+rowStartStr, style)
				if nil != err {
					ginplus.ResError(c, err)
					return
				}
			}
		})
	}
	wg.Wait()
	excelFile.SetActiveSheet(index)
	fileName := time.Now().Format("200601021504") + ".xlsx"

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return

}

// ReserveOrderList 顾客受邀单列表
func (s ReserveHandle) ReserveOrderList(c *gin.Context) {

	var request reserveReq.ReserveOrderListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	res, err := s.reserveSdk.ReserveOrderList(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, res)
	return
}

// ReserveOrderListExport 导出顾客受邀单列表
func (s ReserveHandle) ReserveOrderListExport(c *gin.Context) {
	var request reserveReq.ReserveOrderListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)

	// 获取企业配置

	enterpriseConfigName := helper.GetEnterpriseConfig(c, enterpriseID)
	fileName := time.Now().Format("20060102150405") + ".xlsx"
	// uniqueMap := make(map[string]string)
	excelFile := excelize.NewFile()
	sheetName := "导出顾客受邀单列表"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	excelFile.SetActiveSheet(index)
	err := excelFile.SaveAs(fileName)
	defer os.Remove(fileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	stream, err := excelFile.NewStreamWriter(sheetName)
	if err != nil {
		return
	}

	row := make([]any, 0)
	if enterpriseConfigName == enterpriseNameCoach {
		row = []any{
			"预约单编号", "顾客昵称", "服务导购", "预约门店", "邀请来源", "邀约次数", "预约状态", "是否打开", "是否确认", "顾客类型", "是否成单", "是否到店购买", "导购成单金额",
			"成单时间", "是否评价", "评价时间", "创建时间",
		}

	} else {
		row = []any{
			"预约单编号", "顾客昵称", "服务导购", "预约门店", "邀请来源", "邀约次数", "确认次数", "创建时间",
		}
	}
	err = stream.SetRow("A1", row)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	startRow := int64(2)
	wg := threading.RoutineGroup{}
	request.EnterpriseID = enterpriseID
	request.Size = 10000
	request.Page = 0
	request.Total = int64(0)
	for {
		request.Page++
		list, err := s.reserveSdk.ReserveOrderList(ctx.Wrap(c), request)
		if err != nil {
			s.ResError(c, err)
			return
		}

		jsonStr, _ := json.Marshal(list)
		var data reserveReq.ResultReserveOrderList

		errJson := json.Unmarshal(jsonStr, &data)
		if errJson != nil {
			s.ResError(c, err)
			return
		}
		if len(data.List) == 0 {
			break
		}
		datalist := data.List
		wg.RunSafe(func() {
			for _, item := range datalist {
				rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
				item.CreateTimeStr = time.Unix(item.CreatedAt, 0).Format("2006-01-02 15:04:05")

				if item.OrderOrigin == 1 {
					item.OrderOriginStr = "到店邀约"
				}
				if item.OrderOrigin == 2 {
					item.OrderOriginStr = "邀请函"
				}

				var ReserveOrderPriceFloat decimal.Decimal
				if item.ReserveOrderPrice > 0 {
					ReserveOrderPriceFloat = decimal.NewFromInt(int64(item.ReserveOrderPrice)).Div(decimal.NewFromInt(100))
				}
				priceFloat, _ := ReserveOrderPriceFloat.Float64()

				if item.Status == 1 {
					item.StatusStr = "邀约成功"
				}
				if item.Status == 2 {
					item.StatusStr = "邀约完成"
				}
				var CustomerStatusStr string
				switch item.CustomerStatus {
				case 0:
					CustomerStatusStr = "--"
				case 1:
					CustomerStatusStr = "新客"
				case 2:
					CustomerStatusStr = "老客"
				}

				if enterpriseConfigName == enterpriseNameCoach {
					row = []any{
						item.ReserveNumber, item.UserName, item.MemberName, item.StoreName, item.OrderOriginStr, item.InviteNum,
						item.StatusStr, item.IsOpenStr, item.IsCommitStr, CustomerStatusStr, item.IsPay, item.IsStore, priceFloat, item.ReserveOrderTimeStr,
						item.IsCommentStr, item.CommentTimeStr, item.CreateTimeStr,
					}

				} else {
					row = []any{
						item.ReserveNumber, item.UserName, item.MemberName, item.StoreName, item.OrderOriginStr, item.InviteNum,
						item.IsCommit, item.CreateTimeStr,
					}
				}
				err = stream.SetRow("A"+rowStartStr, row)

				atomic.AddInt64(&startRow, 1)

				// err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "Q"+rowStartStr, style)
				if nil != err {
					ginplus.ResError(c, err)
					return
				}
			}
		})
	}
	wg.Wait()
	err = stream.Flush()
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	err = excelFile.SaveAs(fileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	des, err := os.ReadFile(fileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	_, _ = c.Writer.Write(des)
	return

}

// ServiceCheckDelete 预约服务的删除验证
func (s ReserveHandle) ServiceCheckDelete(c *gin.Context) {
	var request reserveReq.ReserveServiceDeleteReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	res, err := s.reserveSdk.ReserveServiceCheckDelete(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, res)
	return
}

func (s ReserveHandle) ServiceInfo(c *gin.Context) {
	var request reserveReq.ReserveServiceDeleteReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	res, err := s.reserveSdk.ReserveServiceServiceInfo(ctx.Wrap(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, res)
	return
}

// ReserveOrderListExportNew 下载中心顾客受邀单列表导出
func (s ReserveHandle) ReserveOrderListExportNew(c *gin.Context) {
	var request reserveReq.ReserveOrderListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId
	// 获取企业配置
	enterpriseConfigName := helper.GetEnterpriseConfig(c, enterpriseId)
	sourceTitle := make([]string, 0)
	sourceTitleAlias := make(map[string]string, 0)
	sourceTitleValue := make([]string, 0)
	if enterpriseConfigName == enterpriseNameCoach {
		sourceTitle = []string{
			"reserve_number", "user_name", "member_name", "store_name", "order_origin_str", "invite_num", "status_str", "is_open_str",
			"is_commit_str", "customer_status_str", "is_pay", "is_store_str", "order_price", "relation_order_num", "reserve_order_time_str", "create_time_str",
		}
		sourceTitleAlias = map[string]string{
			"reserve_number": "预约单编号", "user_name": "顾客昵称", "member_name": "服务导购", "store_name": "预约门店", "order_origin_str": "邀请来源", "invite_num": "邀约次数",
			"status_str": "预约状态", "is_open_str": "是否打开", "is_commit_str": "是否确认", "customer_status_str": "顾客类型(1新客,0老客)", "is_pay": "是否成单(1是，0否)", "is_store_str": "是否到店购买(1是,0否)",
			"order_price": "导购成单金额", "relation_order_num": "关联订单数", "reserve_order_time_str": "成单时间", "create_time_str": "创建时间",
		}
		sourceTitleValue = []string{
			"order_price", "is_pay", "invite_num", "relation_order_num",
		}
	} else {
		sourceTitle = []string{
			"reserve_number", "user_name", "member_name", "store_name", "order_origin_str", "invite_num", "create_time_str",
		}
		sourceTitleAlias = map[string]string{
			"reserve_number": "预约单编号", "user_name": "顾客昵称", "member_name": "服务导购", "store_name": "预约门店", "order_origin_str": "邀请来源", "invite_num": "邀约次数", "is_commit_str": "是否确认", "create_time_str": "创建时间",
		}
	}
	_, err := downloadcenter.DC.CreateExport(ctx.Wrap(c), downloadcenter.CreateExportRequest{
		FileName:           "受邀统计--受邀列表导出",
		Type:               downloadcenter.FileTypeExcel,
		StoreId:            1,
		EnterpriseId:       enterpriseId,
		OperationId:        adminId,
		OperationName:      "adminId",
		SourceType:         downloadcenter.SourceTypeRestful,
		SourceTitleAlias:   sourceTitleAlias,
		SourceTitle:        sourceTitle,
		SourceTitleValue:   sourceTitleValue,
		SourceHeaderPrefix: nil,
		Label:              downloadcenter.ReserveOrderExport,
		CreateData:         downloadcenter.CreateData{},
		CreateSql:          downloadcenter.CreateSql{},
		CreateRestful: downloadcenter.CreateRestful{
			SourceRestfulPath:             microscrm.GetUrl("%s/micro-scrm/v1/reserve/order/reserve-order-list-new"),
			SourceRestfulMethod:           "POST",
			SourceRestfulHeaders:          map[string]interface{}{"enterprise-id": strconv.FormatInt(enterpriseId, 10)},
			SourceRestfulParams:           helper.StructToMap(request),
			SourceSuffix:                  nil,
			SourceRestfulCountField:       "data.total",
			SourceRestfulDataField:        "data.list",
			SourceRestfulRequestPageField: "page",
			SourceRestfulRequestSizeField: "size",
		},
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, "下载成功")
}

func (s ReserveHandle) ExportOrderRecord(c *gin.Context) {
	var request reserveReq.ReserveOrderListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId
	// 获取企业配置
	enterpriseConfigName := helper.GetEnterpriseConfig(c, enterpriseId)
	sourceTitle := make([]string, 0)
	sourceTitleAlias := make(map[string]string, 0)
	sourceTitleValue := make([]string, 0)
	if enterpriseConfigName == enterpriseNameCoach { // 只有coach有
		sourceTitle = []string{
			"reserve_number", "user_name", "member_name", "store_name", "order_origin_str", "invite_num", "status_str", "is_open_str",
			"is_commit_str", "customer_status_str", "is_pay", "is_store_str", "relation_order_id", "goods_name", "style", "order_price", "pay_time", "store_code",
		}
		sourceTitleAlias = map[string]string{
			"reserve_number": "编号", "user_name": "顾客昵称", "member_name": "邀约导购", "store_name": "邀约门店", "order_origin_str": "邀约来源", "invite_num": "邀约次数",
			"status_str": "邀约状态", "is_open_str": "是否打开", "is_commit_str": "是否确认", "customer_status_str": "顾客类型(1新客,0老客)", "is_pay": "是否成单(1是，0否)", "is_store_str": "是否到店购买(1是,0否)",
			"relation_order_id": "订单单号", "goods_name": "商品名称", "style": "款号", "order_price": "订单金额", "pay_time": "购买日期", "store_code": "购买门店",
		}
		sourceTitleValue = []string{
			"is_pay", "invite_num", "order_price",
		}
	}
	_, err := downloadcenter.DC.CreateExport(ctx.Wrap(c), downloadcenter.CreateExportRequest{
		FileName:           "订单明细数据导出",
		Type:               downloadcenter.FileTypeExcel,
		StoreId:            1,
		EnterpriseId:       enterpriseId,
		OperationId:        adminId,
		OperationName:      "adminId",
		SourceType:         downloadcenter.SourceTypeRestful,
		SourceTitleAlias:   sourceTitleAlias,
		SourceTitle:        sourceTitle,
		SourceTitleValue:   sourceTitleValue,
		SourceHeaderPrefix: nil,
		Label:              downloadcenter.ReserveOrderExport,
		CreateData:         downloadcenter.CreateData{},
		CreateSql:          downloadcenter.CreateSql{},
		CreateRestful: downloadcenter.CreateRestful{
			SourceRestfulPath:             microscrm.GetUrl("%s/micro-scrm/v1/reserve/order/export-order-record"),
			SourceRestfulMethod:           "POST",
			SourceRestfulHeaders:          map[string]interface{}{"enterprise-id": strconv.FormatInt(enterpriseId, 10)},
			SourceRestfulParams:           helper.StructToMap(request),
			SourceSuffix:                  nil,
			SourceRestfulCountField:       "data.total",
			SourceRestfulDataField:        "data.list",
			SourceRestfulRequestPageField: "page",
			SourceRestfulRequestSizeField: "size",
		},
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, "下载成功")
}
