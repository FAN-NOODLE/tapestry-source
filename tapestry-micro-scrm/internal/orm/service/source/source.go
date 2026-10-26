package source

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm/source"
	"time"
)

const (
	SOURCE_ALL_GROUP   = 0  // 全部分组
	SOURCE_UN_GROUP    = -1 // 未分组
	SOURCE_OTHER_GROUP = -2 // 回收站
)

func (s *SourceService) AddSource(info source.GroupModel) error {
	tx := s.Begin()
	if err := tx.Model(&info).Create(&info).Error; err != nil {
		tx.Rollback()
		return err
	}
	newSourceInfo := make([]int64, 0)
	if err := tx.Model(&source.SourceModel{}).Select("id").
		Where("group_id = ?  AND enterprise_id = ?", info.ParentId, info.EnterpriseID).Pluck("id", &newSourceInfo).Error; err != nil {
		tx.Rollback()
		return err
	}

	//if len(newSourceInfo) == 0 {
	//	tx.Commit()
	//	return nil
	//} else {
	//	if err := tx.Model(&source.SourceModel{}).Where("id in (?)", newSourceInfo).Update("group_id", info.ID).Error; err != nil {
	//		tx.Rollback()
	//		return err
	//	}
	//}
	tx.Commit()
	return nil
}

func (s *SourceService) EditSource(sourceID int64, enterpriseID int64, info map[string]interface{}, mold source.GroupModel) error {
	return s.Model(&mold).Where("id = ? AND enterprise_id = ?", sourceID, enterpriseID).Updates(info).Error
}

func (s *SourceService) RemoveSource(id int64, info source.GroupModel) error {
	return s.Model(&info).Where("id = ?", id).Update("deleted_at", time.Now().Unix()).Error
}

func (s *SourceService) RemoveSourceMaterial(id int64) error {
	return s.Model(&source.SourceModel{}).Where("group_id = ?", id).Update("group_id", SOURCE_OTHER_GROUP).Error
}

func (s *SourceService) GetOrder(parentID int64, info source.GroupModel) (result int64, err error) {
	err = s.Model(&info).Select("id").Where("parent_id = ? and enterprise_id = ?", parentID, info.EnterpriseID).Count(&result).Error
	return
}

func (s *SourceService) SourceInfo(info source.GroupModel) (result []source.GroupModel, err error) {
	err = s.Model(&info).Where("deleted_at = ? and enterprise_id = ?", 0, info.EnterpriseID).Order("orders ASC").Find(&result).Error
	return
}

func (s *SourceService) GetGroupID(info source.GroupModel) (result []source.GroupModel, err error) {
	err = s.Model(&info).Where("deleted_at = ? AND enterprise_id = ?", 0, info.EnterpriseID).Find(&result).Error
	return
}

func (s *SourceService) RemoveVerify(ID []int64, info source.SourceModel) (result []int64, err error) {
	err = s.Model(&info).Where("group_id in (?)", ID).Find(&result).Error
	return
}

func (s *SourceService) GetSourceName(enterpriseID, parentID int64) (res []string, err error) {
	err = s.Model(&source.GroupModel{}).
		Where("enterprise_id = ?", enterpriseID).
		Where("parent_id = ? AND deleted_at = 0", parentID).
		Pluck("group_name", &res).Error
	return
}

func (s *SourceService) GetListByIds(IDs []int64, enterpriseId int64) (list []source.GroupModel, err error) {
	err = s.Model(&source.GroupModel{}).
		Where("id in (?) AND enterprise_id = ?", IDs, enterpriseId).
		Find(&list).Error
	return
}

func (s *SourceService) UpdateOrder(id int64, orderStatus int) (err error) {
	err = s.Model(&source.GroupModel{}).Where("id = ?", id).Update("orders", orderStatus).Error
	return
}
func (d *SourceService) IncrementOrder(enterpriseId, parentId int64, order int) error {
	err := d.Model(&source.GroupModel{}).
		Where("enterprise_id = ? AND parent_id = ? AND orders >= ?", enterpriseId, parentId, order).
		UpdateColumn("orders", gorm.Expr(fmt.Sprintf("%s + ?", "orders"), 1)).Error
	return err
}

func (s *SourceService) GetMaxOrder() (result int64, err error) {
	err = s.Model(&source.GroupModel{}).Select("IFNULL(Max(orders),0) AS orders").Find(&result).Error
	return
}

func (s *SourceService) GetSourceLevel(parentID int64) (result int64) {
	s.Model(&source.GroupModel{}).Select("parent_id").Where("id = ?", parentID).First(&result)
	return
}

func (s *SourceService) GetGroupName(groupId int64) (result micro_scrm.ReceiveGroupName) {
	s.Model(&source.GroupModel{}).Select("group_name").Where("id = ?", groupId).Find(&result)
	return
}

func (s *SourceService) SourceInfoById(info source.GroupModel) (result []source.GroupModel, err error) {
	err = s.Model(&info).Where("deleted_at = ? and enterprise_id = ?  and id = ?", 0, info.EnterpriseID, info.ID).Find(&result).Error
	return
}

func (s *SourceService) SourceInfoByParentId(info source.GroupModel) (result []source.GroupModel, err error) {
	err = s.Model(&info).Where("deleted_at = ? and enterprise_id = ?  and parent_id = ?", 0, info.EnterpriseID, info.ID).Find(&result).Error
	return
}

func (s *SourceService) SourceInfoByThisId(info source.GroupModel) (result []source.GroupModel, err error) {
	err = s.Model(&info).Where("deleted_at = ? and enterprise_id = ?  and (id = ? or parent_id = ?)", 0, info.EnterpriseID, info.ID, info.ID).Find(&result).Error
	return
}
