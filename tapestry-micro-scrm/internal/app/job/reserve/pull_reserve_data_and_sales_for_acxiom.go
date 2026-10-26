package reserveJob

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	reserveReq "micro-scrm/internal/app/http/request/reserve"
	excitationModels "micro-scrm/internal/app/models/excitation"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	planService "micro-scrm/internal/app/service/plan"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

const (
	MATCH_SUCCESS = iota + 1
	MATCH_CANCEL
	MATCH_EXPIRE
)

type PullReserveDataAndSalesForAcxiom struct {
	Name       string
	sftpClient *sftp.SFTP
	remotePath string
}

func NewPullReserveDataAndSalesForAcxiom() PullReserveDataAndSalesForAcxiom {
	return PullReserveDataAndSalesForAcxiom{
		Name:       "匹配acxiom的邀约订单销售数据",
		remotePath: "/out/invitation/result",
	}
}

// Run spec的值为时间日期，例如读取20240425文件，则值为20240425，不传获取当天数据
func (job PullReserveDataAndSalesForAcxiom) Run(c ctx.Context, handler *xjob.JobHandler) error {

	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}
	shuffle := time.Now().Format(helper.YYYYMMDD)
	specs := handler.Runner().Specs()
	if len(specs) > 0 {
		shuffle = specs[0]
	}
	for _, v := range enterpriseList {
		enterpriseId := v.ID
		enterpriseName := v.Name
		err = job.handle(c, enterpriseId, enterpriseName, shuffle)
		if err != nil {
			return err
		}
	}
	return nil
}

func (job PullReserveDataAndSalesForAcxiom) handle(c ctx.Context, enterpriseId int64, enterpriseName string, shuffle string) error {
	enterpriseConfig, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
	if err != nil {
		return err
	}
	if enterpriseConfig.Name != "coach" {
		return nil
	}
	config, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfig(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"sftp_type":     micro_scrm.SftpTypeAcxiom,
	})
	if err != nil {
		log.Info(c).Msgf("%s::获取SFTP配置失败: %v", job.Name, err)
	}
	if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
		message := planService.SFTPMessageData{
			EnterpriseName: config.BrandName,
			SFTPName:       "acxiom的sftp连接失败",
			UploadRes:      fmt.Sprintf("失败,错误为：%s", clientErr.Error()),
		}
		planService.NewTaskNotice().SftpSendNotice(c, message)
		return clientErr
	} else {
		job.sftpClient = client
	}

	// 获取任务匹配订单文件
	reserveOrderName := "invitation_result_" + shuffle + ".csv"
	reserveOrderProductName := "invitation_result_product_" + shuffle + ".csv"
	// 获取任务匹配商品文件
	orderFile, err := job.sftpClient.Client.Open(job.remotePath + "/" + reserveOrderName)
	if err != nil {
		log.Info(c).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", job.Name, config.EnterpriseID, err)
		planService.NewTaskNotice().SftpPullNotice(c, planService.SFTPMessageData{
			EnterpriseName: enterpriseConfig.Name,
			SFTPName:       reserveOrderName,
			UploadRes:      "未获取到文件",
			ToAll:          true,
		})
		return err
	}
	orderProductFile, err := job.sftpClient.Client.Open(job.remotePath + "/" + reserveOrderProductName)
	if err != nil {
		log.Info(c).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", job.Name, config.EnterpriseID, err)
		planService.NewTaskNotice().SftpPullNotice(c, planService.SFTPMessageData{
			EnterpriseName: enterpriseConfig.Name,
			SFTPName:       reserveOrderProductName,
			UploadRes:      "未获取到文件",
			ToAll:          true,
		})
		return err
	}
	localFileName := "invitation_result_" + strconv.FormatInt(app.UID(), 10) + `.csv`
	localProductFileName := "invitation_result_product_" + strconv.FormatInt(app.UID(), 10) + `.csv`
	localFilePath := `./storage/acxiom/` + localFileName
	localProductFilePath := `./storage/acxiom/` + localProductFileName
	dstFile, err := os.Create(localFilePath)
	if err != nil {
		log.Info(c).Msgf("%s::创建文件失败: %v", localFilePath, err)
		return err
	}
	if _, err := orderFile.WriteTo(dstFile); err != nil {
		log.Info(c).Msgf("%s::写入文件失败: %v", localFilePath, err)
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Info(c).Msgf("%s::文件删除失败: %v", name, err)
		}
	}(localFilePath)
	dstProductFile, err := os.Create(localProductFilePath)
	if err != nil {
		log.Info(c).Msgf("%s::创建文件失败: %v", localProductFilePath, err)
		return err
	}
	if _, err := orderProductFile.WriteTo(dstProductFile); err != nil {
		log.Info(c).Msgf("%s::写入文件失败: %v", localProductFilePath, err)
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Info(c).Msgf("%s::文件删除失败: %v", name, err)
		}
	}(localProductFilePath)
	databyte, _ := ioutil.ReadFile(localFilePath)
	dataProductbyte, _ := ioutil.ReadFile(localProductFilePath)
	reader := csv.NewReader(bytes.NewBuffer(databyte))
	// 存在本地数据
	var fileChan = make(chan string, 2000)
	go func(reader *csv.Reader) {
		var flag = true
		for {
			row, err := reader.Read()
			if flag {
				flag = false
				continue
			}
			if err != nil && err != io.EOF {
				fmt.Println(err.Error())
				continue
			}
			if err == io.EOF {
				close(fileChan)
				break
			}
			fileChan <- strings.Join(row, ",")
		}
	}(reader)
	modelDao := microScrmDao.NewReservesOrder()
	salesDao := microScrmDao.NewReservesOrderSales()
	nowTime := time.Now().Unix()
	var updatedAt int64
	successReserveNumber := make([]string, 0)
	cancelReserveNumber := make([]string, 0)
	for v := range fileChan {
		var reserveMatchOrder reserveReq.ReserveMatchOrder
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) >= 10 {
			reserveMatchOrder.OrderNo = splitInfo[0]
			reserveMatchOrder.SaCode = splitInfo[1]
			reserveMatchOrder.UnionId = splitInfo[2]
			reserveMatchOrder.PurchaseStoreCode = splitInfo[3]
			reserveMatchOrder.PurchaseCode = splitInfo[4]
			reserveMatchOrder.OriPurchaseCode = splitInfo[5]
			reserveMatchOrder.PurchaseAmount, _ = strconv.ParseFloat(splitInfo[6], 10)
			reserveMatchOrder.PurchaseDate = splitInfo[7]
			reserveMatchOrder.Status, _ = strconv.ParseInt(splitInfo[8], 10, 32)
			reserveMatchOrder.NewCustomer, _ = strconv.ParseInt(splitInfo[9], 10, 64) // 0老客，1新客
			reserveMatchOrder.CreateDate = splitInfo[10]
			purchasePrice := decimal.NewFromFloat(reserveMatchOrder.PurchaseAmount).Mul(decimal.NewFromInt(100)).String()
			orderPrice, _ := strconv.Atoi(purchasePrice)
			updatedAt = helper.GetDateMorningTimestamp(reserveMatchOrder.CreateDate)
			payTime := helper.RuleThree(reserveMatchOrder.Status < 3, helper.GetDateMorningTimestamp(reserveMatchOrder.PurchaseDate), 0)
			createdAt := nowTime
			info, err := salesDao.GetSalesInfo(c, enterpriseId, map[string]interface{}{
				"reserve_number":    reserveMatchOrder.OrderNo,
				"relation_order_id": reserveMatchOrder.PurchaseCode,
			})
			if err != nil {
				fmt.Printf("查询reserve_sales表失败%v\n", err)
				continue
			}
			dataModel := microScrmModel.ReserveOrderSalesModel{
				ReserveNumber:   reserveMatchOrder.OrderNo,
				MemberUserID:    reserveMatchOrder.SaCode,
				UserUnionid:     reserveMatchOrder.UnionId,
				StoreCode:       reserveMatchOrder.PurchaseStoreCode,
				RelationOrderID: reserveMatchOrder.PurchaseCode,
				OriginOrderID:   reserveMatchOrder.OriPurchaseCode,
				OrderPrice:      orderPrice,
				PayTime:         int(payTime),
				OrderType:       int32(reserveMatchOrder.Status),
				NewCustomer:     int32(helper.RuleThree(reserveMatchOrder.NewCustomer == 1, 1, 2)),
			}
			dataModel.UpdatedAt = updatedAt
			if info.ID == 0 {
				dataModel.ID = app.UID()
				dataModel.CreatedAt = createdAt
				_ = salesDao.CreateSaleData(c, enterpriseId, dataModel)
			} else {
				_ = salesDao.UpdateSaleData(c, enterpriseId, info.ID, dataModel)
			}
			// 更新邀约单的数据
			updateMap := map[string]interface{}{
				"updated_at":           time.Now().Unix(),
				"reserve_order_status": reserveMatchOrder.Status,
			}
			if reserveMatchOrder.Status == MATCH_EXPIRE { // 安客成订单匹配过期，需要将邀约单改为已过期, 只处理邀约单为进行中的改为已过期，邀约完成和邀约成功的不更改邀约单状态
				updateMap["status"] = microScrmModel.OVERTIME
				updateMap["is_pay"] = microScrmModel.PAY_FAIL
				modelDao.Update(c, enterpriseId, map[string]interface{}{
					"reserve_number": reserveMatchOrder.OrderNo,
				}, updateMap)
			}
			if reserveMatchOrder.Status == MATCH_SUCCESS {
				//成功
				successReserveNumber = append(successReserveNumber, reserveMatchOrder.OrderNo)
			}
			if reserveMatchOrder.Status == MATCH_CANCEL {
				//取消匹配
				cancelReserveNumber = append(cancelReserveNumber, reserveMatchOrder.OrderNo)
			}
		}
	}
	successReserveNumber = helper.ArrayUniq(successReserveNumber)
	if len(successReserveNumber) > 0 {
		successReserveNumberList := lo.Chunk(successReserveNumber, 1000)
		for _, reserveNumberList := range successReserveNumberList {
			successList, err := salesDao.GetSalesStatusListByReserveNumber(c, enterpriseId, reserveNumberList)
			if err != nil {
				fmt.Printf("查询reserve_sales的匹配成功数据失败%v\n", err)
				continue
			}
			if len(successList) > 0 { // 邀约匹配订单，1匹配成功，2取消匹配，（订单取消匹配，金额扣减，如果金额扣减为0了，则状态回滚）
				updateMap := map[string]interface{}{
					"updated_at":           time.Now().Unix(),
					"reserve_order_status": 1,
				}
				for _, info := range successList {
					orderInfo, _ := modelDao.GetOrderInfo(c, enterpriseId, map[string]interface{}{
						"reserve_number": info.ReserveNumber,
					})
					if orderInfo.CustomerStatus == 0 { // 如果reserve_order表的customer_status状态已更新过，则不更新
						updateMap["customer_status"] = info.NewCustomer
					}
					updateMap["status"] = microScrmModel.FINISH
					updateMap["is_pay"] = microScrmModel.PAY_SUCCESS
					updateMap["reserve_order_time"] = info.LastPayTime
					updateMap["reserve_order_price"] = info.OrderPrice
					if orderInfo.Status != microScrmModel.FINISH {
						excitationSvc.NewPointsProduceService().ProduceQueue(c, excitationReq.PointsProduceSave{
							MemberId:     orderInfo.MemberID,
							Type:         excitationModels.RuleCategoryId4,
							Value:        1,
							UniqueHash:   info.ReserveNumber + fmt.Sprintf("-%d-%d", orderInfo.MemberID, excitationModels.RuleCategoryId4),
							EnterpriseId: enterpriseId,
						})
					}
					_ = modelDao.Update(c, enterpriseId, map[string]interface{}{
						"reserve_number": info.ReserveNumber,
					}, updateMap)
				}
			}
		}
	}

	cancelReserveNumber = helper.ArrayUniq(cancelReserveNumber)
	if len(cancelReserveNumber) > 0 {
		cancelReserveNumberList := lo.Chunk(cancelReserveNumber, 1000)
		for _, reserveNumberList := range cancelReserveNumberList {
			cancelList, err := salesDao.GetSalesStatusListByReserveNumber(c, enterpriseId, reserveNumberList)
			if err != nil {
				fmt.Printf("查询reserve_sales的取消匹配数据失败%v\n", err)
				continue
			}
			if len(cancelList) > 0 {
				updateMap := map[string]interface{}{
					"updated_at":           time.Now().Unix(),
					"reserve_order_status": 2,
				}
				for _, info := range cancelList {
					orderInfo, err := modelDao.GetOrderInfo(c, enterpriseId, map[string]interface{}{
						"reserve_number": info.ReserveNumber,
					})
					if err != nil {
						continue
					}
					updateMap["is_pay"] = microScrmModel.PAY_FAIL
					if orderInfo.IsConfirm == microScrmModel.CONFIRM_YES {
						updateMap["status"] = microScrmModel.WAITFINISH
					} else {
						updateMap["status"] = microScrmModel.RUNNING
					}
					updateMap["reserve_order_time"] = info.LastPayTime
					updateMap["reserve_order_price"] = 0
					_ = modelDao.Update(c, enterpriseId, map[string]interface{}{
						"reserve_number": info.ReserveNumber,
					}, updateMap)
				}

			}
		}
	}

	planService.NewTaskNotice().SftpPullNotice(c, planService.SFTPMessageData{
		EnterpriseName: enterpriseConfig.Name,
		SFTPName:       reserveOrderName,
		UploadRes:      "完成",
	})
	productReader := csv.NewReader(bytes.NewBuffer(dataProductbyte))
	// 存在本地数据
	var productFileChan = make(chan string, 2000)
	go func(reader *csv.Reader) {
		var flag = true
		for {
			row, err := reader.Read()
			if flag {
				flag = false
				continue
			}
			if err != nil && err != io.EOF {
				fmt.Println(err.Error())
				continue
			}
			if err == io.EOF {
				close(productFileChan)
				break
			}
			productFileChan <- strings.Join(row, ",")
		}
	}(productReader)
	goodsSaleMap := make(map[string]int, 0)
	for v := range productFileChan {
		var reserveMatchOrderProduct reserveReq.ReserveMatchOrderProduct
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) >= 7 {
			reserveMatchOrderProduct.PurchaseCode = splitInfo[0]
			reserveMatchOrderProduct.SkuName = splitInfo[1]
			reserveMatchOrderProduct.Style = strings.ToLower(splitInfo[2])
			reserveMatchOrderProduct.Color = strings.ToLower(splitInfo[3])
			reserveMatchOrderProduct.Amount, _ = strconv.ParseFloat(splitInfo[4], 10)
			reserveMatchOrderProduct.PurchaseDate = splitInfo[5]
			reserveMatchOrderProduct.CreateDate = splitInfo[6]
			amount := decimal.NewFromFloat(reserveMatchOrderProduct.Amount).Mul(decimal.NewFromInt(100)).String()
			goodsPrice, _ := strconv.Atoi(amount)
			updatedAt = helper.GetDateMorningTimestamp(reserveMatchOrderProduct.CreateDate)
			payTime := helper.GetDateMorningTimestamp(reserveMatchOrderProduct.PurchaseDate)
			createdAt := nowTime
			// 商品图片拼接：https://s7d2.scene7.com/is/image/Coach/f22237_svm2_a0
			goodsUrl := fmt.Sprintf("https://s7d2.scene7.com/is/image/Coach/%s_%s_a0", reserveMatchOrderProduct.Style, reserveMatchOrderProduct.Color)
			if _, ok := goodsSaleMap[reserveMatchOrderProduct.PurchaseCode]; !ok { // 记录订单号下的所有商品数,读一条记录一次；
				goodsSaleMap[reserveMatchOrderProduct.PurchaseCode] = 0
			} else {
				goodsSaleMap[reserveMatchOrderProduct.PurchaseCode]++
			}
			goodsMapNum := goodsSaleMap[reserveMatchOrderProduct.PurchaseCode]
			info, err := salesDao.GetOrderGoodsData(c, enterpriseId, map[string]interface{}{
				"relation_order_id": reserveMatchOrderProduct.PurchaseCode,
			}) // 查询订单号已记录的所有商品数
			if err != nil {
				fmt.Printf("查询reserve_sales_goods表失败%v\n", err)
				continue
			}
			goodsIndex := len(info) - 1
			if goodsMapNum > goodsIndex { // csv文件中的订单下商品数大于已存表中的数据，则新增数，否则更新数据
				dataModel := microScrmModel.ReserveOrderGoodsModel{
					RelationOrderID: reserveMatchOrderProduct.PurchaseCode,
					GoodsName:       reserveMatchOrderProduct.SkuName,
					Style:           reserveMatchOrderProduct.Style,
					Color:           reserveMatchOrderProduct.Color,
					GoodsPrice:      goodsPrice,
					PayTime:         int(payTime),
					GoodsUrl:        goodsUrl,
				}
				dataModel.ID = app.UID()
				dataModel.CreatedAt = createdAt
				dataModel.UpdatedAt = updatedAt
				_ = salesDao.CreateOrderGoodsData(c, enterpriseId, dataModel)
			} else {
				indexId := info[goodsMapNum].ID
				updateMap := map[string]interface{}{
					"relation_order_id": reserveMatchOrderProduct.PurchaseCode,
					"goods_name":        reserveMatchOrderProduct.SkuName,
					"style":             reserveMatchOrderProduct.Style,
					"color":             reserveMatchOrderProduct.Color,
					"goods_price":       goodsPrice,
					"pay_time":          int(payTime),
					"goods_url":         goodsUrl,
					"updated_at":        updatedAt,
				}
				_ = salesDao.UpdateOrderGoodsData(c, enterpriseId, indexId, updateMap)
			}
		}
	}
	planService.NewTaskNotice().SftpPullNotice(c, planService.SFTPMessageData{
		EnterpriseName: enterpriseConfig.Name,
		SFTPName:       reserveOrderProductName,
		UploadRes:      "完成",
	})
	return nil
}
