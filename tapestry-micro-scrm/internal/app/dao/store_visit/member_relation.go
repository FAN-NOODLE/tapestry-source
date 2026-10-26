package storeVisitDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
)

type MemberRelationDao struct {
	dao.EnterpriseRepo[storeVisitModel.StoreVisitMemberRelationModel]
	model     *storeVisitModel.StoreVisitMemberRelationModel
	tableName string
}

func NewMemberRelationDao() *MemberRelationDao {
	return &MemberRelationDao{
		EnterpriseRepo: dao.EnterpriseRepo[storeVisitModel.StoreVisitMemberRelationModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(storeVisitModel.StoreVisitMemberRelationModel),
		tableName: storeVisitModel.StoreVisitMemberRelationModel{}.TableName(),
	}
}

func (d MemberRelationDao) GetVisitIdsByMemberIds(c ctx.Context, enterpriseId int64, memberIds []int64) ([]int64, error) {
	var visitIds []int64
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Where(map[string]interface{}{
		"member_id": memberIds,
		"type":      storeVisitModel.VisitCreator,
	}).Group("visit_id").Pluck("visit_id", &visitIds).Error
	return visitIds, err
}

func (d MemberRelationDao) GetVisitIdsByDataType(c ctx.Context, req storeVisitReq.DashboardTreeReq) ([]int64, error) {
	var visitIds []int64
	query := d.WithContext(c).Scopes(d.BranchDataBase(req.EnterpriseId, d.model.TableName()))
	switch req.DataType {
	case storeVisitReq.DataTypeCreator:
		//创建者
		query.Where("member_id = ? AND type = ?", req.MemberId, storeVisitModel.VisitCreator)
	case storeVisitReq.DataTypeAll:
		//所有数据
	default:
		//有关联的数据
		query.Where("member_id = ?", req.MemberId)
	}
	err := query.Group("visit_id").Pluck("visit_id", &visitIds).Error
	return visitIds, err
}
