package excitation

import (
	"bytes"
	"encoding/csv"
	"fmt"
	sftp2 "github.com/pkg/sftp"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io"
	"io/ioutil"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	reserveDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SyncWeLearning struct {
	Name string
}

func NewSyncWeLearning() SyncWeLearning {
	return SyncWeLearning{Name: "同步weLearning文件"}
}

func (s SyncWeLearning) Run(c ctx.Context, handler *xjob.JobHandler) error {
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

func (s SyncWeLearning) handle(c ctx.Context, enterpriseId int64) error {
	info, err := microWorkDao.NewStoreSftpConfig().GetOne(c, map[string]interface{}{
		"brand":      "welearning",
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
		log.Info(c).Msgf("%s::SFTP连接失败: %v", config, err)
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
	fileInfoList, err := reserveDao.NewReserveOrderDataAcxiomSftpLogDao().GetExistsFile(c, enterpriseId, 1, "welearning")
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
					err := s.pullLogData(c, enterpriseId, config, fileModel, fileName, info.Brand)
					if err != nil {
						return err
					}
				}
			} else {
				err := s.pullLogData(c, enterpriseId, config, microScrmModel.ReserveOrderDataAcxiomSftpLogModel{}, fileName, info.Brand)
				if err != nil {
					return err
				}
			}
		} else {
			err := s.pullLogData(c, enterpriseId, config, microScrmModel.ReserveOrderDataAcxiomSftpLogModel{}, fileName, info.Brand)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (s SyncWeLearning) pullLogData(ctx ctx.Context, enterpriseId int64, config sftp.SSHConfig, info microScrmModel.ReserveOrderDataAcxiomSftpLogModel, fileName, brandName string) error {
	client, err := sftp.NewSFTPClient(config)
	if err != nil {
		log.Info(ctx).Msgf("%s::SFTP连接失败: %v", fileName, err)
		return err
	}
	sftpFile, err := client.Client.Open(fileName)
	if err != nil {
		log.Info(ctx).Msgf("%s::读取SFTP目录下的文件失败: %v", fileName, err)
		return err
	}
	defer func(Client *sftp2.Client) {
		_ = Client.Close()
	}(client.Client)
	LocalFileName := strconv.FormatInt(app.UID(), 10) + `.csv`
	LocalFilePath := `./storage/acxiom/` + LocalFileName
	dstFile, err := os.Create(LocalFilePath)
	if err != nil {
		log.Info(ctx).Msgf("%s::创建文件失败: %v", LocalFilePath, err)
		return err
	}
	if _, err := sftpFile.WriteTo(dstFile); err != nil {
		log.Info(ctx).Msgf("%s::写入文件失败: %v", LocalFilePath, err)
		return err
	}
	databyte, _ := ioutil.ReadFile(LocalFilePath)
	logId := info.ID
	if logId == 0 {
		// 上传cos
		cs := storage.Init(cos.COS)
		path, err := cs.Upload(ctx, LocalFilePath, "third_sftp/we_learning/"+LocalFileName)
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
			PushType:  "welearning",
		}
		err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().CreateLog(ctx, insertData)
		if err != nil {
			log.Info(ctx).Msgf("创建welearning报错:%v", err)
		}
	}
	// cs2k的csv文件不是以\r\n换行，\n的换行，以\r的换行，csv包处理不了，故换成bufio处理，以\r换行处理
	// reader := bufio.NewReader(bytes.NewBuffer(databyte))
	reader := csv.NewReader(bytes.NewBuffer(databyte))
	// 存在本地数据
	var fileChan = make(chan string, 2000)
	go func(reader *csv.Reader) {
		var flag = true
		for {
			row, errs := reader.Read()
			if flag {
				flag = false
				continue
			}
			if errs != nil && errs != io.EOF {
				fmt.Println(errs)
				continue

			}
			if errs == io.EOF {
				close(fileChan)
				break
			}
			fileChan <- strings.Join(row, ",")
		}
	}(reader)
	MemberLists, err := membersDao.NewMembers().GetMany(ctx, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"deleted_at":    0,
	})
	if err != nil {
		return err
	}
	MemberMap := make(map[string]int64, 0)
	for _, list := range MemberLists {
		info := list
		MemberMap[info.Userid] = info.ID
	}
	err, departmentLists := departmentsDao.NewDepartments().GetAllDepartments(ctx, enterpriseId)
	if err != nil {
		return err
	}
	departmentMap := make(map[string]int64, 0)
	for _, item := range departmentLists {
		value := item
		departmentMap[value.Name] = value.ID
	}

	fileNameStrSlice := strings.Split(fileName, "_")
	if len(fileNameStrSlice) < 3 {
		return nil
	}
	if strings.Contains(fileName, "welearning_data_") {
		s.handlingWeLearningData(ctx, fileChan, enterpriseId, MemberMap, departmentMap)
	} else {
		var fileType string
		var pMonth int
		var fyYear int
		re := regexp.MustCompile(`_data_(\d+)_\d+\.csv`)
		match := re.FindStringSubmatch(fileName)
		if len(match) > 1 {
			//获取文件相关财月作为查询财月的时间条件
			tmpTime, _ := time.ParseInLocation("200601", match[1], time.Local)
			fyYear = tmpTime.Year()
			pMonth = int(tmpTime.Month())
		}
		if strings.Contains(fileName, "welearning_studycourseslist_data") {
			fileType = excitationModels.MedalStudyCoursesList
		} else if strings.Contains(fileName, "welearning_studytimelist_data") {
			fileType = excitationModels.MedalStudyTimeList
		} else if strings.Contains(fileName, "welearning_newcourseslist_data") {
			fileType = excitationModels.MedalNewCoursesList
		}
		if fileType != "" {
			s.handlingWeLearningRank(ctx, fileChan, fileType, fyYear, pMonth, enterpriseId, MemberMap, departmentMap)
		}
	}

	err = reserveDao.NewReserveOrderDataAcxiomSftpLogDao().UpdateLogStatus(ctx, enterpriseId, logId)
	if err != nil {
		return err
	}
	planService.NewTaskNotice().SftpSendNotice(ctx, planService.SFTPMessageData{
		EnterpriseName: brandName,
		SFTPName:       fmt.Sprintf("SFTP同步Welearning文件:%s", fileName),
		UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
	})
	return nil
}

func (s SyncWeLearning) handlingWeLearningData(ctx ctx.Context, fileChan <-chan string, enterpriseId int64, memberMap map[string]int64, departmentMap map[string]int64) {
	salesDao := reserveDao.NewMemberLearnRecord()
	for v := range fileChan {
		// var salesInfo microScrmModel.SftpSalesInfo
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) >= 5 {
			saID := strings.TrimSpace(splitInfo[1])
			storeId := helper.ConvertToString(strings.TrimSpace(splitInfo[2]))
			var memberId, departmentId int64
			if value, ok := memberMap[saID]; ok {
				memberId = value
			}
			if item, ok := departmentMap[storeId]; ok {
				departmentId = item
			}
			startTime := helper.GetDateToTimes(splitInfo[0])
			nowTime := time.Now().Unix()
			titleRemark := helper.ConvertToString(strings.TrimSpace(splitInfo[3]))
			points, _ := strconv.Atoi(splitInfo[4])
			memberLearn := microScrmModel.MemberLearnRecordModel{
				MemberID:     memberId,
				DepartmentID: departmentId,
				StartTime:    startTime,
				StoreCode:    storeId,
				TitleRemark:  titleRemark,
				Points:       int64(points),
				CreatedAt:    nowTime,
				UpdatedAt:    nowTime,
			}
			_ = salesDao.CreateInfo(ctx, enterpriseId, memberLearn)
			// 积分队列埋点
			excitationSvc.NewPointsProduceService().ProduceQueue(ctx, excitationReq.PointsProduceSave{
				MemberId:     memberId,
				Type:         excitationModels.RuleCategoryId15,
				Value:        int64(points),
				UniqueHash:   fmt.Sprintf("%d-%d-%d", app.UIDContext(ctx), memberId, excitationModels.RuleCategoryId15),
				EnterpriseId: enterpriseId,
				SourceRemark: titleRemark,
				AddTime:      startTime,
			})
		}
	}
}

func (s SyncWeLearning) handlingWeLearningRank(ctx ctx.Context, fileChan <-chan string, fileType string, fyYear int, pMonth int, enterpriseId int64, memberMap map[string]int64, departmentMap map[string]int64) {
	rankDao := reserveDao.NewMemberLearnRank()
	var i int
	var insertList []microScrmModel.MemberLearnRankModel
	for v := range fileChan {
		i++
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) >= 3 {
			saID := strings.TrimSpace(splitInfo[0])
			storeId := helper.ConvertToString(strings.TrimSpace(splitInfo[1]))
			var memberId, departmentId int64
			if value, ok := memberMap[saID]; ok {
				memberId = value
			}
			if item, ok := departmentMap[storeId]; ok {
				departmentId = item
			}
			nowTime := time.Now().Unix()
			var remark string
			if len(splitInfo) == 4 {
				remark = strings.TrimSpace(splitInfo[3])
			}
			tmpInfo := microScrmModel.MemberLearnRankModel{
				MemberID:     memberId,
				DepartmentID: departmentId,
				FileType:     fileType,
				PMonth:       pMonth,
				FyYear:       fyYear,
				UserID:       saID,
				StoreCode:    storeId,
				Value:        splitInfo[2],
				Rank:         i,
				Remark:       remark,
				CreatedAt:    nowTime,
				UpdatedAt:    nowTime,
			}
			insertList = append(insertList, tmpInfo)
		}
	}
	_ = rankDao.BatchCreate(ctx, enterpriseId, insertList)
}
