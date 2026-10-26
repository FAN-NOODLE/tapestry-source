package new_cdp

import (
	"archive/zip"
	"database/sql"
	"encoding/csv"
	"fmt"
	sftp2 "github.com/pkg/sftp"
	"github.com/shopspring/decimal"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"io"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	reserveDao "micro-scrm/internal/app/dao/micro_scrm"
	tagsDao "micro-scrm/internal/app/dao/tags"
	"micro-scrm/internal/app/models"
	departmentsModels "micro-scrm/internal/app/models/departments"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	planService "micro-scrm/internal/app/service/plan"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"path"
	"strconv"
	"time"
)

type SyncDataToDEL struct {
	name       string
	startTime  int64
	endTime    int64
	sqlDb      *sql.DB
	suffix     string
	remotePath string
	SftpClient *sftp.SFTP
}

func NewSyncDataToDEL() *SyncDataToDEL {
	startTime, endTime := helper.YesterdayPeriodOfTimestamp(1)
	sqlDb, _ := gormx.GetInstance().DB()
	suffix := helper.GetFormatDate(0, 2)
	return &SyncDataToDEL{
		name:       "同步增量数据上传至DEL",
		startTime:  startTime,
		endTime:    endTime,
		sqlDb:      sqlDb,
		suffix:     suffix,
		remotePath: "inbound/cdp/",
	}
}

func (job *SyncDataToDEL) Run(c ctx.Context, jobHandler *xjob.JobHandler) error {
	sftpConfigList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList(micro_scrm.SftpTypeDEL)
	if err != nil {
		log.Info(c).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	for _, config := range sftpConfigList {
		if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
			log.Info(c).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
			message := planService.SFTPMessageData{
				EnterpriseName: config.BrandName,
				SFTPName:       "NewCdp的sftp连接失败",
				UploadRes:      fmt.Sprintf("失败,错误为：%s", clientErr.Error()),
			}
			planService.NewTaskNotice().SftpSendNotice(c, message)
			// return
		} else {
			job.SftpClient = client
			err = job.handle(c, config.EnterpriseID, client, config)
		}

	}
	return err
}

func (job *SyncDataToDEL) handle(c ctx.Context, enterpriseId int64, client *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail) error {
	dirPath, err := job.CreateDir(client)
	if err != nil {
		log.Info(c).Msgf("%s::企业ID:%d, 生成目录失败: %v", job.name, enterpriseId, err)
		return err
	}
	job.remotePath = dirPath
	brandName := config.BrandName
	departmentsMap := make(map[int]string)
	// 获取所有部门，departments表中name为门店code
	departments, err := departmentsDao.NewDepartments().GetAll(c, departmentsModels.DepartmentsModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		}})
	if err != nil {
		return err
	}
	for _, department := range departments {
		departmentsMap[department.ThirdID] = department.Name
	}
	wg := threading.RoutineGroup{}
	// 生成邀约的数据，并上传COS
	wg.RunSafe(func() {
		err = job.createReserveOrders(c, enterpriseId, dirPath, job.suffix, client, departmentsMap)
		message := planService.SFTPMessageData{
			EnterpriseName: brandName,
			SFTPName:       "sftp增量邀约单上传至NewCdp",
			UploadRes:      "成功",
		}
		if err != nil {
			message.UploadRes = fmt.Sprintf("失败,错误为：%s", err.Error())
		}
		planService.NewTaskNotice().SftpSendNotice(c, message)
	})

	// 生成印象标签的数据，并上传COS
	wg.RunSafe(func() {
		err = job.createImpressionTags(c, enterpriseId, brandName, dirPath, job.suffix, client)
		message := planService.SFTPMessageData{
			EnterpriseName: brandName,
			SFTPName:       "sftp增量印象标签上传至NewCdp",
			UploadRes:      "成功",
		}
		if err != nil {
			message.UploadRes = fmt.Sprintf("失败,错误为：%s", err.Error())
		}
		planService.NewTaskNotice().SftpSendNotice(c, message)
	})
	wg.Wait()

	_ = job.masterTableHandle(c, enterpriseId, "micro_scrm", "sow_task", "id,enterprise_id,work_external_user_id,member_id,member_user_id,content,status,remind_time,is_push,template_id,created_at,updated_at,deleted_at")
	_ = job.masterTableHandle(c, enterpriseId, "micro_scrm", "sow_template", "id,enterprise_id,name,theme,status,admin_name,use_num,created_at,updated_at,deleted_at")
	_ = job.masterTableHandle(c, enterpriseId, "micro_scrm", "sow_template_notice", "id,enterprise_id,template_id,delay_type,delay_value,delay_time,content,created_at,updated_at,deleted_at")
	//活动二维码相关表

	//_ = job.shardTableHandle(c, fmt.Sprintf("micro_scrm_enterprise_%d", enterpriseId), "activity_qrcode_main", "id,name,alias_name,banner_url,time_type,start_time,end_time,category_ids,add_name_tag,add_category_tag,tag_ids,name_tag_id,skip_verify,created_at,updated_at,deleted_at")
	//_ = job.shardTableHandle(c, fmt.Sprintf("micro_scrm_enterprise_%d", enterpriseId), "activity_qrcode_category", "id,name,tag_id,created_at,updated_at,deleted_at")
	//_ = job.shardTableHandle(c, fmt.Sprintf("micro_scrm_enterprise_%d", enterpriseId), "activity_qrcode_member_scope", "id,activity_id,member_id,area_id,created_at,updated_at,deleted_at")
	//_ = job.shardTableHandle(c, fmt.Sprintf("micro_scrm_enterprise_%d", enterpriseId), "activity_qrcode_contact_way", "id,activity_id,state,third_id,store_name,config_id,qr_code,created_at,updated_at,deleted_at")
	//_ = job.shardTableHandle(c, fmt.Sprintf("micro_scrm_enterprise_%d", enterpriseId), "activity_qrcode_point", "id,activity_id,name,is_point,category_id,code_url,link_url,created_at,updated_at,deleted_at")

	message := planService.SFTPMessageData{
		EnterpriseName: brandName,
		SFTPName:       "sftp增量数据上传NewCdp完毕",
		UploadRes:      "成功",
	}
	if err != nil {
		message.UploadRes = "失败"
		return err
	}
	// 生成rdy文件
	err = job.createRdy(c, dirPath, job.suffix, client)
	if err != nil {
		log.Info(c).Msgf("生成rdy文件失败，%v", err)
		return err
	}
	// 消息通知到企微群里
	planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
		EnterpriseName: brandName,
		SFTPName:       "sftp增量数据上传NewCdp完毕",
		UploadRes:      "成功",
	})
	return nil
}

func (job *SyncDataToDEL) createReserveOrders(c ctx.Context, enterpriseId int64, dirPath, timeStr string, client *sftp.SFTP, departmentsMap map[int]string) error {
	list, err := reserveDao.NewReservesOrder().GetMany(c, enterpriseId, "updated_at between ? and ? and deleted_at=0", job.startTime, job.endTime)
	if err != nil {
		return err
	}
	countNum := strconv.Itoa(len(list))
	localFileName := "invitation_order_" + timeStr + `.csv`
	localTextName := "invitation_order_" + timeStr + `.csv.done`
	filePath := `./storage/new_cdp_del/` + localFileName
	textPath := `./storage/new_cdp_del/` + localTextName
	t, _ := os.Create(textPath) // 创建验证上传OK的文件
	_, _ = t.WriteString("\xEF\xBB\xBF")
	_, _ = t.WriteString(fmt.Sprintf("表名,数据量\n"))
	_, _ = t.WriteString(fmt.Sprintf("invitation_order,%s", countNum))
	_ = t.Close()
	f, err := os.Create(filePath) // 创建文件
	if err != nil {
		log.Info(c).Msgf("%s::csv文件创建失败: %v", job.name, err)
		return err
	}

	_, _ = f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	data := [][]string{
		{"最新邀约单时间", "邀约单号", "邀约单状态", "顾客union id", "到店时间", "邀约导购ID", "预约服务名称", "邀约门店code", "订单交易时间", "成交金额", "匹配成功时间", "匹配状态"},
	}
	w := csv.NewWriter(f) // 创建一个新写入文件流
	w.Comma = '\t'
	if len(list) > 0 {
		for _, v := range list {
			var updateTimeStr, reserveTimeStr, reserveOrderTimeStr, successTimeStr, orderPrice string

			updateTimeStr = helper.GetFormatDateByParams(v.UpdatedAt)
			reserveTimeStr = helper.GetFormatDateByParams(int64(v.ReserveTime))
			if v.Status == microScrmModel.FINISH {
				successTimeStr = updateTimeStr
			}
			serveName := v.ReserveService.ReserveName
			storeCode := "101test" // 默认的storeCode为101test
			if item, ok := departmentsMap[int(v.StoreID)]; ok {
				storeCode = item
			}
			if v.ReserveOrderPrice > 0 {
				reserveOrderTimeStr = helper.GetFormatDateByParams(int64(v.ReserveOrderTime))
				reserveOrderPriceFloat := decimal.NewFromInt(int64(v.ReserveOrderPrice)).Div(decimal.NewFromInt(100))
				orderPRiceFloat, _ := reserveOrderPriceFloat.Float64()
				orderPrice = strconv.FormatFloat(orderPRiceFloat, 'f', -1, 64)
			}
			tmp := []string{
				updateTimeStr, v.ReserveNumber, strconv.FormatInt(int64(v.Status), 10), v.UserUnionid, reserveTimeStr, v.MemberUserID, serveName, storeCode, reserveOrderTimeStr, orderPrice, successTimeStr, strconv.Itoa(int(v.ReserveOrderStatus)),
			}
			data = append(data, tmp)
		}
	}
	err = w.WriteAll(data)
	if err != nil {
		fmt.Printf("写入csv文件报错，err:%v\n", err)
		return err
	} // 写入数据
	w.Flush()

	fmt.Printf("%d的企业csv日志文件写入成功"+filePath+",日志数量:"+strconv.Itoa(len(list))+"\n", enterpriseId)
	err = job.uploadFile(c, client, filePath, dirPath)
	if err != nil {
		log.Info(c).Msgf("%s::sftp文件写入失败: %v", job.name, err)
		return err
	}
	log.Info(c).Msgf("%s::写入sftp文件成功: %s", job.name, filePath)

	err = job.uploadFile(c, client, textPath, dirPath)
	if err != nil {
		log.Info(c).Msgf("%s::sftp验证done文件写入失败: %v", job.name, err)
		return err
	}

	// 将csv文件上传直cos存储
	cs := storage.Init(cos.COS)
	cosPath, err := cs.Upload(c, filePath, "new_cdp/"+localFileName)
	if err != nil {
		log.Info(c).Msgf("uploadCloud Error: %v", err)
		return err
	}
	// 删除本地文件
	err = f.Close()
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Info(c).Msgf("%s--文件删除失败 Error: %v", name, err)
		}
	}(textPath)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Info(c).Msgf("%s--文件删除失败 Error: %v", name, err)
		}
	}(filePath)
	zeroTime := helper.GetZeroTimeByTime(time.Now()).Unix()
	fileUrl := cs.GetURL(cosPath)
	insertData := microScrmModel.ReserveOrderDataAcxiomSftpLogModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel:    models.CreateBase(app.UID()),
			EnterpriseID: enterpriseId,
		},
		FileType:  0,
		LocalFile: localFileName,
		FileURL:   fileUrl,
		Status:    1,
		FileTime:  int(zeroTime),
		PushType:  "push_new_cdp_invitation_order",
	}
	err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().CreateLog(c, insertData)
	if err != nil {
		log.Info(c).Msgf("推送数据到NewCdp记录日志失败:%v", err)
	}
	return nil
}

func (job *SyncDataToDEL) createImpressionTags(c ctx.Context, enterpriseId int64, brandName, dirPath, timeStr string, client *sftp.SFTP) error {

	list, err := tagsDao.NewTagsOperateLog().GetMany(c, "enterprise_id= ? and updated_at between ? and ? ", enterpriseId, job.startTime, job.endTime)
	if err != nil {
		return err
	}
	countNum := strconv.Itoa(len(list))
	localFileName := "impression_tags_" + timeStr + `.csv`
	localTextName := "impression_tags_" + timeStr + `.csv.done`
	filePath := `./storage/new_cdp_del/` + localFileName
	textPath := `./storage/new_cdp_del/` + localTextName
	t, _ := os.Create(textPath) // 创建验证上传OK的文件
	_, _ = t.WriteString("\xEF\xBB\xBF")
	_, _ = t.WriteString("\xEF\xBB\xBF")
	_, _ = t.WriteString(fmt.Sprintf("表名,数据量\n"))
	_, _ = t.WriteString(fmt.Sprintf("impression_tags,%s", countNum))
	_ = t.Close()
	_ = t.Close()
	f, err := os.Create(filePath) // 创建文件
	if err != nil {
		log.Info(c).Msgf("%s::csv文件创建失败: %v", job.name, err)
		return err
	}

	_, _ = f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	data := [][]string{
		{"brand", "sa_code", "customer_unionid", "action_type", "tag_id", "unix_time"},
	}
	w := csv.NewWriter(f) // 创建一个新写入文件流
	w.Comma = '\t'
	if len(list) > 0 {
		for _, v := range list {
			tmp := []string{
				brandName, v.OperateUserId, v.ExternalUnionId, strconv.Itoa(int(v.ActionType)), v.ThirdTagID, strconv.FormatInt(v.UpdatedAt, 10),
			}
			data = append(data, tmp)
		}
	}
	err = w.WriteAll(data)
	if err != nil {
		fmt.Printf("写入csv文件报错，err:%v\n", err)
		return err
	} // 写入数据
	w.Flush()
	fmt.Printf("%s::%d的企业csv日志文件写入成功"+filePath+",日志数量:"+strconv.Itoa(len(list))+"\n", brandName, enterpriseId)
	err = job.uploadFile(c, client, filePath, dirPath)
	if err != nil {
		log.Info(c).Msgf("%s::sftp文件写入失败: %v", job.name, err)
		return err
	}
	log.Info(c).Msgf("%s::写入sftp文件成功: %s", job.name, filePath)

	err = job.uploadFile(c, client, textPath, dirPath)
	if err != nil {
		log.Info(c).Msgf("%s::sftp验证done文件写入失败: %v", job.name, err)
		return err
	}

	// 将csv文件上传直cos存储
	cs := storage.Init(cos.COS)
	cosPath, err := cs.Upload(c, filePath, "new_cdp/"+localFileName)
	if err != nil {
		log.Info(c).Msgf("uploadCloud Error: %v", err)
		return err
	}
	// 删除本地文件
	err = f.Close()
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Info(c).Msgf("%s--文件删除失败 Error: %v", name, err)
		}
	}(textPath)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Info(c).Msgf("%s--文件删除失败 Error: %v", name, err)
		}
	}(filePath)
	zeroTime := helper.GetZeroTimeByTime(time.Now()).Unix()
	fileUrl := cs.GetURL(cosPath)
	insertData := microScrmModel.ReserveOrderDataAcxiomSftpLogModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel:    models.CreateBase(app.UID()),
			EnterpriseID: enterpriseId,
		},
		FileType:  0,
		LocalFile: localFileName,
		FileURL:   fileUrl,
		Status:    1,
		FileTime:  int(zeroTime),
		PushType:  "push_new_cdp_impression_tag",
	}
	err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().CreateLog(c, insertData)
	if err != nil {
		log.Info(c).Msgf("推送数据到NewCdp记录日志失败:%v", err)
	}

	return nil
}

func (job *SyncDataToDEL) CreateDir(sftp *sftp.SFTP) (string, error) {
	dirPath := job.remotePath + job.suffix
	_, statErr := sftp.Client.Stat(dirPath)
	if statErr != nil && os.IsNotExist(statErr) {
		err := sftp.Client.Mkdir(dirPath) // 创建目录
		if err != nil {
			return "", err
		}
		err = sftp.Client.Chmod(dirPath, 0777)
	}
	realPath, err := sftp.Client.RealPath(dirPath)
	if err != nil {
		return "", err
	}
	return realPath + "/", nil
}

func (job *SyncDataToDEL) uploadFile(c ctx.Context, sftp *sftp.SFTP, localFilePath string, remotePath string) error {
	srcFile, err := os.Open(localFilePath)
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {
			log.Info(c).Msgf("文件关闭失败 Error: %v", err)
		}
	}(srcFile)
	if err != nil {
		return err
	}
	var remoteFileName = path.Base(localFilePath)
	dstPath := path.Join(remotePath, remoteFileName)
	dstFile, err := sftp.Client.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE)
	defer func(dstFile *sftp2.File) {
		err := dstFile.Close()
		if err != nil {
			log.Info(c).Msgf("sftp远程文件关闭失败 Error: %v", err)
		}
	}(dstFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Info(c).Msgf("sftp远程文件上传失败 Error: %v", err)
		return err
	}
	return err
}

func (job *SyncDataToDEL) createRdy(c ctx.Context, dirPath, timeStr string, client *sftp.SFTP) error {
	localRdyName := timeStr + `.rdy`
	rdyPath := `./storage/new_cdp_del/` + localRdyName
	t, _ := os.Create(rdyPath) // 创建验证上传OK的文件
	err := job.uploadFile(c, client, rdyPath, dirPath)
	err = t.Close()
	return err
}

// 主库表的处理
func (job *SyncDataToDEL) masterTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName, filed string) error {
	rows, err := job.sqlDb.Query("select "+filed+" from "+databaseName+"."+tableName+" where enterprise_id=? and updated_at>=?", enterpriseId, job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "new_cdp的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count` FROM "+databaseName+"."+tableName+" WHERE enterprise_id=?", tableName, enterpriseId)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count",
	}
	//var total int64
	//err = job.sqlDb.QueryRow("SELECT COUNT(1) FROM "+databaseName+"."+tableName+" WHERE enterprise_id=? and updated_at>=?", enterpriseId, job.startTime).Scan(&total)
	//if err != nil {
	//	fmt.Printf("查询"+tableName+"的.done文件数据total报错%v\n", err)
	//	return err
	//}
	err = job.CreateDoneFile(c, tableName, columnsName, sqlRows, incrementCount)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
	}
	return err
}

// 分库的查询逻辑
func (job *SyncDataToDEL) shardTableHandle(c ctx.Context, databaseName, tableName, filed string) error {
	rows, err := job.sqlDb.Query("select "+filed+" from "+databaseName+"."+tableName+" where updated_at>=?", job.startTime)
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count` FROM "+databaseName+"."+tableName, tableName)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count",
	}
	//var total int64
	//err = job.sqlDb.QueryRow("SELECT COUNT(1) FROM "+databaseName+"."+tableName+" WHERE updated_at>=?", job.startTime).Scan(&total)
	//if err != nil {
	//	fmt.Printf("查询"+tableName+"的.done文件数据total报错%v\n", err)
	//	return err
	//}
	err = job.CreateDoneFile(c, tableName, columnsName, sqlRows, incrementCount)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
	}
	return err
}

func (job *SyncDataToDEL) CreateFile(c ctx.Context, tableName string, rows *sql.Rows) (int64, error) {

	fileName := tableName + "_" + job.suffix
	csvFile, incrementCount, err := job.CreateCsvFile(fileName, rows)
	if err != nil {
		return incrementCount, err
	}
	csvPath := "./storage/new_cdp_del/" + fileName + ".csv"
	zipFile, err := job.CreateZipFile(fileName, csvFile)
	if err != nil {
		return incrementCount, err
	}
	localPath := "./storage/new_cdp_del/" + fileName + ".zip"
	err = job.uploadZipFile(c, zipFile, localPath)
	_ = os.Remove(localPath)
	_ = os.Remove(csvPath)
	return incrementCount, err
}

// CreateZipFile 创建ZIP文件
func (job *SyncDataToDEL) CreateZipFile(zipName string, csvFile *os.File) (*os.File, error) {
	// 创建一个新的 ZIP 文件
	zipFileName := "./storage/new_cdp_del/" + zipName + ".zip"
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		fmt.Println("Error creating ZIP file:", err)
		return nil, err
	}
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			fmt.Println("zipFile  close:", err)
		}
	}(zipFile)
	// 创建一个 ZIP 写入器
	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			fmt.Println("zipWriter  close:", err)
		}
	}(zipWriter)

	// 将 CSV 文件添加到 ZIP 文件中
	csvFileInfo, err := os.Stat(csvFile.Name())
	if err != nil {
		fmt.Println("Error getting CSV file info:", err)
		return nil, err
	}

	zipFileHeader, err := zip.FileInfoHeader(csvFileInfo)
	if err != nil {
		fmt.Println("Error creating ZIP file header:", err)
		return nil, err
	}
	zipFileHeader.Method = zip.Deflate

	zipFileWriter, err := zipWriter.CreateHeader(zipFileHeader)
	if err != nil {
		fmt.Println("Error creating ZIP file writer:", err)
		return nil, err
	}
	csvFile, err = os.Open(csvFile.Name())
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return nil, err
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			fmt.Println("csvFile close:", err)
		}
	}(csvFile)
	_, err = io.Copy(zipFileWriter, csvFile)
	if err != nil {
		fmt.Println("Error copying CSV file to ZIP file:", err)
		return nil, err
	}

	return zipFile, nil
}

// CreateCsvFile 创建CSV文件
func (job *SyncDataToDEL) CreateCsvFile(csvName string, rows *sql.Rows) (*os.File, int64, error) {
	incrementCount := int64(0)
	// 创建一个新的 CSV 文件
	csvFileName := "./storage/new_cdp_del/" + csvName + ".csv"
	csvFile, err := os.Create(csvFileName)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return nil, incrementCount, err
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			fmt.Println(" CSV file close:", err)
		}
	}(csvFile)
	// 创建一个 CSV 写入器
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	columnNames, err := rows.Columns()
	if err != nil {
		fmt.Println("Error getting column names:", err)
		return nil, incrementCount, err
	}

	// 写入列名
	err = csvWriter.Write(columnNames)
	if err != nil {
		fmt.Println("Error writing column names:", err)
		return nil, incrementCount, err
	}

	// 遍历数据表中的每一行
	for rows.Next() {
		// 创建一个切片来存储每一行的数据
		rowData := make([]interface{}, len(columnNames))
		rowDataPtrs := make([]interface{}, len(columnNames))
		for i := range rowData {
			rowDataPtrs[i] = &rowData[i]
		}
		// 将每一行的数据读取到切片中
		err := rows.Scan(rowDataPtrs...)
		if err != nil {
			fmt.Println("Error scanning row data:", err)
			return nil, incrementCount, err
		}
		rowDataStr := make([]string, len(rowData))
		for i, v := range rowData {
			switch value := v.(type) {
			case string:
				rowDataStr[i] = value
			case []byte:
				rowDataStr[i] = string(value)
			default:
				rowDataStr[i] = fmt.Sprintf("%v", value)
			}
		}
		// 将每一行的数据写入 CSV 文件
		err = csvWriter.Write(rowDataStr)
		if err != nil {
			fmt.Println("Error writing row data:", err)
			return nil, incrementCount, err
		}
		incrementCount++
	}
	return csvFile, incrementCount, err
}

func (job *SyncDataToDEL) CreateDoneFile(c ctx.Context, tableName string, columnNames []string, rows *sql.Rows, total int64) error {
	// 创建一个新的 CSV 文件
	fileName := tableName + "_" + job.suffix
	csvFileName := "./storage/new_cdp_del/" + fileName + "done.csv"
	doneFileName := "./storage/new_cdp_del/" + fileName + ".done"
	csvFile, err := os.Create(csvFileName)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return err
	}
	// 创建一个 CSV 写入器
	csvWriter := csv.NewWriter(csvFile)
	if err != nil {
		fmt.Println("Error getting column names:", err)
		return err
	}

	// 写入列名
	err = csvWriter.Write(columnNames)
	if err != nil {
		fmt.Println("Error writing column names:", err)
		return err
	}

	rowLength := len(columnNames)
	// 遍历数据表中的每一行
	for rows.Next() {
		// 创建一个切片来存储每一行的数据
		rowData := make([]interface{}, rowLength)
		rowDataPtrs := make([]interface{}, rowLength)
		for i := range rowData {
			if i == 1 {
				rowDataPtrs[1] = &total
				rowData[1] = total
			} else {
				rowDataPtrs[i] = &rowData[i]
			}
		}
		// 将每一行的数据读取到切片中
		err := rows.Scan(rowDataPtrs...)
		if err != nil {
			fmt.Println("Error scanning row data:", err)
			return err
		}
		rowDataStr := make([]string, len(rowDataPtrs))
		for i, v := range rowData {
			switch value := v.(type) {
			case string:
				rowDataStr[i] = value
			case []byte:
				rowDataStr[i] = string(value)
			case float64:
				rowDataStr[i] = strconv.FormatFloat(value, 'f', -1, 64)
			default:
				strValue := fmt.Sprintf("%v", value)
				rowDataStr[i] = strValue
			}
		}
		// 将每一行的数据写入 CSV 文件
		err = csvWriter.Write(rowDataStr)
		if err != nil {
			fmt.Println("Error writing row data:", err)
			return err
		}
	}
	csvWriter.Flush()
	err = csvFile.Close()
	if err != nil {
		fmt.Println(" CSV file close:", err)
	}
	err = os.Rename(csvFileName, doneFileName)
	if err != nil {
		fmt.Println(" CSV file rename err:", err)
	}
	err = job.uploadDoneFile(c, doneFileName)
	_ = os.Remove(doneFileName)
	return err
}

func (job *SyncDataToDEL) uploadDoneFile(c ctx.Context, localFilePath string) error {
	var remoteFileName = path.Base(localFilePath)
	dstPath := path.Join(job.remotePath, remoteFileName)
	dstFile, err := job.SftpClient.Client.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE)
	defer func(dstFile *sftp2.File) {
		err := dstFile.Close()
		if err != nil {
			log.Info(c).Msgf("sftp远程文件关闭失败 Error: %v", err)
		}
	}(dstFile)
	if err != nil {
		return err
	}
	srcFile, err := os.Open(localFilePath)
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {
			fmt.Println("Error close done file:", err)
		}
	}(srcFile)
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Info(c).Msgf("sftp远程文件上传失败 Error: %v", err)
		return err
	}
	return nil
}

func (job *SyncDataToDEL) uploadZipFile(c ctx.Context, zipFile *os.File, localFilePath string) error {
	var remoteFileName = path.Base(localFilePath)
	dstPath := path.Join(job.remotePath, remoteFileName)
	dstFile, err := job.SftpClient.Client.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE)
	defer func(dstFile *sftp2.File) {
		err := dstFile.Close()
		if err != nil {
			log.Info(c).Msgf("sftp远程文件关闭失败 Error: %v", err)
		}
	}(dstFile)
	if err != nil {
		return err
	}
	zipFile, err = os.Open(zipFile.Name())
	if err != nil {
		fmt.Println("Error opening zip file:", err)
		return err
	}
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			fmt.Println("Error close zip file:", err)
		}
	}(zipFile)
	_, err = io.Copy(dstFile, zipFile)
	if err != nil {
		log.Info(c).Msgf("sftp远程文件上传失败 Error: %v", err)
		return err
	}
	return err
}
