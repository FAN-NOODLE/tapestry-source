package externalSearchDao

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	customerSearchReq "micro-scrm/internal/app/http/request/customer_search"
	customerSearchResp "micro-scrm/internal/app/http/response/customer_search"
	"micro-scrm/internal/app/models"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"time"
)

type ExternalCrowd struct {
	dao       dao.Dao
	tableName string
	model     externalSearchModel.ExternalCrowdModel
}

func NewExternalCrowd() *ExternalCrowd {
	return &ExternalCrowd{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: externalSearchModel.ExternalCrowdModel{}.TableName(),
		model:     externalSearchModel.ExternalCrowdModel{},
	}
}

func (d *ExternalCrowd) AddCrowd(c ctx.Context, req customerSearchReq.AddCrowdReq) (int64, error) {
	var total int64
	id := app.UID()
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).Model(&d.model).Where("member_id =?", req.MemberId).Count(&total).Error
	if err != nil {
		return id, err
	}
	if total >= externalSearchModel.ExternalCrowdNum {
		return id, errors.New("最多只能添加5个人群包")
	}
	nowTime := time.Now().Unix()
	data := externalSearchModel.ExternalCrowdModel{
		BaseModel: models.BaseModel{
			ID:        id,
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
			DeletedAt: 0,
		},
		MemberID:       req.MemberId,
		Name:           req.Name,
		SearchConnJson: helper.MarshalToString(req.Search),
		Sort:           req.Sort,
	}
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).Model(&d.model).Create(&data).Error
	return id, err
}

func (d *ExternalCrowd) AddCrowdData(c ctx.Context, enterpriseId int64, crowdId int64, externalIds []int64) error {
	nowTime := time.Now().Unix()
	crowdDatas := make([]externalSearchModel.ExternalCrowdDataModel, 0)
	for _, value := range externalIds {
		crowdData := externalSearchModel.ExternalCrowdDataModel{
			CrowdId:            crowdId,
			WorkExternalUserId: value,
			IsContact:          customerSearchReq.NotContacted,
			CreatedAt:          nowTime,
		}
		crowdDatas = append(crowdDatas, crowdData)
	}

	err := gormx.GetCkInstance().WithContext(c).Transaction(func(tx *gorm.DB) error {
		return tx.Scopes(func(db *gorm.DB) *gorm.DB {
			return db.Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s_local", enterpriseId, externalSearchModel.ExternalCrowdDataModel{}.TableName()))
		}).CreateInBatches(crowdDatas, len(crowdDatas)).Error
	})
	return err
}

func (d *ExternalCrowd) GetOne(c ctx.Context, enterpriseId int64, query interface{}, arg ...interface{}) (externalSearchModel.ExternalCrowdModel, error) {
	var info externalSearchModel.ExternalCrowdModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, arg...).Find(&info).Error
	return info, err
}

func (d *ExternalCrowd) GetMany(c ctx.Context, enterpriseId int64, query interface{}, arg ...interface{}) ([]externalSearchModel.ExternalCrowdModel, error) {
	var list []externalSearchModel.ExternalCrowdModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, arg...).Find(&list).Error
	return list, err
}

func (d *ExternalCrowd) List(c ctx.Context, req customerSearchReq.CrowdListReq) ([]externalSearchModel.ExternalCrowdModel, error) {
	var list []externalSearchModel.ExternalCrowdModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).Where("member_id =?", req.MemberId).Order("sort asc,created_at desc").Find(&list).Error
	return list, err
}
func (d *ExternalCrowd) ListByConn(c ctx.Context, req customerSearchReq.CrowdInfoReq) ([]externalSearchModel.ExternalCrowdDataModel, error) {
	var list []externalSearchModel.ExternalCrowdDataModel
	query := gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", req.EnterpriseID, externalSearchModel.ExternalCrowdDataModel{}.TableName())).
		Where("crowd_id = ?", req.CrowdId)
	if req.ContactedStatus > 0 {
		query = query.Where("is_contact = ?", req.ContactedStatus)
	}
	err := query.Select("work_external_user_id,is_contact").Find(&list).Error
	return list, err
}

func (d *ExternalCrowd) PageListByExternalConn(c ctx.Context, req customerSearchReq.CrowdInfoReq, externalQuery customerSearchResp.SearchQueryResp) ([]externalSearchModel.ExternalCrowdDataModel, int64, error) {
	var (
		total int64
		err   error
	)
	list := make([]externalSearchModel.ExternalCrowdDataModel, 0)
	query := gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", req.EnterpriseID, externalSearchModel.ExternalCrowdDataModel{}.TableName())).
		Where("crowd_id = ?", req.CrowdId)
	if req.ContactedStatus > 0 {
		query = query.Where("is_contact = ?", req.ContactedStatus)
	}
	if len(externalQuery.ExternalExtra.SqlSlice) > 0 {
		externalExtraTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_extra", req.EnterpriseID)
		query = query.Where(fmt.Sprintf("work_external_user_id in (SELECT id FROM %s WHERE %s)", externalExtraTable, externalQuery.ExternalExtra.SqlConn),
			externalQuery.ExternalExtra.SqlSlice...)
	}
	if len(externalQuery.RelationShip.SqlSlice) > 0 {
		relationshipTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.external_relationship", req.EnterpriseID)
		query = query.Where(fmt.Sprintf("work_external_user_id in (SELECT work_external_user_id FROM %s WHERE %s)", relationshipTable, externalQuery.RelationShip.SqlConn),
			externalQuery.RelationShip.SqlSlice...)
	}
	err = query.Count(&total).Error
	if err != nil || req.OnlyCount || total == 0 {
		return list, total, err
	}
	if req.Page > 0 && req.Size > 0 {
		query = query.Offset((req.Page - 1) * req.Size).Limit(req.Size)
	}
	err = query.Select("work_external_user_id,is_contact").Order("is_contact desc,work_external_user_id asc").Find(&list).Error
	return list, total, err
}

func (d *ExternalCrowd) Edit(c ctx.Context, req customerSearchReq.CrowdEditReq) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).Where("id =?", req.CrowdId).Updates(map[string]interface{}{
		"name": req.Name,
	}).Error

}

func (d *ExternalCrowd) Del(c ctx.Context, req customerSearchReq.CrowdEditReq) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).Where("id =?", req.CrowdId).Delete(&d.model).Error

}

// UpdateSort 更新人群包排序
func (d *ExternalCrowd) UpdateSort(c ctx.Context, req customerSearchReq.CrowdUpdateSortReq) error {
	if len(req.SortList) == 0 {
		return nil
	}
	for _, info := range req.SortList {
		d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).Where("id =?", info.CrowdId).Updates(map[string]interface{}{
			"sort": info.Sort,
		})
	}
	return nil
}

func (d *ExternalCrowd) UpdateContact(c ctx.Context, req customerSearchReq.CrowdUpdateContactReq) error {
	return gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s_local", req.EnterpriseId, externalSearchModel.ExternalCrowdDataModel{}.TableName())).
		Where(map[string]interface{}{
			"crowd_id":              req.CrowdId,
			"work_external_user_id": req.WorkExternalUserId,
		}).Updates(map[string]interface{}{
		"is_contact": customerSearchReq.Contacted,
	}).Error
}

func (d *ExternalCrowd) BatchUpdateContact(c ctx.Context, req customerSearchReq.BatchCrowdUpdateContactReq) error {
	return gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s_local", req.EnterpriseId, externalSearchModel.ExternalCrowdDataModel{}.TableName())).
		Where(map[string]interface{}{
			"crowd_id":              req.CrowdId,
			"work_external_user_id": req.WorkExternalUserIds,
		}).Updates(map[string]interface{}{
		"is_contact": customerSearchReq.Contacted,
	}).Error
}

func (d *ExternalCrowd) PageListByExternalConnSw(c ctx.Context, req customerSearchReq.CrowdInfoSwReq, externalQuery customerSearchResp.SearchQuerySwResp) ([]externalSearchModel.ExternalCrowdDataModel, int64, error) {
	var (
		total int64
		err   error
	)
	list := make([]externalSearchModel.ExternalCrowdDataModel, 0)
	query := gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", req.EnterpriseID, externalSearchModel.ExternalCrowdDataModel{}.TableName())).
		Where("crowd_id = ?", req.CrowdId)
	if req.ContactedStatus > 0 {
		query = query.Where("is_contact = ?", req.ContactedStatus)
	}
	if len(externalQuery.External.SqlSlice) > 0 {
		externalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", req.EnterpriseID)
		query = query.Where(fmt.Sprintf("work_external_user_id in (SELECT id FROM %s WHERE %s)", externalTable, externalQuery.External.SqlConn),
			externalQuery.External.SqlSlice...)
	}
	if len(externalQuery.RelationShip.SqlSlice) > 0 {
		relationshipTable := "micro_scrm_ck.work_external_user_relation"
		query = query.Where(fmt.Sprintf("work_external_user_id in (SELECT work_external_user_id FROM %s WHERE %s)", relationshipTable, externalQuery.RelationShip.SqlConn),
			externalQuery.RelationShip.SqlSlice...)
	}
	err = query.Count(&total).Error
	if err != nil || req.OnlyCount || total == 0 {
		return list, total, err
	}
	if req.Page > 0 && req.Size > 0 {
		query = query.Offset((req.Page - 1) * req.Size).Limit(req.Size)
	}
	err = query.Select("work_external_user_id,is_contact").Order("is_contact desc,work_external_user_id asc").Find(&list).Error
	return list, total, err
}

func (d *ExternalCrowd) ListByConnSw(c ctx.Context, req customerSearchReq.CrowdInfoSwReq) ([]externalSearchModel.ExternalCrowdDataModel, error) {
	var list []externalSearchModel.ExternalCrowdDataModel
	query := gormx.GetCkInstance().WithContext(c).
		Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", req.EnterpriseID, externalSearchModel.ExternalCrowdDataModel{}.TableName())).
		Where("crowd_id = ?", req.CrowdId)
	if req.ContactedStatus > 0 {
		query = query.Where("is_contact = ?", req.ContactedStatus)
	}
	err := query.Select("work_external_user_id,is_contact").Find(&list).Error
	return list, err
}
