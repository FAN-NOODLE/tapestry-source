package clickhouse

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gorm.io/gorm"
	groupMember "micro-scrm/internal/modules/group-member"
	"micro-scrm/internal/pkg/db/gormx"
)

type GroupRepo struct {
	*repos.SQLRepo
}

// SelectGroupNum 查询群总数
func (group *GroupRepo) SelectGroupNum(c ctx.Context, param groupMember.GroupParam) (int64, error) {
	var count int64
	err := gormx.GetCkInstance().WithContext(c).
		Table("micro_scrm_ck.work_group_chat").
		Where("created_at <= ?", param.EndAt).
		Where("status = ?", 1).
		Where("enterprise_id = ?", param.EnterpriseId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SelectReGroupUv 查询留存
func (group *GroupRepo) SelectReGroupUv(c ctx.Context, param groupMember.GroupParam) (int64, error) {
	var count int64

	err := gormx.GetCkInstance().WithContext(c).
		Table("micro_scrm_ck.work_group_chat_member").
		Distinct("userid").
		Where("join_time <= ? and status = 1", param.EndAt).
		Where("enterprise_id = ?", param.EnterpriseId).
		Where("deleted_at = 0").
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

// SelectAddGroupUv 查询新增群成员数
func (group *GroupRepo) SelectAddGroupUv(c ctx.Context, param groupMember.GroupParam) (int64, error) {
	var count int64

	err := gormx.GetCkInstance().WithContext(c).
		Table("micro_scrm_ck.work_group_chat_member").
		Distinct("userid").
		Where("enterprise_id = ? and deleted_at = 0", param.EnterpriseId).
		Where("join_time >= ? and join_time <= ? and status = 1", param.StartAt, param.EndAt).
		Where("userid NOT IN (?)", func() *gorm.DB {
			db := gormx.GetCkInstance().WithContext(c).
				Table("micro_scrm_ck.work_group_chat_member").
				Select("userid").
				Distinct("userid").
				Where("join_time < ? and status = 1", param.StartAt).
				Where("enterprise_id = ?", param.EnterpriseId).
				Where("deleted_at = 0")
			return db
		}()).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

// SelectLossGroupUv 查询流失群成员数
func (group *GroupRepo) SelectLossGroupUv(c ctx.Context, param groupMember.GroupParam) (int64, error) {
	var count int64

	err := gormx.GetCkInstance().WithContext(c).
		Table("micro_scrm_ck.work_group_chat_member").
		Distinct("userid").
		Where("enterprise_id = ?", param.EnterpriseId).
		Where("leave_time >= ? and leave_time <= ?", param.StartAt, param.EndAt).
		Where("status = 2 and deleted_at = 0").
		Where("userid NOT IN (?)", func() *gorm.DB {
			db := gormx.GetCkInstance().WithContext(c).
				Table("micro_scrm_ck.work_group_chat_member").
				Select("userid").
				Distinct("userid").
				Where("join_time >= ? and join_time <= ? and `status` = 1", param.StartAt, param.EndAt).
				Where("enterprise_id = ?", param.EnterpriseId).
				Where("deleted_at = 0")
			return db
		}()).Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}
