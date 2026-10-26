package dashboardJob

import (
	"archive/zip"
	"database/sql"
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"io"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	planService "micro-scrm/internal/app/service/plan"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"path"
	"strconv"
)

type SyncDataFabric struct {
	Name       string
	transfer   *dataFabricTransfer
	startTime  int64
	endTime    int64
	sqlDb      *sql.DB
	remotePath string
	prefix     string
	suffix     string
}

func NewSyncDataFabric() *SyncDataFabric {
	startTime, endTime := helper.YesterdayPeriodOfTimestamp(1)
	// startTime, endTime := int64(0), int64(0)
	suffix := helper.GetFormatDate(0, 2)
	sqlDb, _ := gormx.GetInstance().DB()
	return &SyncDataFabric{
		Name:       "同步datafabric文件",
		startTime:  startTime,
		endTime:    endTime,
		sqlDb:      sqlDb,
		remotePath: "inbound/c360/",
		suffix:     suffix,
	}
}

func (job *SyncDataFabric) Run(c ctx.Context, handler *xjob.JobHandler) error {

	enterpriseList, err := orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}
	specs := handler.Runner().Specs()
	if len(specs) > 0 {
		job.startTime = 0
		job.endTime = 0
	}
	for _, v := range enterpriseList {
		enterpriseId := v.ID
		cacheRes, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
		if err != nil {
			continue
		}
		if cacheRes.Name != "coach" {
			continue
		}
		err = job.handle(c, enterpriseId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (job *SyncDataFabric) handle(c ctx.Context, enterpriseId int64) error {
	brands := make(map[int]string)
	job.transfer = newDataFabricTransfer([]int{micro_scrm.SftpTypeDEL, micro_scrm.SftpTypeDELALI}, func(kind int) (*dataFabricTarget, error) {
		brands[kind] = "coach"
		sftpConfig, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfig(c, map[string]interface{}{
			"enterprise_id": enterpriseId,
			"sftp_type":     kind,
		})
		if err != nil {
			return nil, fmt.Errorf("读取SFTP配置失败: %w", err)
		}
		if sftpConfig.ID == 0 {
			return nil, nil
		}
		brands[kind] = sftpConfig.BrandName
		client, err := sftp.NewSFTPClient(sftpConfig.ConfigDetail)
		if err != nil {
			return nil, fmt.Errorf("SFTP连接失败: %w", err)
		}
		remotePath := sftpConfig.PathConfig.DataFabric
		if remotePath == "" {
			remotePath = job.remotePath
		}
		dirPath, err := createDataFabricDir(client, remotePath, job.suffix)
		if err != nil {
			_ = client.Client.Close()
			return nil, fmt.Errorf("创建SFTP目录失败: %w", err)
		}
		return &dataFabricTarget{
			brand: sftpConfig.BrandName,
			upload: func(local string) error {
				return uploadDataFabricFile(client, dirPath, local)
			},
			close: client.Client.Close,
		}, nil
	}, func(kind int, err error) {
		planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
			EnterpriseName: brands[kind],
			SFTPName:       fmt.Sprintf("datafabric sftp_type=%d", kind),
			UploadRes:      fmt.Sprintf("失败,错误为：%s", err.Error()),
		})
	})
	defer job.transfer.close()
	if len(job.transfer.targets) == 0 {
		return job.transfer.err
	}
	// 查询增量的数据，导出数据至CSV文件，压缩CSV文件，上传SFTP
	// 针对每个表写方法执行上传逻辑
	shardDataBase := fmt.Sprintf("micro_scrm_enterprise_%d", enterpriseId)
	for _, target := range job.transfer.targets {
		job.transfer.selectedKind = target.kind
		full := target.kind == micro_scrm.SftpTypeDELALI
		for _, export := range []func() error{
			func() error { return job.memberTableHandle(c, enterpriseId, "micro_scrm", "members", full) },
			func() error {
				return job.masterTableHandleWithMode(c, enterpriseId, "micro_scrm", "member_departments", "id,member_id,department_id,enterprise_id,created_at,updated_at,deleted_at", full)
			},
			func() error { return job.departmentsTableHandle(c, enterpriseId, "micro_scrm", "departments", full) },
		} {
			if target.failed {
				break
			}
			if err := export(); err != nil {
				target.failed = true
				job.transfer.err = fmt.Errorf("datafabric sftp_type=%d: %w", target.kind, err)
				break
			}
		}
	}
	job.transfer.selectedKind = 0
	job.transfer.recordGenerationError(job.storeDepartmentTableHandle(c, enterpriseId, "micro_work", "store_department"))
	job.transfer.recordGenerationError(job.storeEmployeeTableHandle(c, enterpriseId, "micro_work", "store_employee"))
	job.transfer.recordGenerationError(job.storeWorkerTableHandle(c, enterpriseId, "micro_work", "store_worker"))
	job.transfer.recordGenerationError(job.shardBehaviorDataTableHandle(c, shardDataBase, "member_external_behavior_data"))
	job.transfer.recordGenerationError(job.shardExternalTableHandle(c, shardDataBase, "wx_user_corp_external"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "impression_tag", "id,enterprise_id,group_id,third_group_id,third_tag_id,name,description,created_at,updated_at,deleted_at,use_num"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "impression_tag_group", "id,enterprise_id,third_group_id,group_name,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "impression_tag_operate_log", "id,enterprise_id,operate_user_id,operate_name,third_tag_group_id,third_tag_id,external_username,external_unionid,action_type,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "impression_tag_relation", "id,enterprise_id,external_user_id,user_id,tag_id,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.externalTagTableHandle(c, enterpriseId, "micro_scrm", "work_external_tag"))
	job.transfer.recordGenerationError(job.externalTagGroupTableHandle(c, enterpriseId, "micro_scrm", "work_external_tag_group"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "work_external_user_tag", "id,enterprise_id,relation_id,external_tag_id,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.userRelationTableHandle(c, enterpriseId, "micro_scrm", "work_external_user_relation"))
	workInfo, err := microBossDao.NewEnterpriseWorkPlugin().GetOne(c, map[string]interface{}{"enterprise_id": enterpriseId})
	if err != nil {
		job.transfer.recordGenerationError(err)
		planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
			EnterpriseName: "coach",
			SFTPName:       "datafabric的relation_callback_records表获取企业信息corpId失败",
			UploadRes:      "失败" + err.Error(),
		})
	}
	job.transfer.recordGenerationError(job.callbackTableHandle(c, "micro_scrm", "relation_callback_records", workInfo.CorpID))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "drainage_contact_channel_labels", "id,channel_id,label_group_id,label_id,enterprise_id,wx_tag_group_id,wx_tag_id,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "drainage_contact_channels", "id,enterprise_id,name,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "drainage_contact_members", "id,contact_id,enterprise_id,member_id,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "drainage_contacts", "id,name,channel_id,auto_friend,qr_code_url,member_count,enterprise_id,sum_add_member_number,created_at,updated_at,deleted_at,config_id,type"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "drainage_welcomes", "id,content,enterprise_id,is_default,channel_id,type,condition_content,member_level,store_type,custom_type,receive_coupon_status,reserve_status,external_type_data_code,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "welcomes_external_type", "id,enterprise_id,name,code,remark,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "welcomes_external_type_data", "id,enterprise_id,welcomes_external_type_id,name,code,sort,remark,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "customer_inherit", "id,enterprise_id,handover_userid,handover_name,handover_mobile,handover_department_id,handover_department_name,external_userid,external_nickname,external_avatar,takeover_userid,takeover_name,takeover_department_name,dimission_time,status,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.onjobInheritTableHandle(c, enterpriseId, "micro_scrm", "customer_onjob_inherit"))
	job.transfer.recordGenerationError(job.onjobInheritRecordTableHandle(c, enterpriseId, "micro_scrm", "customer_onjob_inherit_record"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "reserve_service", "id,reserve_name,reserve_img,reserve_time,reserve_status,end_time,created_at,updated_at,deleted_at,service_type,start_time,service_role,time_slot"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "reserve_service_relation", "id,service_id,store_id,store_name,member_id,member_name,created_at,updated_at,deleted_at,is_store,store_type,department_id,user_id"))
	job.transfer.recordGenerationError(job.shardReserveOrderTableHandle(c, shardDataBase, "reserve_order"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "reserve_log", "id,reserve_order_id,log_info,order_info,created_at,updated_at,deleted_at,log_type"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "invite_template", "id,reserve_service_id,title_name,status_type,start_time,end_time,background_url,content,share_title,share_url,is_up,user_is_show,member_is_show,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "invite_record", "id,invite_template_id,reserve_order_id,member_id,member_userid,user_name,user_unionid,user_openid,style_type,background_url,send_content,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardPlanTemplateTableHandle(c, shardDataBase, "plan_task_template"))
	job.transfer.recordGenerationError(job.shardPlansTableHandle(c, shardDataBase, "plans"))
	job.transfer.recordGenerationError(job.shardPlanTasksTableHandle(c, shardDataBase, "plan_tasks"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "plan_task_conditions", "id,plan_id,task_id,conditions,content,content_type,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "plan_task_specify_crowd", "id,plan_id,task_id,user_id,user_name,is_ba,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "task_allot_log", "id,plan_id,condition_id,task_id,user_id,start_time,end_time,is_finished,customer_num,group_chat_id,created_at,updated_at,deleted_at,tag_id,tag_name,is_forward,group_msg_id,is_group_send,agent_msg_id,agent_msg_time,chat_num,close_tips,contacted_count"))
	job.transfer.recordGenerationError(job.shardTaskConditionDataTableHandle(c, enterpriseId, shardDataBase))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "source_text", "id,content,enterprise_id,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "source_image", "id,enterprise_id,name,url,size,ext,media_id,media_expire_time,media_status,media_fail_reason,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "source_link", "id,title,description,cover_url,link_type,link_url,author,content,enterprise_id,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "source_mini_program", "id,enterprise_id,mini_program_id,title,path,cover_url,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "source_video", "id,enterprise_id,name,cover_url,url,status,file_id,ext,media_id,media_expire_time,media_status,media_fail_reason,size,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "source_friend", "id,enterprise_id,text_content,name,video_cover_url,original_url,compression_url,ext,size,type,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "source", "id,enterprise_id,group_id,type,status,admin_id,start_time,end_time,created_at,updated_at,deleted_at,member_id,is_deal,compress_id,merge_id,sort,download_times,forwarding_times,copy_times"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "source_category", "id,enterprise_id,name,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "source_category_relation", "id,enterprise_id,category_id,source_id,type,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_integral_rule", "id,enterprise_id,type,rule_name,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_integral_rule_description", "id,enterprise_id,content,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_integral_rule_detail", "id,enterprise_id,rule_id,rule_category_id,validity_type,start_time,end_time,integral,integral_type,is_diff_store,store_type,rule_if,rule_if1,rule_value,is_astrict,astrict_type,max_integral,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_medal", "id,enterprise_id,name,source,rule_type,rule_start_time,rule_end_time,meda_type,meda_description,other_rewards,integral,item_name,status,rule_dimension_id,rule_category_id,invite_top,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_medal_record", "id,enterprise_id,medal_id,medal_style_id,member_id,is_show,is_receive,is_pop,coins_record_id,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_medal_rule", "id,enterprise_id,rule_name,short_name,description,time_type,pid,status,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_medal_style", "id,enterprise_id,medal_id,time_value,time_name,medal_name,light_icon,unlit_icon,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_medal_task_record", "id,enterprise_id,medal_id,time_type,status,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_incentive_config", "id,enterprise_id,rule_name,status,type,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_incentive_relation", "id,enterprise_id,is_diff_store,incentive_config_type,status,ability_index,time_rule,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_incentive_relation_rule", "id,enterprise_id,store_type,incentive_relation_id,incentive_config_type,rule_condition,rule_value,rule_give,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.excitationMemberDataTableHandle(c, enterpriseId, "micro_scrm", "excitation_member_data"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_member_data_record", "id,enterprise_id,member_id,level,invite_value,task_value,interact_value,train_value,online_value,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_level_rule", "id,enterprise_id,level_type,level_name,level_icon,min_value,max_value,rule_invite_id,rule_task_id,rule_interact_id,rule_train_id,rule_online_id,rule_condition_invite,rule_condition_task,rule_condition_interact,rule_condition_train,rule_condition_online,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_achievement", "id,enterprise_id,member_id,type,value,unique_hash,source_remark,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "excitation_points_record", "id,incentive_type,member_id,type,value,condition_value,points,source_type,integral_rule_id,integral_rule_name,integral_rule_detail_name,points_remark,fy_time_str,coins_record_id,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "excitation_coins_record", "id,member_id,ratio_remark,coins,source_type,points_record_id,medal_record_id,coins_remark,order_num,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_coins_expire_record", "id,enterprise_id,expired_at,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_coins_settings", "id,enterprise_id,ratio_remark,ratio,mall_is_open,expire_is_open,expire_at,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_coins_settings_edit_record", "id,enterprise_id,admin_id,admin_name,ratio_remark,mall_is_open,expire_at,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.masterTableHandle(c, enterpriseId, "micro_scrm", "excitation_mall_order", "id,enterprise_id,user_id,member_id,order_num,order_coins,order_status,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "display_task", "id,display_name,store_type,content,class_url,template_info,show_type,status,start_time,end_time,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "display_task_allot", "id,display_name,display_task_id,store_code,store_id,store_name,store_type,store_city,finish_status,grade,owner_member_id,owner_user_id,owner_name,review_remark,submit_time,review_time,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "display_task_allot_data", "id,display_task_id,allot_id,identifier,image_url,status,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "display_task_allot_remark", "id,display_task_id,allot_id,allot_data_id,member_id,member_user_id,member_name,avatar,is_vm,store_code,reply_name,comment_time,content,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "display_task_allot_review", "id,display_task_id,allot_id,review_type,review_id,review_user_id,review_name,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "display_task_allot_submit", "id,display_task_id,allot_id,submit_type,member_id,member_name,member_user_id,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "display_task_look_log", "id,allot_data_id,member_id,allot_id,created_at,updated_at,deleted_at,comment_id"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "display_task_relation", "id,display_task_id,store_id,store_name,member_id,member_name,created_at,updated_at,deleted_at,is_store,store_type,department_id,user_id"))
	//_ = job.masterTableHandle(c, enterpriseId, "micro_scrm", "sow_task", "id,enterprise_id,work_external_user_id,member_id,member_user_id,content,status,remind_time,is_push,template_id,created_at,updated_at,deleted_at")
	//_ = job.masterTableHandle(c, enterpriseId, "micro_scrm", "sow_template", "id,enterprise_id,name,theme,status,admin_name,use_num,created_at,updated_at,deleted_at")
	//_ = job.masterTableHandle(c, enterpriseId, "micro_scrm", "sow_template_notice", "id,enterprise_id,template_id,delay_type,delay_value,delay_time,content,created_at,updated_at,deleted_at")
	//活动二维码相关表
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "activity_qrcode_main", "id,name,alias_name,banner_url,time_type,start_time,end_time,category_ids,add_name_tag,add_category_tag,tag_ids,name_tag_id,skip_verify,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "activity_qrcode_category", "id,name,tag_id,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "activity_qrcode_member_scope", "id,activity_id,member_id,area_id,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "activity_qrcode_contact_way", "id,activity_id,state,third_id,store_name,config_id,qr_code,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.shardTableHandle(c, shardDataBase, "activity_qrcode_point", "id,activity_id,name,is_point,category_id,code_url,link_url,created_at,updated_at,deleted_at"))
	job.transfer.recordGenerationError(job.createRdy(c))
	for _, target := range job.transfer.targets {
		if !target.uploaded[job.suffix+".rdy"] {
			continue
		}
		planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
			EnterpriseName: target.brand,
			SFTPName:       fmt.Sprintf("datafabric sftp_type=%d 文件上传成功", target.kind),
			UploadRes:      "成功",
		})
	}
	return job.transfer.err
}

// members表的处理
func (job *SyncDataFabric) memberTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName string, full bool) error {
	sqlFiled := " id,enterprise_id,userid,name,mobile,position,gender,email,avatar,thumb_avatar,telephone,alias,status,qr_code,address,open_userid,type,qr_poster_code_id,qr_poster_code,created_at,updated_at,deleted_at "
	query := "select " + sqlFiled + " from " + databaseName + "." + tableName + " where enterprise_id=?"
	args := []interface{}{enterpriseId}
	if !full {
		query += " and updated_at>=?"
		args = append(args, job.startTime)
	}
	rows, err := job.sqlDb.Query(query, args...)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as total_count,COUNT(DISTINCT(userid)) AS distinct_userid ,COUNT(DISTINCT(mobile)) as distinct_mobile FROM "+databaseName+"."+tableName+" WHERE enterprise_id=? ", tableName, enterpriseId)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_userid", "distinct_mobile",
	}
	//var total int64
	//err = job.sqlDb.QueryRow("SELECT COUNT(1) FROM "+databaseName+"."+tableName+" WHERE enterprise_id=? and updated_at>=?", enterpriseId, job.startTime).Scan(&total)
	//if err != nil {
	//	fmt.Printf("查询"+tableName+"的.done文件数据total报错%v\n", err)
	//	return err
	//}
	err = job.CreateDoneFile(c, tableName, columnsName, sqlRows, incrementCount)
	if err != nil {
		fmt.Printf("写入"+tableName+"的.done文件数据报错%v\n", err)
	}
	return err
}

// departments表的处理
func (job *SyncDataFabric) departmentsTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName string, full bool) error {
	sqlFiled := " id,parent_id,third_id,third_parent_id,name,name_en,enterprise_id,`order`,created_at,updated_at,deleted_at,child_departments "
	query := "select " + sqlFiled + " from " + databaseName + "." + tableName + " where enterprise_id=?"
	args := []interface{}{enterpriseId}
	if !full {
		query += " and updated_at>=?"
		args = append(args, job.startTime)
	}
	rows, err := job.sqlDb.Query(query, args...)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,COUNT(DISTINCT(name)) AS distinct_name FROM "+databaseName+"."+tableName+" WHERE enterprise_id=?", tableName, enterpriseId)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_name",
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

// store_department表的处理
func (job *SyncDataFabric) storeDepartmentTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName string) error {
	sqlFiled := " id,enterprise_id,third_store_code,uc_department_id,wechat_department_id,store_name,store_type,operate_name,sync_status,sync_at,address,email,popup_store_status,sa_share_order_url,sa_check_order_url,popup_store_keys,sa_check_role_status,sa_check_role,phone,city,inventory_name,closed_at,created_at,updated_at,deleted_at "
	rows, err := job.sqlDb.Query("select "+sqlFiled+" from "+databaseName+"."+tableName+" where enterprise_id=? and updated_at>=?", enterpriseId, job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,COUNT(DISTINCT(enterprise_id)) AS distinct_enterprise_id,COUNT(DISTINCT(store_name)) AS distinct_store_name  FROM "+databaseName+"."+tableName+" WHERE enterprise_id=?", tableName, enterpriseId)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_enterprise_id", "distinct_store_name",
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

// store_employee表的处理
func (job *SyncDataFabric) storeEmployeeTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName string) error {
	sqlFiled := " id,enterprise_id,name,alias_name,employee_code,job_location,job_position,gender,email,status,sync_status,sync_at,third_store_code,wechat_user_id,mobile,created_at,updated_at,deleted_at "
	rows, err := job.sqlDb.Query("select "+sqlFiled+" from "+databaseName+"."+tableName+" where enterprise_id=? and updated_at>=?", enterpriseId, job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,COUNT(DISTINCT(employee_code)) AS distinct_employee_code FROM "+databaseName+"."+tableName+" WHERE enterprise_id=?", tableName, enterpriseId)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_employee_code",
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

// store_worker表的处理
func (job *SyncDataFabric) storeWorkerTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName string) error {
	sqlFiled := " id,enterprise_id,type,user_id,store_id,worker_no,name,mobile,mem_uid,mem_userid,contact_qr_code_id,nickname,avatar,sms_mobile,is_send_sms,qr_code,store_sun_code,is_boss,source,worker_status,entry_time,out_time,recently_logged_store_id,recently_logged_time,created_at,updated_at,deleted_at "
	rows, err := job.sqlDb.Query("select "+sqlFiled+" from "+databaseName+"."+tableName+" where enterprise_id=? and updated_at>=?", enterpriseId, job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,COUNT(DISTINCT(worker_no)) AS distinct_worker_no FROM "+databaseName+"."+tableName+" WHERE enterprise_id=?", tableName, enterpriseId)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_worker_no",
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

// work_external_tag表的处理
func (job *SyncDataFabric) externalTagTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName string) error {
	sqlFiled := " id,enterprise_id,group_id,tag_id,name,type,`order`,member_id,expired_time,create_time,created_at,updated_at,deleted_at,third_tag_id,third_tag_description "
	rows, err := job.sqlDb.Query("select "+sqlFiled+" from "+databaseName+"."+tableName+" where enterprise_id=? and updated_at>=?", enterpriseId, job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,COUNT(DISTINCT(tag_id)) AS distinct_tag_id FROM "+databaseName+"."+tableName+" WHERE enterprise_id=?", tableName, enterpriseId)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_tag_id",
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

// work_external_tag_group表的处理
func (job *SyncDataFabric) externalTagGroupTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName string) error {
	sqlFiled := " id,enterprise_id,group_id,group_name,`order`,type,is_third_tag,create_time,created_at,updated_at,deleted_at "
	rows, err := job.sqlDb.Query("select "+sqlFiled+" from "+databaseName+"."+tableName+" where enterprise_id=? and updated_at>=?", enterpriseId, job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,COUNT(DISTINCT(group_id)) AS distinct_group_id FROM "+databaseName+"."+tableName+" WHERE enterprise_id=?", tableName, enterpriseId)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_group_id",
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

// work_external_user_relation表的处理
func (job *SyncDataFabric) userRelationTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName string) error {
	sqlFiled := " id,enterprise_id,work_external_user_id,member_id,remark,description,createtime,remark_corp_name,remark_mobiles,add_way,state,oper_userid,relation_status,member_deleted_time,external_deleted_time,created_at,updated_at,deleted_at,user_id,unionid,is_key_sa "
	rows, err := job.sqlDb.Query("select "+sqlFiled+" from "+databaseName+"."+tableName+" where enterprise_id=? and updated_at>=?", enterpriseId, job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,COUNT(DISTINCT(member_id)) AS distinct_member_id FROM "+databaseName+"."+tableName+" WHERE enterprise_id=?", tableName, enterpriseId)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_member_id",
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

// customer_onjob_inherit表的处理
func (job *SyncDataFabric) onjobInheritTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName string) error {
	sqlFiled := " id,enterprise_id,member_id,member_name,status,errmsg,customers_numbers,transfer_success_msg,created_at,updated_at,deleted_at "
	rows, err := job.sqlDb.Query("select "+sqlFiled+" from "+databaseName+"."+tableName+" where enterprise_id=? and updated_at>=?", enterpriseId, job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,COUNT(DISTINCT(member_id)) AS distinct_member_id FROM "+databaseName+"."+tableName+" WHERE enterprise_id=?", tableName, enterpriseId)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_member_id",
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

// customer_onjob_inherit_record表的处理
func (job *SyncDataFabric) onjobInheritRecordTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName string) error {
	sqlFiled := " id,inherit_id,handover_member_id,handover_userid,handover_name,handover_mobile,takeover_member_id,takeover_userid,takeover_name,takeover_mobile,enterprise_id,external_id,external_userid,external_nickname,external_name,external_avatar,status,errmsg,takeover_time,request_callback_code,created_at,updated_at,deleted_at "
	rows, err := job.sqlDb.Query("select "+sqlFiled+" from "+databaseName+"."+tableName+" where enterprise_id=? and updated_at>=?", enterpriseId, job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,COUNT(DISTINCT(handover_member_id)) AS distinct_handover_member_id FROM "+databaseName+"."+tableName+" WHERE enterprise_id=?", tableName, enterpriseId)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_handover_member_id",
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

// excitation_member_data表的处理
func (job *SyncDataFabric) excitationMemberDataTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName string) error {
	sqlFiled := " id,enterprise_id,member_id,point_value,point_updated_at,coins,coins_updated_at,spend_coins,total_coins,level,level_type,level_id,level_updated_at,medal_num,medal_updated_at,invite_value,invite_rule_str,invite_rule_value,invite_updated_at,task_value,task_rule_str,task_rule_value,task_updated_at,interact_value,interact_rule_str,interact_rule_value,interact_updated_at,train_value,train_rule_str,train_rule_value,train_updated_at,online_value,online_rule_str,online_rule_value,online_updated_at,created_at,updated_at,deleted_at "
	rows, err := job.sqlDb.Query("select "+sqlFiled+" from "+databaseName+"."+tableName+" where enterprise_id=? and updated_at>=?", enterpriseId, job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,SUM(point_value) AS total_point_value,SUM(total_coins) AS total_coins FROM "+databaseName+"."+tableName+" WHERE enterprise_id=?", tableName, enterpriseId)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "total_point_value", "total_coins",
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

// 主库表的处理
func (job *SyncDataFabric) masterTableHandle(c ctx.Context, enterpriseId int64, databaseName, tableName, filed string) error {
	return job.masterTableHandleWithMode(c, enterpriseId, databaseName, tableName, filed, false)
}

func (job *SyncDataFabric) masterTableHandleWithMode(c ctx.Context, enterpriseId int64, databaseName, tableName, filed string, full bool) error {
	query := "select " + filed + " from " + databaseName + "." + tableName + " where enterprise_id=?"
	args := []interface{}{enterpriseId}
	if !full {
		query += " and updated_at>=?"
		args = append(args, job.startTime)
	}
	rows, err := job.sqlDb.Query(query, args...)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
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
func (job *SyncDataFabric) shardTableHandle(c ctx.Context, databaseName, tableName, filed string) error {
	rows, err := job.sqlDb.Query("select "+filed+" from "+databaseName+"."+tableName+" where updated_at>=?", job.startTime)
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
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

// member_external_behavior_data表的处理
func (job *SyncDataFabric) shardBehaviorDataTableHandle(c ctx.Context, databaseName, tableName string) error {
	filed := "id,member_id,userid,date,stat_time,new_apply_cnt,new_contact_cnt,chat_cnt,message_cnt,reply_percentage,avg_reply_time,negative_feedback_cnt,created_at,updated_at,deleted_at"
	rows, err := job.sqlDb.Query("select "+filed+" from "+databaseName+"."+tableName+" where updated_at>=?", job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,SUM(new_apply_cnt) as total_new_apply_cnt,SUM(chat_cnt) as total_chat_cnt FROM "+databaseName+"."+tableName, tableName)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "total_new_apply_cnt", "total_chat_cnt",
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

// wx_user_corp_external表的处理
func (job *SyncDataFabric) shardExternalTableHandle(c ctx.Context, databaseName, tableName string) error {
	filed := "id,appid,openid,unionid,euid,nick,name,position,avatar,corp_name,corp_full_name,type,sex,external_profile,birthday,union_code,point,grade,store_code,ba_code,created_at,updated_at,deleted_at,birth_year,birth_month,birth_day,mobile,card_level,clt_type,amount_total,customer_code,customer_id,last_purchase_date,last_purchase_store_code,upgrade_need_amount,customer_type,customer_type_updated_at,vip_type,customer_deleted_at,age_tag"
	rows, err := job.sqlDb.Query("select "+filed+" from "+databaseName+"."+tableName+" where updated_at>=?", job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,count(DISTINCT(openid)) as distinct_openid,count(DISTINCT(union_code)) as distinct_union_code FROM "+databaseName+"."+tableName, tableName)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_openid", "distinct_union_code",
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

// reserve_order表的处理
func (job *SyncDataFabric) shardReserveOrderTableHandle(c ctx.Context, databaseName, tableName string) error {
	filed := "id,reserve_number,reserve_id,reserve_time,store_id,store_name,store_address,reserve_type,status,user_name,user_unionid,user_external_id,member_id,member_name,reserve_remark,user_remark,is_open,is_commit,invite_num,member_user_id,reserve_order_sn,reserve_order_price,created_at,updated_at,deleted_at,is_sign,is_comment,comment_time,reserve_order_time,is_store,is_pay,is_confirm,order_origin,customer_status,reserve_order_status,is_vip,reserve_end_time"
	rows, err := job.sqlDb.Query("select "+filed+" from "+databaseName+"."+tableName+" where updated_at>=?", job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,count(DISTINCT(reserve_id)) as distinct_reserve_id,count(DISTINCT(user_external_id)) as distinct_external_id FROM "+databaseName+"."+tableName, tableName)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_reserve_id", "distinct_external_id",
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

// plan_task_template表的处理
func (job *SyncDataFabric) shardPlanTemplateTableHandle(c ctx.Context, databaseName, tableName string) error {
	filed := "id,template_name,template_desc,template_content,template_content_type,status,admin_id,created_at,updated_at,deleted_at,group_send_type"
	rows, err := job.sqlDb.Query("select "+filed+" from "+databaseName+"."+tableName+" where updated_at>=?", job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,count(DISTINCT(template_name)) as distinct_template_name FROM "+databaseName+"."+tableName, tableName)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_template_name",
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

// plans表的处理
func (job *SyncDataFabric) shardPlansTableHandle(c ctx.Context, databaseName, tableName string) error {
	filed := "id,name,task_num,operation_user_id,start_time,end_time,is_open,status,remark,created_at,updated_at,deleted_at,is_push,type,plan_template_id,category"
	rows, err := job.sqlDb.Query("select "+filed+" from "+databaseName+"."+tableName+" where updated_at>=?", job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,count(DISTINCT(name)) as distinct_name FROM "+databaseName+"."+tableName, tableName)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_name",
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

// plan_tasks表的处理
func (job *SyncDataFabric) shardPlanTasksTableHandle(c ctx.Context, databaseName, tableName string) error {
	filed := "id,plan_id,name,dimension,condition_type,is_specify_crowd,is_open,task_type,send_type,status,remark,task_is_push,created_at,task_start_time,task_end_time,task_expire_type,updated_at,deleted_at,birthday_pre_time,consumer_file_id,is_extract,crowd_id,plan_send_time,download_url,task_category,group_send_type,is_add_tag,tag_expire_time,tag_name,qy_tag_id,tag_expire_type,moment_job_id,moment_status,origin_task_id"
	rows, err := job.sqlDb.Query("select "+filed+" from "+databaseName+"."+tableName+" where updated_at>=?", job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,count(DISTINCT(name)) as distinct_name FROM "+databaseName+"."+tableName, tableName)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_name",
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

// plan_tasks表的处理
func (job *SyncDataFabric) shardTaskConditionDataTableHandle(c ctx.Context, enterpriseId int64, databaseName string) (err error) {
	planMany, err := planDao.NewPlanTaskDao().GetPlanMany(c, enterpriseId, "end_time >=? and start_time < ?", job.startTime, job.endTime)
	// planMany, err := planDao.NewPlanTaskDao().GetPlanMany(c, enterpriseId, "updated_at >=?", job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的task_condition_data文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	// 创建csv文件
	fileName := "task_condition_data_" + job.suffix
	// 创建一个新的 CSV 文件
	csvFileName := "./storage/data_fabric/" + fileName + ".csv"
	doneCsvFileName := "./storage/data_fabric/" + fileName + "done.csv"
	doneFileName := "./storage/data_fabric/" + fileName + ".done"
	csvFile, err := os.Create(csvFileName)
	if err != nil {
		return err
	}
	defer csvFile.Close()
	doneCsvFile, err := os.Create(doneCsvFileName)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return err
	}

	defer doneCsvFile.Close()
	// 创建一个 CSV 写入器
	csvWriter := csv.NewWriter(csvFile)
	columnNames := []string{
		"id", "task_allot_log_id", "tag_name", "ww_tag_id", "union_id", "work_external_user_id", "is_union", "mobile", "mars_id", "user_id", "condition_id", "task_id", "plan_id",
		"created_at", "updated_at", "deleted_at", "customer_id", "is_finished", "finished_at", "remark", "is_contact", "is_call", "is_add", "is_match", "user_name", "is_touch", "group_chat_id",
	}
	// 写入列名
	err = csvWriter.Write(columnNames)
	if err != nil {
		fmt.Println("Error writing column names:", err)
		return err
	}
	if _, err = doneCsvFile.WriteString("\xEF\xBB\xBF"); err != nil { // 写入UTF-8 BOM
		return err
	}

	data := [][]string{
		{"table_name", "incremental_count", "total_count"},
	}
	// 创建一个 CSV 写入器
	donCsvWriter := csv.NewWriter(doneCsvFile)
	if err != nil {
		fmt.Println("Error getting column names:", err)
		return err
	}

	if err != nil {
		fmt.Println("Error writing column names:", err)
		return err
	}

	var totalCount int64
	planAllMany, err := planDao.NewPlanTaskDao().GetPlanMany(c, enterpriseId, "updated_at >=?", 0)
	if err != nil {
		return err
	}
	for _, model := range planAllMany {
		planId := model.ID
		var total int64
		shardTableName := fmt.Sprintf("task_condition_data_%d", planId)
		err = job.sqlDb.QueryRow("SELECT COUNT(1) FROM " + databaseName + "." + shardTableName).Scan(&total)
		if err != nil {
			fmt.Printf("查询"+shardTableName+"的total失败%v\n", err)
			return err
		}
		totalCount += total
	}

	var incrementTotal int64
	for _, many := range planMany {
		planId := many.ID
		shardTableName := fmt.Sprintf("task_condition_data_%d", planId)
		rows, err := job.getTaskConditionData(c, databaseName, shardTableName)
		if err != nil {
			fmt.Printf("查询"+shardTableName+"合并数据%v\n", err)
			return err
		}
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
				_ = rows.Close()
				return err
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
				_ = rows.Close()
				return err
			}
			incrementTotal++
		}
		if err = rows.Err(); err != nil {
			_ = rows.Close()
			return err
		}
		_ = rows.Close()
	}
	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		return err
	}
	err = csvFile.Close()
	if err != nil {
		fmt.Println(" CSV file close:", err)
		return err
	}
	zipFile, err := job.CreateZipFile(fileName, csvFile)
	if err != nil {
		return err
	}
	localPath := "./storage/data_fabric/" + fileName + ".zip"
	err = job.uploadFile(c, zipFile, localPath)
	data = append(data, []string{"task_condition_data", strconv.FormatInt(incrementTotal, 10), strconv.FormatInt(totalCount, 10)})
	err = donCsvWriter.WriteAll(data)
	if err != nil {
		fmt.Printf(" doneCSV file write:%v\n", err)
		return err
	}
	donCsvWriter.Flush()
	err = doneCsvFile.Close()
	if err != nil {
		fmt.Println(" CSV file close:", err)
		return err
	}
	err = os.Rename(doneCsvFileName, doneFileName)
	if err != nil {
		fmt.Println(" CSV file rename err:", err)
		return err
	}
	err = job.uploadDoneFile(c, doneFileName)
	_ = os.Remove(csvFileName)
	_ = os.Remove(doneCsvFileName)
	_ = os.Remove(localPath)
	_ = os.Remove(doneFileName)
	return err
}

func (job *SyncDataFabric) getTaskConditionData(c ctx.Context, databaseName, tableName string) (*sql.Rows, error) {
	//var incrmentTotal int64
	//var taskConditionDataInfo taskConditionData
	filed := " CONCAT(plan_id,'_',id) as id,task_allot_log_id,tag_name,ww_tag_id,union_id,work_external_user_id,is_union,mobile,mars_id,user_id,condition_id,task_id,plan_id,created_at,updated_at,deleted_at,customer_id,is_finished,finished_at,remark,is_contact,is_call,is_add,is_match,user_name,is_touch,group_chat_id"
	rows, err := job.sqlDb.Query("select "+filed+" from "+databaseName+"."+tableName+" where updated_at>=?", job.startTime)
	if err != nil {
		//return rows, taskConditionDataInfo, err
		return rows, err
	}
	//err = job.sqlDb.QueryRow("SELECT count(1) FROM "+databaseName+"."+tableName+" WHERE updated_at>=?", job.startTime).Scan(&incrmentTotal)
	//if err != nil {
	//	fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
	//	return rows, taskConditionDataInfo, err
	//}
	//taskConditionDataInfo.IncrementCount = incrmentTotal
	return rows, err
}
func (job *SyncDataFabric) callbackTableHandle(c ctx.Context, databaseName, tableName, toUserName string) error {
	filed := "id,to_user_name,from_user_name,create_time,msg_type,event,change_type,user_id,external_user_id,state,welcome_code,fail_reason,chat_id,status,channel_id,status_trace,add_relation_time"
	rows, err := job.sqlDb.Query("select "+filed+" from "+databaseName+"."+tableName+" where to_user_name=? and create_time>=?", toUserName, job.startTime)
	defer func() {
		if err != nil {
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: "coach",
				SFTPName:       "datafabric的" + tableName + "文件生成失败",
				UploadRes:      "失败" + err.Error(),
			})
		}
	}()
	if err != nil {
		return err
	}
	incrementCount, err := job.CreateFile(c, tableName, rows)
	if err != nil {
		return err
	}
	sqlRows, err := job.sqlDb.Query("SELECT ? as `table_name`,0 as incremental_count,count(1) as `total_count`,count(DISTINCT(user_id)) as distinct_user_id FROM "+databaseName+"."+tableName+" where to_user_name=? ", tableName, toUserName)
	if err != nil {
		fmt.Printf("查询"+tableName+"的.done文件数据报错%v\n", err)
		return err
	}
	columnsName := []string{
		"table_name", "incremental_count", "total_count", "distinct_user_id",
	}
	//var total int64
	//err = job.sqlDb.QueryRow("SELECT COUNT(1) FROM "+databaseName+"."+tableName+" where to_user_name=? and create_time>=?", toUserName, job.startTime).Scan(&total)
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

func createDataFabricDir(client *sftp.SFTP, remotePath, timeStr string) (string, error) {
	dirPath := path.Join(remotePath, timeStr)
	_, err := client.Client.Stat(dirPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		if err = client.Client.Mkdir(dirPath); err != nil {
			return "", err
		}
		if err = client.Client.Chmod(dirPath, 0777); err != nil {
			return "", err
		}
	}
	return client.Client.RealPath(dirPath)
}

func (job *SyncDataFabric) CreateFile(c ctx.Context, tableName string, rows *sql.Rows) (int64, error) {

	fileName := tableName + "_" + job.suffix
	csvFile, incrementCount, err := job.CreateCsvFile(fileName, rows)
	if err != nil {
		return incrementCount, err
	}
	zipFile, err := job.CreateZipFile(fileName, csvFile)
	if err != nil {
		return incrementCount, err
	}
	localPath := "./storage/data_fabric/" + fileName + ".zip"
	err = job.uploadFile(c, zipFile, localPath)
	csvPath := "./storage/data_fabric/" + fileName + ".csv"
	_ = os.Remove(csvPath)
	_ = os.Remove(localPath)
	return incrementCount, err
}

// CreateCsvFile 创建CSV文件
func (job *SyncDataFabric) CreateCsvFile(csvName string, rows *sql.Rows) (csvFile *os.File, incrementCount int64, err error) {
	defer rows.Close()
	// 创建一个新的 CSV 文件
	csvFileName := "./storage/data_fabric/" + csvName + ".csv"
	csvFile, err = os.Create(csvFileName)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return nil, incrementCount, err
	}
	defer func(csvFile *os.File) {
		if closeErr := csvFile.Close(); err == nil {
			err = closeErr
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
	if err = rows.Err(); err != nil {
		return nil, incrementCount, err
	}
	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		return nil, incrementCount, err
	}
	return csvFile, incrementCount, nil
}

// CreateZipFile 创建ZIP文件
func (job *SyncDataFabric) CreateZipFile(zipName string, csvFile *os.File) (zipFile *os.File, err error) {
	// 创建一个新的 ZIP 文件
	zipFileName := "./storage/data_fabric/" + zipName + ".zip"
	zipFile, err = os.Create(zipFileName)
	if err != nil {
		fmt.Println("Error creating ZIP file:", err)
		return nil, err
	}
	defer func(zipFile *os.File) {
		if closeErr := zipFile.Close(); err == nil {
			err = closeErr
		}
	}(zipFile)
	// 创建一个 ZIP 写入器
	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		if closeErr := zipWriter.Close(); err == nil {
			err = closeErr
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

func (job *SyncDataFabric) uploadFile(c ctx.Context, zipFile *os.File, localFilePath string) error {
	job.transfer.upload(localFilePath)
	// Each target records its own failures; keep exporting for the other target.
	return nil
}

func uploadDataFabricFile(client *sftp.SFTP, remotePath, localFilePath string) error {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := client.Client.OpenFile(path.Join(remotePath, path.Base(localFilePath)), os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(dstFile, srcFile)
	closeErr := dstFile.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func (job *SyncDataFabric) CreateDoneFile(c ctx.Context, tableName string, columnNames []string, rows *sql.Rows, total int64) error {
	defer rows.Close()
	// 创建一个新的 CSV 文件
	fileName := tableName + "_" + job.suffix
	csvFileName := "./storage/data_fabric/" + fileName + "done.csv"
	doneFileName := "./storage/data_fabric/" + fileName + ".done"
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
	if err = rows.Err(); err != nil {
		_ = csvFile.Close()
		return err
	}
	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		_ = csvFile.Close()
		return err
	}
	err = csvFile.Close()
	if err != nil {
		fmt.Println(" CSV file close:", err)
		return err
	}
	err = os.Rename(csvFileName, doneFileName)
	if err != nil {
		fmt.Println(" CSV file rename err:", err)
		return err
	}
	err = job.uploadDoneFile(c, doneFileName)
	_ = os.Remove(doneFileName)
	return err
}
func (job *SyncDataFabric) uploadDoneFile(c ctx.Context, localFilePath string) error {
	job.transfer.upload(localFilePath)
	return nil
}

func (job *SyncDataFabric) createRdy(c ctx.Context) error {
	localRdyName := job.suffix + `.rdy`
	rdyPath := `./storage/data_fabric/` + localRdyName
	t, err := os.Create(rdyPath)
	if err != nil {
		return err
	}
	defer os.Remove(rdyPath)
	if err = t.Close(); err != nil {
		return err
	}
	return job.uploadDoneFile(c, rdyPath)
}
