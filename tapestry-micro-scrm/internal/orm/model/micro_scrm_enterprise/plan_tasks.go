package micro_scrm_enterprise

import (
	"fmt"

	"micro-scrm/internal/orm/model"
)

type PlanTaskModel struct {
	model.BranchModel
	PlanId           int64                    `gorm:"column:plan_id;default:0;NOT NULL" json:"plan_id"`               // 子任务对应的计划的ID
	Name             string                   `gorm:"column:name;NOT NULL" json:"name"`                               // 子任务名称
	Dimension        int                      `gorm:"column:dimension;default:1;NOT NULL" json:"dimension"`           // 子任务的维度；1社群数据；2用户画像
	ConditionType    int                      `gorm:"column:condition_type;default:1;NOT NULL" json:"condition_type"` // 子任务的条件类型；1标签；2群人数；3群内新增任务；4参与活动；5购买行为；6群活跃度；7退群名单
	IsSpecifyCrowd   int                      `gorm:"column:is_specify_crowd;default:0" json:"is_specify_crowd"`      // 任务是否指定人群 1 否 2 是
	IsOpen           int                      `gorm:"column:is_open;default:1;NOT NULL" json:"is_open"`               // 任务在小程序中是否可见。0不可见，1可见；默认可见
	Status           int                      `gorm:"column:status;default:0;NOT NULL" json:"status"`                 // 任务状态；0未开始、1进行中、2已完成
	TaskType         int                      `gorm:"column:task_type" json:"task_type"`                              // 任务周期 1 单次 2持续
	SendType         int                      `gorm:"column:send_type" json:"send_type"`                              // 任务推送对象 1 BA 2管理bA
	TaskIsPush       int                      `gorm:"column:task_is_push;default:0" json:"task_is_push"`              // 任务是否执行 0 未执行 1 已执行
	Remark           string                   `gorm:"column:remark;NOT NULL" json:"remark"`                           // 子任务描述
	ConditionList    []PlanTaskConditionModel `gorm:"foreignKey:TaskId;references:ID"`
	BirthdayPreTime  int                      `gorm:"column:birthday_pre_time;default:0;NOT NULL" json:"birthday_pre_time"` // 生日消息推送的提前天数
	ConsumerFileId   int                      `gorm:"column:consumer_file_id;default:0;NOT NULL" json:"consumer_file_id"`
	IsExtract        int                      `gorm:"column:is_extract;type:tinyint(1);not null;default:0;comment:客户名单是否导入完成,0: 未完成，1:已完成"`
	CrowdID          string                   `json:"crowd_id" gorm:"column:crowd_id;type:varchar(200);not null;default:0;comment:MA任务的任务包批次ID"`
	PlanSendTime     int64                    `json:"plan_send_time" gorm:"column:plan_send_time;type:int(10);not null;default:0;comment:MA任务的推送时间"`
	DownloadUrl      string                   `json:"download_url" gorm:"column:download_url;type:varchar(200);not null;default:'';comment:MA人群包匹配COS的地址"`
	TaskCategory     int                      `json:"task_category" gorm:"column:task_category;type:tinyint(1);not null;default:0;comment:子任务的任务类型，沟通任务/朋友圈任务"`
	GroupSendType    int                      `json:"group_send_type" gorm:"column:group_send_type;type:tinyint(1);not null;default:0;comment:子任务群发类型，标签任务群发/企业任务群发"`
	IsAddTag         int                      `json:"is_add_tag" gorm:"column:is_add_tag;type:tinyint(1);not null;default:0;comment:是否生成任务标签"`
	TagExpireTime    int                      `json:"tag_expire_time" gorm:"column:tag_expire_time;type:int(10);not null;default:0;comment:标签任务的过期时间,0为永不过期"`
	TagName          string                   `json:"tag_name" gorm:"column:tag_name;type:varchar(200);not null;default:'';comment:子任务的打标签的标签名称"`
	QyTagId          string                   `json:"qy_tag_id" gorm:"column:qy_tag_id;type:varchar(200);not null;default:'';comment:企微的tag_id"`
	TagExpireType    int                      `json:"tag_expire_type" gorm:"column:tag_expire_type;type:tinyint(1);not null;default:0;comment:企微标签的过期类型,1任务结束后，2自定义删除时间，3长期有效"`
	MomentJobId      string                   `json:"moment_job_id" gorm:"column:moment_job_id;type:varchar(200);not null;default:'';comment:朋友圈的jobid"`
	MomentStatus     int                      `json:"moment_status" gorm:"column:moment_status;type:tinyint(1);not null;default:0;comment:朋友圈的任务创建状态，，1表示开始创建任务，2表示正在创建任务中，3表示创建任务已完成"`
	OriginTaskId     int64                    `json:"origin_task_id" gorm:"column:origin_task_id;type:bigint(20);not null;default:0;comment:针对持续任务一直创建的源子任务id"`
	FinishRate       int                      `json:"finish_rate" gorm:"column:finish_rate;type:tinyint(1);not null;default:100;comment:任务完成率"`
	TaskStartTime    int64                    `json:"task_start_time" gorm:"column:task_start_time;type:int(10);not null;default:0;comment:子任务的开始时间"`
	TaskEndTime      int64                    `json:"task_end_time" gorm:"column:task_end_time;type:int(10);not null;default:0;comment:子任务的结束时间"`
	TaskExpireType   int                      `json:"task_expire_type" gorm:"column:task_expire_type;type:tinyint(1);not null;default:0;comment:子任务有效期，0为根据主任务的有效期一致，2指定时间有效期"`
	IsPush           int32                    `json:"is_push" gorm:"column:is_push;type:tinyint(1);not null;default:0;comment:MA创建的任务推送完毕，值会更新，2为已推送完毕"`
	PlanTemplateId   int64                    `json:"plan_template_id" gorm:"column:plan_template_id;type:bigint(20);not null;default:0;comment:任务模板id"`
	FlowId           string                   `json:"flow_id" gorm:"column:flow_id;type:varchar(200);not null;default:'';comment:MA的自动流ID"`
	Campaign         string                   `json:"campaign" gorm:"column:campaign;type:varchar(200);not null;default:'';comment:MA的营销活动编码"`
	PlanName         string                   `json:"plan_name" gorm:"column:plan_name;type:varchar(200);not null;default:'';comment:MA推送的任务名称"`
	OriginType       int                      `json:"origin_type" gorm:"column:origin_type;type:tinyint(1);not null;default:0;comment:0:默认手动创建，1:convertlab创建，2:startDt创建"`
	AllowSelect      int                      `json:"allow_select" gorm:"column:allow_select;type:tinyint(1);not null;default:0;comment:0:是否允许重新选择发送客户"`
	IsBackupTemplate int                      `json:"is_backup_template" gorm:"column:is_backup_template;type:tinyint(1);not null;default:0;comment:是否使用的是备份模版 0 否 1是"`
}

func (PlanTaskModel) TableName() string {
	return "plan_tasks"
}

func (p PlanTaskModel) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", p.DBPrefix(), enterpriseID, p.TableName())
}
