package db

import (
	"fmt"
	"strings"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/pkg/errorIgnore"
)

// Get 查询封装
func Get(c ctx.Context, o *repos.SQLRepo, tpl string, dest interface{}, args ...interface{}) error {
	stmt, err := o.DB.Template(tpl)
	if err != nil {
		return err
	}
	errStmt := stmt.Get(c, dest, args...)
	if errStmt != nil {
		return errStmt
	}
	return nil
}

// Add 插入封装
func Add(c ctx.Context, o *repos.SQLRepo, tpl string, entity interface{}) (int64, error) {
	return Update(c, o, tpl, entity)
}

// Update 更新封装
func Update(c ctx.Context, o *repos.SQLRepo, tpl string, entity interface{}) (int64, error) {
	stmt, err := o.DB.Template(tpl)
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, entity)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Query 查询封装
func Query(c ctx.Context, o *repos.SQLRepo, tpl string, dest interface{}, args ...interface{}) error {
	stmt, stmtErr := o.DB.Template(tpl)
	if stmtErr != nil {
		return stmtErr
	}
	res, getErr := stmt.Query(c, args...)
	if getErr != nil {
		return getErr
	}
	scanErr := sqlplus.ScanAll(res, dest)
	if scanErr != nil {
		return scanErr
	}
	return nil
}

// NamedSelect NamedSelect封装
func NamedSelect(c ctx.Context, o *repos.SQLRepo, tpl string, dest interface{}, args interface{}, tempArgs ...interface{}) error {
	stmt, err := o.DB.Template(tpl, tempArgs...)
	if err != nil {
		return err
	}
	err = stmt.NamedSelect(c, dest, args)
	if err != nil {
		return err
	}
	return nil
}

// NamedGet NamedGet封装
func NamedGet(c ctx.Context, o *repos.SQLRepo, tpl string, dest interface{}, args interface{}, tempArgs ...interface{}) error {
	stmt, err := o.DB.Template(tpl, tempArgs...)
	if err != nil {
		return err
	}
	err = stmt.NamedGet(c, dest, args)
	if err != nil {
		return err
	}
	return nil
}

// FindPage FindPage封装
func FindPage(c ctx.Context, o *repos.SQLRepo, tpl string, page int64, size int64, condition interface{}, dest interface{}) (int64, interface{}, error) {
	offset := (page - 1) * size
	pageVo := repos.PageVO{
		Offset: int(offset),
		Limit:  int(size),
	}
	option := &repos.FindPageOption{
		TplPrefix: tpl,
		Fn: func(context ctx.Context, i int64) interface{} {
			return dest
		},
	}
	total, list, err := o.FindPage(c, pageVo, condition, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}
	return total, list, nil
}

// Exec Exec方法封装
func Exec(c ctx.Context, o *repos.SQLRepo, tpl string, args ...interface{}) (int64, error) {
	stmt, err := o.DB.Template(tpl)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(c, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

type UpsertColumn struct {
	Name   string
	Insert bool
	Update bool
}

// MakeUpsertSQL 拼接SQL语句
func MakeUpsertSQL(table string, columns []UpsertColumn) string {
	var cols, vals, sets []string
	for _, col := range columns {
		if col.Insert {
			cols = append(cols, "`"+col.Name+"`")
			vals = append(vals, "?")
		}
	}

	for _, col := range columns {
		if col.Update {
			sets = append(sets, fmt.Sprintf("`%s`=?", col.Name))
		}
	}
	query := fmt.Sprintf("insert into %s(%s) values(%s) on duplicate key update %s", table, strings.Join(cols, ","), strings.Join(vals, ","), strings.Join(sets, ","))
	return query
}
