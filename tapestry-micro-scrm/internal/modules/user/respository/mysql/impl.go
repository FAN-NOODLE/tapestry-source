package mysql

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/user"
	"micro-scrm/internal/pkg/db"
)

type UserRepo struct {
	*repos.SQLRepo
}

type ShopUserRepo struct {
	*repos.SQLRepo
}

type EventRepo struct {
	*repos.SQLRepo
}

type OrderRepo struct {
	*repos.SQLRepo
}

func (u *UserRepo) enterpriseUsers(cc ctx.Context, enterpriseID int64, euids []int64, excludes map[int64]bool) ([]*customer.EnterpriseUser, error) {
	for _, euid := range euids {
		excludes[euid] = true
	}
	mapping := make(map[string]interface{})
	mapping["cond"] = &struct {
		EnterpriseID int64
		UserIDs      []int64
	}{
		EnterpriseID: enterpriseID,
		UserIDs:      euids,
	}
	stmt, err := u.DB.Template("t_customer/enterprise_user_by_ids", &mapping)
	if err != nil {
		return nil, err
	}

	var enterpriseUsers []struct {
		customer.EnterpriseUserBindDB
		customer.EnterpriseUserBase
		customer.EnterpriseUserLocation
		customer.EnterpriseUserAttr
		Remark   string `db:"remark"`
		DBTags   []byte `db:"tags"`
		DBGroups []byte `db:"groups"`
	}

	if err = stmt.Select(cc, &enterpriseUsers); err != nil {
		return nil, err
	}

	var bindTos []int64
	var users []*customer.EnterpriseUser
	for _, eu := range enterpriseUsers {
		if eu.BindTo != 0 && eu.BindTo != eu.ID {
			if _, ok := excludes[eu.BindTo]; !ok {
				bindTos = append(bindTos, eu.BindTo)
			}
		} else {
			users = append(users, &customer.EnterpriseUser{
				EnterpriseUserBind:     *eu.ToUser(),
				EnterpriseUserBase:     eu.EnterpriseUserBase,
				EnterpriseUserLocation: eu.EnterpriseUserLocation,
				EnterpriseUserAttr:     eu.EnterpriseUserAttr,
				Remark:                 eu.Remark,
			})
		}
	}

	if len(bindTos) == 0 {
		return users, nil
	}

	bindUsers, err := u.enterpriseUsers(cc, enterpriseID, bindTos, excludes)
	if err != nil {
		return nil, err
	}
	return append(users, bindUsers...), nil
}

// EnterpriseUsers 获取用户列表
func (u *UserRepo) EnterpriseUsers(cc ctx.Context, enterpriseID int64, euids []int64) ([]*customer.EnterpriseUser, error) {
	return u.enterpriseUsers(cc, enterpriseID, euids, make(map[int64]bool))
}

// QueryByID 通过id查询用户信息
func (u *UserRepo) QueryByID(c ctx.Context, enterpriseID int64, customerType int, id int64) (*customer.SourceUserBase, error) {
	var table string
	switch customerType {
	case customer.CustomerTypeOpen:
		table = "wx_user_open"
	case customer.CustomerTypeShop:
		table = "wx_user_shop"
	case customer.CustomerTypeExternal:
		table = "wx_user_corp_external"
	}
	var user customer.SourceUserBase
	stmt, err := u.DB.Template("t_user/select-by-id", &map[string]interface{}{
		"enterprise_id": enterpriseID,
		"table":         table,
	})
	if err != nil {
		return nil, err
	}

	if err := stmt.Get(c, &user, id); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) QueryUserCount(c ctx.Context, cond *customer.ListCond) (int, error) {
	mapping := make(map[string]interface{})
	mapping["enterprise_id"] = cond.EnterpriseID
	mapping["tag_id"] = cond.TagCond.TagIDList[0]
	template := "t_work_external_tag/select-tag-use-num"

	stmt, tmpErr := u.DB.Template(template, &mapping)
	if tmpErr != nil {
		return 0, tmpErr
	}

	var count []struct {
		UseNum int `db:"use_num"`
	}
	resErr := stmt.NamedSelect(c, &count, mapping)

	if resErr != nil {
		return 0, resErr
	}

	return count[0].UseNum, nil
}

// UpdateDescription 更新用户描述
func (u *UserRepo) UpdateDescription(c ctx.Context, enterpriseID int64, id int64, remark string) error {
	stmt, err := u.DB.Template("t_user/update-remark", &map[string]interface{}{"enterprise_id": enterpriseID})
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c, remark, time.Now().Unix(), id)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) insert(c ctx.Context, tpl string, record interface{}) (int64, error) {
	stmt, err := u.DB.Template(tpl)
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, record)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (u *UserRepo) exec(c ctx.Context, tpl string, record interface{}) (int64, error) {
	stmt, err := u.DB.Template(tpl)
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, record)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (u *UserRepo) get(c ctx.Context, tpl string, target interface{}, args ...interface{}) error {
	stmt, err := u.DB.Template(tpl)
	if err != nil {
		return err
	}
	return stmt.Get(c, target, args...)
}

func (u *UserRepo) execMulti(ctx ctx.Context, paras []map[string]interface{}) error {
	return u.DB.Tx(ctx, func(tx *sqlplus.Tx) error {
		for _, para := range paras {
			for tpl, arg := range para {
				stmt, err := tx.Template(tpl)
				if err != nil {
					return err
				}
				_, err = stmt.NamedExec(ctx, arg)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// MergeUser 合并更新用户
func (u *UserRepo) MergeUser(c ctx.Context, source *user.SourceUser) (bool, error) {
	wxuser := NewWxUser(source)
	err := wxuser.SetSource(c, u.DB)
	if err != nil {
		return false, err
	}

	var insertAppUser bool
	if wxuser.NeedInsertUnion() {
		err = wxuser.InsertUnion(c, u.DB)
	} else if !wxuser.Exists() {
		insertAppUser = true
		err = wxuser.Insert(c, u.DB)
	} else if source.UnionID == "" || !wxuser.NeedMerge() {
		err = wxuser.Update(c, u.DB)
	} else {
		err = wxuser.UpdateWithUnion(c, u.DB)
	}

	if mySqlErr, ok := err.(*mysql.MySQLError); ok && mySqlErr.Number == 1062 {
		insertAppUser = false
		if err = wxuser.SetSource(c, u.DB); err != nil {
			return false, err
		}
		if wxuser.NeedInsertUnion() {
			err = wxuser.InsertUnion(c, u.DB)
		} else if !wxuser.Exists() {
			insertAppUser = true
			err = wxuser.Insert(c, u.DB)
		} else if source.UnionID == "" || !wxuser.NeedMerge() {
			err = wxuser.Update(c, u.DB)
		} else {
			err = wxuser.UpdateWithUnion(c, u.DB)
		}
	}
	return insertAppUser, err
}

// GetRobotUserList 获取假用户
func (u *UserRepo) GetRobotUserList(c ctx.Context, limit int) ([]*user.RobotUserEntity, error) {
	var list = make([]*user.RobotUserEntity, 0)
	err := db.Query(c, u.SQLRepo, "t_robot_user/list", &list, limit)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return list, nil
}

// Query 订单查询模板方法
func (o *OrderRepo) query(c ctx.Context, tpl string, dest interface{}, args ...interface{}) error {
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

//UserAdd 实体添加
func (o *OrderRepo) add(cc ctx.Context, entity interface{}, tpl string) (int64, error) {
	stmt, err := o.DB.Template(tpl)
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(cc, entity)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (u *UserRepo) AllEterpriseUserForES(c ctx.Context, ID int64) (*user.UserEntity, error) {
	var user user.UserEntity
	stmt, err := u.DB.Template("t_user/select-by-id")
	if err != nil {
		return nil, err
	}
	errStmt := stmt.Get(c, &user, ID)
	if errStmt != nil {
		log.Info(c).Msgf("通过id查询用户信息错误, ID:%d, err:%v", ID, errStmt)
		return nil, errStmt
	}
	return &user, nil
}
