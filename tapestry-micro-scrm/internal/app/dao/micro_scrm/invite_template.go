package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	inviteRequest "micro-scrm/internal/app/http/request/invite"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type InviteTemplateDao struct {
	dao       dao.Dao
	tableName string
}

func NewInviteTemplateDao() *InviteTemplateDao {
	return &InviteTemplateDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.InviteTemplateModel{}.TableName(),
	}
}

func (d *InviteTemplateDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.InviteTemplateModel, error) {
	var result []microScrmModel.InviteTemplateModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *InviteTemplateDao) Create(c app.Ctx, enterpriseId int64, info microScrmModel.InviteTemplateModel) (err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&info).Error
	return
}

func (d *InviteTemplateDao) Update(c app.Ctx, enterpriseId, id int64, info interface{}) (err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where(map[string]interface{}{
			"id": id,
		}).
		Updates(&info).Error
	return
}

func (d *InviteTemplateDao) UpdateForReserveServiceId(c app.Ctx, enterpriseId, reserveServiceId int64, info interface{}) (err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where(map[string]interface{}{
			"reserve_service_id": reserveServiceId,
		}).
		Updates(&info).Error
	return
}

func (d *InviteTemplateDao) List(c app.Ctx, enterpriseId int64, reqList inviteRequest.TemplateListReq) (list []microScrmModel.InviteTemplateInfo, total int64, err error) {

	reserveServiceTable := d.dao.GetTable("reserve_service", enterpriseId)

	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Joins("LEFT JOIN " + reserveServiceTable + " rs ON invite_template.reserve_service_id=rs.id").
		Select("invite_template.`id`,invite_template.`title_name`,invite_template.`reserve_service_id`,invite_template.`status_type`,invite_template.`start_time`,invite_template.`end_time`,invite_template.`background_url`,invite_template.`content`,invite_template.`share_title`,invite_template.`share_url`,invite_template.`is_up`,invite_template.`user_is_show`,invite_template.`member_is_show`,invite_template.`created_at`,invite_template.`updated_at`," +
			"rs.`reserve_name` AS reserve_service_name,rs.`service_role` AS reserve_service_role").
		Where("invite_template.deleted_at=0")

	if reqList.ReserveServiceId != 0 {
		query.Where("invite_template.reserve_service_id = ?", reqList.ReserveServiceId)
	}
	if reqList.IsUp > 0 {
		query.Where("invite_template.is_up = ?", reqList.IsUp)
	}
	if reqList.Keyword != "" {
		query.Where("invite_template.title_name like ? ", "%"+reqList.Keyword+"%")
	}
	query = query.Group("invite_template.id").
		Order("invite_template.id DESC")

	err = query.Count(&total).
		Offset((reqList.Page - 1) * reqList.Size).Limit(reqList.Size).
		Find(&list).Error

	return
}

func (d *InviteTemplateDao) Info(c app.Ctx, enterpriseId int64, reqInfo inviteRequest.TemplateInfoReq) (info microScrmModel.InviteTemplateInfo, err error) {

	reserveServiceTable := d.dao.GetTable("reserve_service", enterpriseId)

	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Joins("LEFT JOIN "+reserveServiceTable+" rs ON invite_template.reserve_service_id=rs.id").
		Select("invite_template.`id`,invite_template.`title_name`,invite_template.`reserve_service_id`,invite_template.`status_type`,invite_template.`start_time`,invite_template.`end_time`,invite_template.`background_url`,invite_template.`content`,invite_template.`share_title`,invite_template.`share_url`,invite_template.`is_up`,invite_template.`user_is_show`,invite_template.`member_is_show`,"+
			"rs.`reserve_name` AS reserve_service_name,rs.`reserve_time` AS reserve_service_time,rs.`end_time` AS reserve_service_end_time").
		Where("invite_template.id = ?", reqInfo.ID).
		Group("invite_template.id").
		First(&info).Error

	return
}

func (d *InviteTemplateDao) Delete(c app.Ctx, enterpriseId int64, req inviteRequest.TemplateDeleteReq) (err error) {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where(map[string]interface{}{
			"id": req.ID,
		}).
		Updates(map[string]interface{}{
			"deleted_at": time.Now().Unix(),
		}).Error
}

func (d *InviteTemplateDao) ConfigList(c app.Ctx, enterpriseId int64, reqList inviteRequest.TemplateConfigListReq) (list []microScrmModel.InviteTemplateListInfo, total int64, err error) {

	reserveServiceTable := d.dao.GetTable("reserve_service", enterpriseId)
	reserveServiceRelationTable := d.dao.GetTable("reserve_service_relation", enterpriseId)

	nowTime := time.Now().Unix()
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Joins("LEFT JOIN "+reserveServiceTable+" rs ON invite_template.reserve_service_id=rs.id").
		Joins("LEFT JOIN "+reserveServiceRelationTable+" isr ON rs.id=isr.service_id").
		Select("invite_template.`id`,invite_template.`title_name`,invite_template.`reserve_service_id`,invite_template.`status_type`,invite_template.`start_time`,invite_template.`end_time`,invite_template.`background_url`,invite_template.`content`,invite_template.`share_title`,invite_template.`share_url`,invite_template.`is_up`,invite_template.`user_is_show`,invite_template.`member_is_show`,invite_template.`created_at`,invite_template.`updated_at`,"+
			"rs.`reserve_name` AS reserve_service_name,rs.`reserve_time` AS reserve_service_time").
		Where("invite_template.is_up = ?", microScrmModel.IsUp).
		Where("invite_template.start_time < ?", nowTime).
		Where("invite_template.end_time > ?", nowTime).
		Where("rs.service_role = 0 or "+
			"( rs.service_role = 1 and  "+"( isr.store_id = ?  or isr.member_id = ? ) )", reqList.ThirdId, reqList.MemberId)

	err = query.Order("invite_template.id DESC").
		Group("invite_template.id").
		Count(&total).
		Find(&list).Error

	return
}

func (d *InviteTemplateDao) ConfigDetail(c app.Ctx, enterpriseId int64, reqDetail inviteRequest.TemplateConfigDetailReq) (data microScrmModel.InviteTemplateDetailInfo, err error) {

	reserveServiceTable := d.dao.GetTable("reserve_service", enterpriseId)
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Joins("LEFT JOIN "+reserveServiceTable+" rs ON invite_template.reserve_service_id=rs.id").
		Select("invite_template.`id`,invite_template.`title_name`,invite_template.`reserve_service_id`,invite_template.`status_type`,invite_template.`start_time`,invite_template.`end_time`,invite_template.`background_url`,invite_template.`content`,invite_template.`share_title`,invite_template.`share_url`,invite_template.`is_up`,invite_template.`user_is_show`,invite_template.`member_is_show`,invite_template.`created_at`,invite_template.`updated_at`,"+
			"rs.`reserve_name` AS reserve_service_name,rs.`reserve_time` AS reserve_service_time").
		Where("invite_template.`id` = ? ", reqDetail.ID).
		First(&data).Error

	return
}

// UpdateStatus 上下架操作
func (d *InviteTemplateDao) UpdateStatus(c app.Ctx, enterpriseId int64, req inviteRequest.TemplateUpdateStatusReq) (err error) {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where(map[string]interface{}{
			"id": req.ID,
		}).
		Updates(map[string]interface{}{
			"is_up":      req.IsUp,
			"updated_at": time.Now().Unix(),
		}).Error
}
