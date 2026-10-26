package reserveJob

import (
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io/ioutil"
	"micro-scrm/internal/app"
	customerDao "micro-scrm/internal/app/dao/customer"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	reserveDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/models"
	departmentsModels "micro-scrm/internal/app/models/departments"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
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

type ReportDataToAcxiom struct {
	name string
}

func NewReportDataToAcxiom() *ReportDataToAcxiom {
	return &ReportDataToAcxiom{
		name: "Monthly的出报告文件，每个财月结束后上传(Report_Data)",
	}
}

func (job *ReportDataToAcxiom) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	sftpConfigList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList(micro_scrm.SftpTypeAcxiom)
	if err != nil {
		log.Info(ctx).Msgf("%s::获取SFTP配置失败: %v", job.name, err)
	}

	for _, config := range sftpConfigList {

		if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
			log.Info(ctx).Msgf("%s::企业ID:%d, 连接SFTP失败: %v", job.name, config.EnterpriseID, clientErr)
			// return
		} else {
			job.pushData(ctx, config.EnterpriseID, client, config)
		}
	}

	return nil
}

// 下一个财月的第二天会上传上一个财月的数据
func (job *ReportDataToAcxiom) pushData(ctx ctx.Context, enterpriseId int64, sftp *sftp.SFTP, config micro_scrm.BrandConfigSftpDetail) error {
	if len(config.PathConfig.ChangYiReport) == 0 {
		return nil
	}
	// 根据当前时间的零点获取上个财月的信息,
	nowTime := helper.GetNowDateTimestamp()
	zeroTime := helper.GetZeroDateByParams(0)
	fiscalInfo, err := microWorkDao.NewFiscalMonth().GetFiscalMonth(ctx, enterpriseId, nowTime)
	if err != nil {
		return err
	}
	if fiscalInfo.StartTime == zeroTime { // 当天的0点和所在财月的开始时间一致，则不执行
		return nil
	}
	// if fiscalInfo.StartTime+(86400*7) != zeroTime {
	// 	log.Info(ctx).Msgf("%s::每个财月的第7天执行", job.name)
	// 	return nil
	// }
	fiscalPreMonthEndTime := fiscalInfo.StartTime - 1 // 算出上一个财月的结束时间
	fiscalPreInfo, err := microWorkDao.NewFiscalMonth().GetFiscalMonthByEndTime(ctx, enterpriseId, fiscalPreMonthEndTime)
	if err != nil {
		return err
	}
	startTime, endTime := fiscalPreInfo.StartTime, fiscalPreInfo.EndTime

	info, err := reserveDao.NewReserveOrderDataAcxiomSftpLogDao().GetInfo(ctx, enterpriseId, 0, startTime, "Report_Data")
	if err != nil {
		return err
	}
	if info.Status == 1 {
		log.Info(ctx).Msgf("%s::读取SFTP文件已处理完毕", job.name)
		return nil
	}

	list, err := reserveSvc.NewOrderService().GetStoreReportDataLists(ctx, enterpriseId, startTime, endTime)

	if err != nil {
		fmt.Printf("获取Monthly的Store成功邀约,报错%v\n", err)
		return err
	}

	remotePath := config.PathConfig.ChangYiReport[0] // IN/WeClient_Report_PROD/STOREAPPOINTMENTQTY/
	LocalFileName := "Report_Data_" + helper.GetFormatDate(0, 2) + `.csv`
	LocalTextName := "Report_Data_" + helper.GetFormatDate(0, 2) + `OK.txt`
	FilePath := `./storage/acxiom/` + LocalFileName
	TextPath := `./storage/acxiom/` + LocalTextName
	t, _ := os.Create(TextPath) // 创建验证上传OK的文件
	defer t.Close()

	f, err := os.Create(FilePath) // 创建文件
	if err != nil {
		log.Info(ctx).Msgf("%s::csv文件创建失败: %v", job.name, err)
		return err
	}

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	// weClientQTY 导购的好友数; SleepersQTY 沉睡数，不计数，默认0; NewClientQTY 导购的新增好友数; Active 活跃数，不计数，默认0;EviteSent 邀请函发送数;PlantSeed 播种数;OutreachQTY 互动顾客数
	data := [][]string{
		{"year", "period", "saAccount", "storeCode", "appointmentQTY", "weClientQTY", "SleepersQTY", "NewClientQTY", "Active", "EviteSent", "PlantSeed", "OutreachQTY"},
	}
	// 获取所有部门，departments表中name为门店code
	// departments, err := departmentsDao.NewDepartments().GetAll(ctx, departmentsModels.DepartmentsModel{
	//	MicroBaseModel: models.MicroBaseModel{
	//		EnterpriseID: enterpriseId,
	//	}})
	// if err != nil {
	//	return err
	// }
	// departmentsMap := make(map[int]string)
	// for _, department := range departments {
	//	departmentsMap[department.ThirdID] = department.Name
	// }
	w := csv.NewWriter(f) // 创建一个新写入文件流
	defer f.Close()
	// storeQTYMap := make(map[int64]int64)
	// storeWeclientQTYMap := make(map[int64]int64)
	// storeNewClientQTYMap := make(map[int64]int64)
	// storeEviteSentMap := make(map[int64]int64)
	// storePlantSeedMap := make(map[int64]int64)
	// storeChatCntMap := make(map[int64]int64)
	// for _, model := range list {
	//	if _, ok := storeQTYMap[model.StoreId]; !ok {
	//		storeQTYMap[model.StoreId] = 0
	//	}
	//	if _, ok := storeWeclientQTYMap[model.StoreId]; !ok {
	//		storeWeclientQTYMap[model.StoreId] = 0
	//	}
	//	if _, ok := storeNewClientQTYMap[model.StoreId]; !ok {
	//		storeNewClientQTYMap[model.StoreId] = 0
	//	}
	//	if _, ok := storeEviteSentMap[model.StoreId]; !ok {
	//		storeEviteSentMap[model.StoreId] = 0
	//	}
	//	if _, ok := storePlantSeedMap[model.StoreId]; !ok {
	//		storePlantSeedMap[model.StoreId] = 0
	//	}
	//	if _, ok := storeChatCntMap[model.StoreId]; !ok {
	//		storeChatCntMap[model.StoreId] = 0
	//	}
	// }
	// storeMemberIds := make(map[int64][]int64)

	// for k, _ := range storeQTYMap {
	//	var memberIds []int64
	//	memberList, _ := membersDao.NewMembers().GetMemberByThirdIds(ctx, enterpriseId, []int{int(k)}, []int64{})
	//	for _, member := range memberList {
	//		memberIds = append(memberIds, member.MembersId)
	//	}
	//	storeMemberIds[k] = memberIds
	// }
	// for k, _ := range storeWeclientQTYMap {
	//	totalNum, err := reserveDao.NewUserRelation().GetALLFriendsCount(ctx, enterpriseId, endTime, storeMemberIds[k])
	//	if err != nil {
	//		totalNum = 0
	//	}
	//	storeWeclientQTYMap[k] = totalNum
	// }
	//
	// //门店的新增好友数统计
	// for k, _ := range storeNewClientQTYMap {
	//	totalNum, err := reserveDao.NewUserRelation().GetAddFriendsCount(ctx, enterpriseId, startTime, endTime, storeMemberIds[k])
	//	if err != nil {
	//		totalNum = 0
	//	}
	//	storeNewClientQTYMap[k] = totalNum
	// }
	// //邀请函发送数
	// for k, _ := range storeEviteSentMap {
	//	totalNum, err := reserveDao.NewInviteRecordDao().GetRecordNumByMemberIds(ctx, enterpriseId, startTime, endTime, storeMemberIds[k])
	//	if err != nil {
	//		totalNum = 0
	//	}
	//	storeEviteSentMap[k] = totalNum
	// }
	// //播种数
	// for k, _ := range storePlantSeedMap {
	//	totalNum, err := customerDao.NewSowTask().GetSowTaskCountByMemberIds(ctx, enterpriseId, startTime, endTime, storeMemberIds[k])
	//	if err != nil {
	//		totalNum = 0
	//	}
	//	storePlantSeedMap[k] = totalNum
	// }
	// //互动数
	// for k, _ := range storeChatCntMap {
	//	totalNum, err := membersDao.NewMemberExternalBehaviorData().GetChatCntCounts(ctx, enterpriseId, startTime, endTime, storeMemberIds[k])
	//	if err != nil {
	//		totalNum = 0
	//	}
	//	storeChatCntMap[k] = totalNum
	// }
	departments, err := departmentsDao.NewDepartments().GetAll(ctx, departmentsModels.DepartmentsModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel: models.BaseModel{
				DeletedAt: 0,
			},
			EnterpriseID: enterpriseId,
		},
		ThirdParentID: 7,
	})
	if err != nil {
		return err
	}
	departmentsMap := make(map[int]string)
	for _, department := range departments {
		departmentsMap[department.ThirdID] = department.Name
	}
	period := strconv.FormatInt(fiscalPreInfo.Period, 10)
	year := strconv.FormatInt(fiscalPreInfo.FiscalYear, 10)
	for _, v := range list {
		item := v
		weClientQTY := "0"
		allFriendsCount, err := reserveDao.NewUserRelation().GetALLFriendsCount(ctx, enterpriseId, endTime, []int64{item.MemberId})
		if err == nil {
			weClientQTY = strconv.FormatInt(allFriendsCount, 10)
		}

		newClientQTY := "0"
		addFriendsCount, err := reserveDao.NewUserRelation().GetAddFriendsCount(ctx, enterpriseId, startTime, endTime, []int64{item.MemberId})
		if err == nil {
			newClientQTY = strconv.FormatInt(addFriendsCount, 10)
		}
		eviteSent := "0"
		RecordNum, err := reserveDao.NewInviteRecordDao().GetRecordNumByMemberIds(ctx, enterpriseId, startTime, endTime, []int64{item.MemberId})
		if err == nil {
			eviteSent = strconv.FormatInt(RecordNum, 10)
		}
		plantSeed := "0"
		SowTaskCount, err := customerDao.NewSowTask().GetSowTaskCountByMemberIds(ctx, enterpriseId, startTime, endTime, []int64{item.MemberId})
		if err == nil {
			plantSeed = strconv.FormatInt(SowTaskCount, 10)
		}
		chatCnt := "0"
		ChatCntCount, err := membersDao.NewMemberExternalBehaviorData().GetChatCntCounts(ctx, enterpriseId, startTime, endTime, []int64{item.MemberId})
		if err == nil {
			chatCnt = strconv.FormatInt(ChatCntCount, 10)
		}
		// weClientQTY,_ :=membersDao.NewMembers().GetMemberByThirdIds(ctx, enterpriseId, []int{}, )

		// weClientQTY := strconv.FormatInt(storeWeclientQTYMap[v.StoreId], 10)
		// newClientQTY := strconv.FormatInt(storeNewClientQTYMap[v.StoreId], 10)
		// eviteSent := strconv.FormatInt(storeEviteSentMap[v.StoreId], 10)
		// plantSeed := strconv.FormatInt(storePlantSeedMap[v.StoreId], 10)
		// chatCnt := strconv.FormatInt(storeChatCntMap[v.StoreId], 10)
		appointMentQTY := strconv.FormatInt(item.AppointmentQTY, 10)
		tmp := []string{
			year, period, item.MemberUserId, departmentsMap[int(item.StoreId)], appointMentQTY, weClientQTY, "0", newClientQTY, "0", eviteSent, plantSeed, chatCnt,
		}
		data = append(data, tmp)
	}

	err = w.WriteAll(data)
	if err != nil {
		return err
	} // 写入数据
	w.Flush()
	fmt.Println("csv日志文件写入成功" + FilePath + ",日志数量:" + strconv.Itoa(len(list)))
	err = job.uploadFile(sftp, FilePath, remotePath)
	if err != nil {
		log.Info(ctx).Msgf("%s::sftp文件写入失败: %v", job.name, err)
		return err
	}
	err = job.uploadFile(sftp, TextPath, remotePath)
	if err != nil {
		log.Info(ctx).Msgf("%s::sftp验证OK文件写入失败: %v", job.name, err)
		return err
	}
	// 将csv文件上传直cos存储
	cs := storage.Init(cos.COS)
	cosPath, err := cs.Upload(ctx, FilePath, "acxiom/"+LocalFileName)
	if err != nil {
		log.Info(ctx).Msgf("uploadCloud Error: %v", err)
		return err
	}
	_ = os.Remove(TextPath)
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
		FileTime:  int(startTime),
		PushType:  "Report_Data",
		Year:      fiscalPreInfo.FiscalYear,
		Period:    fiscalPreInfo.Period,
	}
	err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().CreateLog(ctx, insertData)

	if err != nil {
		log.Info(ctx).Msgf("推送数据到acxiom记录日志失败:%v", err)
	}
	err = os.Remove(FilePath)
	if err != nil {
		log.Info(ctx).Msgf("delete local file error：%s ,%v", FilePath, err)
	}
	planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
		EnterpriseName: config.BrandName,
		SFTPName:       "SFTP同步财月report_data数据",
		UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
	})
	return nil
}

func (job *ReportDataToAcxiom) uploadFile(sftp *sftp.SFTP, localFilePath string, remotePath string) error {
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
