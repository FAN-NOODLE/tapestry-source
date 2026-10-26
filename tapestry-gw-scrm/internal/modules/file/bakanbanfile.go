package file

import (
	"bytes"
	"io"
	"path"
	"strconv"
	"strings"

	"errors"

	"github.com/tealeg/xlsx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-scrm/internal/pkg/helper"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

func (ba *BaKanBandata) OperationBaKanBanFile(ctx ctx.Context) error {
	filename := ba.File.Filename
	if path.Ext(filename) != ".xlsx" {
		return errors.New("文件格式错误")
	}
	if !strings.Contains(filename, "BA业绩看板") {
		return errors.New("文件名称不匹配")
	}
	fh, err := ba.File.Open()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, fh)
	xlsxfile, err := xlsx.OpenBinary(buf.Bytes())
	if err != nil {
		return err
	}
	BAchan := []BaMonthlyKanban{}
	for index, sheet := range xlsxfile.Sheets {
		if index == 1 { // 门店
			bas := Handsheet1(sheet, ba.Year)
			BAchan = append(BAchan, bas...)
		} else if index == 2 { // 分区
			bas := Handsheet2(sheet, ba.Year)
			BAchan = append(BAchan, bas...)
		} else if index == 3 { // 大区
			bas := Handsheet3(sheet, ba.Year)
			BAchan = append(BAchan, bas...)
		} else if index == 0 { // ba
			bas := Handsheet0(sheet, ba.Year)
			BAchan = append(BAchan, bas...)
		}
	}
	requstDate := BaMonthlyKanbanUpload{
		EnterpriseId: ba.EnterpriseId,
		Data:         []BaMonthlyKanban{},
	}
	for k, value := range BAchan {
		if k%300 == 0 && k != 0 {
			requstDate.Data = append(requstDate.Data, value)
			microscrm.ApiBaKanBan{}.BaKanBanDataUpload(ctx, requstDate)

			requstDate.Data = []BaMonthlyKanban{}
		} else {
			requstDate.Data = append(requstDate.Data, value)
		}
	}
	if len(requstDate.Data) > 0 {
		microscrm.ApiBaKanBan{}.BaKanBanDataUpload(ctx, requstDate)

	}
	log.Info(ctx).Msgf("上传企业%v的BA看板数据一共是:%v", ba.EnterpriseId, len(BAchan))
	return nil
}
func Handsheet0(sheet *xlsx.Sheet, year int) (bas []BaMonthlyKanban) {
	for rowindex, row := range sheet.Rows {
		if rowindex > 2 {
			baMonthlyKanban := BaMonthlyKanban{
				SalesType: 1,
			}

			for cellIndex, cell := range row.Cells {
				switch cellIndex {
				case 1: // 月份
					baMonthlyKanban.Month = cell.Value
				case 2: //柜台
					baMonthlyKanban.CounterName = cell.Value
				case 3: //导购名称
					baMonthlyKanban.BaName = cell.Value
				case 4: //marscode
					baMonthlyKanban.BaMarscode = cell.Value
				case 5: // 电话
					baMonthlyKanban.Mobile = cell.Value
				case 6: // 销售额
					baMonthlyKanban.Sales = helper.StringTofloat64(cell.Value)
				case 7: // 销售目标
					baMonthlyKanban.MonthlySales = helper.StringTofloat64(cell.Value)
				case 8: // 完成率
					baMonthlyKanban.FinishRate = helper.GetRate(cell.Value)
				case 9: // 同比增长
					baMonthlyKanban.GrowthRate = helper.GetRate(cell.Value)
				case 10: //差额（待完成金额）
					baMonthlyKanban.DifferenceSales = helper.StringTofloat64(cell.Value)
				case 11: //订单数
					baMonthlyKanban.OrderNum = helper.StringToInt(cell.Value)
				case 12: // 客单价
					baMonthlyKanban.OrderNumPrice = helper.StringTofloat64(cell.Value)
				case 13: // 客人数
					baMonthlyKanban.CustomerNum = helper.StringToInt(cell.Value)
				case 14: //新客人数
					baMonthlyKanban.NewCustomerNum = helper.StringToInt(cell.Value)
				case 15: // 新客销售额
					baMonthlyKanban.NewCustomerSales = helper.StringTofloat64(cell.Value)
				case 16: // 新客销售额占总销售额
					baMonthlyKanban.NewCustomerRate = helper.GetRate(cell.Value)
				case 17: // 新二购人数
					baMonthlyKanban.NewSecondCustomerNum = helper.StringToInt(cell.Value)
				case 18: // 新二购销售额
					baMonthlyKanban.NewSecondCustomerSales = helper.StringTofloat64(cell.Value)
				case 19: //新二购销售额占总销售额
					baMonthlyKanban.NewSecondCustomerSalesRate = helper.GetRate(cell.Value)
				case 20: //老客人数
					baMonthlyKanban.OldCustomerNum = helper.StringToInt(cell.Value)
				case 21: //老客销售额
					baMonthlyKanban.OldCustomerSales = helper.StringTofloat64(cell.Value)
				case 22: //老客销售额占总销售额
					baMonthlyKanban.OldCustomerSalesRate = helper.GetRate(cell.Value)
				case 23: // 新二购且一码三关注的指标人数
					baMonthlyKanban.MonthlySecondBuyAndThreeScan = helper.StringToInt(cell.Value)
				case 24: // 新二购一码三关人数
					baMonthlyKanban.NewSecondAndThreeScan = helper.StringToInt(cell.Value)
				case 25: // 新二购人数
					baMonthlyKanban.NewSecondCustomerYNum = helper.StringToInt(cell.Value)
				case 26: // 新二购一码三关率
					baMonthlyKanban.NewSecondAndThreeScanRate = helper.GetRate(cell.Value)
				case 27: //新二购工资完成进度
					baMonthlyKanban.NewSecondAndThreeScanWageProgress = helper.StringTofloat64(cell.Value) * 100
				case 28: //老客且一码三关注的指标人数
					baMonthlyKanban.MonthlyOldCustomerAndThreeScan = helper.StringToInt(cell.Value)
				case 29: //老客且一码三关注的完成人数
					baMonthlyKanban.OldAndThreeScan = helper.StringToInt(cell.Value)
				case 30: //老客人数
					baMonthlyKanban.OldXCustomerNum = helper.StringToInt(cell.Value)
				case 31: // 老客一码三观率
					baMonthlyKanban.MonthlyOldCustomerAndThreeScanRate = helper.GetRate(cell.Value)
				case 32: // 老客工资完成进度
					baMonthlyKanban.MonthlyOldCustomerAndThreeWageProgress = helper.StringTofloat64(cell.Value) * 100
				case 33: // 新客人数
					baMonthlyKanban.NewCustomerAndThreeScan = helper.StringToInt(cell.Value)
				case 34: // 新客扫码完成人数
					baMonthlyKanban.MonthlyNewCustomerAndThreeScan = helper.StringToInt(cell.Value)
				case 35: // 新客关注率
					baMonthlyKanban.NewCustomerAndThreeScanRate = helper.GetRate(cell.Value)
				case 36: // 游客人数
					baMonthlyKanban.VisitCustomerNum = helper.StringToInt(cell.Value)
				case 37: //游客销售额
					baMonthlyKanban.VisitCustomerSales = helper.StringTofloat64(cell.Value)
				case 38: //游客销售额占总销售额占比
					baMonthlyKanban.VisitCustomerSalesRate = helper.GetRatefloat64(cell.Value)
				}
			}
			if baMonthlyKanban.BaName == "TTL" {
				continue
			}
			bas = append(bas, baMonthlyKanban)
		}
	}
	return bas
}
func Handsheet1(sheet *xlsx.Sheet, year int) (bas []BaMonthlyKanban) {
	for rowindex, row := range sheet.Rows {
		if rowindex > 1 {
			baMonthlyKanban := BaMonthlyKanban{
				SalesType: 2,
			}

			for cellIndex, cell := range row.Cells {
				switch cellIndex {
				case 1: // 月份
					baMonthlyKanban.Month = cell.Value
				case 2: //柜台
					baMonthlyKanban.CounterName = cell.Value
				case 3: //导购名称
					baMonthlyKanban.BaName = cell.Value
				case 4: //marscode
					baMonthlyKanban.Mobile = cell.Value
				case 5: // 销售额
					baMonthlyKanban.Sales = helper.StringTofloat64(cell.Value)
				case 6: // 销售目标
					baMonthlyKanban.MonthlySales = helper.StringTofloat64(cell.Value)
				case 7: // 完成率
					baMonthlyKanban.FinishRate = helper.GetRate(cell.Value)
				case 8: // 同比增长
					baMonthlyKanban.GrowthRate = helper.GetRate(cell.Value)
				case 9: //差额（待完成金额）
					baMonthlyKanban.DifferenceSales = helper.StringTofloat64(cell.Value)
				case 10: //订单数
					baMonthlyKanban.OrderNum = helper.StringToInt(cell.Value)
				case 11: // 客单价
					baMonthlyKanban.OrderNumPrice = helper.StringTofloat64(cell.Value)
				case 12: // 客人数
					baMonthlyKanban.CustomerNum = helper.StringToInt(cell.Value)
				case 13: //新客人数
					baMonthlyKanban.NewCustomerNum = helper.StringToInt(cell.Value)
				case 14: // 新客销售额
					baMonthlyKanban.NewCustomerSales = helper.StringTofloat64(cell.Value)
				case 15: // 新客销售额占总销售额
					baMonthlyKanban.NewCustomerRate = helper.GetRate(cell.Value)
				case 16: // 新二购人数
					baMonthlyKanban.NewSecondCustomerNum = helper.StringToInt(cell.Value)
				case 17: // 新二购销售额
					baMonthlyKanban.NewSecondCustomerSales = helper.StringTofloat64(cell.Value)
				case 18: //新二购销售额占总销售额
					baMonthlyKanban.NewSecondCustomerSalesRate = helper.GetRate(cell.Value)
				case 19: //老客人数
					baMonthlyKanban.OldCustomerNum = helper.StringToInt(cell.Value)
				case 20: //老客销售额
					baMonthlyKanban.OldCustomerSales = helper.StringTofloat64(cell.Value)
				case 21: //老客销售额占总销售额
					baMonthlyKanban.OldCustomerSalesRate = helper.GetRate(cell.Value)
				case 22: // 新二购且一码三关注的指标人数
					baMonthlyKanban.MonthlySecondBuyAndThreeScan = helper.StringToInt(cell.Value)
				case 23: // 新二购一码三关人数
					baMonthlyKanban.NewSecondAndThreeScan = helper.StringToInt(cell.Value)
				case 24: // 新二购人数
					baMonthlyKanban.NewSecondCustomerYNum = helper.StringToInt(cell.Value)
				case 25: // 新二购一码三关率
					baMonthlyKanban.NewSecondAndThreeScanRate = helper.GetRate(cell.Value)
				case 26: //新二购工资完成进度
					baMonthlyKanban.NewSecondAndThreeScanWageProgress = helper.StringTofloat64(cell.Value) * 100
				case 27: //老客且一码三关注的指标人数
					baMonthlyKanban.MonthlyOldCustomerAndThreeScan = helper.StringToInt(cell.Value)
				case 28: //老客且一码三关注的完成人数
					baMonthlyKanban.OldAndThreeScan = helper.StringToInt(cell.Value)
				case 29: //老客人数
					baMonthlyKanban.OldXCustomerNum = helper.StringToInt(cell.Value)
				case 30: // 老客一码三观率
					baMonthlyKanban.MonthlyOldCustomerAndThreeScanRate = helper.GetRate(cell.Value)
				case 31: // 老客工资完成进度
					baMonthlyKanban.MonthlyOldCustomerAndThreeWageProgress = helper.StringTofloat64(cell.Value) * 100
				case 32: // 新客人数
					baMonthlyKanban.NewCustomerAndThreeScan = helper.StringToInt(cell.Value)
				case 33: // 新客扫码完成人数
					baMonthlyKanban.MonthlyNewCustomerAndThreeScan = helper.StringToInt(cell.Value)
				case 34: // 新客关注率
					baMonthlyKanban.NewCustomerAndThreeScanRate = helper.GetRate(cell.Value)
				case 35: // 游客人数
					baMonthlyKanban.VisitCustomerNum = helper.StringToInt(cell.Value)
				case 36: //游客销售额
					baMonthlyKanban.VisitCustomerSales = helper.StringTofloat64(cell.Value)
				case 37: //游客销售额占总销售额占比
					baMonthlyKanban.VisitCustomerSalesRate = helper.GetRatefloat64(cell.Value)
				}
			}
			bas = append(bas, baMonthlyKanban)
		}
	}
	return bas
}
func Handsheet2(sheet *xlsx.Sheet, year int) (bas []BaMonthlyKanban) {
	for rowindex, row := range sheet.Rows {
		if rowindex > 1 {
			baMonthlyKanban := BaMonthlyKanban{
				SalesType: 3,
			}

			for cellIndex, cell := range row.Cells {
				switch cellIndex {
				case 1: // 月份
					baMonthlyKanban.Month = cell.Value
				case 2: //导购名称
					baMonthlyKanban.BaName = cell.Value
				case 3: //marscode
					baMonthlyKanban.Mobile = cell.Value
				case 4: // 销售额
					baMonthlyKanban.Sales = helper.StringTofloat64(cell.Value)
				case 5: // 销售目标
					baMonthlyKanban.MonthlySales = helper.StringTofloat64(cell.Value)
				case 6: // 完成率
					baMonthlyKanban.FinishRate = helper.GetRate(cell.Value)
				case 7: // 同比增长
					baMonthlyKanban.GrowthRate = helper.GetRate(cell.Value)
				case 8: //差额（待完成金额）
					baMonthlyKanban.DifferenceSales = helper.StringTofloat64(cell.Value)
				case 9: //订单数
					baMonthlyKanban.OrderNum = helper.StringToInt(cell.Value)
				case 10: // 客单价
					baMonthlyKanban.OrderNumPrice = helper.StringTofloat64(cell.Value)
				case 11: // 客人数
					baMonthlyKanban.CustomerNum = helper.StringToInt(cell.Value)
				case 12: //新客人数
					baMonthlyKanban.NewCustomerNum = helper.StringToInt(cell.Value)
				case 13: // 新客销售额
					baMonthlyKanban.NewCustomerSales = helper.StringTofloat64(cell.Value)
				case 14: // 新客销售额占总销售额
					baMonthlyKanban.NewCustomerRate = helper.GetRate(cell.Value)
				case 15: // 新二购人数
					baMonthlyKanban.NewSecondCustomerNum = helper.StringToInt(cell.Value)
				case 16: // 新二购销售额
					baMonthlyKanban.NewSecondCustomerSales = helper.StringTofloat64(cell.Value)
				case 17: //新二购销售额占总销售额
					baMonthlyKanban.NewSecondCustomerSalesRate = helper.GetRate(cell.Value)
				case 18: //老客人数
					baMonthlyKanban.OldCustomerNum = helper.StringToInt(cell.Value)
				case 19: //老客销售额
					baMonthlyKanban.OldCustomerSales = helper.StringTofloat64(cell.Value)
				case 20: //老客销售额占总销售额
					baMonthlyKanban.OldCustomerSalesRate = helper.GetRate(cell.Value)
				case 21: // 新二购且一码三关注的指标人数
					baMonthlyKanban.MonthlySecondBuyAndThreeScan = helper.StringToInt(cell.Value)
				case 22: // 新二购一码三关人数
					baMonthlyKanban.NewSecondAndThreeScan = helper.StringToInt(cell.Value)
				case 23: // 新二购人数
					baMonthlyKanban.NewSecondCustomerYNum = helper.StringToInt(cell.Value)
				case 24: // 新二购一码三关率
					baMonthlyKanban.NewSecondAndThreeScanRate = helper.GetRate(cell.Value)
				case 25: //新二购工资完成进度
					baMonthlyKanban.NewSecondAndThreeScanWageProgress = helper.StringTofloat64(cell.Value) * 100
				case 26: //老客且一码三关注的指标人数
					baMonthlyKanban.MonthlyOldCustomerAndThreeScan = helper.StringToInt(cell.Value)
				case 27: //老客且一码三关注的完成人数
					baMonthlyKanban.OldAndThreeScan = helper.StringToInt(cell.Value)
				case 28: //老客人数
					baMonthlyKanban.OldXCustomerNum = helper.StringToInt(cell.Value)
				case 29: // 老客一码三观率
					baMonthlyKanban.MonthlyOldCustomerAndThreeScanRate = helper.GetRate(cell.Value)
				case 30: // 老客工资完成进度
					baMonthlyKanban.MonthlyOldCustomerAndThreeWageProgress = helper.StringTofloat64(cell.Value) * 100
				case 31: // 新客人数
					baMonthlyKanban.NewCustomerAndThreeScan = helper.StringToInt(cell.Value)
				case 32: // 新客扫码完成人数
					baMonthlyKanban.MonthlyNewCustomerAndThreeScan = helper.StringToInt(cell.Value)
				case 33: // 新客关注率
					baMonthlyKanban.NewCustomerAndThreeScanRate = helper.GetRate(cell.Value)
				case 34: // 游客人数
					baMonthlyKanban.VisitCustomerNum = helper.StringToInt(cell.Value)
				case 35: //游客销售额
					baMonthlyKanban.VisitCustomerSales = helper.StringTofloat64(cell.Value)
				case 36: //游客销售额占总销售额占比
					baMonthlyKanban.VisitCustomerSalesRate = helper.GetRatefloat64(cell.Value)
				}
			}
			bas = append(bas, baMonthlyKanban)
		}
	}
	return bas
}
func Handsheet3(sheet *xlsx.Sheet, year int) (bas []BaMonthlyKanban) {
	for rowindex, row := range sheet.Rows {
		if rowindex > 1 {
			baMonthlyKanban := BaMonthlyKanban{
				SalesType: 4,
			}

			for cellIndex, cell := range row.Cells {
				switch cellIndex {
				case 1: // 月份
					baMonthlyKanban.Month = cell.Value
				case 2: //导购名称
					baMonthlyKanban.BaName = cell.Value
				case 3: //marscode
					baMonthlyKanban.Mobile = cell.Value
				case 4: // 销售额
					baMonthlyKanban.Sales = helper.StringTofloat64(cell.Value)
				case 5: // 销售目标
					baMonthlyKanban.MonthlySales = helper.StringTofloat64(cell.Value)
				case 6: // 完成率
					baMonthlyKanban.FinishRate = helper.GetRate(cell.Value)
				case 7: // 同比增长
					baMonthlyKanban.GrowthRate = helper.GetRate(cell.Value)
				case 8: //差额（待完成金额）
					baMonthlyKanban.DifferenceSales = helper.StringTofloat64(cell.Value)
				case 9: //订单数
					baMonthlyKanban.OrderNum = helper.StringToInt(cell.Value)
				case 10: // 客单价
					baMonthlyKanban.OrderNumPrice = helper.StringTofloat64(cell.Value)
				case 11: // 客人数
					baMonthlyKanban.CustomerNum = helper.StringToInt(cell.Value)
				case 12: //新客人数
					baMonthlyKanban.NewCustomerNum = helper.StringToInt(cell.Value)
				case 13: // 新客销售额
					baMonthlyKanban.NewCustomerSales = helper.StringTofloat64(cell.Value)
				case 14: // 新客销售额占总销售额
					baMonthlyKanban.NewCustomerRate = helper.GetRate(cell.Value)
				case 15: // 新二购人数
					baMonthlyKanban.NewSecondCustomerNum = helper.StringToInt(cell.Value)
				case 16: // 新二购销售额
					baMonthlyKanban.NewSecondCustomerSales = helper.StringTofloat64(cell.Value)
				case 17: //新二购销售额占总销售额
					baMonthlyKanban.NewSecondCustomerSalesRate = helper.GetRate(cell.Value)
				case 18: //老客人数
					baMonthlyKanban.OldCustomerNum = helper.StringToInt(cell.Value)
				case 19: //老客销售额
					baMonthlyKanban.OldCustomerSales = helper.StringTofloat64(cell.Value)
				case 20: //老客销售额占总销售额
					baMonthlyKanban.OldCustomerSalesRate = helper.GetRate(cell.Value)
				case 21: // 新二购且一码三关注的指标人数
					baMonthlyKanban.MonthlySecondBuyAndThreeScan = helper.StringToInt(cell.Value)
				case 22: // 新二购一码三关人数
					baMonthlyKanban.NewSecondAndThreeScan = helper.StringToInt(cell.Value)
				case 23: // 新二购人数
					baMonthlyKanban.NewSecondCustomerYNum = helper.StringToInt(cell.Value)
				case 24: // 新二购一码三关率
					baMonthlyKanban.NewSecondAndThreeScanRate = helper.GetRate(cell.Value)
				case 25: //新二购工资完成进度
					baMonthlyKanban.NewSecondAndThreeScanWageProgress = helper.StringTofloat64(cell.Value) * 100
				case 26: //老客且一码三关注的指标人数
					baMonthlyKanban.MonthlyOldCustomerAndThreeScan = helper.StringToInt(cell.Value)
				case 27: //老客且一码三关注的完成人数
					baMonthlyKanban.OldAndThreeScan = helper.StringToInt(cell.Value)
				case 28: //老客人数
					baMonthlyKanban.OldXCustomerNum = helper.StringToInt(cell.Value)
				case 29: // 老客一码三观率
					baMonthlyKanban.MonthlyOldCustomerAndThreeScanRate = helper.GetRate(cell.Value)
				case 30: // 老客工资完成进度
					baMonthlyKanban.MonthlyOldCustomerAndThreeWageProgress = helper.StringTofloat64(cell.Value) * 100
				case 31: // 新客人数
					baMonthlyKanban.NewCustomerAndThreeScan = helper.StringToInt(cell.Value)
				case 32: // 新客扫码完成人数
					baMonthlyKanban.MonthlyNewCustomerAndThreeScan = helper.StringToInt(cell.Value)
				case 33: // 新客关注率
					baMonthlyKanban.NewCustomerAndThreeScanRate = helper.GetRate(cell.Value)
				case 34: // 游客人数
					baMonthlyKanban.VisitCustomerNum = helper.StringToInt(cell.Value)
				case 35: //游客销售额
					baMonthlyKanban.VisitCustomerSales = helper.StringTofloat64(cell.Value)
				case 36: //游客销售额占总销售额占比
					baMonthlyKanban.VisitCustomerSalesRate = helper.GetRatefloat64(cell.Value)
				}
			}
			bas = append(bas, baMonthlyKanban)
		}
	}
	return bas
}
func YearMouth(yearInt int, params string) string {
	mouth := helper.StringToInt(params)
	year := strconv.Itoa(yearInt)
	if mouth < 10 && mouth > 0 {
		return year + "0" + strconv.Itoa(mouth)
	} else if mouth <= 12 && mouth >= 10 {
		return year + strconv.Itoa(mouth)
	} else {
		return params
	}
}
