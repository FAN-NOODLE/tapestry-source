package microScrmDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	inviteRequest "micro-scrm/internal/app/http/request/invite"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type InviteRecordDao struct {
	dao       dao.Dao
	tableName string
}

func NewInviteRecordDao() *InviteRecordDao {
	return &InviteRecordDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.InviteRecordModel{}.TableName(),
	}
}

// GetRecordNumByInviteIds 邀请函关联的预约记录数
func (d *InviteRecordDao) GetRecordNumByInviteIds(c app.Ctx, enterpriseId int64, inviteId []int64) (recordNum []inviteRequest.OrderCountNum, err error) {
	reserveOrderTable := d.dao.GetTable("reserve_order", enterpriseId)
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("count(ro.id) AS num, invite_record.invite_template_id").
		Joins("LEFT JOIN "+reserveOrderTable+" ro ON invite_record.reserve_order_id=ro.id").
		Where("invite_record.invite_template_id in ?", inviteId).
		Where("ro.reserve_time > 0").
		Where("ro.is_commit = 1").
		Where("invite_record.deleted_at = 0").
		Where("ro.deleted_at = 0").
		Group("invite_record.invite_template_id").
		Find(&recordNum).Error
	return
}

func (d *InviteRecordDao) Create(c app.Ctx, enterpriseId int64, recordModel microScrmModel.InviteRecordModel) (microScrmModel.InviteRecordModel, error) {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&recordModel).Error
	return recordModel, err
}

func (d *InviteRecordDao) Update(c app.Ctx, enterpriseId int64, cond interface{}, data interface{}) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(cond).UpdateColumns(data).Error
}

func (d *InviteRecordDao) GetRecordInfo(c app.Ctx, enterpriseId, id int64) (data inviteRequest.RecordDetailInfo, err error) {

	reserveServiceTable := d.dao.GetTable("reserve_service", enterpriseId)
	inviteTemplateTable := d.dao.GetTable("invite_template", enterpriseId)

	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Joins("LEFT JOIN "+inviteTemplateTable+" it ON invite_record.invite_template_id=it.id").
		Joins("LEFT JOIN "+reserveServiceTable+" rs ON it.reserve_service_id=rs.id").
		Select("invite_record.`id`,invite_record.`member_id`,invite_record.`member_userid`,invite_record.`user_name`,invite_record.`user_unionid`,invite_record.`user_openid`,invite_record.`style_type`,invite_record.`background_url`,invite_record.`reserve_order_id`,invite_record.`send_content`,"+
			"it.`title_name`,it.`content`,it.`share_title`,it.`share_url`,it.`is_up`,it.`user_is_show`,it.`member_is_show`,"+
			"rs.`id` AS reserve_service_id,rs.`reserve_name` AS reserve_service_name,rs.`reserve_time` AS reserve_service_time,rs.`time_slot` AS time_slot_str").
		Where("invite_record.`id` = ? ", id).
		Scan(&data).Error

	return
}

func (d *InviteRecordDao) GetRecordOrderList(c app.Ctx, enterpriseId int64, reqList inviteRequest.RecordOrderListReq) (list []*inviteRequest.RecordOrder, total int64, err error) {

	reserveServiceTable := d.dao.GetTable("reserve_service", enterpriseId)
	inviteTemplateTable := d.dao.GetTable("invite_template", enterpriseId)
	reserveOrderTable := d.dao.GetTable("reserve_order", enterpriseId)

	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Joins("LEFT JOIN "+inviteTemplateTable+" it ON invite_record.invite_template_id=it.id").
		Joins("LEFT JOIN "+reserveServiceTable+" rs ON it.reserve_service_id=rs.id").
		Joins("LEFT JOIN "+reserveOrderTable+" ro ON invite_record.reserve_order_id=ro.id").
		Select("invite_record.`id`,invite_record.`member_id`,invite_record.`member_userid`,invite_record.`user_unionid`,invite_record.`user_openid`,invite_record.`style_type`,invite_record.`background_url`,"+
			"it.`title_name`,it.`content`,it.`share_title`,it.`share_url`,"+
			"rs.`id` AS reserve_service_id,rs.`reserve_name` AS reserve_service_name,rs.`reserve_time` AS reserve_service_time,"+
			"ro.member_name,ro.store_id,ro.store_name,ro.user_name,ro.reserve_time,ro.created_at").
		Where("invite_record.`deleted_at` = ? ", 0).
		Where("invite_record.`reserve_order_id` > ? ", 0)

	if reqList.InviteTemplateId != 0 {
		query.Where("it.id = ?", reqList.InviteTemplateId)
	}
	if reqList.ReserveServiceID != 0 {
		query.Where("it.reserve_service_id = ?", reqList.ReserveServiceID)
	}
	if len(reqList.ThirdIds) > 0 {
		query.Where("ro.store_id in (?) ", reqList.ThirdIds)
	}

	if reqList.ReserveServiceStartTime > 0 && reqList.ReserveServiceEndTime > 0 {
		query.Where("ro.reserve_time > ?  and ro.reserve_time < ?", reqList.ReserveServiceStartTime, reqList.ReserveServiceEndTime)
	}
	if reqList.OrderCreateStartTime > 0 && reqList.OrderCreateEndTime > 0 {
		query.Where("ro.created_at > ?  and ro.created_at < ?", reqList.OrderCreateStartTime, reqList.OrderCreateEndTime)
	}
	if reqList.TitleName != "" {
		query.Where("it.title_name like ? ", "%"+reqList.TitleName+"%")
	}
	if reqList.MemberName != "" {
		query.Where("ro.member_name like ? ", "%"+reqList.MemberName+"%")
	}
	if reqList.UserName != "" {
		query.Where("ro.user_name like ? ", "%"+reqList.UserName+"%")
	}

	query = query.Group("invite_record.id").Order("invite_record.id DESC")

	err = query.Count(&total).
		Offset((reqList.Page - 1) * reqList.Size).Limit(reqList.Size).
		Find(&list).Error

	return
}

// GetRecordNumByMemberIds  邀请函的发送数
func (d *InviteRecordDao) GetRecordNumByMemberIds(c app.Ctx, enterpriseId, startTime, endTime int64, memberIds []int64) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Where("member_id in ?", memberIds).
		Count(&total).Error
	return
}

func (d *InviteRecordDao) RecordDelete(c ctx.Context, enterpriseId, id int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("id=?", id).Delete(&microScrmModel.InviteRecordModel{}).Error
}
