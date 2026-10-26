package click_house

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	workexternal "micro-scrm/internal/modules/work-external"
	"micro-scrm/internal/pkg/db/gormx"
)

// MemberRepo 定义成员存储
type WorkRepo struct {
	*repos.SQLRepo
}

func (w WorkRepo) HistoryUcRelation(c ctx.Context, page repos.PageVO, status int8) ([]*workexternal.HistoryUcRelationLog, error) {
	panic("implement me")
}

func (w WorkRepo) UpdateHistoryStatus(c ctx.Context, id int64, status int8) error {
	panic("implement me")
}

func (w WorkRepo) SelectReUv(c ctx.Context, param workexternal.SelectKeep) (int64, error) {
	// stmt, err := w.DB.Template("t_external_relation_day_local/select-re-uv", w.Cond2Args(param))
	// if err != nil {
	// 	return nil, err
	// }

	// addRes := &workexternal.SelectUvRes{}
	// // err = stmt.NamedGet(c, addRes, param)
	// // if err != nil && err != sql.ErrNoRows {
	// // 	return nil, err
	// // }
	// return addRes, nil

	var count int64
	gormx.GetCkInstance().WithContext(c).
		Table("micro_scrm_ck.work_external_user_relation").
		Distinct("work_external_user_id").
		Where("enterprise_id = ? and created_at <= ?", param.EnterpriseId, param.EndAt).
		Where("deleted_at = 0 AND relation_status = 0").
		Count(&count)
	return count, nil
}

func (w WorkRepo) SelectAddUv(c ctx.Context, param workexternal.SelectKeep) (int64, error) {
	// stmt, err := w.DB.Template("t_external_relation_day_local/select-add-uv", w.Cond2Args(param))
	// if err != nil {
	// 	return nil, err
	// }

	// addRes := &workexternal.SelectUvRes{}
	// err = stmt.NamedGet(c, addRes, param)
	// if err != nil && err != sql.ErrNoRows {
	// 	return nil, err
	// }
	// return addRes, nil

	var count int64
	gormx.GetCkInstance().WithContext(c).
		Table("micro_scrm_ck.work_external_user_relation").
		Distinct("work_external_user_id").
		Where("enterprise_id = ? and createtime BETWEEN ? AND ?", param.EnterpriseId, param.StartAt, param.EndAt).
		Where("deleted_at = 0 AND relation_status = 0").
		Count(&count)
	return count, nil
}

func (w WorkRepo) SelectLossUv(c ctx.Context, param workexternal.SelectKeep) (int64, error) {
	// stmt, err := w.DB.Template("t_external_relation_day_local/select-loss-uv", w.Cond2Args(param))
	// if err != nil {
	// 	return nil, err
	// }

	// addRes := &workexternal.SelectUvRes{}
	// err = stmt.NamedGet(c, addRes, param)
	// if err != nil && err != sql.ErrNoRows {
	// 	return nil, err
	// }

	// return addRes, nil
	var count int64
	gormx.GetCkInstance().WithContext(c).
		Table("micro_scrm_ck.work_external_user_relation").
		Distinct("work_external_user_id").
		Where("enterprise_id = ? and (member_deleted_time BETWEEN ? AND ? OR external_deleted_time BETWEEN ? AND ?)",
			param.EnterpriseId, param.StartAt, param.EndAt, param.StartAt, param.EndAt).
		Where("relation_status in (1,2,3)").
		Count(&count)
	return count, nil
}

func (w WorkRepo) ExistsJoinExternal(c ctx.Context, userId string, enterpriseId int64) (*workexternal.ExternalTrace, error) {
	panic("implement me")
}

func (w WorkRepo) LastExternalTrace(c ctx.Context, userId string, enterpriseId int64) (*workexternal.ExternalTrace, error) {
	panic("implement me")
}

func (w WorkRepo) JoinExternal(c ctx.Context, req *workexternal.ExternalTrace, lastTrace *workexternal.ExternalTrace) error {
	panic("implement me")
}

func (w WorkRepo) LeaveExternalByUserId(c ctx.Context, trace *workexternal.ExternalTrace) error {
	panic("implement me")
}

func (w WorkRepo) ExistsJoinExternalInChannel(c ctx.Context, userId string, enterpriseId, channelId int64) (*workexternal.ExternalChannelTrace, error) {
	panic("implement me")
}

func (w WorkRepo) LastExternalTraceInChannel(c ctx.Context, userId string, enterpriseId, channelId int64) (*workexternal.ExternalChannelTrace, error) {
	panic("implement me")
}

func (w WorkRepo) JoinExternalInChannel(c ctx.Context, req, existsTrace *workexternal.ExternalChannelTrace) error {
	panic("implement me")
}

func (w WorkRepo) LeaveChannelExternalByUserId(c ctx.Context, trace *workexternal.ExternalChannelTrace) error {
	panic("implement me")
}

func (w WorkRepo) AddDelExternalContact(c ctx.Context, trace *workexternal.DelRelationChannelTrace) error {
	panic("implement me")
}

func (w WorkRepo) FindDelRelation(c ctx.Context, condition *workexternal.DelRelationChannelTrace) (*workexternal.DelRelationChannelTrace, error) {
	panic("implement me")
}

func (w WorkRepo) RecoveryDelRelation(c ctx.Context, id int64, update *workexternal.DelRelationChannelTrace) error {
	panic("implement me")
}

func (w WorkRepo) AddRelationRecord(c ctx.Context, record *workexternal.WorkRelationAddRecord) error {
	panic("implement me")
}

func (w WorkRepo) FindLastAddRelationRecord(c ctx.Context, record *workexternal.WorkRelationAddRecord) (*workexternal.WorkRelationAddRecord, error) {
	panic("implement me")
}

func (w WorkRepo) FinishAddRelation(c ctx.Context, record *workexternal.WorkRelationAddRecord) error {
	panic("implement me")
}
