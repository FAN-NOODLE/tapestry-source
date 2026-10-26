package new_cdp

import (
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
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"path"
	"strconv"
	"time"
)

type SyncDataToDELAll struct {
	name      string
	startTime int64
	endTime   int64
}

func NewSyncDataToDELALL() *SyncDataToDELAll {
	startTime, endTime := helper.YesterdayPeriodOfTimestamp(1)
	return &SyncDataToDELAll{
		name:      "同步全量数据上传至DEL",
		startTime: startTime,
		endTime:   endTime,
	}
}

func (job *SyncDataToDELAll) Run(c ctx.Context, jobHandler *xjob.JobHandler) error {
	sftpConfigList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList(micro_scrm.SftpTypeDEL)
	if err != nil {
		log.Info(c).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	for _, config := range sftpConfigList {
		if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
			log.Info(c).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
			// return
		} else {
			err = job.handle(c, config.EnterpriseID, client, config)
		}

	}
	return err
}

func (job *SyncDataToDELAll) handle(c ctx.Context, enterpriseId int64, client *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail) error {
	timeStr := helper.GetFormatDate(0, 2)
	dirPath, err := job.CreateDir(client, timeStr)
	if err != nil {
		log.Info(c).Msgf("%s::企业ID:%d, 生成目录失败: %v", job.name, enterpriseId, err)
		return err
	}
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
		err = job.createReserveOrders(c, enterpriseId, dirPath, timeStr, client, departmentsMap)
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
		err = job.createImpressionTags(c, enterpriseId, brandName, dirPath, timeStr, client)
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
	err = job.createRdy(c, dirPath, timeStr, client)
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

func (job *SyncDataToDELAll) createReserveOrders(c ctx.Context, enterpriseId int64, dirPath, timeStr string, client *sftp.SFTP, departmentsMap map[int]string) error {
	localFileName := "invitation_order_" + timeStr + `.csv`
	localTextName := "invitation_order_" + timeStr + `.csv.done`
	filePath := `./storage/new_cdp_del/` + localFileName
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
	id := 0
	count := 0
	for {
		list, err := reserveDao.NewReservesOrder().GetManyById(c, enterpriseId, "id > ? and updated_at <=? and deleted_at=0 ", id, job.endTime)
		if err != nil {
			return err
		}
		index := len(list)
		if index == 0 {
			break
		}
		for _, v := range list {
			count++
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
		id = int(list[index-1].ID)
	}
	err = w.WriteAll(data)
	if err != nil {
		fmt.Printf("写入csv文件报错，err:%v\n", err)
		return err
	} // 写入数据
	w.Flush()
	countNum := strconv.Itoa(count)
	textPath := `./storage/new_cdp_del/` + localTextName
	t, _ := os.Create(textPath) // 创建验证上传OK的文件
	_, _ = t.WriteString("\xEF\xBB\xBF")
	_, _ = t.WriteString(fmt.Sprintf("表名,数据量\n"))
	_, _ = t.WriteString(fmt.Sprintf("invitation_order,%s", countNum))
	_ = t.Close()
	fmt.Printf("%d的企业csv日志文件写入成功"+filePath+",日志数量:"+countNum+"\n", enterpriseId)
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
	if err != nil {
		log.Info(c).Msgf("delete local file error：%s ,%v", filePath, err)
	}
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

func (job *SyncDataToDELAll) createImpressionTags(c ctx.Context, enterpriseId int64, brandName, dirPath, timeStr string, client *sftp.SFTP) error {

	localFileName := "impression_tags_" + timeStr + `.csv`
	localTextName := "impression_tags_" + timeStr + `.csv.done`
	filePath := `./storage/new_cdp_del/` + localFileName
	f, err := os.Create(filePath) // 创建文件
	if err != nil {
		log.Info(c).Msgf("%s::csv文件创建失败: %v", job.name, err)
		return err
	}
	id := 0
	count := 0
	_, _ = f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	data := [][]string{
		{"brand", "sa_code", "customer_unionid", "action_type", "tag_id", "unix_time"},
	}
	w := csv.NewWriter(f) // 创建一个新写入文件流
	w.Comma = '\t'

	many, err := tagsDao.NewImpression().GetMany(c, "enterprise_id= ? and deleted_at=0", enterpriseId)
	if err != nil {
		return err
	}
	impressionTagMap := make(map[int64]string, 0)
	for _, model := range many {
		impressionTagMap[model.ID] = model.ThirdTagID
	}
	for {
		list, err := tagsDao.NewImpressionRelation().GetManyById(c, "enterprise_id= ? and id> ? and updated_at <= ? ", enterpriseId, id, job.endTime)
		if err != nil {
			return err
		}
		index := len(list)
		if index == 0 {
			break
		}
		externalIds := make([]int64, 0)
		for _, v := range list {
			externalIds = append(externalIds, v.ExternalUserID)
		}
		externalMany, err := reserveDao.NewWxUserCorpExternal().GetMany(c, enterpriseId, "id in ?", externalIds)
		if err != nil {
			return err
		}
		externalUnionIds := make(map[int64]string, 0)
		for _, v := range externalMany {
			externalUnionIds[v.ID] = v.Unionid
		}
		for _, v := range list {
			count++
			tmp := []string{
				brandName, v.UserID, externalUnionIds[v.ExternalUserID], "1", impressionTagMap[v.TagID], strconv.FormatInt(v.UpdatedAt, 10),
			}
			data = append(data, tmp)
		}
		id = int(list[index-1].ID)
	}
	err = w.WriteAll(data)
	if err != nil {
		fmt.Printf("写入csv文件报错，err:%v\n", err)
		return err
	} // 写入数据
	w.Flush()

	countNum := strconv.Itoa(count)
	textPath := `./storage/new_cdp_del/` + localTextName
	t, _ := os.Create(textPath) // 创建验证上传OK的文件
	_, _ = t.WriteString("\xEF\xBB\xBF")
	_, _ = t.WriteString("\xEF\xBB\xBF")
	_, _ = t.WriteString(fmt.Sprintf("表名,数据量\n"))
	_, _ = t.WriteString(fmt.Sprintf("impression_tags,%s", countNum))
	_ = t.Close()
	fmt.Printf("%s::%d的企业csv日志文件写入成功"+filePath+",日志数量:"+countNum+"\n", brandName, enterpriseId)
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
	if err != nil {
		log.Info(c).Msgf("delete local file error：%s ,%v", filePath, err)
	}
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

func (job *SyncDataToDELAll) CreateDir(sftp *sftp.SFTP, timeStr string) (string, error) {
	dirPath := "inbound/cdp/" + timeStr
	err := sftp.Client.Mkdir(dirPath) // 创建目录
	if err != nil {
		return "", err
	}
	err = sftp.Client.Chmod(dirPath, 0777)
	realPath, err := sftp.Client.RealPath(dirPath)
	if err != nil {
		return "", err
	}
	return realPath + "/", nil
}

func (job *SyncDataToDELAll) uploadFile(c ctx.Context, sftp *sftp.SFTP, localFilePath string, remotePath string) error {
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

func (job *SyncDataToDELAll) createRdy(c ctx.Context, dirPath, timeStr string, client *sftp.SFTP) error {
	localRdyName := timeStr + `.rdy`
	rdyPath := `./storage/new_cdp_del/` + localRdyName
	t, _ := os.Create(rdyPath) // 创建验证上传OK的文件
	err := job.uploadFile(c, client, rdyPath, dirPath)
	err = t.Close()
	return err
}
