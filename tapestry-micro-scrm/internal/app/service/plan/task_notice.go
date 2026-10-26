package planService

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	internalHelper "micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/notice"
	"strings"
	"time"
)

type TaskNotice struct {
	planTaskDao *planDao.PlanTaskDao
}

func NewTaskNotice() *TaskNotice {
	return &TaskNotice{
		planTaskDao: planDao.NewPlanTaskDao(),
	}
}

type TaskMessageData struct {
	EnterpriseId      int64  `json:"enterprise_id"`
	PlanId            int64  `json:"plan_id"`
	TaskId            int64  `json:"task_id"`
	TaskTitle         string `json:"task_title"`
	SendUserCount     int    `json:"send_user_count"`
	TaskCategory      int    `json:"task_category"`
	Dimension         int    `json:"dimension"`
	RealSendUserCount int    `json:"real_send_user_count"`
	StartTime         int    `json:"start_time"`
	EndTime           int    `json:"end_time"`
}

type SFTPMessageData struct {
	EnterpriseName string `json:"enterprise_name" comment:"企业名称"`
	SFTPName       string `json:"sftp_name" comment:"sftp名称"`
	UploadRes      string `json:"upload_res" comment:"上传结果"`
	ToAll          bool   `json:"to_all" comment:"是否通知所有人"`
}

type MaterialUploadFailMessageData struct {
	EnterpriseId int64  `json:"enterprise_id"`
	PlanId       int64  `json:"plan_id"`
	TaskId       int64  `json:"task_id"`
	TaskTitle    string `json:"task_title"`
	TaskCategory int    `json:"task_category"`
	Dimension    int    `json:"dimension"`
	MaterialType string `json:"material_type"`
	MaterialUrl  string `json:"material_url"`
	ErrorMessage string `json:"error_message"`
}

// CreateTaskNotice 任务创建通知
func (p TaskNotice) CreateTaskNotice(c app.Ctx, data TaskMessageData) {
	enterpriseMap := make(map[int64]string, 6)
	enterpriseMap = map[int64]string{3360610327843968: "coach", 3374939108972352: "sw", 3383774945595776: "ks", 3363603667487424: "coach", 3373738709328576: "sw", 3376048056335040: "ks"}
	content := fmt.Sprintf("任务计划通知：\n"+
		">企业名称:<font color='comment'>%s</font>\n"+
		">主任务ID:<font color='comment'>%d</font>\n"+
		">子任务ID:<font color='comment'>%d</font>\n"+
		">任务类型:<font color='comment'>%s</font>\n"+
		">数据来源:<font color='comment'>%s</font>\n"+
		">子任务标题:<font color='comment'>%s</font>\n"+
		">任务开始时间：<font color='red'>%s</font>\n"+
		">任务结束时间：<font color='comment'>%s</font>\n"+
		">RequestId：<font color='comment'>%v</font>\n",
		enterpriseMap[data.EnterpriseId],
		data.PlanId,
		data.TaskId,
		internalHelper.RuleThree(data.TaskId == 0, "", p.getCategoryName(data.TaskCategory)),
		internalHelper.RuleThree(data.TaskId == 0, "", p.getDimensionStr(data.Dimension)),
		data.TaskTitle,
		time.Unix(int64(data.StartTime), 0).Format("2006-01-02 15:04:05"),
		time.Unix(int64(data.EndTime), 0).Format("2006-01-02 15:04:05"),
		c.ID(),
	)
	notice.NewWechat().SendTaskMarkdownMessage(c, content)
}

// TaskResult 任务结果获取
func (p TaskNotice) TaskResult(ctx app.Ctx, taskInfo microScrmModel.PlanTasksModel, enterpriseId int64, sendUserIds []string, realUserIds []string) {
	// 企业任务群发(群发助手)
	if taskInfo.GroupSendType == microScrmModel.PLANTASKQYSEND {
		var realTotal int64
		if taskInfo.TaskCategory == microScrmModel.PLANTASKMOMENT {
			// 朋友圈任务的群发类型
			realTotal = int64(len(internalHelper.ArrayUniq(realUserIds)))
		} else if taskInfo.TaskCategory == microScrmModel.PLANTASKCONTACT {
			// 沟通任务需要单个导购发群发助手
			conStr := "task_id = ? AND group_msg_id != '' AND deleted_at = 0"
			realTotal, _ = planDao.NewPlanTaskAllotLogDao().GetCount(ctx, enterpriseId, conStr, taskInfo.ID)
		}
		// 任务发送结果通知
		p.TaskSendNotice(ctx, TaskMessageData{
			EnterpriseId:      enterpriseId,
			TaskId:            taskInfo.ID,
			PlanId:            taskInfo.PlanID,
			TaskCategory:      taskInfo.TaskCategory,
			Dimension:         int(taskInfo.Dimension),
			TaskTitle:         taskInfo.Name,
			SendUserCount:     len(internalHelper.ArrayUniq(sendUserIds)),
			RealSendUserCount: int(realTotal),
		})
	}
}

// TaskSendNotice 任务发送结果通知
func (p TaskNotice) TaskSendNotice(c app.Ctx, data TaskMessageData) {
	enterpriseMap := make(map[int64]string, 6)
	enterpriseMap = map[int64]string{3360610327843968: "coach", 3374939108972352: "sw", 3383774945595776: "ks", 3363603667487424: "coach", 3373738709328576: "sw", 3376048056335040: "ks"}
	content := fmt.Sprintf("任务发送结果通知：\n"+
		">企业名称:<font color='comment'>%s</font>\n"+
		">主任务ID:<font color='comment'>%d</font>\n"+
		">子任务ID:<font color='comment'>%d</font>\n"+
		">任务类型:<font color='comment'>%s</font>\n"+
		">数据来源:<font color='comment'>%s</font>\n"+
		">子任务标题:<font color='comment'>%s</font>\n"+
		">应发送导购人数：<font color='comment'>%d</font>\n"+
		">实际发送导购人数：<font color='comment'>%d</font>\n"+
		">发送结果：<font color='red'>%s</font>\n"+
		">RequestId：<font color='comment'>%v</font>\n",
		enterpriseMap[data.EnterpriseId],
		data.PlanId,
		data.TaskId,
		p.getCategoryName(data.TaskCategory),
		p.getDimensionStr(data.Dimension),
		data.TaskTitle,
		data.SendUserCount,
		data.RealSendUserCount,
		internalHelper.RuleThree(data.SendUserCount == data.RealSendUserCount, "成功", "异常"),
		c.ID(),
	)
	notice.NewWechat().SendTaskMarkdownMessage(c, content)
}

func (p TaskNotice) TaskMaterialUploadFailNotice(c app.Ctx, data MaterialUploadFailMessageData) {
	enterpriseMap := make(map[int64]string, 6)
	enterpriseMap = map[int64]string{3360610327843968: "coach", 3374939108972352: "sw", 3383774945595776: "ks", 3363603667487424: "coach", 3373738709328576: "sw", 3376048056335040: "ks"}
	content := fmt.Sprintf("MA任务素材上传失败告警：\n"+
		"企业名称：%s\n"+
		"主任务ID：%d\n"+
		"子任务ID：%d\n"+
		"任务类型：%s\n"+
		"数据来源：%s\n"+
		"子任务标题：%s\n"+
		"素材类型：%s\n"+
		"素材地址：%s\n"+
		"错误信息：%s\n"+
		"处理结果：已标记子任务为下发失败，停止自动下发\n"+
		"RequestId：%v",
		enterpriseMap[data.EnterpriseId],
		data.PlanId,
		data.TaskId,
		p.getCategoryName(data.TaskCategory),
		p.getDimensionStr(data.Dimension),
		data.TaskTitle,
		data.MaterialType,
		data.MaterialUrl,
		data.ErrorMessage,
		c.ID(),
	)
	notice.NewWechat().SendTaskTextMessage(c, content)
}

func (p TaskNotice) getCategoryName(taskCategory int) string {
	var typeStr string
	if taskCategory == microScrmModel.PLANTASKMOMENT {
		typeStr = "朋友圈任务"
	} else if taskCategory == microScrmModel.PLANTASKCONTACT {
		typeStr = "沟通任务"
	}
	return typeStr
}

func (p TaskNotice) getDimensionStr(dimension int) string {
	var typeStr string
	if dimension == microScrmModel.TaskDimensionTag {
		typeStr = "用户画像"
	} else if dimension == microScrmModel.TaskDimensionImportUnionId || dimension == microScrmModel.TaskDimensionImportSaUnionId {
		typeStr = "人群包"
	}
	return typeStr
}

// SftpSendNotice  sftp文件通知
func (p TaskNotice) SftpSendNotice(c app.Ctx, data SFTPMessageData) {
	content := fmt.Sprintf("sftp上传结果通知：\n"+
		">企业名称:<font color='comment'>%s</font>\n"+
		">sftp名称:<font color='comment'>%s</font>\n"+
		">上传结果：<font color='red'>%s</font>\n"+
		">RequestId：<font color='comment'>%v</font>\n",
		data.EnterpriseName,
		data.SFTPName,
		data.UploadRes,
		c.ID(),
	)
	if strings.Contains(data.UploadRes, "失败") || strings.Contains(data.UploadRes, "异常") {
		notice.NewWechat().SendTaskTextMessage(c, internalHelper.StripHTML(content))
	} else {
		notice.NewWechat().SendTaskMarkdownMessage(c, content)
	}
}

// NfsUploadNotice nfs文件上传的通知
func (p TaskNotice) NfsUploadNotice(c ctx.Context, enterpriseId int64, fileName, cosUrl string, num int) {
	enterpriseMap := make(map[int64]string, 6)
	enterpriseMap = map[int64]string{3360610327843968: "coach", 3374939108972352: "sw", 3383774945595776: "ks", 3363603667487424: "coach", 3373738709328576: "sw", 3376048056335040: "ks"}
	content := fmt.Sprintf("DashBoard文件生成：\n"+
		">企业名称:<font color='comment'>%s</font>\n"+
		">文件名称:<font color='comment'>%s</font>\n"+
		">cos地址:<font color='comment'>%s</font>\n"+
		">文件数量:<font color='comment'>%d</font>\n"+
		">RequestId：<font color='comment'>%v</font>\n",
		enterpriseMap[enterpriseId],
		fileName,
		cosUrl,
		num,
		c.ID(),
	)
	notice.NewWechat().SendTaskMarkdownMessage(c, content)
}

// SftpPullNotice  sftp文件获取通知
func (p TaskNotice) SftpPullNotice(c app.Ctx, data SFTPMessageData) {
	content := fmt.Sprintf("sftp获取文件结果通知：\n"+
		">企业名称:<font color='comment'>%s</font>\n"+
		">sftp名称:<font color='comment'>%s</font>\n"+
		">获取结果：<font color='red'>%s</font>\n"+
		">RequestId：<font color='comment'>%v</font>\n",
		data.EnterpriseName,
		data.SFTPName,
		data.UploadRes,
		c.ID(),
	)
	if data.ToAll {
		notice.NewWechat().SendTaskTextMessage(c, internalHelper.StripHTML(content))
	} else {
		notice.NewWechat().SendTaskMarkdownMessage(c, content)
	}
}
