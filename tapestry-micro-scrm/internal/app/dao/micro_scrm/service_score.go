package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type ServiceScoreDao struct {
	dao       dao.Dao
	tableName string
}

func NewServiceScoreDao() *ServiceScoreDao {
	return &ServiceScoreDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.ReserveServiceScoreModel{}.TableName(),
	}
}

func (d *ServiceScoreDao) GetList(c app.Ctx, enterpriseId int64, serviceScoreListReq reserveRequest.ServiceScoreListReq) (serviceScoreList []microScrmModel.ReserveServiceScoreModel, total int64, err error) {

	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId))
	//Joins(fmt.Sprintf("LEFT JOIN  micro_scrm_enterprise_%d.reserve_service_score rss ON reserve_order.reserve_number=rss.reserve_number", enterpriseId)).

	if len(serviceScoreListReq.MemberIds) > 0 {
		query.Where("member_id in ?", serviceScoreListReq.MemberIds)
	}
	if len(serviceScoreListReq.StoreIds) > 0 {
		query.Where("store_id in ?", serviceScoreListReq.StoreIds)
	}
	if len(serviceScoreListReq.Ids) > 0 {
		query.Where("id in ?", serviceScoreListReq.Ids)
	}
	if serviceScoreListReq.ReserveStartTime > 0 && serviceScoreListReq.ReserveEndTime > 0 {
		query.Where("reserve_time >= ? AND reserve_time <= ?", serviceScoreListReq.ReserveStartTime, serviceScoreListReq.ReserveEndTime)

	}
	if serviceScoreListReq.ScoreStartTime > 0 && serviceScoreListReq.ScoreEndTime > 0 {
		query.Where("created_at >= ? AND created_at <= ?", serviceScoreListReq.ScoreStartTime, serviceScoreListReq.ScoreEndTime)

	}

	err = query.Count(&total).
		Order("created_at DESC").
		Offset((serviceScoreListReq.Page - 1) * serviceScoreListReq.Size).
		Limit(serviceScoreListReq.Size).
		Find(&serviceScoreList).Error

	return serviceScoreList, total, err

}

// UserMakeScore 用户评分，存在的就更改
func (d *ServiceScoreDao) UserMakeScore(c app.Ctx, enterpriseId int64, data microScrmModel.ReserveServiceScoreModel) error {
	info, err := d.GetScoreInfo(c, enterpriseId, reserveRequest.UserOrderBase{
		OrderBase: reserveRequest.OrderBase{
			ID:            data.ReserveOrderID,
			ReserveNumber: data.ReserveNumber,
		},
		EnterpriseBase: reserveRequest.EnterpriseBase{
			EnterpriseId: enterpriseId,
		},
		UnionId: data.UserUnionid,
	})
	if err != nil {
		return err
	}
	if info.ID > 0 {
		updateData := map[string]interface{}{
			"service_score": data.ServiceScore,
			"store_score":   data.StoreScore,
			"remark":        data.Remark,
			"flow_score":    data.FlowScore,
			"reserve_time":  data.ReserveTime,
		}
		return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("id=?", info.ID).Updates(updateData).Error
	}
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
}

// GetScoreInfo 获取最新评分信息一条信息
func (d *ServiceScoreDao) GetScoreInfo(c app.Ctx, enterpriseId int64, request reserveRequest.UserOrderBase) (microScrmModel.ReserveServiceScoreModel, error) {
	var result microScrmModel.ReserveServiceScoreModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"reserve_order_id": request.ID,
		"reserve_number":   request.ReserveNumber,
		"user_unionid":     request.UnionId,
	}).Order("created_at DESC").Find(&result).Error
	return result, err
}

// UserEditScore 用户编辑预约单评分
func (d *ServiceScoreDao) UserEditScore(c app.Ctx, enterpriseId int64, request reserveRequest.UserScoreReq) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"reserve_order_id": request.ID,
		"reserve_number":   request.ReserveNumber,
		"user_unionid":     request.UnionId,
	}).Updates(map[string]interface{}{
		"service_score": request.ServiceScore,
		"store_score":   request.StoreScore,
		"remark":        request.Remark,
	}).Error
}
