package mysql

import (
	"database/sql"
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/fx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	workexternal "micro-scrm/internal/modules/work-external"
	"micro-scrm/internal/pkg/db"
)

type WorkRepo struct {
	*repos.SQLRepo
}

func (work *WorkRepo) SelectReUv(c ctx.Context, param workexternal.SelectKeep) (int64, error) {
	stmt, err := work.DB.Template("t_work_external_user_relation/select-re-uv", work.Cond2Args(param))
	if err != nil {
		return 0, err
	}

	addRes := &workexternal.SelectUvRes{}
	err = stmt.NamedGet(c, addRes, param)
	if err != nil {
		return 0, err
	}
	return addRes.RelationNum, nil
}

func (work *WorkRepo) SelectLossUv(c ctx.Context, param workexternal.SelectKeep) (int64, error) {
	stmt, err := work.DB.Template("t_work_external_user_relation/select-loss-uv", work.Cond2Args(param))
	if err != nil {
		return 0, err
	}

	addRes := &workexternal.SelectUvRes{}
	err = stmt.NamedGet(c, addRes, param)
	if err != nil {
		return 0, err
	}
	return addRes.RelationNum, nil
}

func (work *WorkRepo) SelectAddUv(c ctx.Context, param workexternal.SelectKeep) (int64, error) {
	stmt, err := work.DB.Template("t_work_external_user_relation/select-add-uv", work.Cond2Args(param))
	if err != nil {
		return 0, err
	}

	addRes := &workexternal.SelectUvRes{}
	err = stmt.NamedGet(c, addRes, param)
	if err != nil {
		return 0, err
	}
	return addRes.RelationNum, nil
}

// ExistsJoinExternal 查询是否存在正在加入的外部联系人状态
func (work *WorkRepo) ExistsJoinExternal(c ctx.Context, userId string, enterpriseId int64) (*workexternal.ExternalTrace, error) {
	records := make(map[string]interface{})
	records["EnterpriseID"] = enterpriseId
	stmt, err := work.DB.Template("t_work_external_trace/select_external_by_userid", work.Cond2Args(records))
	if err != nil {
		return nil, err
	}

	data := new(workexternal.ExternalTrace)
	err = stmt.Get(c, data, userId, enterpriseId)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return data, nil
	}
}

// LastExternalTrace 查询上一条记录
func (work *WorkRepo) LastExternalTrace(c ctx.Context, userId string, enterpriseId int64) (*workexternal.ExternalTrace, error) {
	records := make(map[string]interface{})
	records["EnterpriseID"] = enterpriseId
	stmt, err := work.DB.Template("t_work_external_trace/select_last_leave_by_userid", work.Cond2Args(records))
	if err != nil {
		return nil, err
	}

	trace := new(workexternal.ExternalTrace)
	err = stmt.Get(c, trace, userId, enterpriseId)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return trace, nil
	}
}

// JoinExternal 插入外部联系人状态
func (work *WorkRepo) JoinExternal(c ctx.Context, req, lastTrace *workexternal.ExternalTrace) error {
	// 使用事务
	err := work.DB.Tx(c, func(tx *sqlplus.Tx) error {
		// 防止 crash
		e, crash := fx.RunSafe(func() error {
			if lastTrace != nil {
				cond := work.Cond2Args(lastTrace)
				// 更新上一条记录的下次加入时间
				stmt, err := work.DB.Template("t_work_external_trace/update_next_join_time", cond)
				if err != nil {
					return err
				}
				res, err := stmt.NamedExec(c, map[string]interface{}{
					"update_id":      lastTrace.ID,
					"next_join_time": req.JoinTime,
				})
				if err != nil || res == nil {
					return err
				}
			}

			// 插入记录
			stmt, err := work.DB.Template("t_work_external_trace/join_external_trace", work.Cond2Args(req))
			if err != nil {
				return err
			}
			res, err := stmt.NamedExec(c, req)
			if err != nil || res == nil {
				return err
			}
			return nil
		})

		if crash != nil {
			return crash
		}

		return e
	})

	return err
}

// LeaveExternalByUserId 外部联系人离开系统
func (work *WorkRepo) LeaveExternalByUserId(c ctx.Context, trace *workexternal.ExternalTrace) error {
	cond := work.Cond2Args(trace)
	// 更新记录的离开时间，要求取第一次的离开时间
	stmt, err := work.DB.Template("t_work_external_trace/update_leave_external", cond)
	if err != nil {
		return err
	}
	res, err := stmt.NamedExec(c, trace)
	if err != nil || res == nil {
		return err
	}

	return nil
}

// ExistsJoinExternalInChannel 查询正常外部联系人是否存在该渠道
func (work *WorkRepo) ExistsJoinExternalInChannel(c ctx.Context, userId string, enterpriseId int64, channelId int64) (*workexternal.ExternalChannelTrace, error) {
	records := make(map[string]interface{})
	records["EnterpriseID"] = enterpriseId
	stmt, err := work.DB.Template("t_work_external_channel_trace/select_external_by_userid", work.Cond2Args(records))
	if err != nil {
		return nil, err
	}

	data := new(workexternal.ExternalChannelTrace)
	err = stmt.Get(c, data, userId, enterpriseId, channelId)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return data, nil
	}
}

// LastExternalTraceInChannel 查询上一条记录
func (work *WorkRepo) LastExternalTraceInChannel(c ctx.Context, userId string, enterpriseId, channelId int64) (*workexternal.ExternalChannelTrace, error) {
	records := make(map[string]interface{})
	records["EnterpriseID"] = enterpriseId
	stmt, err := work.DB.Template("t_work_external_channel_trace/select_last_leave_by_userid", work.Cond2Args(records))
	if err != nil {
		return nil, err
	}

	trace := new(workexternal.ExternalChannelTrace)
	err = stmt.Get(c, trace, userId, enterpriseId, channelId)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return trace, nil
	}
}

func (work *WorkRepo) JoinExternalInChannel(c ctx.Context, req, lastTrace *workexternal.ExternalChannelTrace) error {
	// 使用事务
	err := work.DB.Tx(c, func(tx *sqlplus.Tx) error {
		// 防止 crash
		e, crash := fx.RunSafe(func() error {
			if lastTrace != nil {
				cond := work.Cond2Args(lastTrace)
				// 更新上一条记录的上次加入时间
				stmt, err := work.DB.Template("t_work_external_channel_trace/update_next_join_time", cond)
				if err != nil {
					return err
				}
				res, err := stmt.NamedExec(c, map[string]interface{}{
					"update_id":      lastTrace.ID,
					"next_join_time": req.JoinTime,
				})
				if err != nil || res == nil {
					return err
				}
			}

			// 插入记录
			stmt, err := work.DB.Template("t_work_external_channel_trace/join_external_trace", work.Cond2Args(req))
			if err != nil {
				return err
			}
			res, err := stmt.NamedExec(c, req)
			if err != nil || res == nil {
				return err
			}
			return nil
		})

		if crash != nil {
			return crash
		}

		return e
	})

	return err
}

func (work *WorkRepo) LeaveChannelExternalByUserId(c ctx.Context, trace *workexternal.ExternalChannelTrace) error {
	cond := work.Cond2Args(trace)
	// 更新记录的离开时间，要求取第一次的离开时间
	stmt, err := work.DB.Template("t_work_external_channel_trace/update_leave_external", cond)
	if err != nil {
		return err
	}
	res, err := stmt.NamedExec(c, trace)
	if err != nil || res == nil {
		return err
	}

	return nil
}

func (work *WorkRepo) AddDelExternalContact(c ctx.Context, trace *workexternal.DelRelationChannelTrace) error {
	// 插入记录
	stmt, err := work.DB.Template("t_work_del_relation_channel_trace/insert_del_external_trace", work.Cond2Args(trace))
	if err != nil {
		return err
	}
	res, err := stmt.NamedExec(c, trace)
	if err != nil {
		return err
	}

	if res == nil {
		return errors.New("更新未返回数据")
	}

	return nil
}

func (work *WorkRepo) FindDelRelation(c ctx.Context, condition *workexternal.DelRelationChannelTrace) (*workexternal.DelRelationChannelTrace, error) {
	stmt, err := work.DB.Template("t_work_del_relation_channel_trace/select_del_relation_by_userid", work.Cond2Args(condition))
	if err != nil {
		return nil, err
	}

	trace := new(workexternal.DelRelationChannelTrace)
	err = stmt.Get(c, trace, condition.UserID, condition.ExternalUserID, condition.EnterpriseID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return trace, nil
	}
}

func (work *WorkRepo) RecoveryDelRelation(c ctx.Context, id int64, update *workexternal.DelRelationChannelTrace) error {
	// 插入记录
	stmt, err := work.DB.Template("t_work_del_relation_channel_trace/recovery_del_relation", work.Cond2Args(update))
	if err != nil {
		return err
	}
	update.ID = id
	res, err := stmt.NamedExec(c, update)
	if err != nil {
		return err
	}

	if res == nil {
		return errors.New("更新未返回数据")
	}

	return nil
}

// AddRelationRecord 添加记录
func (work *WorkRepo) AddRelationRecord(c ctx.Context, record *workexternal.WorkRelationAddRecord) error {
	stmt, err := work.DB.Template("t_work_relation_add_record/insert")
	if err != nil {
		return err
	}
	res, err := stmt.NamedExec(c, record)
	if err != nil {
		return err
	}

	if res == nil {
		return errors.New("更新未返回数据")
	}

	return nil
}

// FindLastAddRelationRecord 查询上次添加的记录
func (work *WorkRepo) FindLastAddRelationRecord(c ctx.Context, record *workexternal.WorkRelationAddRecord) (*workexternal.WorkRelationAddRecord, error) {
	stmt, err := work.DB.Template("t_work_relation_add_record/select_last_relation", work.Cond2Args(record))
	if err != nil {
		return nil, err
	}
	trace := new(workexternal.WorkRelationAddRecord)
	err = stmt.NamedGet(c, trace, record)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return trace, nil
	}
}

func (work *WorkRepo) FinishAddRelation(c ctx.Context, record *workexternal.WorkRelationAddRecord) error {
	stmt, err := work.DB.Template("t_work_relation_add_record/finish_add_relation", work.Cond2Args(record))
	if err != nil {
		return err
	}

	res, err := stmt.NamedExec(c, record)
	if err != nil {
		return err
	}

	if res == nil {
		return errors.New("更新未返回数据")
	}

	return nil
}

// HistoryUcRelation 查询 UC 历史数据
func (work *WorkRepo) HistoryUcRelation(c ctx.Context, page repos.PageVO, status int8) ([]*workexternal.HistoryUcRelationLog, error) {
	dest := make([]*workexternal.HistoryUcRelationLog, 0, page.Limit)
	err := db.Query(c, work.SQLRepo, "t_history_uc_relation_log/select-history-relation", &dest, status, page.Offset, page.Limit)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (work *WorkRepo) UpdateHistoryStatus(c ctx.Context, id int64, status int8) error {
	stmt, err := work.DB.Template("t_history_uc_relation_log/update-status")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c, status, id)
	if err != nil {
		return err
	}
	return nil
}
