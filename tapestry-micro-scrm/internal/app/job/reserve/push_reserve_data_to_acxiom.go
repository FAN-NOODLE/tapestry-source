package reserveJob

import (
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io/ioutil"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	reserveDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	"micro-scrm/internal/app/models"
	departmentsModels "micro-scrm/internal/app/models/departments"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	planService "micro-scrm/internal/app/service/plan"
	reserveSvc "micro-scrm/internal/app/service/reserve"
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
)

type pushReserveDataToAcxiom struct {
	name string
}

func NewPushReserveDataToAcxiom() *pushReserveDataToAcxiom {
	return &pushReserveDataToAcxiom{
		name: "推送昨天的未完成，邀约完成的邀约单给到安客诚（acxiom）",
	}
}

func (job *pushReserveDataToAcxiom) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

	sftpConfigList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList(micro_scrm.SftpTypeAcxiom)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	for _, config := range sftpConfigList {

		if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
			log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
			// return
		} else {
			job.pushLogData(ctx, config.EnterpriseID, client, config)
		}

	}

	return nil

	// var (
	//	err            error
	//	enterpriseList []boss.EnterpriseModel
	// )
	// enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	// if err != nil {
	//	return err
	// }
	// var wg sync.WaitGroup
	// for _, enterprise := range enterpriseList {
	//	wg.Add(1)
	//	go func(enterpriseID int64) {
	//		defer wg.Done()
	//		job.getTask(ctx, enterpriseID)
	//	}(enterprise.ID)
	// }
	// wg.Wait()
	// return nil

}

func (job *pushReserveDataToAcxiom) pushLogData(ctx ctx.Context, enterpriseId int64, sftp *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail) error {
	// 判断当前企业是否为coach
	if len(config.PathConfig.Acxiom) == 0 {
		log.Info(ctx).Msgf("%s::安客诚的sftp信息未配置", job.name)
		return nil
	}

	zeroTime := helper.GetDateMorningTimestamp(helper.GetNowDate())
	info, err := reserveDao.NewReserveOrderDataAcxiomSftpLogDao().GetInfo(ctx, enterpriseId, 0, zeroTime, "push_day_order")
	if err != nil {
		return err
	}
	if info.Status == 1 {
		log.Info(ctx).Msgf("%s::读取SFTP文件已处理完毕", job.name)
		return nil
	}
	enterpriseConfig, err := cacheSvc.NewRedisData().GetEnterpriseConfig(ctx, enterpriseId, "store_name")
	if err != nil {
		return err
	}
	startTime, endTime := helper.YesterdayPeriodOfTimestamp(1)
	var list []microScrmModel.ReserveOrderModel
	if enterpriseConfig.Name == "coach" {
		list, err = reserveSvc.NewOrderService().GetOrderListsByStatus(ctx, enterpriseId, startTime, endTime)
	} else {
		list, err = reserveSvc.NewOrderService().GetReserveTimeOrderLists(ctx, enterpriseId, startTime, endTime)
	}
	if len(list) == 0 {
		planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
			EnterpriseName: enterpriseConfig.Name,
			SFTPName:       "sftp邀约单上传至cdp",
			UploadRes:      "今天没有邀约单",
		})
		return err
	}
	listReq := reserveRequest.ReserveServiceListReq{}
	serviceList, _, err := reserveDao.NewReserveServiceDao().GetListAll(ctx, enterpriseId, listReq)
	if err != nil {
		return err
	}
	serviceListMap := make(map[int64]microScrmModel.ReserveServiceModel, 0)
	for _, item := range serviceList {
		serviceListMap[item.ID] = item
	}
	var filePrefix string
	departmentsMap := make(map[int]string)
	if enterpriseConfig.Name == "coach" {
		// 获取所有部门，departments表中name为门店code
		departments, err := departmentsDao.NewDepartments().GetAll(ctx, departmentsModels.DepartmentsModel{
			MicroBaseModel: models.MicroBaseModel{
				EnterpriseID: enterpriseId,
			}})
		if err != nil {
			return err
		}
		for _, department := range departments {
			departmentsMap[department.ThirdID] = department.Name
		}
		filePrefix = "invitation_order_"
	} else {
		// 获取所有部门，departments表中name为门店code
		departments, err := microWorkDao.NewStoreDepartment().GetMany(ctx, map[string]interface{}{
			"enterprise_id": enterpriseId,
			"deleted_at":    0,
		})
		if err != nil {
			return err
		}
		for _, department := range departments {
			departmentsMap[int(department.WechatDepartmentID)] = department.ThirdStoreCode
		}
		filePrefix = "invitation_sa_"
	}
	localFilePath := `./storage/acxiom/`
	LocalFileName := filePrefix + helper.GetFormatDate(0, 2) + ".csv"
	FilePath := localFilePath + LocalFileName
	f, err := os.Create(FilePath) // 创建文件
	if err != nil {
		log.Info(ctx).Msgf("%s::csv文件创建失败: %v", job.name, err)
		return err
	}

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	data := [][]string{
		{"最新邀约单时间", "邀约单号", "邀约单状态,(0-进行中,2-确认用户到店,3-预约单取消)", "顾客unionid", "到店时间", "邀约导购ID", "预约服务名称", "邀约门店code"},
	}
	swSaCodeMap := make(map[string]string)
	if enterpriseConfig.Name == "sw" {
		swMembers, _ := membersDao.NewMembers().GetMany(ctx, "enterprise_id= ?", enterpriseId)
		if len(swMembers) > 0 {
			for _, member := range swMembers {
				if member.OpenUserid != "" {
					swSaCodeMap[member.Userid] = member.OpenUserid
				}
			}
		}
	}

	w := csv.NewWriter(f) // 创建一个新写入文件流
	for _, v := range list {
		var updateTimeStr, reserveTimeStr string
		if enterpriseConfig.Name == "coach" {
			updateTimeStr = helper.GetFormatDateDay(v.CreatedAt)
			reserveTimeStr = helper.GetFormatDateDay(int64(v.ReserveTime))
		} else {
			updateTimeStr = helper.GetFormatDateByParams(v.CreatedAt)
			reserveTimeStr = helper.GetFormatDateByParams(int64(v.ReserveTime))
			if swSaCode, ok := swSaCodeMap[v.MemberUserID]; ok {
				v.MemberUserID = swSaCode
			}
		}

		serveName := serviceListMap[v.ReserveID].ReserveName
		storeCode := "101test" // 默认的storeCode为101test
		if item, ok := departmentsMap[int(v.StoreID)]; ok {
			storeCode = item
		}
		tmp := []string{
			updateTimeStr, v.ReserveNumber, strconv.FormatInt(int64(v.Status), 10), v.UserUnionid, reserveTimeStr, v.MemberUserID, serveName, storeCode,
		}
		data = append(data, tmp)
	}
	err = w.WriteAll(data)
	if err != nil {
		return err
	} // 写入数据
	w.Flush()
	fmt.Printf("%s::%d的企业csv日志文件写入成功"+FilePath+",日志数量:"+strconv.Itoa(len(list))+"\n", enterpriseConfig.Name, enterpriseId)
	err = job.uploadFile(sftp, FilePath, config.PathConfig.Acxiom[0])
	if err != nil {
		log.Info(ctx).Msgf("%s::sftp文件写入失败: %v", job.name, err)
		return err
	}
	log.Info(ctx).Msgf("%s::写入sftp文件成功: %s", job.name, FilePath)

	// 将csv文件上传直cos存储
	cs := storage.Init(cos.COS)
	cosPath, err := cs.Upload(ctx, FilePath, "acxiom/"+LocalFileName)
	if err != nil {
		log.Info(ctx).Msgf("uploadCloud Error: %v", err)
		return err
	}
	// 删除本地文件
	f.Close()
	defer os.Remove(FilePath)
	if err != nil {
		log.Info(ctx).Msgf("delete local file error：%s ,%v", FilePath, err)
	}
	fileUrl := cs.GetURL(cosPath)
	insertData := microScrmModel.ReserveOrderDataAcxiomSftpLogModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel:    models.CreateBase(app.UIDContext(ctx)),
			EnterpriseID: enterpriseId,
		},
		FileType:  0,
		LocalFile: LocalFileName,
		FileURL:   fileUrl,
		Status:    1,
		FileTime:  int(zeroTime),
	}
	err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().CreateLog(ctx, insertData)
	if err != nil {
		log.Info(ctx).Msgf("推送数据到acxiom记录日志失败:%v", err)
	}
	planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
		EnterpriseName: enterpriseConfig.Name,
		SFTPName:       "sftp邀约单上传至cdp",
		UploadRes:      "成功",
	})
	return err
}

func (job *pushReserveDataToAcxiom) uploadFile(sftp *sftp.SFTP, localFilePath string, remotePath string) error {
	srcFile, err := os.Open(localFilePath)
	defer srcFile.Close()
	if err != nil {
		return err
	}
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftp.Client.Create(path.Join(remotePath, remoteFileName))
	defer dstFile.Close()
	if err != nil {
		return err
	}
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return err
	}
	dstFile.Write(ff)
	return nil
}
