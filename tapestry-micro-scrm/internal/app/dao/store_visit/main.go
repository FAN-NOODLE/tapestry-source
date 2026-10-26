package storeVisitDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	storeVisitResp "micro-scrm/internal/app/http/response/store_visit"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
)

type MainDao struct {
	dao.EnterpriseRepo[storeVisitModel.StoreVisitMainModel]
	model     *storeVisitModel.StoreVisitMainModel
	tableName string
}

func NewMainDao() *MainDao {
	return &MainDao{
		EnterpriseRepo: dao.EnterpriseRepo[storeVisitModel.StoreVisitMainModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(storeVisitModel.StoreVisitMainModel),
		tableName: storeVisitModel.StoreVisitMainModel{}.TableName(),
	}
}

func (d *MainDao) GetStoreVisitList(c ctx.Context, req storeVisitReq.PlanListReq) ([]storeVisitResp.StoreVisitMember, int64, error) {
	var (
		result []storeVisitResp.StoreVisitMember
		total  int64
	)
	relTableName := d.Dao.GetTable(storeVisitModel.StoreVisitMemberRelationModel{}.TableName(), req.EnterpriseId)

	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model).
		Joins("left join " + relTableName + " as rel on store_visit_main.id=rel.visit_id").
		Select("store_visit_main.id,store_visit_main.date,store_visit_main.time_span,store_visit_main.third_id,store_visit_main.remark,store_visit_main.status,store_visit_main.updated_at")
	if req.StartTime > 0 {
		query.Where("unix_time > ? ", req.StartTime)
	}
	if req.EndTime > 0 {
		query.Where("unix_time < ? ", req.EndTime)
	}
	if req.Status == storeVisitReq.StoreVisitStatusVisiting {
		query.Where("status in ? ", []int{storeVisitModel.VisitStatusPlan, storeVisitModel.VisitVisiting})
	} else if req.Status == storeVisitReq.StoreVisitStatusVisited {
		query.Where("status = ? ", []int{storeVisitModel.VisitComplete})
	}
	if req.ThirdId > 0 {
		query.Where("third_id = ? ", req.ThirdId)
	}
	switch req.DataType {
	case storeVisitReq.DataTypeCreator:
		//创建者
		query.Where("rel.member_id = ? AND rel.type = ?", req.MemberId, storeVisitModel.VisitCreator)
	case storeVisitReq.DataTypeAll:
		//所有数据
	default:
		//有关联的数据
		query.Where("rel.member_id = ?", req.MemberId)
	}
	offset := (req.Page - 1) * req.Size
	err := query.Group("store_visit_main.id").
		Order("store_visit_main.date desc").Count(&total).
		Offset(offset).Limit(req.Size).
		Scan(&result).Error
	return result, total, err
}

func (d *MainDao) TotalCountByStore(c ctx.Context, req storeVisitReq.DashboardTreeReq) ([]storeVisitResp.TotalCountByStore, error) {
	var result = make([]storeVisitResp.TotalCountByStore, 0)
	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model)
	if req.StartTime > 0 {
		query.Where("unix_time > ? ", req.StartTime)
	}
	if req.EndTime > 0 {
		query.Where("unix_time < ? ", req.EndTime)
	}
	if len(req.VisitIds) > 0 {
		query.Where("id in ?", req.VisitIds)
	}
	err := query.Select("third_id,count(1) as total_count").
		Group("third_id").
		Scan(&result).Error
	return result, err
}

func (d *MainDao) TotalCountByStoreAndCreator(c ctx.Context, req storeVisitReq.DashboardExportEmail, thirdIds []int) ([]storeVisitResp.TotalCountByStoreAndMember, error) {
	var result = make([]storeVisitResp.TotalCountByStoreAndMember, 0)
	relTableName := d.Dao.GetTable(storeVisitModel.StoreVisitMemberRelationModel{}.TableName(), req.EnterpriseId)
	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model).
		Joins("left join "+relTableName+" as rel on store_visit_main.id = rel.visit_id").
		Where("rel.type = ?", storeVisitModel.VisitCreator)
	if req.StartTime > 0 {
		query.Where("store_visit_main.unix_time > ? ", req.StartTime)
	}
	if req.EndTime > 0 {
		query.Where("store_visit_main.unix_time < ? ", req.EndTime)
	}
	if len(req.VisitIds) > 0 {
		query.Where("store_visit_main.id in ?", req.VisitIds)
	}
	if len(thirdIds) > 0 {
		query.Where("store_visit_main.third_id in ?", thirdIds)
	}
	err := query.Select("store_visit_main.third_id,rel.member_id,count(1) as total_count").
		Group("store_visit_main.third_id,rel.member_id").
		Scan(&result).Error
	return result, err
}

func (d *MainDao) GetAdminStoreVisitList(c ctx.Context, req storeVisitReq.AdminStoreVisitListReq, thirdIds []int, ids []int64) ([]storeVisitModel.StoreVisitMainModel, int64, error) {
	result := make([]storeVisitModel.StoreVisitMainModel, 0)
	var total int64
	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(thirdIds) > 0 {
			db.Where("third_id in ?", thirdIds)
		}
		if len(ids) > 0 {
			db.Where("id in ?", ids)
		}
		if req.StartTime > 0 && req.EndTime > 0 {
			db.Where("created_at between ? and ?", req.StartTime, req.EndTime)
		}
		if req.VisitStartTime > 0 && req.VisitEndTime > 0 {
			db.Where("unix_time between ? and ?", req.VisitStartTime, req.VisitEndTime)
		}
		if req.Status > 0 {
			db.Where("status = ? ", req.Status)
		}
		return db
	}).Unscoped()
	offset := (req.Page - 1) * req.Size
	err := query.Count(&total).
		Offset(offset).Limit(req.Size).
		Scan(&result).Error
	return result, total, err
}

func (d *MainDao) GetStoreVisitDashboardByStore(c ctx.Context, req storeVisitReq.AdminStoreDashboardReq, thirdIds []int, memberIds []int64, isAll bool) ([]storeVisitResp.AdminDashboardStore, int64, error) {
	var (
		result []storeVisitResp.AdminDashboardStore
		total  int64
		err    error
	)
	relTableName := d.Dao.GetTable(storeVisitModel.StoreVisitMemberRelationModel{}.TableName(), req.EnterpriseId)

	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model).
		Joins("left join " + relTableName + " as rel on store_visit_main.id=rel.visit_id").
		Select("rel.member_id,store_visit_main.third_id,SUM(store_visit_main.status=3) as finish_num,SUM(1) as visit_num").
		Scopes(func(db *gorm.DB) *gorm.DB {
			db.Where("rel.type =?", storeVisitModel.VisitCreator)
			if len(thirdIds) > 0 {
				db.Where("store_visit_main.third_id in ?", thirdIds)
			}
			if len(memberIds) > 0 {
				db.Where("rel.member_id in ?", memberIds)
			}
			if req.StartTime > 0 && req.EndTime > 0 {
				db.Where("unix_time between ? and ?", req.StartTime, req.EndTime)
			}
			return db
		}).Group("store_visit_main.third_id")
	offset := (req.Page - 1) * req.Size
	if !isAll {
		err = query.Count(&total).
			Offset(offset).Limit(req.Size).Order("visit_num desc").
			Scan(&result).Error
	} else {
		err = query.Scan(&result).Error
	}
	return result, total, err
}

func (d *MainDao) GetStoreVisitDashboardByStoreVisitNum(c ctx.Context, req storeVisitReq.AdminStoreDashboardReq, thirdIds []int, memberIds []int64) (int64, error) {
	var (
		result []storeVisitResp.AdminDashboardStore
		total  int64
		err    error
	)
	relTableName := d.Dao.GetTable(storeVisitModel.StoreVisitMemberRelationModel{}.TableName(), req.EnterpriseId)
	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model).
		Joins("left join " + relTableName + " as rel on store_visit_main.id=rel.visit_id").
		Select("SUM(1) as visit_num").Scopes(func(db *gorm.DB) *gorm.DB {
		db.Where("rel.type =?", storeVisitModel.VisitCreator)
		if len(thirdIds) > 0 {
			db.Where("store_visit_main.third_id in ?", thirdIds)
		}
		if len(memberIds) > 0 {
			db.Where("rel.member_id in ?", memberIds)
		}
		if req.StartTime > 0 && req.EndTime > 0 {
			db.Where("unix_time between ? and ?", req.StartTime, req.EndTime)
		}
		return db
	}).Group("store_visit_main.third_id")
	err = query.Scan(&result).Error
	if err != nil {
		return 0, err
	}
	for _, v := range result {
		total += v.VisitNum
	}
	return total, err
}

func (d *MainDao) GetInfoPreloadCreator(c ctx.Context, req storeVisitReq.VisitInfoReq) (storeVisitModel.StoreVisitMainModel, error) {
	var result storeVisitModel.StoreVisitMainModel
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model).
		Where("id = ?", req.ID).
		Preload("MemberList", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.BranchTable(storeVisitModel.StoreVisitMemberRelationModel{}.TableName(), req.EnterpriseId)).
				Where("type = ?", storeVisitModel.VisitCreator)
		}).
		Find(&result).Error
	return result, err
}
