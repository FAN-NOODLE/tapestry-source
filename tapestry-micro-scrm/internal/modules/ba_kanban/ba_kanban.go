package ba_kanban

import (
	"errors"
	"math"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/logic/member"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
)

type baKanBan struct {
	Name  string
	Month string
	MiniBaKanBanDataRequest
}

func NewBaKanBan(request MiniBaKanBanDataRequest) *baKanBan {
	// 计算月份
	// 暂定EndTime

	return &baKanBan{
		Name:                    "BA业绩看板",
		MiniBaKanBanDataRequest: request,
		Month:                   time.Unix(request.EndTime, 0).Format(helper.YYYYMM),
	}
}

// GetBaKanBanDataMini 小程序BA业绩看板
func (kan *baKanBan) GetBaKanBanDataMini(ctx ctx.Context) (data BaKanBanDataResponse, err error) {
	mobileList, mobileErr := kan.getBaMobileList(ctx)
	if mobileErr != nil {
		err = mobileErr
		return
	}
	if len(mobileList) == 0 {
		return
	}

	// 根据手机号统计导购的信息
	baData, dataErr := orm.Service.BaPerformance.GetBaSaleDataByMobile(kan.EnterpriseID, mobileList, kan.Month)
	if dataErr != nil {
		if !errors.Is(dataErr, gorm.ErrRecordNotFound) {
			err = dataErr
		}
	}

	lastSales, _ := orm.Service.BaPerformance.GetBATotalSales(kan.EnterpriseID, mobileList, kan.Month)

	data = kan.makeData(baData, lastSales)

	return
}

// GetBaKanBanDataConsole 管理后台业绩看板
func (kan *baKanBan) GetBaKanBanDataConsole(ctx ctx.Context) (data BaKanBanDataResponse, err error) {
	mobileList, mobileErr := kan.getBaMobileList(ctx)
	if mobileErr != nil {
		err = mobileErr
		return
	}
	if len(mobileList) == 0 {
		return
	}

	// 根据手机号统计导购的信息
	baData, dataErr := orm.Service.BaPerformance.GetBaSaleDataByMobile(kan.EnterpriseID, mobileList, kan.Month)
	if dataErr != nil {
		if !errors.Is(dataErr, gorm.ErrRecordNotFound) {
			err = dataErr
		}
		return
	}

	lastSales, _ := orm.Service.BaPerformance.GetBATotalSales(kan.EnterpriseID, mobileList, kan.Month)

	data = kan.makeData(baData, lastSales)

	return
}

// GetSalesDataDetail 销售数据明细
func (kan *baKanBan) GetSalesDataDetail(ctx ctx.Context) (result SalesDataDetailResponse, err error) {
	mobileList, mobileErr := kan.getBaMobileList(ctx)
	if mobileErr != nil {
		err = mobileErr
		return
	}
	baListData, baListErr := orm.Service.BaPerformance.GetBaSaleDataByMobileList(kan.EnterpriseID, mobileList, kan.Month, service.PaginateParam{
		PageSize: kan.PageSize,
		Page:     kan.Page,
	})
	if baListErr != nil {
		err = baListErr
		return
	}
	// 整理数据
	result.PageSize = kan.PageSize
	result.Page = kan.Page
	result.Total, _ = orm.Service.BaPerformance.BaDataCount(kan.EnterpriseID, mobileList, kan.Month)

	result.List = []SalesDataDetailList{}
	for _, item := range baListData {
		result.List = append(result.List, SalesDataDetailList{
			Name: item.Name,
			// 销售数据
			Sales:        item.Sales,        // 销售额
			TargetSales:  item.TargetSales,  // 销售目标
			FinishRate:   item.FinishRate,   // 销售完成率
			SurplusSales: item.SurplusSales, // 销售差额
			OrderNum:     item.OrderNum,     // 订单数
			OrderPrice:   item.OrderPrice,   // 客单价

			CustomerNum:          item.CustomerNum,          // 客人数量
			NewCustomerNum:       item.NewCustomerNum,       // 新客数量
			NewCustomerSales:     item.NewCustomerSales,     // 新客月度销售额
			NewCustomerSalesRate: item.NewCustomerSalesRate, // 新客销售占比

			OldCustomerNum:       item.OldCustomerNum,       // 老客数量
			OldCustomerSales:     item.OldCustomerSales,     // 老客月度销售额
			OldCustomerSalesRate: item.OldCustomerSalesRate, // 老客销售占比

			NewSecondCustomerNum:       item.NewSecondCustomerNum,       // 新二购数量
			NewSecondCustomerSales:     item.NewSecondCustomerSales,     // 新二购月度销售额
			NewSecondCustomerSalesRate: item.NewSecondCustomerSalesRate, // 新二购销售占比

			VisitCustomerNum:       item.VisitCustomerNum,       // 老客数量
			VisitCustomerSales:     item.VisitCustomerSales,     // 老客月度销售额
			VisitCustomerSalesRate: item.VisitCustomerSalesRate, // 老客销售占比
		})
	}

	return
}

// GetOperationDataDetail 运营数据明细
func (kan *baKanBan) GetOperationDataDetail(ctx ctx.Context) (result OperationDataDetailResponse, err error) {
	mobileList, mobileErr := kan.getBaMobileList(ctx)
	if mobileErr != nil {
		err = mobileErr
		return
	}
	baListData, baListErr := orm.Service.BaPerformance.GetBaSaleDataByMobileList(kan.EnterpriseID, mobileList, kan.Month, service.PaginateParam{
		PageSize: kan.PageSize,
		Page:     kan.Page,
	})
	if baListErr != nil {
		err = baListErr
		return
	}
	// 整理数据
	result.PageSize = kan.PageSize
	result.Page = kan.Page
	result.Total, _ = orm.Service.BaPerformance.BaDataCount(kan.EnterpriseID, mobileList, kan.Month)

	result.List = []OperationDataDetailList{}
	for _, item := range baListData {
		result.List = append(result.List, OperationDataDetailList{
			Name: item.Name,
			// 运营数据
			MonthlySecondBuyAndThreeScanNum: item.MonthlySecondBuyAndThreeScanNum, // 新二购指标人数
			NewSecondAndThreeScanNum:        item.NewSecondAndThreeScanNum,        // 新二购一码三关注人数
			NewSecondAndThreeScanRate:       item.NewSecondAndThreeScanRate,       // 新二购一码三关注率

			MonthlyOldCustomerAndThreeScanNum: item.MonthlyOldCustomerAndThreeScanNum, // 老客指标人数
			OldAndThreeScanNum:                item.OldAndThreeScanNum,                // 老顾客一码三关注人数
			OldAndThreeScanRate:               item.OldAndThreeScanRate,               // 老顾客一码三关注率

			MonthlyNewCustomerAndThreeScanNum: item.MonthlyNewCustomerAndThreeScanNum, // 新顾客指标人数
			NewCustomerAndThreeScanNum:        item.NewCustomerAndThreeScanNum,        // 新顾客一码三关注人数
			NewCustomerAndThreeScanRate:       item.NewCustomerAndThreeScanRate,       // 新顾客一码三关注率
		})
	}

	return
}

func (kan *baKanBan) getBaMobileList(ctx ctx.Context) (mobileList []string, err error) {
	// 查询所有的BA数据
	if kan.StoresID == 0 && kan.UserID == 0 {
		// 获取全部的BA的user_id电话号码
		memberList, listErr := orm.Service.Member.QueryAllMemberInfo(micro_scrm.MemberModel{
			MicroScrmBaseModel: model.MicroScrmBaseModel{
				EnterpriseID: kan.EnterpriseID,
			},
		})
		if listErr != nil {
			err = errors.New("没有查询到导购的信息")
			return
		}
		for _, item := range memberList {
			mobileList = append(mobileList, item.Mobile)
		}
	} else {
		// 只有门店ID时
		if kan.UserID == 0 && kan.StoresID != 0 {
			if memberList, memberErr := member.NewMemberLogic(kan.EnterpriseID).GetMemberListByDepartmentIDScrm(ctx, kan.StoresID); memberErr != nil {
				err = memberErr
				return
			} else {
				for _, item := range memberList {
					mobileList = append(mobileList, item.Member.Mobile)
				}
			}
		} else {
			// 获取某个BA的电话
			memberInfo, infoErr := orm.Service.Member.QueryMemberInfo(micro_scrm.MemberModel{
				MicroScrmBaseModel: model.MicroScrmBaseModel{
					EnterpriseID: kan.EnterpriseID,
					BaseModel: model.BaseModel{
						ID: kan.UserID,
					},
				},
			})
			if infoErr != nil {
				err = errors.New("没有查询到该导购的信息")
				return
			}
			mobileList = append(mobileList, memberInfo.Mobile)
		}
	}
	return
}

func (ken baKanBan) makeData(baData micro_scrm_enterprise.BAMonthlyKanBanResult, lastSales float32) BaKanBanDataResponse {
	return BaKanBanDataResponse{
		SaleData: SaleData{
			SaleAmount:   baData.Sales,
			FinishRate:   baData.FinishRate,
			TargetAmount: baData.TargetSales,
			SurplusSales: baData.SurplusSales,
			SurplusRate:  baData.SurplusRate,
			OrderNum:     baData.OrderNum,
			OrderPrice:   baData.OrderPrice,
		},
		GrowthData: GrowthData{
			CurrentSaleAmount: baData.Sales,
			LastSaleAmount:    lastSales,
			GrowthRate:        int8(math.Float32bits((baData.Sales/lastSales - 0.5) * 100)),
		},
		CustomerNumData: CustomerNumData{
			TotalCustomer:   baData.CustomerNum,
			TotalSaleAmount: baData.Sales,
			NewCustomerData: NewCustomerData{
				CustomerData{
					Num:        baData.NewCustomerNum,
					SaleAmount: baData.NewCustomerSales,
					Rate:       baData.NewCustomerSalesRate,
				},
			},
			OldCustomerData: OldCustomerData{
				CustomerData{
					Num:        baData.OldCustomerNum,
					SaleAmount: baData.OldCustomerSales,
					Rate:       baData.OldCustomerSalesRate,
				},
			},
			NewSecondCustomerData: NewSecondCustomerData{
				CustomerData{
					Num:        baData.NewSecondCustomerNum,
					SaleAmount: baData.NewSecondCustomerSales,
					Rate:       baData.NewSecondCustomerSalesRate,
				},
			},
			VisitCustomerData: VisitCustomerData{
				CustomerData{
					Num:        baData.VisitCustomerNum,
					SaleAmount: baData.VisitCustomerSales,
					Rate:       baData.VisitCustomerSalesRate,
				},
			},
		},
		OneQrThreeFocusData: OneQrThreeFocusData{
			NewSecondOneQrThreeFocusData: NewSecondOneQrThreeFocusData{
				OneQrThreeFocusBaseData{
					Num:        baData.NewSecondCustomerNum,
					QrNum:      baData.NewSecondAndThreeScanNum,
					Rate:       baData.NewSecondAndThreeScanRate,
					FinishRate: baData.NewSecondAndThreeScanFinishRate,
				},
			},
			OldCustomerOneQrThreeFocusData: OldCustomerOneQrThreeFocusData{
				OneQrThreeFocusBaseData{
					Num:        baData.OldCustomerNum,
					QrNum:      baData.OldAndThreeScanNum,
					Rate:       baData.OldAndThreeScanRate,
					FinishRate: baData.OldAndThreeScanFinishRate,
				},
			},
			NewCustomerOneQrThreeFocusData: NewCustomerOneQrThreeFocusData{
				OneQrThreeFocusBaseData{
					Num:        baData.NewCustomerNum,
					QrNum:      baData.NewCustomerAndThreeScanNum,
					Rate:       baData.NewCustomerAndThreeScanRate,
					FinishRate: baData.OldAndThreeScanFinishRate,
				},
			},
		},
	}
}
func BaSourceDepte(params BaSourceDeptRequst) map[string]interface{} {
	result := map[string]interface{}{
		"counter_name": "全部",
		"leader_ba":    []interface{}{},
		"ba_level":     0,
		"ba":           0,
	}
	// 根据电话号码查询BA的职位级别
	salesType, err := orm.Service.BaPerformance.BaSourceMaxDept(params.EnterpriseId, params.Mobile)
	if err != nil {
		return result
	}
	switch salesType {
	case 1:
		if bainfo, err := orm.Service.BaPerformance.BaKanbaninfo(params.EnterpriseId, 1, params.Mobile); err == nil && len(bainfo) > 0 {
			result["counter_name"] = bainfo[0].CounterName
			result["leader_ba"] = bainfo
			result["ba_level"] = 1
		}
	case 2:
		if bainfo, err := orm.Service.BaPerformance.BaKanbaninfo(params.EnterpriseId, 2, params.Mobile); err == nil && len(bainfo) > 0 {
			result["counter_name"] = bainfo[0].ManagerName
			result["leader_ba"] = bainfo
			result["ba_level"] = 2
		}
	case 3:
		if bainfo, err := orm.Service.BaPerformance.BaKanbaninfo(params.EnterpriseId, 3, params.Mobile); err == nil && len(bainfo) > 0 {
			result["counter_name"] = bainfo[0].PartitionLeaderName
			result["leader_ba"] = bainfo
			result["ba_level"] = 3
		}
	case 4:
		if bainfo, err := orm.Service.BaPerformance.BaKanbaninfo(params.EnterpriseId, 4, params.Mobile); err == nil && len(bainfo) > 0 {
			result["counter_name"] = bainfo[0].RegionalLeaderName
			result["leader_ba"] = bainfo
			result["ba_level"] = 4
		}
	}
	return result
}
func GetBaKanBanDataMiniV1(params MiNiBaKanBanDataRequest) (micro_scrm_enterprise.BaMonthlyKanbanModel, error) {
	result := micro_scrm_enterprise.BaMonthlyKanbanModel{}
	selectMobile := params.Default_mobile
	month := time.Unix(params.Month, 0).Format("200601")
	// 区分以哪个电话查询
	if params.Dept_leader_mobile != "" {
		selectMobile = params.Dept_leader_mobile
	}
	if params.Ba_mobile != "" {
		selectMobile = params.Ba_mobile
	}
	// 查询该电话的最高权限
	salesType, err := orm.Service.BaPerformance.BaSourceMaxDept(params.EnterpriseId, selectMobile)
	if err != nil {
		return result, err
	}
	if salesType == 0 { // 表示暂未查询到该角色的任何信息
		return result, errors.New("暂未查询到该手机下的BA信息")
	}
	return orm.Service.BaPerformance.FindBAkanBanDateBymobie(params.EnterpriseId, salesType, selectMobile, month)

}
func GetSalesDataDetailV1(params BaKanBanDataRequest) (micro_scrm_enterprise.BaMonthlyKanbanModel, error) {
	result := micro_scrm_enterprise.BaMonthlyKanbanModel{}
	selectMobile := params.Default_mobile
	month := time.Unix(params.Month, 0).Format("200601")
	// 区分以哪个电话查询
	if params.Dept_leader_mobile != "" {
		selectMobile = params.Dept_leader_mobile
	}
	if params.Ba_mobile != "" {
		selectMobile = params.Ba_mobile
	}
	// 查询该电话的最高权限
	baLevel, err := orm.Service.BaPerformance.BaSourceMaxDept(params.EnterpriseId, selectMobile)
	if err != nil {
		return result, err
	}
	if baLevel == 0 { // 表示暂未查询到该角色的任何信息
		return result, errors.New("暂未查询到该手机下的BA信息")
	}
	return orm.Service.BaPerformance.FindBAkanBanDateBymobie(params.EnterpriseId, baLevel, selectMobile, month)
}

func GrowthRate(params BaKanBanDataRequest) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"nowMonth": map[string]interface{}{ // 查询月的信息
			"sales":       0, // 查询月的销售额
			"growth_rate": 0, // 查询月的同比增长
		},
		"beforeyearMonth": map[string]interface{}{
			"sales":       0, // 上一年该月的销售额
			"growth_rate": 0, // 上一年该月的同比增长
		},
	}
	selectMobile := params.Default_mobile
	month := time.Unix(params.Month, 0).Format("200601")
	year, err := strconv.Atoi(time.Unix(params.Month, 0).Format("2006"))
	if err != nil {
		return result, errors.New("获取上年的数据错误")
	}
	beforeYearMonth := strconv.Itoa(year-1) + month[5:]

	// 区分以哪个电话查询
	if params.Dept_leader_mobile != "" {
		selectMobile = params.Dept_leader_mobile
	}
	if params.Ba_mobile != "" {
		selectMobile = params.Ba_mobile
	}
	// 查询该电话的最高权限
	baLevel, err := orm.Service.BaPerformance.BaSourceMaxDept(params.EnterpriseId, selectMobile)
	if err != nil {
		return result, err
	}
	if baLevel == 0 { // 表示暂未查询到该角色的任何信息
		return result, errors.New("暂未查询到该手机下的BA信息")
	}
	if nowMonthinfo, err := orm.Service.BaPerformance.FindBAkanBanDateBymobie(params.EnterpriseId, baLevel, selectMobile, month); err == nil {
		result["nowMonth"] = map[string]interface{}{
			"sales":       nowMonthinfo.Sales,
			"growth_rate": nowMonthinfo.GrowthRate,
			"month":       month,
		}
	}
	if beforeYearMonthinfo, err := orm.Service.BaPerformance.FindBAkanBanDateBymobie(params.EnterpriseId, baLevel, selectMobile, beforeYearMonth); err == nil {
		result["nowMonth"] = map[string]interface{}{
			"sales":       beforeYearMonthinfo.Sales,
			"growth_rate": beforeYearMonthinfo.GrowthRate,
			"month":       beforeYearMonth,
		}
	}
	return result, nil
}
func GetOperationDataDetailv1(params BaKanBanDataRequest) ([]micro_scrm_enterprise.BaMonthlyKanbanModel, int64) {
	result := []micro_scrm_enterprise.BaMonthlyKanbanModel{}
	selectMobile := params.Default_mobile
	month := time.Unix(params.Month, 0).Format("200601")
	// 区分以哪个电话查询
	if params.Dept_leader_mobile != "" {
		selectMobile = params.Dept_leader_mobile
	}
	if params.Ba_mobile != "" {
		selectMobile = params.Ba_mobile
	}
	// 查询该电话的最高权限
	baLevel, err := orm.Service.BaPerformance.BaSourceMaxDept(params.EnterpriseId, selectMobile)
	if err != nil {
		return result, 0
	}
	// 查询BA下的电话号码
	baMobile := []string{}
	if mobiles, err := orm.Service.BaPerformance.BaKanbanBaMobiles(params.EnterpriseId, baLevel, selectMobile); err == nil {
		baMobile = mobiles
	}

	// 查询BA的看板
	return orm.Service.BaPerformance.FindBAkanBanDateBymobies(params.EnterpriseId, month, baMobile, params.Page, params.PageSize)
}

func BaKanbandateupload(params BaKanDateUploadRequst) error {
	creatList := make([]micro_scrm_enterprise.BaMonthlyKanbanModel, 0)
	for _, value := range params.Data {
		if value.Month == "" {
			continue
		}
		if value.SalesType == 1 {
			if (value.BaName == "" || value.Mobile == "") && value.BaMarscode != "" {
				// 根据marscode查询电话号码
				membersModel, _ := orm.Service.BaPerformance.GetBaInfoBYmarscode(params.EnterpriseId, value.BaMarscode)
				if value.BaName == "" {
					value.BaName = membersModel.Name
				}
				if value.Mobile == "" {
					value.Mobile = membersModel.Mobile
				}
			}
		}
		if infos, err := orm.Service.BaPerformance.FindBaMonthlyKan(params.EnterpriseId, value.BaName, value.Mobile, value.Month, value.SalesType); err == nil && len(infos) > 0 {
			orm.Service.BaPerformance.DeleteBaMonthlyKan(params.EnterpriseId, value.BaName, value.Mobile, value.Month, value.SalesType)
		}
		creatList = append(creatList, value)
	}
	return orm.Service.BaPerformance.UpoladBathCreateBaMonthlyKanbanModel(params.EnterpriseId, creatList)
}
