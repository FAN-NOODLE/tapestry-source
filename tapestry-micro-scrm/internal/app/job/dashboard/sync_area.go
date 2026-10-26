package dashboardJob

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/tealeg/xlsx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io/ioutil"
	"micro-scrm/internal/app"
	dashboardDao "micro-scrm/internal/app/dao/dashboard"
	reserveDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	dashboardResp "micro-scrm/internal/app/http/response/dashboard"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/models"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	SFTPROOT = `/`
)

type SyncArea struct {
	Name string
}

func NewSyncArea() SyncArea {
	return SyncArea{Name: "同步core_area文件"}
}

func (s SyncArea) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}
	for _, id := range enterpriseIds {
		enterpriseId := id
		err := s.handle(c, enterpriseId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s SyncArea) handle(c ctx.Context, enterpriseId int64) error {
	info, err := microWorkDao.NewStoreSftpConfig().GetOne(c, map[string]interface{}{
		"brand":      "core_area",
		"deleted_at": 0,
	})
	if err != nil {
		return err
	}
	port, _ := strconv.Atoi(info.Port)
	config := sftp.SSHConfig{
		Host:     info.Host,
		Port:     int64(port),
		UserName: info.User,
		Password: info.Password,
	}
	client, err := sftp.NewSFTPClient(config)
	if err != nil {
		return err
	}
	fileLists, err := client.Client.ReadDir(SFTPROOT + info.RootPath)
	if err != nil {
		return err
	}
	if err != nil {
		log.Info(c).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", s.Name, enterpriseId, err)
		return err
	}
	fileInfoList, err := reserveDao.NewReserveOrderDataAcxiomSftpLogDao().GetExistsFile(c, enterpriseId, 1, "area")
	if err != nil {
		log.Info(c).Msgf("%s::企业ID:%d, 获取目录下文件失败: %v", s.Name, enterpriseId, err)
		return err
	}
	fileInfoMap := make(map[string]microScrmModel.ReserveOrderDataAcxiomSftpLogModel)
	if len(fileInfoList) > 0 {
		for _, v := range fileInfoList {
			value := v
			fileInfoMap[value.LocalFile] = value
		}
	}

	for _, fileInfo := range fileLists {
		if fileInfo.IsDir() {
			continue
		}
		fileName := SFTPROOT + info.RootPath + fileInfo.Name()
		if len(fileInfoMap) > 0 {
			if fileModel, ok := fileInfoMap[fileName]; ok {
				if fileModel.Status == 0 {
					err := s.pullLogData(c, enterpriseId, client, fileModel, fileName, info.Brand)
					if err != nil {
						return err
					}
				}
			} else {
				err := s.pullLogData(c, enterpriseId, client, microScrmModel.ReserveOrderDataAcxiomSftpLogModel{}, fileName, info.Brand)
				if err != nil {
					return err
				}
			}
		} else {
			err := s.pullLogData(c, enterpriseId, client, microScrmModel.ReserveOrderDataAcxiomSftpLogModel{}, fileName, info.Brand)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (s SyncArea) pullLogData(ctx ctx.Context, enterpriseId int64, client *sftp.SFTP, info microScrmModel.ReserveOrderDataAcxiomSftpLogModel, fileName, brandName string) error {
	sftpFile, err := client.Client.Open(fileName)
	if err != nil {
		log.Info(ctx).Msgf("%s::读取SFTP目录下的文件失败: %v", fileName, err)
		return err
	}

	LocalFileName := strconv.FormatInt(app.UID(), 10) + `.csv`
	LocalFilePath := `./storage/core_area/` + LocalFileName
	dstFile, err := os.Create(LocalFilePath)
	if err != nil {
		log.Info(ctx).Msgf("%s::创建文件失败: %v", LocalFilePath, err)
		return err
	}
	if _, err := sftpFile.WriteTo(dstFile); err != nil {
		log.Info(ctx).Msgf("%s::写入文件失败: %v", LocalFilePath, err)
		return err
	}
	dataByte, _ := ioutil.ReadFile(LocalFilePath)
	logId := info.ID
	if logId == 0 {
		// 上传cos
		cs := storage.Init(cos.COS)
		path, err := cs.Upload(ctx, LocalFilePath, "third_sftp/core_area/"+LocalFileName)
		if err != nil {
			log.Info(ctx).Msgf("uploadCloud Error: %v", err)
			return err
		}
		// 删除本地文件
		err = os.Remove(LocalFilePath)
		if err != nil {
			log.Info(ctx).Msgf("delete local file error：%s ,%v", LocalFilePath, err)
		}
		fileUrl := cs.GetURL(path)
		logId = app.UIDContext(ctx)
		insertData := microScrmModel.ReserveOrderDataAcxiomSftpLogModel{
			MicroBaseModel: models.MicroBaseModel{
				BaseModel:    models.CreateBase(logId),
				EnterpriseID: enterpriseId,
			},
			FileType:  1,
			LocalFile: sftpFile.Name(),
			FileURL:   fileUrl,
			Status:    0,
			FileTime:  int(time.Now().Unix()),
			PushType:  "area",
		}
		err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().CreateLog(ctx, insertData)
		if err != nil {
			log.Info(ctx).Msgf("创建core_area报错:%v", err)
		}
	}
	// 读取 CSV 记录
	zr, err := zip.NewReader(bytes.NewReader(dataByte), int64(len(dataByte)))
	if err != nil {
		return err
	}
	fileData, err := xlsx.ReadZipReader(zr)
	if err != nil {
		return err
	}
	records, _ := readXlsxData(fileData.Sheets)

	if err != nil {
		return err
	}

	//开始处理业务逻辑
	unixTime := time.Now().Unix()
	var fileSourceFiscalYearMap = make(map[string]string) //财年标识
	var fileSourceNameMap = make(map[string]string)       //渠道
	var fileRegionNameMap = make(map[string]string)       //大区
	var fileAreaNameMap = make(map[string]string)         //分区
	var storeNodeNumMap = make(map[string]int)            //门店出现的次数标识
	for _, record := range records {
		storeCode := strings.TrimSpace(record.RetailPOSNumber)
		isVirtualIndex := strings.Index(record.Location, "VIRTUAL")
		if isVirtualIndex != -1 {
			// 忽略虚拟门店
			continue
		}
		if storeCode == "N/A" {
			// 忽略无编码门店
			continue
		}
		if _, ok := storeNodeNumMap[storeCode]; !ok {
			//门店编码第一次出现时是渠道数据
			fileSourceNameMap[storeCode] = strings.TrimSpace(record.CustomGroupElement)
			fileSourceFiscalYearMap[storeCode] = strings.TrimSpace(record.CustomGroup)
			storeNodeNumMap[storeCode] = 1
		} else {
			if storeNodeNumMap[storeCode] == 1 {
				//门店编码第2次出现时是大区数据
				fileRegionNameMap[storeCode] = strings.TrimSpace(record.CustomGroupElement)
			} else if storeNodeNumMap[storeCode] == 2 {
				//门店编码第3次出现时是分区数据
				fileAreaNameMap[storeCode] = strings.TrimSpace(record.CustomGroupElement)
			}
			storeNodeNumMap[storeCode]++
		}
	}
	//处理数据不完整的门店（hainan)
	for storeCode, num := range storeNodeNumMap {
		if num == 2 {
			isRegion := strings.Index(fileRegionNameMap[storeCode], "Total")
			if isRegion != -1 {
				// 正常大区中如果不存在分区，则认为已经闭店
				continue
			}
			// 没有分区的门店，将大区数据复制一份给分区
			fileAreaNameMap[storeCode] = fileRegionNameMap[storeCode]
		}
	}
	//获取db中所有已存在节点
	allArea, err := dashboardDao.NewArea().GetAll(ctx, enterpriseId)
	if err != nil {
		return err
	}
	var allAreaMap = make(map[int]dashboardModel.DashboardAreaModel)
	var maxId int64
	for _, area := range allArea {
		if area.ID > maxId {
			maxId = area.ID //获取表中最大的id
		}
		allAreaMap[int(area.ID)] = area
	}
	//使用节点的整个path的名称组合起来作为节点的唯一key，映射节点id
	var pathNamePrimaryIdMap = make(map[string]int64)
	for _, area := range allArea {
		pathNameStr := s.getPathNameStr(area, allAreaMap)
		pathNamePrimaryIdMap[pathNameStr] = area.ID
	}
	//获取企微门店id
	thirdIdMap, err := microWorkDao.NewStoreDepartment().GetThirdIdMap(ctx, enterpriseId)
	if err != nil {
		return err
	}
	var insertData []dashboardModel.DashboardAreaModel //要新增数据
	var updatePrimaryIds []int64                       //要更新的节点
	//循环中的门店的渠道id、大区id、分区id、门店id
	var regionId int64
	var sourceId int64
	var areaId int64
	var storeId int64
	//循环所有分区数据的map（自动忽略无分区的数据）
	for storeCode := range fileAreaNameMap {
		if _, ok := thirdIdMap[storeCode]; !ok {
			//门店在企微系统内不存在直接跳过
			continue
		}
		//处理该门店渠道
		if _, ok := pathNamePrimaryIdMap[fileSourceNameMap[storeCode]]; ok {
			//门店渠道存在，获取渠道主键id
			sourceId = pathNamePrimaryIdMap[fileSourceNameMap[storeCode]]
			updatePrimaryIds = append(updatePrimaryIds, sourceId)
		} else {
			//门店渠道不存在，新增渠道
			maxId++
			sourceId = maxId
			tmp := dashboardModel.DashboardAreaModel{
				Name:       fileSourceNameMap[storeCode],
				Sort:       sourceId,
				Type:       dashboardModel.AreaSource,
				FiscalYear: fileSourceFiscalYearMap[storeCode],
				BaseModel: models.BaseModel{
					ID:        sourceId,
					CreatedAt: unixTime,
					UpdatedAt: unixTime,
				},
			}
			insertData = append(insertData, tmp)
			pathNamePrimaryIdMap[fileSourceNameMap[storeCode]] = sourceId
		}

		//处理该门店大区
		if _, ok := fileRegionNameMap[storeCode]; !ok {
			//大区不存在，跳出
			continue
		}
		//门店大区存在，根据path获取大区主键id
		pathName := fileSourceNameMap[storeCode] + "-" + fileRegionNameMap[storeCode]
		if _, ok := pathNamePrimaryIdMap[pathName]; ok {
			regionId = pathNamePrimaryIdMap[pathName]
			updatePrimaryIds = append(updatePrimaryIds, regionId)
		} else {
			//门店大区不存在，新增大区
			maxId++
			regionId = maxId
			tmp := dashboardModel.DashboardAreaModel{
				Name:     fileRegionNameMap[storeCode],
				Sort:     regionId,
				ParentID: int(sourceId),
				Path:     helper.Int64ToString(sourceId),
				Type:     dashboardModel.AreaRegion,
				BaseModel: models.BaseModel{
					ID:        regionId,
					CreatedAt: unixTime,
					UpdatedAt: unixTime,
				},
			}
			insertData = append(insertData, tmp)
			pathNamePrimaryIdMap[pathName] = regionId
		}

		//处理该门店分区
		if _, ok := fileAreaNameMap[storeCode]; !ok {
			//分区不存在，跳出
			continue
		}
		//门店分区存在，根据path获取分区主键id
		pathName = fileSourceNameMap[storeCode] + "-" + fileRegionNameMap[storeCode] + "-" + fileAreaNameMap[storeCode]
		if _, ok := pathNamePrimaryIdMap[pathName]; ok {
			areaId = pathNamePrimaryIdMap[pathName]
			updatePrimaryIds = append(updatePrimaryIds, areaId)
		} else {
			//门店分区不存在，新增分区
			maxId++
			areaId = maxId
			tmp := dashboardModel.DashboardAreaModel{
				Name:     fileAreaNameMap[storeCode],
				Sort:     areaId,
				Path:     helper.Int64ToString(sourceId) + "," + helper.Int64ToString(regionId),
				ParentID: int(regionId),
				Type:     dashboardModel.AreaArea,
				BaseModel: models.BaseModel{
					ID:        areaId,
					CreatedAt: unixTime,
					UpdatedAt: unixTime,
				},
			}
			insertData = append(insertData, tmp)
			pathNamePrimaryIdMap[pathName] = areaId
		}

		//处理该门店
		//门店存在，根据path获取门店主键id
		pathName = fileSourceNameMap[storeCode] + "-" + fileRegionNameMap[storeCode] + "-" + fileAreaNameMap[storeCode] + "-" + storeCode
		if _, ok := pathNamePrimaryIdMap[pathName]; ok {
			storeId = pathNamePrimaryIdMap[pathName]
			updatePrimaryIds = append(updatePrimaryIds, storeId)
		} else {
			//门店不存在，新增门店
			maxId++
			storeId = maxId
			tmp := dashboardModel.DashboardAreaModel{
				Name:     storeCode,
				Sort:     storeId,
				Path:     helper.Int64ToString(sourceId) + "," + helper.Int64ToString(regionId) + "," + helper.Int64ToString(areaId),
				ParentID: int(areaId),
				Type:     dashboardModel.AreaStore,
				BaseModel: models.BaseModel{
					ID:        storeId,
					CreatedAt: unixTime,
					UpdatedAt: unixTime,
				},
			}
			if thirdId, ok := thirdIdMap[storeCode]; ok {
				tmp.ThirdID = thirdId
			}
			insertData = append(insertData, tmp)
			pathNamePrimaryIdMap[pathName] = storeId
		}
	}
	// 写入所有新增节点
	err = dashboardDao.NewArea().CreateInBatchesData(ctx, enterpriseId, insertData)
	if err != nil {
		return err
	}
	// 把所有需要更新的节点更新时间标记为当前时间
	err = dashboardDao.NewArea().UpdateData(ctx, enterpriseId, map[string]interface{}{
		"id": updatePrimaryIds,
	}, map[string]interface{}{
		"updated_at": unixTime,
	})
	if err != nil {
		return err
	}
	// 把所有更新时间在当前时间之前的节点都删除
	err = dashboardDao.NewArea().DeleteOldData(ctx, enterpriseId, unixTime)
	if err != nil {
		return err
	}
	err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().UpdateLogStatus(ctx, enterpriseId, logId)
	if err != nil {
		return err
	}
	planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
		EnterpriseName: brandName,
		SFTPName:       "SFTP同步core管理架构文件",
		UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
	})
	return nil
}

// GetMemberAuthAreaIds 获取导购组织架构权限节点
func (s SyncArea) getPathNameStr(area dashboardModel.DashboardAreaModel, areaMap map[int]dashboardModel.DashboardAreaModel) string {
	// 组装节点的path
	var pathNameStr string
	for _, tmp := range append(strings.Split(area.Path, ","), helper.Int64ToString(area.ID)) {
		if areaInfo, ok := areaMap[helper.StringToInt(tmp)]; ok {
			pathNameStr += areaInfo.Name + "-"
		}
	}
	return strings.TrimRight(pathNameStr, "-")
}

func readXlsxData(sheets []*xlsx.Sheet) ([]*dashboardResp.FileRow, error) {
	riMin := 1

	var records []*dashboardResp.FileRow
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
			var record = &dashboardResp.FileRow{}

			if row.Cells[0].String() != "" {
				record.CustomGroup = row.Cells[0].String()
				record.CustomGroupElement = row.Cells[1].String()
				record.Location = row.Cells[2].String()
				record.RetailPOSNumber = row.Cells[4].String()
				records = append(records, record)
			}

		}
	}
	return records, nil
}
