package excitationDao

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"strings"
)

type PointsRecord struct {
	dao       dao.Dao
	model     *microScrmModel.ExcitationPointsRecordModel
	tableName string
}

func NewPointsRecord() *PointsRecord {
	return &PointsRecord{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(microScrmModel.ExcitationPointsRecordModel),
		tableName: microScrmModel.ExcitationPointsRecordModel{}.TableName(),
	}
}

// WithTx 切换db连接,用于事物
func (d *PointsRecord) WithTx(tx *gorm.DB) *PointsRecord {
	d.dao.DB = tx
	return d
}

func (d *PointsRecord) ReleaseTx() {
	d.dao.DB = gormx.GetInstance()
}

func (d *PointsRecord) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.ExcitationPointsRecordModel, error) {
	var result microScrmModel.ExcitationPointsRecordModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *PointsRecord) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.ExcitationPointsRecordModel, error) {
	var result []microScrmModel.ExcitationPointsRecordModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *PointsRecord) GetPagination(c app.Ctx, cond excitationReq.RecordListReq) ([]microScrmModel.ExcitationPointsRecordModel, int64, error) {
	var (
		total int64
		resp  []microScrmModel.ExcitationPointsRecordModel
		err   error
	)

	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).
		Scopes(func(db *gorm.DB) *gorm.DB {
			// 导购维度
			if len(cond.MemberIdArr) > 0 {
				db.Where("member_id in ?", cond.MemberIdArr)
			}
			// 积分维度
			if cond.IncentiveType > 0 {
				db.Where("incentive_type = ?", cond.IncentiveType)
			}
			// 积分来源
			if cond.SourceType > 0 {
				db.Where("source_type = ?", cond.SourceType)
			}
			// 规则名称
			if cond.IntegralRuleName != "" {
				db.Where("integral_rule_name like ?", "%"+strings.Trim(cond.IntegralRuleName, " ")+"%")
			}
			// 获得时间
			if cond.StartTime > 0 && cond.EndTime > 0 {
				db.Where("created_at between ? and ?", cond.StartTime, cond.EndTime)
			}
			return db
		})
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size
	if cond.HasCoins {
		query.Scopes(d.dao.BranchTable(d.tableName, cond.EnterpriseId)).
			Preload("CoinsRecord", func(db *gorm.DB) *gorm.DB {
				return db.Scopes(d.dao.BranchTable(microScrmModel.ExcitationCoinsRecordModel{}.TableName(), cond.EnterpriseId))
			})
	}
	err = query.Order("created_at DESC").Order("id DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return resp, total, err
}

func (d *PointsRecord) CreateData(c app.Ctx, enterpriseId int64, data microScrmModel.ExcitationPointsRecordModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Create(&data).Error
	return err
}

func (d *PointsRecord) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []microScrmModel.ExcitationPointsRecordModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *PointsRecord) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *PointsRecord) DeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model microScrmModel.ExcitationPointsRecordModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).Delete(&model).Error
	return err
}

// GetTopList 导购积分排行
func (d *PointsRecord) GetTopList(c app.Ctx, enterpriseId int64, queryType int, startTime, endTime int64) (map[int64]float64, error) {
	var (
		resMap = make(map[int64]float64)
	)
	var result []excitationResp.MemberSumPointsList
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("member_id,SUM(points) AS sum_points").
		Where("created_at between ? and ?", startTime, endTime).
		Where("type = ?", queryType).
		Group("member_id").
		Order("sum_points desc").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}
	for _, tmp := range result {
		resMap[tmp.MemberId] = tmp.SumPoints
	}
	return resMap, err
}

// GetSumPoints 获取某导购某时间段内的积分总和
func (d *PointsRecord) GetSumPoints(c app.Ctx, req excitationReq.SumPointsReq) (float64, error) {
	var sum float64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model).
		Where("member_id = ?", req.MemberId).
		Where("type = ?", req.Type).
		Where("created_at between ? and ?", req.StartTime, req.EndTime).
		Select("COALESCE(sum(points),0) as sum").
		Pluck("sum", &sum).Error
	return sum, err
}

// GetSumPointsByCategoryIds  获取某导购某时间段内多个分类的的积分总和
func (d *PointsRecord) GetSumPointsByCategoryIds(c app.Ctx, req excitationReq.SumPointsCategoryReq) (res []excitationResp.SumPointsResp, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model).
		Where("member_id = ?", req.MemberId).
		Where("type in ?", req.CategoryIds).
		Where("created_at between ? and ?", req.StartTime, req.EndTime).
		Select("COALESCE(sum(points),0) as sum_points, type as category_id").
		Group("type").
		Find(&res).Error
	return
}

// GetTopListByRemark 获取top
func (d *PointsRecord) GetTopListByRemark(c app.Ctx, enterpriseId int64, queryType int, inviteTop int, startTime, endTime int64, remark []string) ([]excitationResp.MemberSumPointsList, error) {
	var result []excitationResp.MemberSumPointsList
	var lastTopTotal int64
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("SUM(points) AS sum_points").
		Where("created_at between ? and ?", startTime, endTime).
		Where("type = ?", queryType).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(remark) > 0 {
				db.Where("points_remark in ?", remark)
			}
			return db
		}).
		Group("member_id").
		Order("sum_points desc")

	//获取top以及和top最后一位数量相同的其他导购
	newOrm := orm.Session(&gorm.Session{NewDB: true}).WithContext(c)
	err := orm.Offset(inviteTop-1).Limit(1).Pluck("sum_points", &lastTopTotal).Error
	err = newOrm.Select("member_id,SUM(points) AS sum_points").
		Having("total >= ?", lastTopTotal).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, err
}

// GetPointsTopByCategoryIds 获取所有导购时间段内部分类别的积分
func (d *PointsRecord) GetPointsTopByCategoryIds(c app.Ctx, enterpriseId int64, categoryIds []int, startTime int64, endTime int64) (list []models.MemberTopList, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.model.TableName(), enterpriseId)).
		Where("created_at BETWEEN ? AND ?  ", startTime, endTime).
		Where("type in ?", categoryIds).
		Where("deleted_at = 0").
		Select("member_id,SUM(points) AS total").
		Group("member_id").
		Order("total desc").
		Scan(&list).Error
	return list, err
}
