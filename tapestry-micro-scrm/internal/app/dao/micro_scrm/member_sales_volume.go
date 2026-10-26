package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	memberReq "micro-scrm/internal/app/http/request/member"
	"micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type MemberSalesVolume struct {
	dao       dao.Dao
	tableName string
}

func NewMemberSalesVolume() *MemberSalesVolume {
	return &MemberSalesVolume{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.MemberSalesVolumeModel{}.TableName(),
	}
}

// GetMemberSalesVolume 获取导购销售额
func (d *MemberSalesVolume) GetMemberSalesVolume(c app.Ctx, req memberReq.MemberSalesVolumeReq) (map[int64]float64, error) {
	var returnMap = make(map[int64]float64)
	var resList []memberReq.MemberSalesVolume
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).
		Where("start_time BETWEEN ? AND ?", req.StartTime, req.EndTime).
		Group("member_id")
	selectStr := "ll_eshop+ll_livestream+pp_livestream+pp_post+remote_sales+cloud_store+link_to_ec"
	if req.SearchType == microScrmModel.ALL_SALES_VOLUME {
		selectStr += "+home_store+o2o+try_at_home"
	}
	err := orm.Select("member_id,sum(" + selectStr + ") as total").
		Find(&resList).Error
	for _, data := range resList {
		returnMap[data.MemberId] = data.Total
	}

	return returnMap, err
}

// GetMemberSalesVolumeByReserve  获取导购销售额(邀约相关只包含cs2000数据不能计算link_to_ec)
func (d *MemberSalesVolume) GetMemberSalesVolumeByReserve(c app.Ctx, req memberReq.MemberSalesVolumeReq) (map[int64]float64, error) {
	var returnMap = make(map[int64]float64)
	var resList []memberReq.MemberSalesVolume
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).
		Where("start_time BETWEEN ? AND ?", req.StartTime, req.EndTime).
		Group("member_id")
	selectStr := "ll_eshop+ll_livestream+pp_livestream+pp_post+remote_sales+cloud_store+home_store+o2o+try_at_home"
	err := orm.Select("member_id,sum(" + selectStr + ") as total").
		Find(&resList).Error
	for _, data := range resList {
		returnMap[data.MemberId] = data.Total
	}

	return returnMap, err
}

// GetOnlineRatio 获取导购线上销售占比
func (d *MemberSalesVolume) GetOnlineRatio(c app.Ctx, req memberReq.MemberSalesVolumeReq) (map[int64]float64, error) {
	var returnMap = make(map[int64]float64)
	var resList []memberReq.MemberSalesVolume
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).
		Where("start_time BETWEEN ? AND ?", req.StartTime, req.EndTime).
		Group("member_id")
	onlineStr := "ll_eshop+ll_livestream+pp_livestream+pp_post+remote_sales+cloud_store+link_to_ec"
	allStr := onlineStr + "+home_store+o2o+try_at_home"
	err := orm.Select("member_id,sum(" + onlineStr + ")/sum(" + allStr + ") as total").
		Find(&resList).Error
	for _, data := range resList {
		returnMap[data.MemberId] = data.Total
	}

	return returnMap, err
}

// GetAllVolumeByMemberIds 获取导购总销售额
func (d *MemberSalesVolume) GetAllVolumeByMemberIds(c app.Ctx, req memberReq.AchievementParams) (total int64, err error) {
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).
		Where("start_time BETWEEN ? AND ?", req.StartTime, req.EndTime).
		Where("member_id in ?", req.MemberIds)
	selectStr := "ll_eshop+ll_livestream+pp_livestream+pp_post+remote_sales+cloud_store+link_to_ec"
	if req.SearchType == microScrmModel.ALL_SALES_VOLUME {
		selectStr += "+home_store+o2o+try_at_home"
	}
	err = orm.Select("round(sum("+selectStr+")) as total").
		Pluck("total", &total).Error
	return
}

func (d *MemberSalesVolume) CreateInfo(c app.Ctx, enterpriseId int64, info microScrmModel.MemberSalesVolumeModel) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&info).Error
}

func (d *MemberSalesVolume) GetOne(c app.Ctx, enterpriseId int64, query interface{}) (microScrmModel.MemberSalesVolumeModel, error) {
	var result microScrmModel.MemberSalesVolumeModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).Find(&result).Error
	return result, err
}

func (d *MemberSalesVolume) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).UpdateColumns(data).Error
	return err
}
