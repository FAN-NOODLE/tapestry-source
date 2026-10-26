package plan_tasks

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io"
	"io/ioutil"
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	planReq "micro-scrm/internal/app/http/request/plan"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
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
)

type TaskOrderForAcxiom struct {
	Name       string
	sftpClient *sftp.SFTP
	remotePath string
}

func NewTaskOrderForAcxiom() TaskOrderForAcxiom {
	return TaskOrderForAcxiom{
		Name:       "匹配acxiom的任务订单销售数据",
		remotePath: "/out/task",
	}
}

func (job TaskOrderForAcxiom) Run(c ctx.Context, handler *xjob.JobHandler) error {

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

func (job TaskOrderForAcxiom) handle(c ctx.Context, enterpriseId int64, enterpriseName, shuffle string) error {
	// todo 获取sftp的配置信息，拉取sftp的文件内容
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
		log.Info(c).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.Name, config.EnterpriseID, clientErr)
		// return
	} else {
		job.sftpClient = client
	}
	// 获取任务匹配订单文件
	taskOrderName := "plan_task_match_" + shuffle + ".csv"
	taskOrderProductName := "plan_task_match_product_" + shuffle + ".csv"
	// 获取任务匹配商品文件
	orderFile, err := job.sftpClient.Client.Open(job.remotePath + "/" + taskOrderName)
	if err != nil {
		log.Info(c).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", job.Name, config.EnterpriseID, err)
		planService.NewTaskNotice().SftpPullNotice(c, planService.SFTPMessageData{
			EnterpriseName: enterpriseConfig.Name,
			SFTPName:       taskOrderName,
			UploadRes:      "未获取到文件",
			ToAll:          true,
		})
		return err
	}
	orderProductFile, err := job.sftpClient.Client.Open(job.remotePath + "/" + taskOrderProductName)
	if err != nil {
		log.Info(c).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", job.Name, config.EnterpriseID, err)
		planService.NewTaskNotice().SftpPullNotice(c, planService.SFTPMessageData{
			EnterpriseName: enterpriseConfig.Name,
			SFTPName:       taskOrderProductName,
			UploadRes:      "未获取到文件",
			ToAll:          true,
		})
		return err
	}
	localFileName := "task_order_" + strconv.FormatInt(app.UID(), 10) + `.csv`
	localProductFileName := "task_order_product_" + strconv.FormatInt(app.UID(), 10) + `.csv`
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
	dstProductFile, err := os.Create(localProductFilePath)
	if err != nil {
		log.Info(c).Msgf("%s::创建文件失败: %v", localProductFilePath, err)
		return err
	}
	if _, err := orderProductFile.WriteTo(dstProductFile); err != nil {
		log.Info(c).Msgf("%s::写入文件失败: %v", localProductFilePath, err)
		return err
	}
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
	// execNum := 0
	modelDao := microScrmDao.NewTaskConditionDataSaleDao()
	nowTime := time.Now().Unix()
	var updatedAt int64
	for v := range fileChan {
		var taskMatchOrder planReq.TaskMatchOrder
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) >= 12 {
			taskMatchOrder.TaskId = splitInfo[0]
			taskMatchOrder.PrimaryTaskId, _ = strconv.ParseInt(splitInfo[1], 10, 64)
			taskMatchOrder.SubTaskId, _ = strconv.ParseInt(splitInfo[2], 10, 64)
			taskMatchOrder.SaCode = splitInfo[3]
			taskMatchOrder.UnionId = splitInfo[4]
			taskMatchOrder.CounterCode = splitInfo[5]
			taskMatchOrder.PurchaseCode = splitInfo[6]
			taskMatchOrder.OriPurchaseCode = splitInfo[7]
			taskMatchOrder.PurchaseAmount, _ = strconv.ParseFloat(splitInfo[8], 10)
			taskMatchOrder.PurchaseDate = splitInfo[9]
			taskMatchOrder.CreateDate = splitInfo[10]
			taskMatchOrder.MatchStatus, _ = strconv.ParseInt(splitInfo[11], 10, 32)
			conditionIds := strings.Split(taskMatchOrder.TaskId, "_")
			purchasePrice := decimal.NewFromFloat(taskMatchOrder.PurchaseAmount).Mul(decimal.NewFromInt(100)).String()
			orderPrice, _ := strconv.Atoi(purchasePrice)
			updatedAt = helper.GetDateMorningTimestamp(taskMatchOrder.CreateDate)
			payTime := helper.RuleThree(taskMatchOrder.MatchStatus < 3, helper.GetDateMorningTimestamp(taskMatchOrder.PurchaseDate), 0)
			createdAt := nowTime
			conditionID, _ := strconv.Atoi(conditionIds[1])
			info, err := modelDao.GetSalesInfo(c, enterpriseId, map[string]interface{}{
				"condition_id": conditionID,
				"task_id":      taskMatchOrder.SubTaskId,
				"order_no":     taskMatchOrder.PurchaseCode,
			})
			if err != nil {
				fmt.Printf("查询task_condition_sales表失败%v\n", err)
				continue
			}
			dataModel := microScrmModel.TaskConditionDataSaleModel{
				PlanID:        taskMatchOrder.PrimaryTaskId,
				ConditionID:   conditionID,
				TaskID:        taskMatchOrder.SubTaskId,
				UnionID:       taskMatchOrder.UnionId,
				UserID:        taskMatchOrder.SaCode,
				StoreCode:     taskMatchOrder.CounterCode,
				OrderNo:       taskMatchOrder.PurchaseCode,
				OriginOrderNo: taskMatchOrder.OriPurchaseCode,
				OrderPrice:    orderPrice,
				OrderPayTime:  int(payTime),
				OrderStatus:   int32(taskMatchOrder.MatchStatus),
			}
			dataModel.UpdatedAt = updatedAt
			if info.ID == 0 {
				dataModel.ID = app.UID()
				dataModel.CreatedAt = createdAt
				_ = modelDao.CreateData(c, enterpriseId, dataModel)
			} else {
				_ = modelDao.UpdateSaleData(c, enterpriseId, info.ID, dataModel)
			}
		}
	}
	planService.NewTaskNotice().SftpPullNotice(c, planService.SFTPMessageData{
		EnterpriseName: enterpriseConfig.Name,
		SFTPName:       taskOrderName,
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
	for v := range productFileChan {
		var taskMatchOrderProduct planReq.TaskMatchOrderProduct
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) >= 7 {
			taskMatchOrderProduct.PurchaseCode = splitInfo[0]
			taskMatchOrderProduct.SkuName = splitInfo[1]
			taskMatchOrderProduct.Style = strings.ToLower(splitInfo[2])
			taskMatchOrderProduct.Color = strings.ToLower(splitInfo[3])
			taskMatchOrderProduct.Amount, _ = strconv.ParseFloat(splitInfo[4], 10)
			taskMatchOrderProduct.PurchaseDate = splitInfo[5]
			taskMatchOrderProduct.CreateDate = splitInfo[6]
			amount := decimal.NewFromFloat(taskMatchOrderProduct.Amount).Mul(decimal.NewFromInt(100)).String()
			goodsPrice, _ := strconv.Atoi(amount)
			updatedAt = helper.GetDateMorningTimestamp(taskMatchOrderProduct.CreateDate)
			payTime := helper.GetDateMorningTimestamp(taskMatchOrderProduct.PurchaseDate)
			createdAt := nowTime
			// 商品图片拼接：https://s7d2.scene7.com/is/image/Coach/f22237_svm2_a0
			goodsUrl := fmt.Sprintf("https://s7d2.scene7.com/is/image/Coach/%s_%s_a0", taskMatchOrderProduct.Style, taskMatchOrderProduct.Color)

			info, err := modelDao.GetOrderGoodsData(c, enterpriseId, map[string]interface{}{
				"order_no":   taskMatchOrderProduct.PurchaseCode,
				"goods_name": taskMatchOrderProduct.SkuName,
			})
			if err != nil {
				fmt.Printf("查询reserve_sales_goods表失败%v\n", err)
				continue
			}
			dataModel := microScrmModel.TaskConditionDataOrderGoodsModel{
				OrderNo:    taskMatchOrderProduct.PurchaseCode,
				GoodsName:  taskMatchOrderProduct.SkuName,
				Style:      taskMatchOrderProduct.Style,
				Color:      taskMatchOrderProduct.Color,
				GoodsPrice: goodsPrice,
				PayTime:    int(payTime),
				GoodsURL:   goodsUrl,
			}
			dataModel.UpdatedAt = updatedAt
			if info.ID == 0 {
				dataModel.ID = app.UID()
				dataModel.CreatedAt = createdAt
				_ = modelDao.CreateGoodsData(c, enterpriseId, dataModel)
			} else {
				_ = modelDao.UpdateOrderGoodsData(c, enterpriseId, info.ID, dataModel)
			}
		}
	}
	planService.NewTaskNotice().SftpPullNotice(c, planService.SFTPMessageData{
		EnterpriseName: enterpriseConfig.Name,
		SFTPName:       taskOrderProductName,
		UploadRes:      "完成",
	})
	return nil
}
