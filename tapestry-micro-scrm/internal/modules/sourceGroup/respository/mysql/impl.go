package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/sourceGroup"
)

type SourceGroupRepo struct {
	*repos.SQLRepo
}

// RepeatName 分组名称是否重复 true 是
func (r *SourceGroupRepo) IsRepeatName(cc ctx.Context, params *sourceGroup.SourceGroup) (bool, error) {
	cond := r.Cond2Args(params)
	countStmt, err := r.DB.Template("t_source_group/repeat-name", cond)
	if err != nil {
		return false, err
	}

	total := new(int64)
	countRes, err := countStmt.NamedQuery(cc, params)
	if err != nil {
		return false, err
	}
	defer countRes.Close()
	if countRes.Next() {
		err = countRes.Scan(total)
	}
	if err != nil {
		return false, err
	}
	if *total > 0 {
		return true, nil
	}

	return false, nil
}

// BatchDel 批量删除
func (r *SourceGroupRepo) BatchDel(cc ctx.Context, ids []int64) (int64, error) {
	cond := &map[string]interface{}{
		"IDs": ids,
	}

	stmt, err := r.DB.Template("t_source_group/batch-delete", cond)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(cc)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
