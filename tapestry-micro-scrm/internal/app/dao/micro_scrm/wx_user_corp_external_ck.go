package microScrmDao

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	customerReq "micro-scrm/internal/app/http/request/customer"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type WxUserCorpExternalCk struct {
	dao       dao.Dao
	tableName string
}

func NewWxUserCorpExternalCk() *WxUserCorpExternalCk {
	return &WxUserCorpExternalCk{
		dao:       dao.Dao{DB: gormx.GetCkInstance()},
		tableName: microScrmModel.WxUserCorpExternalModel{}.TableName(),
	}
}

func (d *WxUserCorpExternalCk) GetPagination(c app.Ctx, cond customerReq.UserList) ([]microScrmModel.WxUserCorpExternalModel, int64, error) {
	var (
		resp  []microScrmModel.WxUserCorpExternalModel
		total int64
		err   error
	)

	userCorpExternalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", cond.EnterpriseID)

	query := gormx.GetCkInstance().WithContext(c).Table(userCorpExternalTable).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if cond.Name != "" {
				db.Where("nick Like ?", "%"+cond.Name+"%").Or("name Like ?", "%"+cond.Name+"%")
			}
			if cond.Mobile != "" {
				db.Where("mobile = ?", cond.Mobile)
			}

			if cond.Remark != "" {
				db.Where("id in (?)", func() *gorm.DB {
					return gormx.GetCkInstance().WithContext(c).
						Table("micro_scrm_ck.work_external_user_relation").
						Select("work_external_user_id").
						Where("enterprise_id = ?", cond.EnterpriseID).
						Where("remark like ?", "%"+cond.Remark+"%")
				})
			}

			if len(cond.TagList) > 0 {
				db.Where("id in (?)", func() *gorm.DB {
					return gormx.GetCkInstance().WithContext(c).
						Table("micro_scrm_ck.work_external_user_relation").
						Distinct("work_external_user_id").
						Where("id in (?)", func() *gorm.DB {
							return gormx.GetCkInstance().WithContext(c).
								Distinct("relation_id").
								Table("micro_scrm_ck.work_external_user_tag").
								Where("external_tag_id in (?)", cond.TagList).
								Where("enterprise_id = ?", cond.EnterpriseID).
								Where("deleted_at = 0")
						}())
				}())
			}

			if len(cond.ImpressionTagList) > 0 {
				db.Where("id in (?)", func() *gorm.DB {
					return gormx.GetCkInstance().WithContext(c).
						Table("micro_scrm_ck.impression_tag_relation").
						Distinct("external_user_id").
						Where("tag_id in (?)", cond.ImpressionTagList).
						Where("enterprise_id = ?", cond.EnterpriseID)
				}())
			}

			return db
		})

	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size
	err = query.Order("id DESC,created_at DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error

	return resp, total, err
}

func (d *WxUserCorpExternalCk) GetOne(c app.Ctx, query interface{}, enterpriseId int64) (microScrmModel.WxUserCorpExternalModel, error) {
	var result microScrmModel.WxUserCorpExternalModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Table(fmt.Sprintf("%s%d_ck.%s", app.Cfg().DB.DBPrefix, enterpriseId, d.tableName))
	}).Where(query).Find(&result).Error
	return result, err
}

func (d *WxUserCorpExternalCk) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.WxUserCorpExternalModel, error) {
	var result []microScrmModel.WxUserCorpExternalModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Table(fmt.Sprintf("%s%d_ck.%s", app.Cfg().DB.DBPrefix, enterpriseId, d.tableName))
	}).Where(query, args...).Find(&result).Error
	return result, err
}

// GetTagList 标签详情.标签列表
func (d *WxUserCorpExternalCk) GetTagList(c app.Ctx, req tagsReq.InfoReq) ([]microScrmModel.WxUserCorpExternalModel, int64, error) {
	var (
		resp  []microScrmModel.WxUserCorpExternalModel
		total int64
		err   error
	)

	userCorpExternalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", req.EnterpriseID)
	query := gormx.GetCkInstance().WithContext(c).Table(userCorpExternalTable).
		Scopes(func(db *gorm.DB) *gorm.DB {
			db.Where("deleted_at = 0")
			if req.TagType == 0 {
				db.Where("id in (?)", func() *gorm.DB {
					return gormx.GetCkInstance().WithContext(c).
						Table("micro_scrm_ck.work_external_user_relation").
						Select("DISTINCT work_external_user_id").
						Where("id in (?)", func() *gorm.DB {
							return gormx.GetCkInstance().WithContext(c).
								Table("micro_scrm_ck.work_external_user_tag").
								Select("relation_id").
								Where("external_tag_id", req.TagID).
								Where("deleted_at = 0")
						}())
				}())
			} else {
				db.Where("id in (?)", func() *gorm.DB {
					return gormx.GetCkInstance().WithContext(c).
						Table("micro_scrm_ck.impression_tag_relation").
						Select("DISTINCT external_user_id").
						Where("tag_id", req.TagID).
						Where("deleted_at = 0")
				}())
			}
			return db
		})

	query.Count(&total)

	offset := (req.Page - 1) * req.Size
	err = query.Limit(req.Size).Offset(offset).Find(&resp).Error

	return resp, total, err
}

func (d *WxUserCorpExternalCk) GetManyByPage(c app.Ctx, enterpriseId int64, size int, query interface{}, args ...interface{}) ([]microScrmModel.WxUserCorpExternalModel, error) {
	var result []microScrmModel.WxUserCorpExternalModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Table(fmt.Sprintf("%s%d_ck.%s", app.Cfg().DB.DBPrefix, enterpriseId, d.tableName))
	}).Where(query, args...).Limit(size).Find(&result).Error
	return result, err
}
