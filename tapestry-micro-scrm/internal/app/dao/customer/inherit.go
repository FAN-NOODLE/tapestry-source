package customerDao

import (
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	customerReq "micro-scrm/internal/app/http/request/customer"
	customerModels "micro-scrm/internal/app/models/customer"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type Inherit struct {
	dao       dao.Dao
	tableName string
}

func NewInherit() *Inherit {
	return &Inherit{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: customerModels.CustomerInheritModel{}.TableName(),
	}
}

func (d *Inherit) GetExternalAllotList(c app.Ctx, request customerReq.ExternalListReq) ([]customerReq.ExternalListRes, int64, error) {
	var (
		total           int64
		ExternalListRes []customerReq.ExternalListRes
		err             error
	)
	offset := (request.Page - 1) * request.Size
	query := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id":   request.EnterpriseID,
			"status":          0,
			"handover_userid": request.HandoverUserId,
			"deleted_at":      0,
		})
	query = query.Where("handover_mobile <> ?", " ")
	if request.SearchName != "" {
		query = query.Where("external_nickname like ?", "%"+request.SearchName+"%")
	}
	err = query.Count(&total).Order("created_at DESC").Offset(offset).Limit(request.Size).Find(&ExternalListRes).Error

	return ExternalListRes, total, err
}

func (d *Inherit) GetExternalAllotCount(c app.Ctx, request customerReq.ExternalListReq) (count int64, err error) {
	query := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id":   request.EnterpriseID,
			"status":          0,
			"handover_userid": request.HandoverUserId,
			"deleted_at":      0,
		})
	query = query.Where("handover_mobile <> ?", " ")
	if request.SearchName != "" {
		query = query.Where("external_nickname like ?", "%"+request.SearchName+"%")
	}
	err = query.Count(&count).Error
	return
}

// GetInheritList 获取继承列表 - 优化版本，分步查询，单表操作
func (d *Inherit) GetInheritList(c app.Ctx, request customerReq.CustomerListReq) ([]customerModels.ResponseCustomerInheritList, int64, error) {
	var (
		total                       int64
		ResponseCustomerInheritList []customerModels.ResponseCustomerInheritList
		err                         error
	)

	// 1. 构建基础查询条件
	query := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id = ?", request.EnterpriseID).
		Where("handover_mobile <> ?", " ").
		Where("deleted_at = 0")

	if request.HandoverName != "" {
		query = query.Where("handover_name like ?", "%"+request.HandoverName+"%")
	}
	if request.HandoverUserId != "" {
		query = query.Where("handover_userid = ?", request.HandoverUserId)
	}
	if request.HandoverDepartmentId > 0 {
		query = query.Where("handover_department_id = ?", request.HandoverDepartmentId)
	}

	// 2. 根据状态过滤
	if request.Status == 0 { // 未分配
		query = query.Where("status = 0")
	} else { // 已分配
		query = query.Where("status = 1")
	}

	// 3. 查询 handover_userid 列表（用于后续查询）
	type HandoverUserInfo struct {
		HandoverUserid         string
		HandoverMobile         string
		HandoverName           string
		HandoverDepartmentId   int64
		HandoverDepartmentName string
		DimissionTime          int
	}
	var handoverUsers []HandoverUserInfo

	// 子查询获取去重的 handover_userid
	subQuery := d.dao.WithContext(c).Table(d.tableName).
		Select("handover_userid, handover_mobile, handover_name, handover_department_id, handover_department_name, dimission_time").
		Where("enterprise_id = ?", request.EnterpriseID).
		Where("handover_mobile <> ?", " ").
		Where("deleted_at = 0").
		Group("handover_userid")

	// 应用额外的过滤条件到子查询
	if request.HandoverName != "" {
		subQuery = subQuery.Where("handover_name like ?", "%"+request.HandoverName+"%")
	}
	if request.HandoverUserId != "" {
		subQuery = subQuery.Where("handover_userid = ?", request.HandoverUserId)
	}
	if request.HandoverDepartmentId > 0 {
		subQuery = subQuery.Where("handover_department_id = ?", request.HandoverDepartmentId)
	}

	// 统计总数
	err = d.dao.WithContext(c).Table("(?) as sub", subQuery).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询 handover 用户信息
	offset := (request.Page - 1) * request.Size
	err = subQuery.Order("id DESC").Offset(offset).Limit(request.Size).Find(&handoverUsers).Error
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回
	if len(handoverUsers) == 0 {
		return ResponseCustomerInheritList, total, nil
	}

	// 4. 提取所有 handover_userid
	var handoverUserIds []string
	for _, user := range handoverUsers {
		handoverUserIds = append(handoverUserIds, user.HandoverUserid)
	}

	// 5. 查询每个 handover_userid 的统计信息（已分配和未分配数量）
	type StatusCount struct {
		HandoverUserid string
		Status         int
		Count          int64
	}
	var statusCounts []StatusCount

	err = d.dao.WithContext(c).Table(d.tableName).
		Select("handover_userid, status, COUNT(*) as count").
		Where("enterprise_id = ?", request.EnterpriseID).
		Where("handover_userid IN ?", handoverUserIds).
		Where("deleted_at = 0").
		Group("handover_userid, status").
		Find(&statusCounts).Error
	if err != nil {
		return nil, 0, err
	}

	// 6. 构建统计信息的 map
	statusMap := make(map[string]map[int]int64)
	for _, sc := range statusCounts {
		if _, ok := statusMap[sc.HandoverUserid]; !ok {
			statusMap[sc.HandoverUserid] = make(map[int]int64)
		}
		statusMap[sc.HandoverUserid][sc.Status] = sc.Count
	}

	// 7. 组装结果
	for _, user := range handoverUsers {
		result := customerModels.ResponseCustomerInheritList{
			HandoverUserid:         user.HandoverUserid,
			HandoverMobile:         user.HandoverMobile,
			HandoverName:           user.HandoverName,
			HandoverDepartmentId:   user.HandoverDepartmentId,
			HandoverDepartmentName: user.HandoverDepartmentName,
			DimissionTime:          user.DimissionTime,
		}
		if counts, ok := statusMap[user.HandoverUserid]; ok {
			result.UnassignedUserNum = int(counts[0])
			result.AssignedUserNum = int(counts[1])
		}
		ResponseCustomerInheritList = append(ResponseCustomerInheritList, result)
	}

	return ResponseCustomerInheritList, total, nil
}

func (d *Inherit) GetTakeOverList(c app.Ctx, request customerReq.GroupChatList) (data []customerReq.TakeOverInfo, err error) {
	offset := (request.Page - 1) * request.Size
	err = d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id":   request.EnterpriseID,
			"handover_userid": request.HandoverUserId,
			"status":          1,
			"deleted_at":      0,
		}).Where("takeover_userid <> ?", " ").Offset(offset).Limit(request.Size).Group("takeover_userid").Distinct("takeover_userid").Select("takeover_userid,takeover_name,takeover_department_name").Find(&data).Error
	return
}

func (d *Inherit) ChangeTakeoverUserByHandoverIdsExternalUserIds(c app.Ctx, enterpriseID int64, handoverUserId string, externalUserIds []string, data customerModels.CustomerInheritModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ? and handover_userid = ? and external_userid In(?)", enterpriseID, handoverUserId, externalUserIds).
		UpdateColumns(data).Error
	return err
}

func (d *Inherit) GetUniqueStr(c app.Ctx, enterpriseID int64, newList [][]string) (data []string, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id", enterpriseID).
		Where("deleted_at", 0).
		Where("(handover_userid, external_userid, dimission_time) in ?", newList).
		Pluck("CONCAT(`handover_userid`,`external_userid`,`dimission_time`) unique_str ", &data).Error
	return
}

func (d *Inherit) GetListByUserIds(c app.Ctx, enterpriseID int64, userIds []string) (list []customerModels.CustomerInheritModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id", enterpriseID).
		Where("deleted_at", 0).
		Where("handover_userid in ?", userIds).
		Find(&list).Error
	return
}

// InsertNewData 入库
func (d *Inherit) InsertNewData(c app.Ctx, data []customerModels.CustomerInheritModel) error {
	return d.dao.WithContext(c).Table(d.tableName).CreateInBatches(&data, len(data)).Error

}

// UpdateByIds 根据id更新时间
func (d *Inherit) UpdateByIds(c app.Ctx, enterpriseId int64, ids []int64) error {
	return d.dao.WithContext(c).Table(d.tableName).Where("deleted_at", 0).
		Where("enterprise_id", enterpriseId).
		Where("id in ?", ids).
		UpdateColumn("updated_at", time.Now().Unix()).Error

}

// UpdateDataByUnique 分批更新数据
func (d *Inherit) UpdateDataByUnique(c app.Ctx, uniqueData [][]string) error {
	const batchSize = 100 // 每批处理的数量

	// 获取数据库连接
	db := d.dao.WithContext(c).Table(d.tableName)

	// 如果数据量小于等于批次大小，直接处理
	if len(uniqueData) <= batchSize {
		return db.Where("(handover_userid, external_userid, dimission_time) in ?", uniqueData).
			Where("deleted_at", 0).
			UpdateColumn("updated_at", time.Now().Unix()).Error
	}

	// 分批处理
	for i := 0; i < len(uniqueData); i += batchSize {
		end := i + batchSize
		if end > len(uniqueData) {
			end = len(uniqueData)
		}

		batch := uniqueData[i:end]

		err := db.Where("(handover_userid, external_userid, dimission_time) in ?", batch).
			Where("deleted_at", 0).
			UpdateColumn("updated_at", time.Now().Unix()).Error

		if err != nil {
			return fmt.Errorf("处理第 %d-%d 条数据时出错: %w", i+1, end, err)
		}
	}

	return nil
}

// DelInvalidData 移除失效数据
func (d *Inherit) DelInvalidData(c app.Ctx, enterpriseId int64) error {
	nowUnix := time.Now().Unix()
	return d.dao.WithContext(c).Table(d.tableName).
		Where("status", 0).
		Where("enterprise_id = ?", enterpriseId).
		Where("deleted_at", 0).
		Where("updated_at < ?", nowUnix-86400).
		UpdateColumn("deleted_at", nowUnix).Error
}

// UpdateHandoverData 数据补全
func (d *Inherit) UpdateHandoverData(c app.Ctx, enterpriseID int64) error {
	return d.dao.WithContext(c).
		Exec("UPDATE customer_inherit inherit" +
			" LEFT JOIN members ON inherit.handover_userid = members.userid" +
			" AND members.enterprise_id =" + fmt.Sprintf("%d", enterpriseID) +
			" LEFT JOIN member_departments md ON md.member_id = members.id" +
			" LEFT JOIN departments ON departments.id = md.department_id" +
			" SET inherit.handover_name = if(members.id > 0, members.`name`, inherit.handover_userid)" +
			",inherit.handover_mobile = if( members.id > 0, members.`mobile`, '')" +
			",inherit.handover_department_id = if( departments.id > 0, departments.`id`, 0)" +
			",inherit.handover_department_name = if( departments.id > 0, departments.`name`, '')" +
			" WHERE inherit.enterprise_id=" + fmt.Sprintf("%d", enterpriseID) +
			" AND inherit.handover_name=''").Error
}

// GetCustomerLimit  数据补全
func (d *Inherit) GetCustomerLimit(c app.Ctx, enterpriseID, LastID int64, limit int) ([]customerModels.CustomerInheritModel, error) {
	var list = make([]customerModels.CustomerInheritModel, 0)
	err := d.dao.WithContext(c).Where("enterprise_id = ?", enterpriseID).
		Where("external_nickname = ''").
		Where("id > ?", LastID).
		Order("id asc").
		Limit(limit).
		Find(&list).Error
	return list, err
}

func (d *Inherit) UpdateData(c app.Ctx, enterpriseID int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id = ?", enterpriseID).Where(query).UpdateColumns(data).Error
	return err
}

// UpdateExternalData 数据补全
func (d *Inherit) UpdateExternalData(c app.Ctx, enterpriseID int64) error {
	return d.dao.WithContext(c).
		Exec("UPDATE customer_inherit inherit LEFT JOIN micro_scrm_enterprise_" +
			fmt.Sprintf("%d", enterpriseID) + ".wx_user_corp_external external ON inherit.external_userid = external.openid" +
			" SET inherit.external_avatar = if(external.id > 0, external.avatar, 'https://oss-crs.vchangyi.com/standard/lottery_default_avatar.png')" +
			",inherit.external_nickname = if(external.id > 0, external.name, inherit.external_userid)" +
			" WHERE inherit.enterprise_id=" + fmt.Sprintf("%d", enterpriseID) +
			" AND inherit.external_nickname=''").Error
}

func (d *Inherit) GetTakeOverExternalList(c app.Ctx, request customerReq.ExternalListReq, ids []string) (data []customerReq.TakeOverExternalInfo, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"handover_userid": request.HandoverUserId,
			"takeover_userid": request.TakeoverUserId,
			"external_userid": ids,
		}).Select("external_nickname,external_avatar,external_userid").Find(&data).Error
	return
}

func (d *Inherit) GetTakeOverCount(c app.Ctx, enterpriseID int64, UserID string) (total int64, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id":   enterpriseID,
			"handover_userid": UserID,
			"status":          1,
			"deleted_at":      0,
		}).Where("takeover_userid <> ?", " ").Group("takeover_userid").Distinct("takeover_userid").Count(&total).Error
	return
}

func (d *Inherit) GetTakeOverExternalListByLocal(c app.Ctx, request customerReq.ExternalListReq) (total int64, data []customerReq.ExternalListRes, err error) {
	offset := (request.Page - 1) * request.Size
	err = d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id":   request.EnterpriseID,
			"handover_userid": request.HandoverUserId,
			"status":          1,
			"takeover_userid": request.TakeoverUserId,
		}).Count(&total).Offset(offset).Limit(request.Size).Find(&data).Error
	return
}
