package appJob

import (
	"fmt"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
)

type CosToCdn struct {
	name      string
	originUri string
	cdnUri    string
	db        *gorm.DB
}

func NewCosToCdn() *CosToCdn {
	env := config.CfgData.Env
	originUri := helper.RuleThree(env != "prod", "tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com", "tapestry-prod-1308608494.cos.ap-shanghai.myqcloud.com")
	cdnUri := helper.RuleThree(env != "prod", "cdn-cos-uat.coach.com.cn", "cdn-cos.coach.com.cn")
	db := gormx.GetInstance()
	return &CosToCdn{
		name:      "cos转cdn",
		originUri: originUri,
		cdnUri:    cdnUri,
		db:        db,
	}

}

func (job *CosToCdn) Run(c ctx.Context, jobHandler *xjob.JobHandler) error {
	microScrmTableMap := map[string][]string{
		"chat_group_msg_task":                  {"content"},
		"crm_sftp_downloaded_file":             {"file_url"},
		"dashboard_log":                        {"file_url"},
		"download_resource":                    {"source_data"},
		"drainage_welcomes":                    {"content"},
		"excitation_integral_rule_description": {"content"},
		"excitation_level_rule":                {"level_icon"},
		"excitation_medal_style":               {"light_icon", "unlit_icon"},
		"marketing_strategy_channels":          {"content"},
		"members":                              {"qr_poster_code"},
		"order_achievement_sftp_log":           {"file_url"},
		"reserve_order_data_acxiom_sftp_log":   {"file_url"},
		"source_file":                          {"url"},
		"source_friend":                        {"video_cover_url", "original_url", "compression_url"},
		"source_image":                         {"url"},
		"source_link":                          {"cover_url", "link_url"},
		"source_mini_program":                  {"cover_url"},
		"source_video":                         {"cover_url", "url"},
		"template":                             {"template_picture_url"},
	}
	for key, item := range microScrmTableMap {
		for _, v := range item {
			err := job.replace(c, "micro_scrm", key, v)
			if err != nil {
				fmt.Printf("执行更新数据库micro_scrm下的%s表报错%s", key, err.Error())
				continue
			}
		}
	}

	microScrmEnterpriseTableMap := map[string][]string{
		"activity_qrcode_main":            {"banner_url"},
		"activity_qrcode_point":           {"code_url"},
		"activity_qrcode_welcomes":        {"content"},
		"display_task":                    {"review_file"},
		"display_task_allot_data":         {"image_url"},
		"display_task_allot_remark":       {"avatar", "content"},
		"drainage_contact_group":          {"zip_url"},
		"drainage_contact_group_welcomes": {"content"},
		"external_pay_order":              {"tickets"},
		"invite_record":                   {"background_url"},
		"invite_template":                 {"background_url", "share_url"},
		"notification_crowd":              {"file_url"},
		"notification_info":               {"attachment", "page_url", "pdf_url", "zip_url"},
		"notification_info_merge":         {"cover"},
		"online_warehouse_order":          {"transfer_file", "share_poster_url"},
		"plan_task_conditions":            {"content"},
		"plan_task_import_file":           {"file_url", "download_url"},
		"plan_task_template":              {"template_content"},
		"plan_tasks":                      {"download_url"},
		"reserve_service":                 {"reserve_img"},
		"store_visit_main":                {"pdf_url"},
		"store_visit_media":               {"url", "video_image_url"},
		"user_group_file":                 {"file_url"},
	}
	enterpriseIds, err := GetEnterpriseIds(jobHandler)
	if err != nil {
		return err
	}

	for _, enterpriseId := range enterpriseIds {
		dataBaseName := fmt.Sprintf("%s_%d", "micro_scrm_enterprise", enterpriseId)
		for key, item := range microScrmEnterpriseTableMap {
			for _, v := range item {
				err := job.replace(c, dataBaseName, key, v)
				if err != nil {
					fmt.Printf("执行更新数据库%s下的%s表报错%s", dataBaseName, key, err.Error())
					continue
				}
			}
		}
	}

	microWorkTableMap := map[string][]string{
		"contact_qr_code": {"transform_url"},
		"store_worker":    {"qr_code"},
	}
	for key, item := range microWorkTableMap {
		for _, v := range item {
			err := job.replace(c, "micro_work", key, v)
			if err != nil {
				fmt.Printf("执行更新数据库micro_work下的%s表报错%s", key, err.Error())
				continue
			}
		}
	}
	return nil
}

func (job *CosToCdn) replace(c ctx.Context, databaseName, tableName, fieldName string) error {
	// 查询满足条件的数据总数
	var totalCount int
	countQuery := "SELECT COUNT(1) FROM " + databaseName + "." + tableName + " WHERE " + fieldName + " LIKE CONCAT('%', ?, '%')"
	err := job.db.Raw(countQuery, job.originUri).Scan(&totalCount).Error

	if err != nil {
		fmt.Printf("查询%s.%s的字段%s总数报错，%s\n", databaseName, tableName, fieldName, err.Error())
		return err
	}

	if totalCount == 0 {
		fmt.Printf("表%s.%s中无满足条件的数据\n", databaseName, tableName)
		return nil
	}

	// 计算分批次数
	batchSize := 1000
	batches := (totalCount + batchSize - 1) / batchSize
	nowUnix := time.Now().Unix()
	// 分批更新
	for i := 0; i < batches; i++ {
		offset := i * batchSize
		updateQuery := "UPDATE " + databaseName + "." + tableName + " SET " + fieldName + " = REPLACE(" + fieldName + ", ?, ?) ,updated_at = ? WHERE id IN (SELECT id FROM (SELECT id FROM " + databaseName + "." + tableName + " WHERE " + fieldName + " LIKE CONCAT('%', ?, '%') ORDER BY id DESC LIMIT ? OFFSET ?) AS temp)"
		err = job.db.Exec(updateQuery, job.originUri, job.cdnUri, nowUnix, job.originUri, batchSize, offset).Error
		if err != nil {
			fmt.Printf("替换%s.%s的字段%s第%d批报错，%s\n", databaseName, tableName, fieldName, i+1, err.Error())
			return err
		}
		fmt.Printf("替换%s.%s的字段%s第%d批完成\n", databaseName, tableName, fieldName, i+1)
	}

	return nil
}
