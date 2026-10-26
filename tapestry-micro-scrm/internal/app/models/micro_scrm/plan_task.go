package microScrmModel

import "micro-scrm/internal/app/models"

const (
	PLANTASKTAGSEND = iota // 计划任务标签群发类型
	PLANTASKQYSEND         // 计划任务企业群发类型
	PLANTASK1V1SEND        // 计划任务1v1企业群发类型
)

const (
	PLANTASKCONTACT = iota // 计划任务沟通任务类型
	PLANTASKMOMENT         // 计划任务朋友圈任务类型
)
const (
	PLANTSINGLE = iota + 1 // 单次任务
	PLANTCOUNT             // 持续任务
)

const (
	PlanTaskPushNotStarted int32 = iota // 未执行
	PlanTaskPushSuccess                 // 已执行
	PlanTaskPushRunning                 // 推送中
	PlanTaskPushFailed                  // 下发失败
)

const (
	TAGTASKFINISHEXPIRE  = iota + 1 // 标签任务的标签有效期，任务结束后多长时间过期
	TAGTASKCUSTOMEXPIRE             // 标签任务的标签自定义有效期
	TAGTASKNOTEXPIRE                // 标签任务的标签永不过期
	TAG_TASK_SEND_EXPIRE            //任务发送立刻删除
)

// 子任务的条件类型
const (
	TaskConditionTypeTag           = 1 // 标签
	TaskConditionTypeGroupNum      = 2 // 群人数
	TaskConditionTypeAddGroupNum   = 3 // 群内新增人数
	TaskConditionTypeLeaveGroupNum = 7 // 群内新增人数
	TaskConditionTypeAllCustomer   = 9 // 全部顾客
)

// 子任务的维度
const (
	TaskDimensionCommunity       = iota + 1 // 社群数据
	TaskDimensionTag                        // 用户画像
	TaskDimensionImportUnionId              // 导入客户名单
	TaskDimensionImportSaUnionId            // 导入客户&&sa名单
)

const (
	SendTypeBA    = iota + 1 // BA
	SendTypeKeySa            // 2管理bA
)

const (
	BackupTemplateCLose = iota
	BackupTemplateOpen
)

const START_DT_TASK_ORIGIN_TYPE = 2
const SHENYAN_TASK_ORIGIN_TYPE = 3

type PlanTasksModel struct {
	models.BaseModel
	PlanID           int64      `json:"plan_id" gorm:"column:plan_id;type:bigint(20);not null;default:0;comment:子任务对应的计划的ID"`
	Name             string     `json:"name" gorm:"column:name;type:varchar(64);not null;comment:子任务名称"`
	Dimension        uint32     `json:"dimension" gorm:"column:dimension;type:tinyint(2) unsigned;not null;default:1;comment:子任务的维度；1社群数据；2用户画像；3导入模板名单,4导入顾客&&saCode名单"`
	ConditionType    int32      `json:"condition_type" gorm:"column:condition_type;type:tinyint(2);not null;default:1;comment:子任务的条件类型；1标签；2群人数；3群内新增人数；4参与活动；5购买行为；6群活跃度；7退群名单；8顾客生日,9全部顾客"`
	IsSpecifyCrowd   int32      `json:"is_specify_crowd" gorm:"column:is_specify_crowd;type:tinyint(1);default:0;comment:任务是否指定人群 1 否 2 是"`
	IsOpen           int32      `json:"is_open" gorm:"column:is_open;type:tinyint(2);not null;default:1;comment:任务在小程序中是否可见。0不可见，1可见；默认可见"`
	TaskType         int32      `json:"task_type" gorm:"column:task_type;type:tinyint(1);default:0;comment:任务周期 1 单次 2持续"`
	SendType         int32      `json:"send_type" gorm:"column:send_type;type:tinyint(1);default:0;comment:任务推送对象 1 BA 2管理bA"`
	Status           int32      `json:"status" gorm:"column:status;type:tinyint(2);not null;default:0;comment:任务状态；0未开始、1进行中、2已结束"`
	Remark           string     `json:"remark" gorm:"column:remark;type:varchar(200);not null;comment:子任务描述"`
	TaskIsPush       int32      `json:"task_is_push" gorm:"column:task_is_push;type:tinyint(1);default:0;comment:任务是否推送 0 未执行 1 已执行 2 推送中 3 下发失败"`
	BirthdayPreTime  int        `json:"birthday_pre_time" gorm:"column:birthday_pre_time;type:int(10);not null;comment:生日消息推送的提前天数"`
	ConsumerFileID   int        `json:"consumer_file_id" gorm:"column:consumer_file_id;type:int(10);not null;default:0;comment:导入模板名单的文件id"`
	IsExtract        int32      `json:"is_extract" gorm:"column:is_extract;type:tinyint(1);not null;default:0;comment:客户名单是否导入完成,0: 未完成，1:已完成"`
	CrowdID          string     `json:"crowd_id" gorm:"column:crowd_id;type:varchar(200);not null;default:0;comment:MA任务的任务包批次ID"`
	PlanSendTime     int64      `json:"plan_send_time" gorm:"column:plan_send_time;type:int(10);not null;default:0;comment:MA任务的推送时间"`
	DownloadUrl      string     `json:"download_url" gorm:"column:download_url;type:varchar(200);not null;default:'';comment:MA人群包匹配COS的地址"`
	TaskCategory     int        `json:"task_category" gorm:"column:task_category;type:tinyint(1);not null;default:0;comment:子任务的任务类型，沟通任务/朋友圈任务"`
	GroupSendType    int        `json:"group_send_type" gorm:"column:group_send_type;type:tinyint(1);not null;default:0;comment:子任务群发类型，标签任务群发/企业任务群发"`
	IsAddTag         int        `json:"is_add_tag" gorm:"column:is_add_tag;type:tinyint(1);not null;default:0;comment:是否生成任务标签"`
	TagExpireTime    int        `json:"tag_expire_time" gorm:"column:tag_expire_time;type:int(10);not null;default:0;comment:标签任务的过期时间,0为永不过期"`
	TagName          string     `json:"tag_name" gorm:"column:tag_name;type:varchar(200);not null;default:'';comment:子任务的打标签的标签名称"`
	QyTagId          string     `json:"qy_tag_id" gorm:"column:qy_tag_id;type:varchar(200);not null;default:'';comment:企微的tag_id"`
	TagExpireType    int        `json:"tag_expire_type" gorm:"column:tag_expire_type;type:tinyint(1);not null;default:0;comment:企微标签的过期类型,1任务结束后，2自定义删除时间，3长期有效"`
	MomentJobId      string     `json:"moment_job_id" gorm:"column:moment_job_id;type:varchar(200);not null;default:'';comment:朋友圈的jobid"`
	MomentStatus     int        `json:"moment_status" gorm:"column:moment_status;type:tinyint(1);not null;default:0;comment:朋友圈的任务创建状态，，1表示开始创建任务，2表示正在创建任务中，3表示创建任务已完成"`
	OriginTaskId     int64      `json:"origin_task_id" gorm:"column:origin_task_id;type:bigint(20);not null;default:0;comment:针对持续任务一直创建的源子任务id"`
	FinishRate       int        `json:"finish_rate" gorm:"column:finish_rate;type:tinyint(1);not null;default:100;comment:任务完成率"`
	TaskStartTime    int64      `json:"task_start_time" gorm:"column:task_start_time;type:int(10);not null;default:0;comment:子任务的开始时间"`
	TaskEndTime      int64      `json:"task_end_time" gorm:"column:task_end_time;type:int(10);not null;default:0;comment:子任务的结束时间"`
	TaskExpireType   int        `gorm:"column:task_expire_type;type:tinyint(1);not null;default:0;comment:子任务有效期，0为根据主任务的有效期一致，2指定时间有效期"`
	IsPush           int32      `json:"is_push" gorm:"column:is_push;type:tinyint(1);not null;default:0;comment:MA创建的任务推送完毕，值会更新，2为已推送完毕"`
	PlanTemplateId   int64      `json:"plan_template_id" gorm:"column:plan_template_id;type:bigint(20);not null;default:0;comment:任务模板id"`
	FlowId           string     `json:"flow_id" gorm:"column:flow_id;type:varchar(200);not null;default:'';comment:MA的自动流ID"`
	Campaign         string     `json:"campaign" gorm:"column:campaign;type:varchar(200);not null;default:'';comment:MA的营销活动编码"`
	PlanName         string     `json:"plan_name" gorm:"column:plan_name;type:varchar(200);not null;default:'';comment:MA推送的任务名称"`
	OriginType       int        `json:"origin_type" gorm:"column:origin_type;type:tinyint(1);not null;default:0;comment:0:默认手动创建，1:convertlab创建，2:startDt创建"`
	AllowSelect      int        `json:"allow_select" gorm:"column:allow_select;type:tinyint(1);not null;default:0;comment:0:是否允许重新选择发送客户"`
	IsBackupTemplate int        `json:"is_backup_template" gorm:"column:is_backup_template;type:tinyint(1);not null;default:0;comment:是否使用的是备份模版 0 否 1是"`
	Plan             PlansModel `json:"plan" gorm:"foreignKey:PlanID;references:ID"`
}

// TableName returns the table name of the PlanTasks model
func (PlanTasksModel) TableName() string {
	return "plan_tasks"
}
