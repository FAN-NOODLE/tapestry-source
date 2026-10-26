package plan_tasks

import (
	"encoding/csv"
	"fmt"
	sftp2 "github.com/pkg/sftp"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io"
	membersDao "micro-scrm/internal/app/dao/members"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"path"
	"strconv"
	"time"
)

type SyncTaskInfoToAcxiom struct {
	Name       string
	startTime  int64
	endTime    int64
	sftpClient *sftp.SFTP
	remotePath string
}

func NewTaskInfoToAcxiom() SyncTaskInfoToAcxiom {
	startTime, endTime := helper.YesterdayPeriodOfTimestamp(1)
	return SyncTaskInfoToAcxiom{
		Name:       "进行中的已联系的任务数据推送给acxiom",
		startTime:  startTime,
		endTime:    endTime,
		remotePath: "in/task",
	}
}

func (job SyncTaskInfoToAcxiom) Run(c ctx.Context, handler *xjob.JobHandler) error {

	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}

	for _, v := range enterpriseList {
		enterpriseId := v.ID
		enterpriseName := v.Name
		err = job.handle(c, enterpriseId, enterpriseName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (job SyncTaskInfoToAcxiom) handle(c ctx.Context, enterpriseId int64, enterpriseName string) error {
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
	// nowTime := time.Now().Unix()
	// 查询进行中的子任务，并且是1V1沟通的任务
	taskLists, err := planDao.NewPlanTaskDao().GetMany(c, enterpriseId, "task_end_time >=? and task_start_time < ?  and task_category=? and task_is_push=? and deleted_at=0", job.startTime, job.endTime, microScrmModel.PLANTASKCONTACT, 1)
	if err != nil {
		return err
	}

	memberList, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"deleted_at":    0,
	})
	if err != nil {
		return err
	}
	memberNameMap := make(map[string]string, 0)
	for _, info := range memberList {
		memberNameMap[info.Userid] = info.Name
	}
	filePrefix := "plan_task_"
	// 创建CSV文件
	localFilePath := `./storage/acxiom/`
	LocalFileName := filePrefix + time.Now().Format(helper.YYYYMMDD) + ".csv"
	FilePath := localFilePath + LocalFileName
	f, err := os.Create(FilePath) // 创建文件
	if err != nil {
		log.Info(c).Msgf("%s::csv文件创建失败: %v", job.Name, err)
		return err
	}

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	data := [][]string{
		{"唯一ID", "主任务ID", "任务名称", "子任务ID", "任务开始时间", "任务结束时间", "子任务开始时间", "子任务结束时间", "子任务创建时间", "任务创建方式", "触达导购", "导购ID", "客户ID(unionid)", "customer ID", "子任务是否完成", "完成时间", "顾客联系状态", "联系状态变更时间"},
	}
	w := csv.NewWriter(f) // 创建一个新写入文件流
	// yesterdayZeroTime, yesterdayEndTime := helper.YesterdayPeriodOfTimestamp(1)
	yesterdayZeroTime, yesterdayEndTime := job.startTime, job.endTime
	// taskPlanIdMap := make(map[int64][]microScrmModel.PlanTasksModel, 0)
	planInfoMap := make(map[int64]struct {
		category  int
		startTime int
		endTime   int
	}, 0)
	planIds := make([]int64, 0)
	for _, task := range taskLists {
		if !helper.InArray(task.PlanID, planIds) {
			planIds = append(planIds, task.PlanID)
		}
	}
	planLists, err := planDao.NewPlanTaskDao().GetPlanMany(c, enterpriseId, map[string]interface{}{
		"id": planIds,
	})
	if err != nil {
		return err
	}
	for _, plan := range planLists {
		if _, ok := planInfoMap[plan.ID]; !ok {
			planInfoMap[plan.ID] = struct {
				category  int
				startTime int
				endTime   int
			}{category: plan.Category, startTime: plan.StartTime, endTime: plan.EndTime}
		}
	}
	for _, task := range taskLists {
		planInfo := planInfoMap[task.PlanID]
		planCreateType := helper.RuleThree(planInfo.category == microScrmModel.PlanMaTypeYes, "MA模板", "手动创建")
		planStartTime, planEndTime := helper.GetFormatDateByParams(int64(planInfo.startTime)), helper.GetFormatDateByParams(int64(planInfo.endTime))
		taskAllotMap := make(map[int64]string, 0)
		startTime, endTime := helper.GetFormatDateByParams(task.TaskStartTime), helper.GetFormatDateByParams(task.TaskEndTime)
		createTime := helper.GetFormatDateByParams(task.CreatedAt)
		planId := task.PlanID
		taskId := task.ID
		taskAllotMany, err := planDao.NewPlanTaskAllotLogDao().GetMany(c, enterpriseId, map[string]interface{}{
			"plan_id":     task.PlanID,
			"task_id":     task.ID,
			"is_finished": microScrmModel.PLAN_FINISHED,
		})
		if err != nil {
			return err
		}
		// conditionList, _ := planDao.NewPlanTaskConditionDataDao().GetMany(c, enterpriseId, plan.ID, `finished_at BETWEEN ? AND ? AND is_finished=? AND task_id in ?`, yesterdayZeroTime, yesterdayEndTime, microScrmModel.PLAN_FINISHED, taskIds)
		conditionList, _ := planDao.NewPlanTaskConditionDataDao().GetMany(c, enterpriseId, planId, `updated_at BETWEEN ? AND ? AND task_id = ?`, yesterdayZeroTime, yesterdayEndTime, taskId)
		if len(conditionList) == 0 {
			continue
		}
		for _, v := range taskAllotMany {
			if _, ok := taskAllotMap[v.ID]; !ok {
				taskAllotMap[v.ID] = helper.GetFormatDateByParams(int64(v.EndTime))
			}
		}
		for _, item := range conditionList {
			uniqueId := fmt.Sprintf("%d_%d", planId, item.ID)
			memberName := item.UserId
			if value, ok := memberNameMap[item.UserId]; ok {
				memberName = value
			}
			taskFinishStr := "否"
			taskFinishTime := ""
			key := item.TaskAllotLogId
			if value, isOk := taskAllotMap[key]; isOk {
				taskFinishStr = "是"
				taskFinishTime = value
			}
			finishStr := "否"
			finishTime := ""
			if item.IsFinished == microScrmModel.PLAN_FINISHED {
				finishStr = "是"
				finishTime = helper.GetFormatDateByParams(int64(item.FinishedAt))
			}

			temp := []string{
				uniqueId, strconv.Itoa(int(planId)), task.Name, strconv.Itoa(int(taskId)), planStartTime, planEndTime, startTime, endTime, createTime, planCreateType, memberName, item.UserId, item.UnionId, item.CustomerID, taskFinishStr, taskFinishTime, finishStr, finishTime,
			}
			data = append(data, temp)
		}
	}
	err = w.WriteAll(data)
	if err != nil {
		return err
	} // 写入数据
	w.Flush()
	fmt.Printf("%s::%d的企业csv任务完成数据写入成功"+FilePath, enterpriseName, enterpriseId)
	// todo 文件上传值sftp,并发送消息通知
	err = job.uploadFile(c, job.sftpClient, FilePath, job.remotePath)
	err = os.Remove(FilePath)
	return err
}

func (job *SyncTaskInfoToAcxiom) uploadFile(c ctx.Context, sftp *sftp.SFTP, localFilePath string, remotePath string) error {
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
