package planService

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/dao/micro_scrm"
	planReq "micro-scrm/internal/app/http/request/plan"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	"micro-scrm/internal/app/models"
	"micro-scrm/internal/app/models/micro_scrm"
	tagsSvc "micro-scrm/internal/app/service/tags"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/notice"
	"micro-scrm/internal/pkg/sdk/microwork"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shakinm/xlsReader/xls"
	"github.com/shakinm/xlsReader/xls/structure"
	"github.com/tealeg/xlsx"
)

type PlanTaskImport struct {
	planTaskImport *microScrmDao.PlanImportDao
}

func NewPlanTaskImport() *PlanTaskImport {
	return &PlanTaskImport{
		planTaskImport: microScrmDao.NewPlanImportDao(),
	}
}

func (p *PlanTaskImport) FileDetection(c app.Ctx, enterpriseID int64, FileName string) (count int64, err error) {
	return p.planTaskImport.FileDetection(c, enterpriseID, FileName)

}

func (p *PlanTaskImport) SaveImportFile(c app.Ctx, enterpriseID int64, info microScrmModel.TaskImportFileModel) (Id int64, err error) {
	return p.planTaskImport.SaveImportFile(c, enterpriseID, info)
}

// SyncTaskImportData 同步数据到plan_task_import_data表中
func (p *PlanTaskImport) SyncTaskImportData(c app.Ctx, enterpriseID, fileID int64, sendType int, taskName string, taskId int64, storeName string) {
	fileInfo, errs := p.planTaskImport.GetImportFile(c, enterpriseID, fileID)
	if errs != nil {
		return
	}
	if fileInfo.IsExtract == 1 {
		return
	}
	b, err := readTaskImportRemoteFile(c, fileInfo.FileURL)
	if err != nil {
		return
	}
	urlByte := strings.Split(fileInfo.FileURL, ".")
	ext := urlByte[len(urlByte)-1]
	var data []*planReq.UnionIDData
	switch ext {
	case "xlsx":
		zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
		if err != nil {
			return
		}
		fileData, err := xlsx.ReadZipReader(zr)
		if err != nil {
			return
		}
		data, _ = readXlsxData(fileData.Sheets, fileInfo.TaskType)
	case "xls":
		var fileName = fmt.Sprintf("%d", time.Now().Unix()) + "_taskConsumerFile.xls"
		newFile, _ := os.Create(fileName)
		if _, err0 := newFile.Write(b); err0 != nil {
			return
		}
		if err := newFile.Close(); err != nil {
			return
		}
		file, err := os.Open(fileName)
		if err != nil {
			return
		}
		wb, err := xls.OpenReader(file)
		if err != nil {
			return
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				return
			}
		}(fileName)
		data, _ = readXlsData(wb, fileInfo.TaskType)
		if err := file.Close(); err != nil {
			return
		}
	default:
		return
	}
	fileName := strconv.FormatInt(taskId, 10) + ".csv"
	// FilePath := `./storage/upload/` + fileName
	create, err := os.Create(fileName)
	if err != nil {
		return
	}
	create.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	csvFile := csv.NewWriter(create)
	var header []string
	header = []string{
		"unionid", "customer_id", "sacode", "匹配结果",
	}
	err = csvFile.Write(header)
	if err != nil {
		return
	}
	members, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseID,
		"deleted_at":    0,
	})
	if err != nil {
		return
	}
	var userIds []string
	for _, member := range members {
		userIds = append(userIds, member.Userid)
	}
	limit := 1000
	dataLen := len(data)
	if dataLen > limit {
		dataPage := int(math.Ceil(float64(dataLen / limit)))
		dataPage += 1
		for i := 1; i <= dataPage; i++ {
			start := (i - 1) * limit
			end := i * limit
			var dataList []*planReq.UnionIDData
			if end >= dataLen {
				dataList = data[start:dataLen]
			} else {
				dataList = data[start:end]
			}
			if errInsert := p.createImportData(c, enterpriseID, dataList, fileInfo, sendType, csvFile, userIds, taskId, storeName); errInsert != nil {
				continue
			}
		}
	} else {
		_ = p.createImportData(c, enterpriseID, data, fileInfo, sendType, csvFile, userIds, taskId, storeName)
	}
	csvFile.Flush()
	cs := storage.Init(cos.COS)
	cosPath, err := cs.Upload(c, fileName, "plan_task/"+fileName)
	if err != nil {
		return
	}
	err = create.Close()
	if err != nil {
		return
	}
	defer os.Remove(fileName)
	_ = p.planTaskImport.UpdateImportFileExtract(c, enterpriseID, fileInfo.ID, cs.GetURL(cosPath))
	// 修改任务的isExtract状态
	_ = microScrmDao.NewPlanTaskDao().UpdatePlanTaskISExtract(c, enterpriseID, fileInfo.ID)
}

func readTaskImportRemoteFile(c context.Context, fileURL string) ([]byte, error) {
	return helper.ReadRemoteFile(c, fileURL)
}

func readXlsxData(sheets []*xlsx.Sheet, taskType int) ([]*planReq.UnionIDData, error) {
	riMin := 1
	if taskType != microScrmModel.TaskDimensionImportUnionId {
		riMin = 2
	}
	var settings []*planReq.UnionIDData
	for index, sheet := range sheets {
		if index > 0 {
			return nil, errors.New("文件格式错误")
		}
		for ri, row := range sheet.Rows {
			if ri < riMin {
				continue
			}
			if len(row.Cells) == 0 {
				continue
			}
			var setting = &planReq.UnionIDData{}
			setting.UnionID = strings.TrimSpace(row.Cells[0].String())
			if len(row.Cells) >= 2 {
				if row.Cells[0].String() == "" && row.Cells[1].String() == "" {
					continue
				}
				setting.CustomerId = strings.TrimSpace(row.Cells[1].String())
				if taskType == microScrmModel.TaskDimensionImportSaUnionId {
					for i, cell := range row.Cells {
						if i == 2 {
							setting.SaCode = cell.String()
						}
					}
				}
			}
			settings = append(settings, setting)
		}
	}
	return settings, nil
}

func readXlsData(wb xls.Workbook, taskType int) ([]*planReq.UnionIDData, error) {
	var settings []*planReq.UnionIDData
	riMin := 1
	if taskType != microScrmModel.TaskDimensionImportUnionId {
		riMin = 2
	}
	for index, sheet := range wb.GetSheets() {
		if index > 0 {
			return nil, errors.New("文件格式错误")
		}
		for ri, row := range sheet.GetRows() {
			if ri < riMin {
				continue
			}
			if len(row.GetCols()) == 0 {
				continue
			}
			if _, err := row.GetCol(0); err != nil {
				continue
			}
			getColStr := func(c structure.CellData, err error) string {
				if err != nil {
					return ""
				}
				return c.GetString()
			}
			var setting = &planReq.UnionIDData{}
			setting.UnionID = strings.TrimSpace(getColStr(row.GetCol(0)))
			if len(row.GetCols()) >= 2 {
				if getColStr(row.GetCol(0)) == "" && getColStr(row.GetCol(1)) == "" {
					continue
				}
				setting.CustomerId = strings.TrimSpace(getColStr(row.GetCol(1)))
				if taskType == microScrmModel.TaskDimensionImportSaUnionId {
					for i := range row.GetCols() {
						if i == 2 {
							setting.SaCode = getColStr(row.GetCol(2))
						}
					}
				}
			}
			settings = append(settings, setting)
		}
	}
	return settings, nil
}

// CreateImportData 插入数据到plan_task_import_data表中
func (p *PlanTaskImport) createImportData(c app.Ctx, enterpriseID int64, data []*planReq.UnionIDData, fileInfo microScrmModel.TaskImportFileModel, sendType int, csvFile *csv.Writer, userIds []string, taskId int64, storeName string) error {
	lists := make(map[string]planReq.ImportInsertData, 0)
	var datalist []planReq.ImportInsertData
	var importData []microScrmModel.PlanTaskImportDataModel
	if fileInfo.TaskType == microScrmModel.TaskDimensionImportUnionId { // 专属导购
		var importDataLog []microScrmModel.PlanTaskImportDataLogModel
		var customerIds []string
		externalMaps := make(map[string]string)
		externalUserIdMaps := make(map[string]string)
		for _, v := range data {
			info := v
			if info.UnionID == "" && info.CustomerId != "" {
				customerIds = append(customerIds, info.CustomerId)
			}
		}
		if len(customerIds) > 0 { // 计算出unionid为空，customerid有值的数据
			whereMap := make(map[string]interface{})
			if storeName == "sw" {
				whereMap = map[string]interface{}{
					"customer_code": customerIds,
				}
			} else {
				whereMap = map[string]interface{}{
					"union_code": customerIds,
				}
			}
			externalResult, err := microScrmDao.NewWxUserCorpExternal().GetMany(c, enterpriseID, whereMap)
			if err != nil {
				return err
			}
			if len(externalResult) > 0 { // 根据customerid找到库中的unionid用户
				for _, model := range externalResult {
					if storeName == "sw" {
						if _, ok := externalMaps[model.CustomerCode]; !ok {
							externalMaps[model.CustomerCode] = model.Unionid
						}
						externalUserIdMaps[model.CustomerCode] = model.Openid
					} else {
						if _, ok := externalMaps[model.UnionCode]; !ok {
							externalMaps[model.UnionCode] = model.Unionid
						}
						externalUserIdMaps[model.UnionCode] = model.Openid
					}

				}
			}
			for k, v := range data {
				info := v
				if info.UnionID == "" && info.CustomerId != "" {
					info.UnionID = externalMaps[info.CustomerId]
				}
				data[k] = info
			}
		}
		list, err := p.planTaskImport.GetInsertImportDataInfoByUnionId(c, enterpriseID, data, sendType)
		if err != nil {
			return err
		}
		keySaRelation := make(map[string][]planReq.ImportInsertData, 0)
		for _, value := range list {
			item := value
			keySaIndex := item.UnionID
			keySaRelation[keySaIndex] = append(keySaRelation[keySaIndex], item)
		}
		for _, v := range data {
			ImportRes := "匹配失败"
			item := v
			keySaIndex := item.UnionID
			saName := ""
			userId := item.SaCode
			unionId := item.UnionID
			customerId := item.CustomerId
			relationLists := keySaRelation[keySaIndex]
			if len(relationLists) == 0 {
				logId := app.UID()
				var LogData = microScrmModel.PlanTaskImportDataLogModel{
					Name:       saName,
					Said:       userId,
					UnionId:    unionId,
					ImportRes:  ImportRes,
					FileID:     fileInfo.ID,
					BaseModel:  models.CreateBase(logId),
					CustomerID: customerId,
				}
				importDataLog = append(importDataLog, LogData)
				_ = csvFile.Write([]string{
					unionId, item.CustomerId, userId, ImportRes,
				})
			}
			for _, value := range relationLists {
				info := value
				saName = info.Name
				userId = info.UserID
				unionId = info.UnionID
				if helper.InArray[string](info.UserID, userIds) {
					logId := app.UID()
					var newData = microScrmModel.PlanTaskImportDataModel{
						Name:               info.Name,
						Said:               info.UserID,
						UnionId:            info.UnionID,
						FileID:             fileInfo.ID,
						WorkExternalUserId: info.WorkExternalUserId,
						BaseModel:          models.CreateBase(logId),
						IsMatch:            1,
						UserName:           info.UserName,
					}
					importData = append(importData, newData)
					ImportRes = "匹配成功"
					var LogData = microScrmModel.PlanTaskImportDataLogModel{
						Name:       saName,
						Said:       userId,
						UnionId:    unionId,
						ImportRes:  ImportRes,
						FileID:     fileInfo.ID,
						BaseModel:  models.CreateBase(app.UID()),
						CustomerID: customerId,
					}
					importDataLog = append(importDataLog, LogData)
				}
				_ = csvFile.Write([]string{
					unionId, item.CustomerId, userId, ImportRes,
				})
			}
		}
		_ = microScrmDao.NewPlanTaskImportDataDao().InsertImportDataBatch(c, importData, enterpriseID)

		_ = microScrmDao.NewPlanTaskImportDataLogDao().InsertImportDataLogBatch(c, importDataLog, enterpriseID)
	}
	if fileInfo.TaskType == microScrmModel.TaskDimensionImportSaUnionId {
		var importDataLog []microScrmModel.PlanTaskImportDataLogModel
		var customerIds []string
		externalMaps := make(map[string]string)
		externalUserIdMaps := make(map[string]string)
		for _, v := range data {
			info := v
			if info.UnionID == "" && info.CustomerId != "" {
				customerIds = append(customerIds, info.CustomerId)
			}
		}
		if len(customerIds) > 0 { // 计算出unionid为空，customerid有值的数据
			whereMap := make(map[string]interface{})
			if storeName == "sw" {
				whereMap = map[string]interface{}{
					"customer_code": customerIds,
				}
			} else {
				whereMap = map[string]interface{}{
					"union_code": customerIds,
				}
			}
			externalResult, err := microScrmDao.NewWxUserCorpExternal().GetMany(c, enterpriseID, whereMap)
			if err != nil {
				return err
			}
			if len(externalResult) > 0 { // 根据customerid找到库中的unionid用户
				for _, model := range externalResult {
					if storeName == "sw" {
						if _, ok := externalMaps[model.CustomerCode]; !ok {
							externalMaps[model.CustomerCode] = model.Unionid
						}
						externalUserIdMaps[model.CustomerCode] = model.Openid
					} else {
						if _, ok := externalMaps[model.UnionCode]; !ok {
							externalMaps[model.UnionCode] = model.Unionid
						}
						externalUserIdMaps[model.UnionCode] = model.Openid
					}
				}
				for k, value := range data {
					info := value
					if info.UnionID == "" && info.CustomerId != "" {
						if unionId, ok := externalMaps[info.CustomerId]; ok {
							info.UnionID = unionId
						}
					}
					data[k] = info
				}
			}
		}
		list, err := p.planTaskImport.GetInsertImportDataInfoBySaCodeAndUnionId(c, enterpriseID, data)
		if err != nil {
			return err
		}
		datalist = list
		for _, item := range datalist {
			key := item.UnionID + "_" + item.UserID
			lists[key] = item
		}

		for _, item := range data {
			ImportRes := "匹配失败"
			var name, workExternalUserId string
			said := item.SaCode
			unionId := item.UnionID
			customerId := item.CustomerId
			if said != "" { // 优化处理said为空的数据兼容
				if unionId == "" && customerId != "" { // unionid为空，customer_id不为空，并且没有在库中匹配到unionid的用户
					isMatch := 1
					customerInfo, err2 := microwork.NewCrm().GetCustomerInfoByUnionId(c, microwork.UserInfoRequest{
						EnterpriseID: enterpriseID,
						CustomerCode: customerId,
					})
					if err2 != nil {
						isMatch = 2
					}
					if customerInfo.Data.Mobile == "" {
						isMatch = 2
					}
					if !helper.InArray[string](said, userIds) {
						isMatch = 2
					}
					userName := customerInfo.Data.Name
					unionId = customerInfo.Data.UnionId
					workExternalUserId = externalUserIdMaps[customerId]
					if helper.InArray[string](said, userIds) {
						var newData = microScrmModel.PlanTaskImportDataModel{
							Name:               name,
							Said:               said,
							UnionId:            unionId,
							FileID:             fileInfo.ID,
							WorkExternalUserId: workExternalUserId,
							BaseModel:          models.CreateBase(app.UID()),
							CustomerID:         customerId,
							IsMatch:            isMatch,
							UserName:           userName,
						}
						if isMatch == 1 {
							ImportRes = "匹配成功"
						}
						importData = append(importData, newData)
						var LogData = microScrmModel.PlanTaskImportDataLogModel{
							Name:       name,
							Said:       said,
							UnionId:    unionId,
							ImportRes:  ImportRes,
							FileID:     fileInfo.ID,
							BaseModel:  models.CreateBase(app.UID()),
							CustomerID: customerId,
						}
						importDataLog = append(importDataLog, LogData)
					}
				} else {
					pickKey := unionId + "_" + said
					if _, ok := lists[pickKey]; ok {
						name = lists[pickKey].Name
						workExternalUserId = lists[pickKey].WorkExternalUserId
						if helper.InArray[string](said, userIds) {
							var newData = microScrmModel.PlanTaskImportDataModel{
								Name:               name,
								Said:               said,
								UnionId:            unionId,
								FileID:             fileInfo.ID,
								WorkExternalUserId: workExternalUserId,
								BaseModel:          models.CreateBase(app.UID()),
								CustomerID:         customerId,
								IsMatch:            1,
								UserName:           lists[pickKey].UserName,
							}
							ImportRes = "匹配成功"
							importData = append(importData, newData)
							var LogData = microScrmModel.PlanTaskImportDataLogModel{
								Name:       name,
								Said:       said,
								UnionId:    unionId,
								ImportRes:  ImportRes,
								FileID:     fileInfo.ID,
								BaseModel:  models.CreateBase(app.UID()),
								CustomerID: customerId,
							}
							importDataLog = append(importDataLog, LogData)
						}
					}
				}
			} else {
				ImportRes = "数据异常"
			}
			_ = csvFile.Write([]string{
				unionId, customerId + "\t", said, ImportRes,
			})
		}
		err = microScrmDao.NewPlanTaskImportDataDao().InsertImportDataBatch(c, importData, enterpriseID)

		err = microScrmDao.NewPlanTaskImportDataLogDao().InsertImportDataLogBatch(c, importDataLog, enterpriseID)

	}
	// 如果需要给人群包打标签，开始处理标签
	taskInfo, err := microScrmDao.NewPlanTaskDao().GetOne(c, enterpriseID, map[string]interface{}{"id": taskId})
	if err != nil {
		return err
	}
	if (taskInfo.GroupSendType == microScrmModel.PLANTASKTAGSEND && taskInfo.IsAddTag == 1) || (taskInfo.TaskCategory == microScrmModel.PLANTASKMOMENT && taskInfo.GroupSendType == microScrmModel.PLANTASKQYSEND) {
		var externalUserTags []tagsReq.ExternalUserTag
		for _, value := range importData {
			item := value
			if item.UnionId == "" {
				continue
			}
			externalUserTags = append(externalUserTags, tagsReq.ExternalUserTag{
				TagId:   taskInfo.QyTagId,
				UnionId: item.UnionId,
				UserId:  item.Said,
			})

		}
		err = tagsSvc.NewWorkTags().BatchMarking(c, "计划任务人群包", enterpriseID, externalUserTags, 1)
		if err != nil {
			app.Logger().Error(c).Msgf("导入人群数据打标签失败，错误原因:%s", err.Error())
			notice.NewWechat().SendTextMessage(c, fmt.Sprintf("计划任务人群包-批量打标异常,错误:[%+v]", err))
			return err
		}
	}
	return nil
}
